package modules

import (
	"barista/pkg/models"
	"barista/pkg/repo"
	"context"
)

type PaymentHandler struct {
	PaymentRepo repo.Transaction
}

func (h PaymentHandler) Transfer(ctx context.Context, userID int32, r *models.RequestTransfer) error {
	return h.PaymentRepo.Create(ctx, &models.Transaction{
		SenderID:   userID,
		ReceiverID: r.To,
		Amount:     r.Amount,
		Type:       models.Transfer,
	})
}

func (h PaymentHandler) Deposit(ctx context.Context, userID int32, r *models.RequestDeposit) error {
	return h.PaymentRepo.Create(ctx, &models.Transaction{
		SenderID:   userID,
		ReceiverID: userID,
		Amount:     r.Amount,
		Type:       models.Deposit,
	})
}

func (h PaymentHandler) Withdraw(ctx context.Context, userID int32, r *models.RequestWithdraw) error {
	return h.PaymentRepo.Create(ctx, &models.Transaction{
		SenderID:   userID,
		ReceiverID: userID,
		Amount:     r.Amount,
		Type:       models.Withdraw,
	})
}
