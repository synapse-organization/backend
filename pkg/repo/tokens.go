package repo

import (
	"barista/pkg/log"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

func init() {
	log.GetLog().Info("Init TokensRepo")
}

type TokensRepo interface {
	Create(ctx context.Context, tokenID int32, token string, userID int32, expiredAt time.Time) error
	GetIDByTokenString(ctx context.Context, token string) (int32, error)
	DeleteByID(ctx context.Context, tokenID int32) error
}

type TokenRepoImp struct {
	postgres *pgxpool.Pool
}

func NewTokenRepoImp(postgres *pgxpool.Pool) *TokenRepoImp {
	_, err := postgres.Exec(context.Background(),
		`CREATE TABLE IF NOT EXISTS tokens (
				token_id INT PRIMARY KEY,
    			token TEXT,
				expired_at TIMESTAMP,
				user_id	INT)`)
	if err != nil {
		panic(fmt.Sprintf("Unable to create tokens table. error: %v", err))
	}
	return &TokenRepoImp{postgres: postgres}
}

func CheckTokenExistence(postgres *pgxpool.Pool, tokenID int32) (bool, error) {
	var exists bool
	err := postgres.QueryRow(context.Background(),
		`SELECT EXISTS
		(SELECT 1 FROM tokens WHERE token_id = $1)`,
		tokenID).Scan(&exists)
	if err != nil {
		log.GetLog().Errorf("Unable to check token existence. error: %v", err)
	}

	return exists, err
}

func DeleteByID(postgres *pgxpool.Pool, tokenID int32) error {
	_, err := postgres.Exec(context.Background(),
		`DELETE FROM tokens
		WHERE token_id = $1`, tokenID)
	if err != nil {
		log.GetLog().Errorf("Unable to delete token by ID. error: %v", err)
	}

	return err
}

func (t *TokenRepoImp) Create(ctx context.Context, tokenID int32, token string, userID int32, expiredAt time.Time) error {
	_, err := t.postgres.Exec(ctx,
		`INSERT INTO tokens (token_id, token, expired_at, user_id)
		VALUES ($1, $2, $3, $4)`,
		tokenID, token, expiredAt, userID)
	if err != nil {
		log.GetLog().Errorf("Unable to insert new token. error: %v", err)
	}

	return nil
}

func (t *TokenRepoImp) GetIDByTokenString(ctx context.Context, token string) (int32, error) {
	var userID int32
	err := t.postgres.QueryRow(ctx,
		`SELECT user_id
		FROM tokens
		WHERE token = $1`, token).Scan(&userID)
	if err != nil {
		log.GetLog().Errorf("Unable to get user ID by tokens string. error: %v", err)
	}

	return userID, err
}

func (t *TokenRepoImp) DeleteByID(ctx context.Context, tokenID int32) error {
	_, err := t.postgres.Exec(ctx,
		`DELETE FROM tokens
		WHERE token_id = $1`, tokenID)
	if err != nil {
		log.GetLog().Errorf("Unable to delete token by ID. error: %v", err)
	}

	return err
}
