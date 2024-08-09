package model

type Response struct {
	Data    TransactionHash `json:"data"`
	Status  string          `json:"status"`
	Message string          `json:"message"`
}
