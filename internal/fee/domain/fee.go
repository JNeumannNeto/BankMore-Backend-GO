package domain

import (
	"time"

	"github.com/google/uuid"
)

type Fee struct {
	ID          string    `json:"id" gorm:"column:idtarifa;primaryKey"`
	AccountID   string    `json:"accountId" gorm:"column:idcontacorrente"`
	Date        time.Time `json:"date" gorm:"column:datamovimento"`
	Amount      float64   `json:"amount" gorm:"column:valor"`
	Type        string    `json:"type" gorm:"-"`
	Description string    `json:"description" gorm:"-"`
	RequestID   string    `json:"requestId" gorm:"-"`
}

func (Fee) TableName() string {
	return "tarifa"
}

func NewFee(accountID string, amount float64) *Fee {
	return &Fee{
		ID:        uuid.New().String(),
		AccountID: accountID,
		Date:      time.Now(),
		Amount:    amount,
	}
}

const (
	FeeTypeTransfer = "TRANSFER"
)
