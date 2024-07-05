package models

import "time"

type TransactionType int

const (
	InvalidTransaction TransactionType = iota
	Deposit
	Withdraw
	Transfer
)

type Transaction struct {
	ID          string          `json:"id"`
	SenderID    int32           `json:"sender_id"`
	ReceiverID  int32           `json:"receiver_id"`
	Amount      int64           `json:"amount"`
	Description string          `json:"description"`
	Type        TransactionType `json:"transaction_type"`
	CreatedAt   time.Time       `json:"created_at"`
}
