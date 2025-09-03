package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"bankmore/internal/account/domain"
	"bankmore/internal/account/repository"
	"bankmore/internal/shared/middleware"
	"bankmore/internal/shared/models"
	"bankmore/internal/shared/utils"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AccountService interface {
	Register(request RegisterRequest) (*RegisterResponse, error)
	Login(request LoginRequest) (*LoginResponse, error)
	Deactivate(accountID, password string) error
	CreateMovement(request MovementRequest) error
	GetBalance(accountID string) (*BalanceResponse, error)
	GetBalanceByAccountNumber(accountNumber string) (*BalanceResponse, error)
	AccountExists(accountNumber string) (bool, error)
}

type accountService struct {
	repo   repository.AccountRepository
	logger *logrus.Logger
}

func NewAccountService(repo repository.AccountRepository, logger *logrus.Logger) AccountService {
	return &accountService{
		repo:   repo,
		logger: logger,
	}
}

type RegisterRequest struct {
	CPF      string `json:"cpf" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterResponse struct {
	AccountNumber string `json:"accountNumber"`
}

type LoginRequest struct {
	CPF      string `json:"cpf" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token         string `json:"token"`
	AccountNumber string `json:"accountNumber"`
}

type MovementRequest struct {
	RequestID     string  `json:"requestId" binding:"required"`
	AccountNumber string  `json:"accountNumber" binding:"required"`
	Amount        float64 `json:"amount" binding:"required"`
	Type          string  `json:"type" binding:"required"`
}

type BalanceResponse struct {
	AccountNumber string  `json:"accountNumber"`
	Balance       float64 `json:"balance"`
}

func (s *accountService) Register(request RegisterRequest) (*RegisterResponse, error) {
	cleanCPF := utils.CleanCPF(request.CPF)
	
	if !utils.ValidateCPF(cleanCPF) {
		return nil, fmt.Errorf("CPF inválido")
	}

	existingAccount, err := s.repo.GetByCPF(cleanCPF)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.WithError(err).Error("Error checking existing account")
		return nil, fmt.Errorf("erro interno do servidor")
	}
	if existingAccount != nil {
		return nil, fmt.Errorf("CPF já cadastrado")
	}

	salt, err := utils.GenerateSalt()
	if err != nil {
		s.logger.WithError(err).Error("Error generating salt")
		return nil, fmt.Errorf("erro interno do servidor")
	}

	passwordHash := utils.HashPassword(request.Password, salt)

	accountNumber, err := s.repo.GetNextAccountNumber()
	if err != nil {
		s.logger.WithError(err).Error("Error getting next account number")
		return nil, fmt.Errorf("erro interno do servidor")
	}

	account := domain.NewAccount(request.Name, cleanCPF, passwordHash, salt, accountNumber)

	if err := s.repo.Create(account); err != nil {
		s.logger.WithError(err).Error("Error creating account")
		return nil, fmt.Errorf("erro interno do servidor")
	}

	s.logger.WithFields(logrus.Fields{
		"accountId":     account.ID,
		"accountNumber": account.Number,
		"cpf":           cleanCPF,
	}).Info("Account created successfully")

	return &RegisterResponse{
		AccountNumber: strconv.Itoa(account.Number),
	}, nil
}

func (s *accountService) Login(request LoginRequest) (*LoginResponse, error) {
	cleanCPF := utils.CleanCPF(request.CPF)

	account, err := s.repo.GetByCPF(cleanCPF)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("credenciais inválidas")
		}
		s.logger.WithError(err).Error("Error getting account by CPF")
		return nil, fmt.Errorf("erro interno do servidor")
	}

	if !account.Active {
		return nil, fmt.Errorf("conta inativa")
	}

	if !utils.VerifyPassword(request.Password, account.Salt, account.PasswordHash) {
		return nil, fmt.Errorf("credenciais inválidas")
	}

	token, err := middleware.GenerateJWT(account.ID, strconv.Itoa(account.Number), account.CPF)
	if err != nil {
		s.logger.WithError(err).Error("Error generating JWT token")
		return nil, fmt.Errorf("erro interno do servidor")
	}

	s.logger.WithFields(logrus.Fields{
		"accountId":     account.ID,
		"accountNumber": account.Number,
		"cpf":           cleanCPF,
	}).Info("User logged in successfully")

	return &LoginResponse{
		Token:         token,
		AccountNumber: strconv.Itoa(account.Number),
	}, nil
}

func (s *accountService) Deactivate(accountID, password string) error {
	account, err := s.repo.GetByID(accountID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("conta não encontrada")
		}
		s.logger.WithError(err).Error("Error getting account by ID")
		return fmt.Errorf("erro interno do servidor")
	}

	if !utils.VerifyPassword(password, account.Salt, account.PasswordHash) {
		return fmt.Errorf("senha inválida")
	}

	account.Deactivate()

	if err := s.repo.Update(account); err != nil {
		s.logger.WithError(err).Error("Error deactivating account")
		return fmt.Errorf("erro interno do servidor")
	}

	s.logger.WithFields(logrus.Fields{
		"accountId":     account.ID,
		"accountNumber": account.Number,
	}).Info("Account deactivated successfully")

	return nil
}

func (s *accountService) CreateMovement(request MovementRequest) error {
	if request.Type != domain.MovementTypeCredit && request.Type != domain.MovementTypeDebit {
		return fmt.Errorf("tipo de movimentação inválido")
	}

	if request.Amount <= 0 {
		return fmt.Errorf("valor deve ser positivo")
	}

	idempotency, err := s.repo.CheckIdempotency(request.RequestID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.WithError(err).Error("Error checking idempotency")
		return fmt.Errorf("erro interno do servidor")
	}
	if idempotency != nil {
		s.logger.WithField("requestId", request.RequestID).Info("Duplicate request ignored")
		return nil
	}

	account, err := s.repo.GetByNumber(request.AccountNumber)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("conta não encontrada")
		}
		s.logger.WithError(err).Error("Error getting account by number")
		return fmt.Errorf("erro interno do servidor")
	}

	if !account.Active {
		return fmt.Errorf("conta inativa")
	}

	if request.Type == domain.MovementTypeDebit {
		balance, err := s.repo.GetBalance(account.ID)
		if err != nil {
			s.logger.WithError(err).Error("Error getting account balance")
			return fmt.Errorf("erro interno do servidor")
		}
		if balance < request.Amount {
			return fmt.Errorf("saldo insuficiente")
		}
	}

	movement := domain.NewMovement(account.ID, request.Type, request.Amount, &request.RequestID)

	if err := s.repo.CreateMovement(movement); err != nil {
		s.logger.WithError(err).Error("Error creating movement")
		return fmt.Errorf("erro interno do servidor")
	}

	requestData, _ := json.Marshal(request)
	idempotencyRecord := &domain.Idempotency{
		Key:     request.RequestID,
		Request: string(requestData),
		Result:  "SUCCESS",
	}
	s.repo.SaveIdempotency(idempotencyRecord)

	s.logger.WithFields(logrus.Fields{
		"accountId":     account.ID,
		"accountNumber": account.Number,
		"movementType":  request.Type,
		"amount":        request.Amount,
		"requestId":     request.RequestID,
	}).Info("Movement created successfully")

	return nil
}

func (s *accountService) GetBalance(accountID string) (*BalanceResponse, error) {
	account, err := s.repo.GetByID(accountID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("conta não encontrada")
		}
		s.logger.WithError(err).Error("Error getting account by ID")
		return nil, fmt.Errorf("erro interno do servidor")
	}

	balance, err := s.repo.GetBalance(accountID)
	if err != nil {
		s.logger.WithError(err).Error("Error getting account balance")
		return nil, fmt.Errorf("erro interno do servidor")
	}

	return &BalanceResponse{
		AccountNumber: strconv.Itoa(account.Number),
		Balance:       balance,
	}, nil
}

func (s *accountService) GetBalanceByAccountNumber(accountNumber string) (*BalanceResponse, error) {
	account, err := s.repo.GetByNumber(accountNumber)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		s.logger.WithError(err).Error("Error getting account by number")
		return nil, fmt.Errorf("erro interno do servidor")
	}

	balance, err := s.repo.GetBalance(account.ID)
	if err != nil {
		s.logger.WithError(err).Error("Error getting account balance")
		return nil, fmt.Errorf("erro interno do servidor")
	}

	return &BalanceResponse{
		AccountNumber: strconv.Itoa(account.Number),
		Balance:       balance,
	}, nil
}

func (s *accountService) AccountExists(accountNumber string) (bool, error) {
	_, err := s.repo.GetByNumber(accountNumber)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		s.logger.WithError(err).Error("Error checking account existence")
		return false, fmt.Errorf("erro interno do servidor")
	}
	return true, nil
}
