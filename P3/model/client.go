package model

type EndpointRequest struct {
	TxHash   string `json:"tx_hash"`
	TxStatus string `json:"tx_status"`
}
