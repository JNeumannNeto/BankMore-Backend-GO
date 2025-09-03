package models

type Result[T any] struct {
	IsSuccess    bool   `json:"isSuccess"`
	Data         T      `json:"data,omitempty"`
	ErrorType    string `json:"errorType,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

func Success[T any](data T) Result[T] {
	return Result[T]{
		IsSuccess: true,
		Data:      data,
	}
}

func Error[T any](errorType, errorMessage string) Result[T] {
	return Result[T]{
		IsSuccess:    false,
		ErrorType:    errorType,
		ErrorMessage: errorMessage,
	}
}
