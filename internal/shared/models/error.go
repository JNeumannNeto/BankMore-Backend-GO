package models

type ErrorResponse struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

const (
	ErrorInvalidDocument     = "INVALID_DOCUMENT"
	ErrorUserUnauthorized    = "USER_UNAUTHORIZED"
	ErrorInvalidAccount      = "INVALID_ACCOUNT"
	ErrorInactiveAccount     = "INACTIVE_ACCOUNT"
	ErrorInvalidValue        = "INVALID_VALUE"
	ErrorInvalidType         = "INVALID_TYPE"
	ErrorInsufficientBalance = "INSUFFICIENT_BALANCE"
	ErrorInvalidAmount       = "INVALID_AMOUNT"
	ErrorInvalidTransfer     = "INVALID_TRANSFER"
	ErrorAccountNotFound     = "ACCOUNT_NOT_FOUND"
	ErrorInvalidOperation    = "INVALID_OPERATION"
	ErrorInvalidArgument     = "INVALID_ARGUMENT"
	ErrorInternalError       = "INTERNAL_ERROR"
	ErrorInvalidData         = "INVALID_DATA"
)
