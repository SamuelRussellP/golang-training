package model

import "github.com/google/uuid"

type Transaction struct {
	TransactionId     uuid.UUID `json:"transaction_id"`
	BankId            string    `json:"bank_id"`
	TransactionStatus string    `json:"transaction_status"`
	TransactionAmount float32   `json:"transaction_amount"`
	TransactionFee    float32   `json:"transaction_fee"`
	TransactionTime   string    `json:"transaction_time"`
}
