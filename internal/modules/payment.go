package modules

import (
	"barista/pkg/models"
	"barista/pkg/repo"
	"context"
)

type PaymentHandler struct {
	PaymentRepo repo.Transaction
	UserRepo    repo.UsersRepo
}

func (h PaymentHandler) Transfer(ctx context.Context, userID int32, r *models.RequestTransfer) error {
	_, err := h.PaymentRepo.Create(ctx, &models.Transaction{
		SenderID:   userID,
		ReceiverID: r.To,
		Amount:     r.Amount,
		Type:       models.Transfer,
	})
	return err
}

func (h PaymentHandler) Deposit(ctx context.Context, userID int32, r *models.RequestDeposit) error {
	_, err := h.PaymentRepo.Create(ctx, &models.Transaction{
		SenderID:   userID,
		ReceiverID: userID,
		Amount:     r.Amount,
		Type:       models.Deposit,
	})
	return err
}

func (h PaymentHandler) Withdraw(ctx context.Context, userID int32, r *models.RequestWithdraw) error {
	_, err := h.PaymentRepo.Create(ctx, &models.Transaction{
		SenderID:   userID,
		ReceiverID: userID,
		Amount:     r.Amount,
		Type:       models.Withdraw,
	})
	return err
}

func (h PaymentHandler) Balance(ctx context.Context, userID int32) int64 {
	balance, err := h.UserRepo.GetBalance(ctx, userID)
	if err != nil {
		return 0
	}
	return balance
}

func (h PaymentHandler) TransactionsList(ctx context.Context, userID int32) ([]models.Transaction, error) {
	return h.PaymentRepo.GetBySenderID(ctx, userID)
}
