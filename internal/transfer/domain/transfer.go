package domain

import (
	"time"

	"github.com/google/uuid"
)

type Transfer struct {
	ID                    string     `json:"id" gorm:"column:idtransferencia;primaryKey"`
	OriginAccountID       string     `json:"originAccountId" gorm:"column:idcontacorrente_origem"`
	DestinationAccountID  string     `json:"destinationAccountId" gorm:"column:idcontacorrente_destino"`
	Date                  time.Time  `json:"date" gorm:"column:datamovimento"`
	Amount                float64    `json:"amount" gorm:"column:valor"`
	Status                int        `json:"status" gorm:"column:status"`
	CompletionDate        *time.Time `json:"completionDate" gorm:"column:data_conclusao"`
	Description           string     `json:"description" gorm:"column:descricao"`
	IdempotencyKey        *string    `json:"idempotencyKey" gorm:"column:idempotencia_key"`
}

func (Transfer) TableName() string {
	return "transferencia"
}

func NewTransfer(originAccountID, destinationAccountID string, amount float64, description string, idempotencyKey *string) *Transfer {
	return &Transfer{
		ID:                   uuid.New().String(),
		OriginAccountID:      originAccountID,
		DestinationAccountID: destinationAccountID,
		Date:                 time.Now(),
		Amount:               amount,
		Status:               TransferStatusPending,
		Description:          description,
		IdempotencyKey:       idempotencyKey,
	}
}

func (t *Transfer) Complete() {
	t.Status = TransferStatusCompleted
	now := time.Now()
	t.CompletionDate = &now
}

func (t *Transfer) Fail() {
	t.Status = TransferStatusFailed
	now := time.Now()
	t.CompletionDate = &now
}

const (
	TransferStatusPending   = 0
	TransferStatusCompleted = 1
	TransferStatusFailed    = 2
)
