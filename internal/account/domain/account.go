package domain

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID           string    `json:"id" gorm:"column:idcontacorrente;primaryKey"`
	Number       int       `json:"number" gorm:"column:numero;unique"`
	Name         string    `json:"name" gorm:"column:nome"`
	CPF          string    `json:"cpf" gorm:"column:cpf;unique"`
	Active       bool      `json:"active" gorm:"column:ativo"`
	PasswordHash string    `json:"-" gorm:"column:senha"`
	Salt         string    `json:"-" gorm:"column:salt"`
	CreatedAt    time.Time `json:"createdAt" gorm:"-"`
	UpdatedAt    time.Time `json:"updatedAt" gorm:"-"`
}

func (Account) TableName() string {
	return "contacorrente"
}

func NewAccount(name, cpf, passwordHash, salt string, number int) *Account {
	return &Account{
		ID:           uuid.New().String(),
		Number:       number,
		Name:         name,
		CPF:          cpf,
		Active:       true,
		PasswordHash: passwordHash,
		Salt:         salt,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

func (a *Account) Deactivate() {
	a.Active = false
	a.UpdatedAt = time.Now()
}

func (a *Account) Activate() {
	a.Active = true
	a.UpdatedAt = time.Now()
}

type Movement struct {
	ID              string    `json:"id" gorm:"column:idmovimento;primaryKey"`
	AccountID       string    `json:"accountId" gorm:"column:idcontacorrente"`
	Date            time.Time `json:"date" gorm:"column:datamovimento"`
	Type            string    `json:"type" gorm:"column:tipomovimento"`
	Amount          float64   `json:"amount" gorm:"column:valor"`
	IdempotencyKey  *string   `json:"idempotencyKey" gorm:"column:idempotencia_key"`
}

func (Movement) TableName() string {
	return "movimento"
}

func NewMovement(accountID, movementType string, amount float64, idempotencyKey *string) *Movement {
	return &Movement{
		ID:             uuid.New().String(),
		AccountID:      accountID,
		Date:           time.Now(),
		Type:           movementType,
		Amount:         amount,
		IdempotencyKey: idempotencyKey,
	}
}

type Idempotency struct {
	Key     string `json:"key" gorm:"column:chave_idempotencia;primaryKey"`
	Request string `json:"request" gorm:"column:requisicao"`
	Result  string `json:"result" gorm:"column:resultado"`
}

func (Idempotency) TableName() string {
	return "idempotencia"
}

const (
	MovementTypeCredit = "C"
	MovementTypeDebit  = "D"
)
