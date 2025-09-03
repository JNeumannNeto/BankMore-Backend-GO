package repository

import (
	"bankmore/internal/fee/domain"
	"strconv"

	"gorm.io/gorm"
)

type FeeRepository interface {
	Create(fee *domain.Fee) error
	GetByID(id int) (*domain.Fee, error)
	GetByAccountNumber(accountNumber string) ([]domain.Fee, error)
}

type feeRepository struct {
	db *gorm.DB
}

func NewFeeRepository(db *gorm.DB) FeeRepository {
	return &feeRepository{db: db}
}

func (r *feeRepository) Create(fee *domain.Fee) error {
	return r.db.Create(fee).Error
}

func (r *feeRepository) GetByID(id int) (*domain.Fee, error) {
	var fee domain.Fee
	err := r.db.Where("idtarifa = ?", id).First(&fee).Error
	if err != nil {
		return nil, err
	}

	fee.Type = domain.FeeTypeTransfer
	fee.Description = "Taxa de transferência"

	return &fee, nil
}

func (r *feeRepository) GetByAccountNumber(accountNumber string) ([]domain.Fee, error) {
	var fees []domain.Fee
	
	err := r.db.Table("tarifa t").
		Select("t.idtarifa, t.idcontacorrente, t.datamovimento, t.valor").
		Joins("JOIN contacorrente c ON t.idcontacorrente = c.idcontacorrente").
		Where("c.numero = ?", accountNumber).
		Order("t.datamovimento DESC").
		Find(&fees).Error

	if err != nil {
		return nil, err
	}

	for i := range fees {
		fees[i].Type = domain.FeeTypeTransfer
		fees[i].Description = "Taxa de transferência"
	}

	return fees, nil
}
