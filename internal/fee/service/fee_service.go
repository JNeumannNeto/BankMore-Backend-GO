package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"bankmore/internal/fee/domain"
	"bankmore/internal/fee/repository"
	"bankmore/internal/shared/kafka"

	"github.com/sirupsen/logrus"
)

type FeeService interface {
	GetFeesByAccountNumber(accountNumber string) ([]domain.Fee, error)
	GetFeeByID(id int) (*domain.Fee, error)
	HandleTransferEvent(event kafka.TransferEvent) error
}

type feeService struct {
	repo   repository.FeeRepository
	logger *logrus.Logger
}

func NewFeeService(repo repository.FeeRepository, logger *logrus.Logger) FeeService {
	return &feeService{
		repo:   repo,
		logger: logger,
	}
}

func (s *feeService) GetFeesByAccountNumber(accountNumber string) ([]domain.Fee, error) {
	fees, err := s.repo.GetByAccountNumber(accountNumber)
	if err != nil {
		s.logger.WithError(err).Error("Error getting fees by account number")
		return nil, fmt.Errorf("erro interno do servidor")
	}

	for i := range fees {
		fees[i].Type = domain.FeeTypeTransfer
		fees[i].Description = fmt.Sprintf("Taxa de transferência")
		fees[i].RequestID = ""
	}

	return fees, nil
}

func (s *feeService) GetFeeByID(id int) (*domain.Fee, error) {
	fee, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.WithError(err).Error("Error getting fee by ID")
		return nil, err
	}

	fee.Type = domain.FeeTypeTransfer
	fee.Description = "Taxa de transferência"

	return fee, nil
}

func (s *feeService) HandleTransferEvent(event kafka.TransferEvent) error {
	feeAmount := s.getTransferFeeAmount()

	fee := domain.NewFee(event.OriginAccountID, feeAmount)

	if err := s.repo.Create(fee); err != nil {
		s.logger.WithError(err).Error("Error creating fee")
		return fmt.Errorf("erro ao criar tarifa")
	}

	if err := s.debitFeeFromAccount(event.OriginAccountID, feeAmount, event.RequestID); err != nil {
		s.logger.WithError(err).Error("Error debiting fee from account")
		return fmt.Errorf("erro ao debitar tarifa da conta")
	}

	s.logger.WithFields(logrus.Fields{
		"feeId":           fee.ID,
		"accountId":       event.OriginAccountID,
		"amount":          feeAmount,
		"transferId":      event.TransferID,
		"requestId":       event.RequestID,
	}).Info("Transfer fee processed successfully")

	return nil
}

func (s *feeService) getTransferFeeAmount() float64 {
	feeAmountStr := os.Getenv("TRANSFER_FEE_AMOUNT")
	if feeAmountStr == "" {
		return 2.00
	}

	feeAmount, err := strconv.ParseFloat(feeAmountStr, 64)
	if err != nil {
		s.logger.WithError(err).Error("Error parsing transfer fee amount")
		return 2.00
	}

	return feeAmount
}

func (s *feeService) debitFeeFromAccount(accountID string, amount float64, requestID string) error {
	accountAPIURL := os.Getenv("ACCOUNT_API_URL")
	if accountAPIURL == "" {
		accountAPIURL = "http://localhost:8001"
	}

	accountNumber, err := s.getAccountNumberByID(accountID)
	if err != nil {
		return err
	}

	request := map[string]interface{}{
		"requestId":     requestID + "-fee",
		"accountNumber": accountNumber,
		"amount":        amount,
		"type":          "D",
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/api/account/movement", accountAPIURL)
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

func (s *feeService) getAccountNumberByID(accountID string) (string, error) {
	var accountNumber int
	err := s.repo.(*feeRepository).db.Table("contacorrente").
		Select("numero").
		Where("idcontacorrente = ?", accountID).
		Scan(&accountNumber).Error
	if err != nil {
		return "", err
	}
	return strconv.Itoa(accountNumber), nil
}

type feeRepository struct {
	db interface{}
}
