package repo

import (
	"barista/pkg/errors"
	"barista/pkg/log"
	"barista/pkg/models"
	"barista/pkg/utils"
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Transaction interface {
	Create(ctx context.Context, transaction *models.Transaction) error
	GetByID(ctx context.Context, id string) (*models.Transaction, error)
	GetBySenderID(ctx context.Context, senderID int32) ([]models.Transaction, error)
	GetByReceiverID(ctx context.Context, receiverID int32) ([]models.Transaction, error)
	GetBySenderAndReceiverID(ctx context.Context, senderID int32, receiverID int32) ([]models.Transaction, error)
	GetBySenderOrReceiverID(ctx context.Context, accountID int32) ([]models.Transaction, error)
}

type TransactionImp struct {
	postgres *pgxpool.Pool
}

func NewTransactionImp(postgres *pgxpool.Pool) Transaction {
	_, err := postgres.Exec(context.Background(),
		`CREATE TABLE IF NOT EXISTS transactions (
			id TEXT PRIMARY KEY,
			sender_id INT,
			receiver_id INT,
			amount INT,
			description TEXT,
			transaction_type INT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`)

	if err != nil {
		log.GetLog().WithError(err).WithField("table", "transactions").Fatal("Unable to create table")
	}

	return &TransactionImp{postgres: postgres}
}

func (t *TransactionImp) Create(ctx context.Context, transaction *models.Transaction) (e error) {
	transaction.ID = utils.GenerateRandomStr(12)
	// transactions
	tx, e := t.postgres.BeginTx(ctx, pgx.TxOptions{})
	if e != nil {
		return
	}
	defer func() {
		if e != nil {
			tx.Rollback(ctx)
			return
		}
	}()

	// get sender balance
	var senderBalance int64
	row := tx.QueryRow(ctx, "SELECT balance FROM users WHERE id = $1", transaction.SenderID)
	e = row.Scan(&senderBalance)
	if e != nil {
		return
	}

	if senderBalance < transaction.Amount && transaction.Type != models.Deposit {
		return errors.ErrNotEnoughBalance.Error()
	}

	_, e = tx.Exec(ctx, "INSERT INTO transactions (id, sender_id, receiver_id, amount, description, transaction_type) VALUES ($1, $2, $3, $4, $5, $6)", transaction.ID, transaction.SenderID, transaction.ReceiverID, transaction.Amount, transaction.Description, transaction.Type)
	if e != nil {
		return
	}

	// update sender balance
	_, e = tx.Exec(ctx, "UPDATE users SET balance = balance - $1 WHERE id = $2", transaction.Amount, transaction.SenderID)
	if e != nil {
		return
	}

	// update receiver balance
	_, e = tx.Exec(ctx, "UPDATE users SET balance = balance + $1 WHERE id = $2", transaction.Amount, transaction.ReceiverID)
	if e != nil {
		return
	}

	e = tx.Commit(ctx)
	return e
}

func (t *TransactionImp) GetByID(ctx context.Context, id string) (transaction *models.Transaction, e error) {
	transaction = &models.Transaction{}
	e = t.postgres.QueryRow(ctx, "SELECT id, sender_id, receiver_id, amount, description, created_at, transaction_type FROM transactions WHERE id = $1", id).Scan(&transaction.ID, &transaction.SenderID, &transaction.ReceiverID, &transaction.Amount, &transaction.Description, &transaction.CreatedAt, &transaction.Type)
	return
}

func (t *TransactionImp) GetBySenderID(ctx context.Context, senderID int32) (transactions []models.Transaction, e error) {
	rows, e := t.postgres.Query(ctx, "SELECT id, sender_id, receiver_id, amount, description, created_at, transaction_type FROM transactions WHERE sender_id = $1", senderID)
	if e != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var transaction models.Transaction
		e = rows.Scan(&transaction.ID, &transaction.SenderID, &transaction.ReceiverID, &transaction.Amount, &transaction.Description, &transaction.CreatedAt, &transaction.Type)
		if e != nil {
			return
		}
		transactions = append(transactions, transaction)
	}
	return
}

func (t *TransactionImp) GetByReceiverID(ctx context.Context, receiverID int32) (transactions []models.Transaction, e error) {
	rows, e := t.postgres.Query(ctx, "SELECT id, sender_id, receiver_id, amount, description, created_at, transaction_type FROM transactions WHERE receiver_id = $1", receiverID)
	if e != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var transaction models.Transaction
		e = rows.Scan(&transaction.ID, &transaction.SenderID, &transaction.ReceiverID, &transaction.Amount, &transaction.Description, &transaction.CreatedAt, &transaction.Type)
		if e != nil {
			return
		}
		transactions = append(transactions, transaction)
	}
	return
}

func (t *TransactionImp) GetBySenderAndReceiverID(ctx context.Context, senderID int32, receiverID int32) (transactions []models.Transaction, e error) {
	rows, e := t.postgres.Query(ctx, "SELECT id, sender_id, receiver_id, amount, description, created_at, transaction_type FROM transactions WHERE sender_id = $1 AND receiver_id = $2", senderID, receiverID)
	if e != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var transaction models.Transaction
		e = rows.Scan(&transaction.ID, &transaction.SenderID, &transaction.ReceiverID, &transaction.Amount, &transaction.Description, &transaction.CreatedAt, &transaction.Type)
		if e != nil {
			return
		}
		transactions = append(transactions, transaction)
	}
	return
}

func (t *TransactionImp) GetBySenderOrReceiverID(ctx context.Context, accountID int32) (transactions []models.Transaction, e error) {
	rows, e := t.postgres.Query(ctx, "SELECT id, sender_id, receiver_id, amount, description, created_at, transaction_type FROM transactions WHERE sender_id = $1 OR receiver_id = $1", accountID)
	if e != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var transaction models.Transaction
		e = rows.Scan(&transaction.ID, &transaction.SenderID, &transaction.ReceiverID, &transaction.Amount, &transaction.Description, &transaction.CreatedAt, &transaction.Type)
		if e != nil {
			return
		}
		transactions = append(transactions, transaction)
	}
	return
}
