package model

type TransactionRequest struct {
	Symbol    string `json:"symbol" binding:"required"`
	Price     int    `json:"price" binding:"required"`
	Timestamp int    `json:"timestamp,omitempty"`
}

type TransactionHash struct {
	TxHash string `json:"tx_hash"`
}

type TransactionStatus struct {
	TxStatus string `json:"tx_status"`
}
