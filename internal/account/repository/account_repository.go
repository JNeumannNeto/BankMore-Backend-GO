package repository

import (
	"bankmore/internal/account/domain"
	"strconv"

	"gorm.io/gorm"
)

type AccountRepository interface {
	Create(account *domain.Account) error
	GetByCPF(cpf string) (*domain.Account, error)
	GetByID(id string) (*domain.Account, error)
	GetByNumber(number string) (*domain.Account, error)
	Update(account *domain.Account) error
	GetBalance(accountID string) (float64, error)
	CreateMovement(movement *domain.Movement) error
	GetNextAccountNumber() (int, error)
	CheckIdempotency(key string) (*domain.Idempotency, error)
	SaveIdempotency(idempotency *domain.Idempotency) error
}

type accountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) Create(account *domain.Account) error {
	return r.db.Create(account).Error
}

func (r *accountRepository) GetByCPF(cpf string) (*domain.Account, error) {
	var account domain.Account
	err := r.db.Where("cpf = ?", cpf).First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *accountRepository) GetByID(id string) (*domain.Account, error) {
	var account domain.Account
	err := r.db.Where("idcontacorrente = ?", id).First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *accountRepository) GetByNumber(number string) (*domain.Account, error) {
	var account domain.Account
	err := r.db.Where("numero = ?", number).First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *accountRepository) Update(account *domain.Account) error {
	return r.db.Save(account).Error
}

func (r *accountRepository) GetBalance(accountID string) (float64, error) {
	var balance float64
	err := r.db.Model(&domain.Movement{}).
		Select("COALESCE(SUM(CASE WHEN tipomovimento = 'C' THEN valor ELSE -valor END), 0)").
		Where("idcontacorrente = ?", accountID).
		Scan(&balance).Error
	return balance, err
}

func (r *accountRepository) CreateMovement(movement *domain.Movement) error {
	return r.db.Create(movement).Error
}

func (r *accountRepository) GetNextAccountNumber() (int, error) {
	var maxNumber int
	err := r.db.Model(&domain.Account{}).
		Select("COALESCE(MAX(numero), 100000)").
		Scan(&maxNumber).Error
	if err != nil {
		return 0, err
	}
	return maxNumber + 1, nil
}

func (r *accountRepository) CheckIdempotency(key string) (*domain.Idempotency, error) {
	var idempotency domain.Idempotency
	err := r.db.Where("chave_idempotencia = ?", key).First(&idempotency).Error
	if err != nil {
		return nil, err
	}
	return &idempotency, nil
}

func (r *accountRepository) SaveIdempotency(idempotency *domain.Idempotency) error {
	return r.db.Create(idempotency).Error
}
