package repository

import (
	"bankmore/internal/transfer/domain"

	"gorm.io/gorm"
)

type TransferRepository interface {
	Create(transfer *domain.Transfer) error
	GetByID(id string) (*domain.Transfer, error)
	Update(transfer *domain.Transfer) error
	GetByAccountID(accountID string) ([]domain.Transfer, error)
}

type transferRepository struct {
	db *gorm.DB
}

func NewTransferRepository(db *gorm.DB) TransferRepository {
	return &transferRepository{db: db}
}

func (r *transferRepository) Create(transfer *domain.Transfer) error {
	return r.db.Create(transfer).Error
}

func (r *transferRepository) GetByID(id string) (*domain.Transfer, error) {
	var transfer domain.Transfer
	err := r.db.Where("idtransferencia = ?", id).First(&transfer).Error
	if err != nil {
		return nil, err
	}
	return &transfer, nil
}

func (r *transferRepository) Update(transfer *domain.Transfer) error {
	return r.db.Save(transfer).Error
}

func (r *transferRepository) GetByAccountID(accountID string) ([]domain.Transfer, error) {
	var transfers []domain.Transfer
	err := r.db.Where("idcontacorrente_origem = ? OR idcontacorrente_destino = ?", accountID, accountID).
		Order("datamovimento DESC").
		Find(&transfers).Error
	return transfers, err
}
