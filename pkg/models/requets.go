package models

type RequestTransfer struct {
	To     int32 `json:"to"`
	Amount int64 `json:"amount"`
}

type RequestWithdraw struct {
	To     string `json:"to"`
	Amount int64  `json:"amount"`
}

type RequestDeposit struct {
	Amount int64 `json:"amount"`
}
