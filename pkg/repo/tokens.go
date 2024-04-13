package repo

import (
	"barista/pkg/log"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

func init() {
	log.GetLog().Info("Init TokensRepo")
}

type TokensRepo interface {
	Create(ctx context.Context, token string, refreshToken string, userID int32, updatedAt time.Time) error
	GetIDByTokenString(ctx context.Context, token string) (int32, error)
	// GetByTokenString(ctx context.Context, token string)
}

type TokenRepoImp struct {
	postgres *pgx.Conn
}

func NewTokenRepoImp(postgres *pgx.Conn) *TokenRepoImp {
	_, err := postgres.Exec(context.Background(),
		`CREATE TABLE IF NOT EXISTS tokens (
    			token TEXT PRIMARY KEY,
				refresh_token TEXT,
				updated_at TIMESTAMP,
				user_id	INT)`)
	if err != nil {
		panic(fmt.Sprintf("Unable to create tokens table. error: %v", err))
	}
	return &TokenRepoImp{postgres: postgres}
}

func CheckTokenExistence(postgres *pgx.Conn, token string) (bool, error) {
	var exists bool
	err := postgres.QueryRow(context.Background(),
		`SELECT EXISTS
		(SELECT 1 FROM tokens WHERE token = $1)`,
		token).Scan(&exists)
	if err != nil {
		log.GetLog().Errorf("Unable to check token existence. error: %v", err)
	}

	return exists, err
}

func (t *TokenRepoImp) Create(ctx context.Context, token string, refreshToken string, userID int32, updatedAt time.Time) error {
	_, err := t.postgres.Exec(ctx,
		`INSERT INTO tokens (token, refresh_token, updated_at, user_id)
		VALUES ($1, $2, $3, $4)`,
		token, refreshToken, updatedAt, userID)
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

// func (t *TokenRepoImp) GetByTokenString(ctx context.Context, token string) (*models.JWTToken, error) {
// 	var JWTtoken models.JWTToken
// 	err := t.postgres.QueryRow(ctx,
// 		`SELECT token, refresh_token, updated_at, user_id
// 		FROM tokens
// 		WHERE token = $1`,
// 		token).Scan(&JWTtoken.Token, &JWTtoken.RefreshToken, &JWTtoken.UpdatedAt, &JWTtoken.UserID)
// 	if err != nil {
// 		log.GetLog().Errorf("Unable to get token by it's string. error: %v", err)
// 	}

// 	return &JWTtoken, err
// }