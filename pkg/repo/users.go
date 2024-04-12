package repo

import (
	"barista/pkg/log"
	"barista/pkg/models"
	"context"
	"github.com/jackc/pgx/v5"
)

func init() {
	log.GetLog().Info("Init UsersRepo")
}

type UsersRepo interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id int32) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	DeleteByID(ctx context.Context, id int32) error
	UpdatePassword(ctx context.Context, id int32, newPassword string) error
}

type UserRepoImp struct {
	postgres *pgx.Conn
}

func NewUserRepoImp(postgres *pgx.Conn) *UserRepoImp {
	_, err := postgres.Exec(context.Background(),
		`CREATE TABLE IF NOT EXISTS users (
    			id INT PRIMARY KEY, 
    			first_name TEXT, 
    			last_name TEXT, 
    			email TEXT, 
    			password TEXT, 
    			phone BIGINT, 
    			sex INT,
    			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    			UNIQUE(email),
    			UNIQUE(phone))`)
	if err != nil {
		log.GetLog().WithError(err).WithField("table", "users").Fatal("Unable to create table")
	}
	return &UserRepoImp{postgres: postgres}
}

func (u *UserRepoImp) Create(ctx context.Context, user *models.User) error {
	_, err := u.postgres.Exec(ctx, "INSERT INTO users (id, first_name, last_name, email, password, phone, sex) VALUES ($1, $2, $3, $4, $5, $6, $7)", user.ID, user.FirstName, user.LastName, user.Email, user.Password, user.Phone, user.Sex)
	if err != nil {
		log.GetLog().Errorf("Unable to intser user. error: %v", err)
	}
	return err
}

func (u *UserRepoImp) GetByID(ctx context.Context, id int32) (*models.User, error) {
	var user models.User
	err := u.postgres.QueryRow(ctx, "SELECT id, first_name, last_name, email, password, phone, sex FROM users WHERE id = $1", id).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Phone, &user.Sex)
	if err != nil {
		log.GetLog().Errorf("Unable to get user by id. error: %v", err)
	}
	return &user, err
}

func (u *UserRepoImp) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := u.postgres.QueryRow(ctx, "SELECT id, first_name, last_name, email, password, phone, sex FROM users WHERE email = $1", email).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Phone, &user.Sex)
	if err != nil {
		log.GetLog().Errorf("Unable to get user by id. error: %v", err)
	}
	return &user, err
}

func (u *UserRepoImp) DeleteByID(ctx context.Context, id int32) error {
	_, err := u.postgres.Exec(ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		log.GetLog().Errorf("Unable to delete user by id. error: %v", err)
	}
	return err
}

func (u *UserRepoImp) UpdatePassword(ctx context.Context, id int32, newPassword string) error {
	_, err := u.postgres.Exec(ctx, "UPDATE users SET password = $1 WHERE id = $2", newPassword, id)
	if err != nil {
		log.GetLog().Errorf("Unable to update user's password. error: %v", err)
	}
	return nil
}
