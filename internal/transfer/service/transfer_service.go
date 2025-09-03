package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"bankmore/internal/shared/kafka"
	"bankmore/internal/shared/models"
	"bankmore/internal/transfer/domain"
	"bankmore/internal/transfer/repository"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TransferService interface {
	CreateTransfer(request CreateTransferRequest, originAccountID string) (*models.Result[TransferResponse], error)
}

type transferService struct {
	repo     repository.TransferRepository
	producer *kafka.Producer
	logger   *logrus.Logger
}

func NewTransferService(repo repository.TransferRepository, producer *kafka.Producer, logger *logrus.Logger) TransferService {
	return &transferService{
		repo:     repo,
		producer: producer,
		logger:   logger,
	}
}

type CreateTransferRequest struct {
	RequestID                string  `json:"requestId" binding:"required"`
	DestinationAccountNumber string  `json:"destinationAccountNumber" binding:"required"`
	Amount                   float64 `json:"amount" binding:"required"`
}

type TransferResponse struct {
	TransferID string `json:"transferId"`
	Message    string `json:"message"`
}

type AccountAPIResponse struct {
	AccountNumber string  `json:"accountNumber"`
	Balance       float64 `json:"balance"`
}

func (s *transferService) CreateTransfer(request CreateTransferRequest, originAccountID string) (*models.Result[TransferResponse], error) {
	if request.Amount <= 0 {
		return &models.Result[TransferResponse]{
			IsSuccess:    false,
			ErrorType:    models.ErrorInvalidAmount,
			ErrorMessage: "Valor deve ser positivo",
		}, nil
	}

	destinationAccountID, err := s.getAccountIDByNumber(request.DestinationAccountNumber)
	if err != nil {
		s.logger.WithError(err).Error("Error getting destination account")
		return &models.Result[TransferResponse]{
			IsSuccess:    false,
			ErrorType:    models.ErrorAccountNotFound,
			ErrorMessage: "Conta de destino não encontrada",
		}, nil
	}

	if originAccountID == destinationAccountID {
		return &models.Result[TransferResponse]{
			IsSuccess:    false,
			ErrorType:    models.ErrorInvalidTransfer,
			ErrorMessage: "Não é possível transferir para a mesma conta",
		}, nil
	}

	originBalance, err := s.getAccountBalance(originAccountID)
	if err != nil {
		s.logger.WithError(err).Error("Error getting origin account balance")
		return &models.Result[TransferResponse]{
			IsSuccess:    false,
			ErrorType:    models.ErrorInternalError,
			ErrorMessage: "Erro interno do servidor",
		}, nil
	}

	if originBalance < request.Amount {
		return &models.Result[TransferResponse]{
			IsSuccess:    false,
			ErrorType:    models.ErrorInsufficientBalance,
			ErrorMessage: "Saldo insuficiente",
		}, nil
	}

	description := fmt.Sprintf("Transferência para conta %s", request.DestinationAccountNumber)
	transfer := domain.NewTransfer(originAccountID, destinationAccountID, request.Amount, description, &request.RequestID)

	if err := s.repo.Create(transfer); err != nil {
		s.logger.WithError(err).Error("Error creating transfer")
		return &models.Result[TransferResponse]{
			IsSuccess:    false,
			ErrorType:    models.ErrorInternalError,
			ErrorMessage: "Erro interno do servidor",
		}, nil
	}

	if err := s.processTransferMovements(transfer); err != nil {
		s.logger.WithError(err).Error("Error processing transfer movements")
		transfer.Fail()
		s.repo.Update(transfer)
		return &models.Result[TransferResponse]{
			IsSuccess:    false,
			ErrorType:    models.ErrorInternalError,
			ErrorMessage: "Erro ao processar transferência",
		}, nil
	}

	transfer.Complete()
	if err := s.repo.Update(transfer); err != nil {
		s.logger.WithError(err).Error("Error updating transfer status")
	}

	event := kafka.TransferEvent{
		RequestID:               request.RequestID,
		OriginAccountID:         originAccountID,
		DestinationAccountID:    destinationAccountID,
		DestinationAccountNumber: request.DestinationAccountNumber,
		Amount:                  request.Amount,
		TransferID:              transfer.ID,
	}

	if err := s.producer.PublishTransferEvent(event); err != nil {
		s.logger.WithError(err).Error("Error publishing transfer event")
	}

	s.logger.WithFields(logrus.Fields{
		"transferId":         transfer.ID,
		"originAccountId":    originAccountID,
		"destinationAccountId": destinationAccountID,
		"amount":             request.Amount,
		"requestId":          request.RequestID,
	}).Info("Transfer completed successfully")

	return &models.Result[TransferResponse]{
		IsSuccess: true,
		Data: TransferResponse{
			TransferID: transfer.ID,
			Message:    "Transferência realizada com sucesso",
		},
	}, nil
}

func (s *transferService) getAccountIDByNumber(accountNumber string) (string, error) {
	accountAPIURL := os.Getenv("ACCOUNT_API_URL")
	if accountAPIURL == "" {
		accountAPIURL = "http://localhost:8001"
	}

	url := fmt.Sprintf("%s/api/account/balance/%s", accountAPIURL, accountNumber)
	
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", errors.New("account not found")
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to get account")
	}

	var accountResp AccountAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&accountResp); err != nil {
		return "", err
	}

	return s.getAccountIDByNumberFromDB(accountNumber)
}

func (s *transferService) getAccountIDByNumberFromDB(accountNumber string) (string, error) {
	var accountID string
	err := s.repo.(*transferRepository).db.Table("contacorrente").
		Select("idcontacorrente").
		Where("numero = ?", accountNumber).
		Scan(&accountID).Error
	return accountID, err
}

func (s *transferService) getAccountBalance(accountID string) (float64, error) {
	var balance float64
	err := s.repo.(*transferRepository).db.Table("movimento").
		Select("COALESCE(SUM(CASE WHEN tipomovimento = 'C' THEN valor ELSE -valor END), 0)").
		Where("idcontacorrente = ?", accountID).
		Scan(&balance).Error
	return balance, err
}

func (s *transferService) processTransferMovements(transfer *domain.Transfer) error {
	accountAPIURL := os.Getenv("ACCOUNT_API_URL")
	if accountAPIURL == "" {
		accountAPIURL = "http://localhost:8001"
	}

	debitRequest := map[string]interface{}{
		"requestId":     transfer.ID + "-debit",
		"accountNumber": s.getAccountNumberByID(transfer.OriginAccountID),
		"amount":        transfer.Amount,
		"type":          "D",
	}

	if err := s.callAccountMovementAPI(accountAPIURL, debitRequest); err != nil {
		return fmt.Errorf("failed to debit origin account: %w", err)
	}

	creditRequest := map[string]interface{}{
		"requestId":     transfer.ID + "-credit",
		"accountNumber": s.getAccountNumberByID(transfer.DestinationAccountID),
		"amount":        transfer.Amount,
		"type":          "C",
	}

	if err := s.callAccountMovementAPI(accountAPIURL, creditRequest); err != nil {
		rollbackRequest := map[string]interface{}{
			"requestId":     transfer.ID + "-rollback",
			"accountNumber": s.getAccountNumberByID(transfer.OriginAccountID),
			"amount":        transfer.Amount,
			"type":          "C",
		}
		s.callAccountMovementAPI(accountAPIURL, rollbackRequest)
		return fmt.Errorf("failed to credit destination account: %w", err)
	}

	return nil
}

func (s *transferService) getAccountNumberByID(accountID string) string {
	var accountNumber int
	s.repo.(*transferRepository).db.Table("contacorrente").
		Select("numero").
		Where("idcontacorrente = ?", accountID).
		Scan(&accountNumber)
	return strconv.Itoa(accountNumber)
}

func (s *transferService) callAccountMovementAPI(baseURL string, request map[string]interface{}) error {
	url := fmt.Sprintf("%s/api/account/movement", baseURL)
	
	jsonData, err := json.Marshal(request)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", strings.NewReader(string(jsonData)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("account API returned status %d", resp.StatusCode)
	}

	return nil
}

type transferRepository struct {
	db *gorm.DB
}
