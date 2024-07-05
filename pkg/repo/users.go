package repo

import (
	"barista/pkg/log"
	"barista/pkg/models"
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
)

func init() {
	log.GetLog().Info("Init UsersRepo")
}

type UsersRepo interface {
	Create(ctx context.Context, user *models.User) error
	Verify(ctx context.Context, email string) error
	GetByID(ctx context.Context, id int32) (*models.User, error)
	GetByEmail(ctx context.Context, email string) ([]*models.User, error)
	GetBalance(ctx context.Context, id int32) (int64, error)
	DeleteByID(ctx context.Context, id int32) error
	UpdateFirstName(ctx context.Context, id int32, newFirstName string) error
	UpdateLastName(ctx context.Context, id int32, newLastName string) error
	UpdatePassword(ctx context.Context, id int32, newPassword string) error
	UpdateSex(ctx context.Context, id int32, newSex string) error
	UpdatePhone(ctx context.Context, id int32, newPhone int32) error
	UpdateRole(ctx context.Context, id int32, newRole int32) error
	UpdateExtraInfo(ctx context.Context, id int32, newExtraInfo map[string]interface{}) error
}

type UserRepoImp struct {
	postgres *pgxpool.Pool
}

func NewUserRepoImp(postgres *pgxpool.Pool) *UserRepoImp {
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
    			is_verified BOOLEAN DEFAULT FALSE,
    			user_role INT DEFAULT 1,
    			balance BIGINT DEFAULT 0,
    			extra_info JSONB,
    			UNIQUE(email, user_role))`)

	if err != nil {
		log.GetLog().WithError(err).WithField("table", "users").Fatal("Unable to create table")
	}

	_, err = postgres.Exec(context.Background(), `INSERT INTO users (id, first_name, last_name, email, password, phone, sex, is_verified, user_role, balance, extra_info)
			VALUES
			(31, 'Dan', 'Adams', 'dan.adams@gmail.com', '$2a$10$examplehash30', 34509845678, 1, true, 2, 1000, '{}'),
			(32, 'Alice', 'Smith', 'alice.smith@gmail.com', '$2a$10$examplehash1', 12345678901, 2, true, 2, 1000, '{}'),
			(33, 'Bob', 'Johnson', 'bob.johnson@gmail.com', '$2a$10$examplehash2', 23456789012, 1, true, 2, 1000, '{}'),
			(34, 'Carol', 'Williams', 'carol.williams@gmail.com', '$2a$10$examplehash3', 34567890123, 2, true, 2, 1000, '{}'),
			(35, 'David', 'Brown', 'david.brown@gmail.com', '$2a$10$examplehash4', 45678901234, 1, true, 2, 1000, '{}'),
			(36, 'Eve', 'Jones', 'eve.jones@gmail.com', '$2a$10$examplehash5', 56789012345, 2, true, 2, 1000, '{}'),
			(37, 'Frank', 'Garcia', 'frank.garcia@gmail.com', '$2a$10$examplehash6', 67890123456, 1, true, 2, 1000, '{}'),
			(38, 'Grace', 'Martinez', 'grace.martinez@gmail.com', '$2a$10$examplehash7', 78901234567, 2, true, 2, 1000, '{}'),
			(39, 'Hank', 'Davis', 'hank.davis@gmail.com', '$2a$10$examplehash8', 89012345678, 1, true, 2, 1000, '{}'),
			(40, 'Ivy', 'Rodriguez', 'ivy.rodriguez@gmail.com', '$2a$10$examplehash9', 90123456789, 2, true, 2, 1000, '{}'),
			(41, 'Jack', 'Martinez', 'jack.martinez@gmail.com', '$2a$10$examplehash10', 12309845678, 1, true, 2, 1000, '{}'),
			(42, 'Karen', 'Hernandez', 'karen.hernandez@gmail.com', '$2a$10$examplehash11', 23410956789, 2, true, 2, 1000, '{}'),
			(43, 'Leo', 'Lopez', 'leo.lopez@gmail.com', '$2a$10$examplehash12', 34521067890, 1, true, 2, 1000, '{}'),
			(44, 'Mia', 'Gonzalez', 'mia.gonzalez@gmail.com', '$2a$10$examplehash13', 45632178901, 2, true, 2, 1000, '{}'),
			(45, 'Nate', 'Wilson', 'nate.wilson@gmail.com', '$2a$10$examplehash14', 56743289012, 1, true, 2, 1000, '{}')
			ON CONFLICT DO NOTHING;
	`)
	if err != nil {
		log.GetLog().Errorf("Unable to insert users. error: %v", err)
	}

	return &UserRepoImp{postgres: postgres}
}

func (u *UserRepoImp) Create(ctx context.Context, user *models.User) error {
	newExtraInfo := map[string]interface{}{
		"bank_account": user.BankAccount,
		"national_id":  user.NationalID,
	}
	data, err := json.Marshal(newExtraInfo)
	if err != nil {
		log.GetLog().Errorf("Unable to marshal extra info. error: %v", err)
		return err
	}

	_, err = u.postgres.Exec(ctx, "INSERT INTO users (id, first_name, last_name, email, password, phone, sex, user_role, extra_info) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)", user.ID, user.FirstName, user.LastName, user.Email, user.Password, user.Phone, user.Sex, user.Role, data)
	if err != nil {
		log.GetLog().Errorf("Unable to intser user. error: %v", err)
	}
	return err
}

func (u *UserRepoImp) Verify(ctx context.Context, email string) error {
	_, err := u.postgres.Exec(ctx, "UPDATE users SET is_verified = TRUE WHERE email = $1", email)
	if err != nil {
		log.GetLog().Errorf("Unable to verify user. error: %v", err)
	}
	return err
}

func (u *UserRepoImp) GetByID(ctx context.Context, id int32) (*models.User, error) {
	var user models.User
	err := u.postgres.QueryRow(ctx, "SELECT id, first_name, last_name, email, password, phone, sex, user_role, balance, extra_info->>'bank_account', extra_info->>'national_id' FROM users WHERE id = $1", id).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Phone, &user.Sex, &user.Role, &user.Balance, &user.BankAccount, &user.NationalID)
	if err != nil {
		log.GetLog().Errorf("Unable to get user by id. error: %v", err)
	}
	return &user, err
}

func (u *UserRepoImp) GetByEmail(ctx context.Context, email string) ([]*models.User, error) {
	var users []*models.User
	rows, err := u.postgres.Query(ctx, "SELECT id, first_name, last_name, email, password, phone, sex, user_role, balance, is_verified, extra_info->>'bank_account', extra_info->>'national_id' FROM users WHERE email = $1", email)
	if err != nil {
		log.GetLog().Errorf("Unable to get user by email. error: %v", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user models.User
		err = rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Phone, &user.Sex, &user.Role, &user.Balance, &user.IsVerified, &user.BankAccount, &user.NationalID)
		if err != nil {
			log.GetLog().Errorf("Unable to scan user. error: %v", err)
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}
func (u *UserRepoImp) GetBalance(ctx context.Context, id int32) (int64, error) {
	var balance int64
	err := u.postgres.QueryRow(ctx, "SELECT balance FROM users WHERE id = $1", id).Scan(&balance)
	if err != nil {
		log.GetLog().Errorf("Unable to get balance by id. error: %v", err)
	}
	return balance, err
}

func (u *UserRepoImp) DeleteByID(ctx context.Context, id int32) error {
	_, err := u.postgres.Exec(ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		log.GetLog().Errorf("Unable to delete user by id. error: %v", err)
	}
	return err
}

func (u *UserRepoImp) UpdateFirstName(ctx context.Context, id int32, newFirstName string) error {
	_, err := u.postgres.Exec(ctx, "UPDATE users SET first_name = $1 WHERE id = $2", newFirstName, id)
	if err != nil {
		log.GetLog().Errorf("Unable to update user's first name. error: %v", err)
	}
	return nil
}

func (u *UserRepoImp) UpdateLastName(ctx context.Context, id int32, newLastName string) error {
	_, err := u.postgres.Exec(ctx, "UPDATE users SET last_name = $1 WHERE id = $2", newLastName, id)
	if err != nil {
		log.GetLog().Errorf("Unable to update user's last name. error: %v", err)
	}
	return nil
}

func (u *UserRepoImp) UpdatePassword(ctx context.Context, id int32, newPassword string) error {
	_, err := u.postgres.Exec(ctx, "UPDATE users SET password = $1 WHERE id = $2", newPassword, id)
	if err != nil {
		log.GetLog().Errorf("Unable to update user's password. error: %v", err)
	}
	return nil
}

func (u *UserRepoImp) UpdateSex(ctx context.Context, id int32, newSex string) error {
	_, err := u.postgres.Exec(ctx, "UPDATE users SET sex = $1 WHERE id = $2", newSex, id)
	if err != nil {
		log.GetLog().Errorf("Unable to update user's sex. error: %v", err)
	}
	return nil
}

func (u *UserRepoImp) UpdatePhone(ctx context.Context, id int32, newPhone int32) error {
	_, err := u.postgres.Exec(ctx, "UPDATE users SET phone = $1 WHERE id = $2", newPhone, id)
	if err != nil {
		log.GetLog().Errorf("Unable to update user's phone number. error: %v", err)
	}
	return nil
}

func (u *UserRepoImp) UpdateRole(ctx context.Context, id int32, newRole int32) error {
	_, err := u.postgres.Exec(ctx, "UPDATE users SET user_role = $1 WHERE id = $2", newRole, id)
	if err != nil {
		log.GetLog().Errorf("Unable to update user's role. error: %v", err)
	}
	return nil
}

func (u *UserRepoImp) UpdateExtraInfo(ctx context.Context, id int32, newExtraInfo map[string]interface{}) error {

	data, err := json.Marshal(newExtraInfo)
	if err != nil {
		log.GetLog().Errorf("Unable to marshal extra info. error: %v", err)
		return err
	}

	_, err = u.postgres.Exec(ctx, "UPDATE users SET extra_info = $1 WHERE id = $2", string(data), id)
	if err != nil {
		log.GetLog().Errorf("Unable to update user's extra info. error: %v", err)
	}
	return nil
}
