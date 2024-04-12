package repo

import (
	"barista/pkg/log"
	"barista/pkg/models"
	"context"
	"github.com/jackc/pgx/v5"
)

func init() {
	log.GetLog().Info("Init CafesRepo")
}

type CafesRepo interface {
	Create(ctx context.Context, cafe *models.Cafe) error
	GetByID(ctx context.Context, id int32) (*models.Cafe, error)
	SearchCafe(ctx context.Context, name string, address string, location string, categories []string) ([]models.Cafe, error)
}

type CafesRepoImp struct {
	postgres *pgx.Conn
}

func NewCafeRepoImp(postgres *pgx.Conn) *CafesRepoImp {
	_, err := postgres.Exec(context.Background(),
		`CREATE TABLE IF NOT EXISTS cafes (
				id INTEGER PRIMARY KEY,
				owner_id INTEGER,
				name TEXT,
				description TEXT,
				opening_time TIMESTAMP,
				closing_time TIMESTAMP,
				capacity INTEGER,
				phone_number TEXT,
				email TEXT,
				address TEXT,
				location TEXT,
				catagoires TEXT[],
				FOREIGN KEY (owner_id) REFERENCES users(id)
			);`)
	if err != nil {
		log.GetLog().WithError(err).WithField("table", "cafes").Fatal("Unable to create table")
	}
	return &CafesRepoImp{postgres: postgres}
}

func (c *CafesRepoImp) Create(ctx context.Context, cafe *models.Cafe) error {
	_, err := c.postgres.Exec(ctx, "INSERT INTO cafes (id, owner_id, name, description, opening_time, closing_time, capacity, phone_number, email, address, location, catagoires) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)", cafe.ID, cafe.OwnerID, cafe.Name, cafe.Description, cafe.OpeningTime, cafe.ClosingTime, cafe.Capacity, cafe.ContactInfo.Phone, cafe.ContactInfo.Email, cafe.ContactInfo.Address, cafe.ContactInfo.Location, cafe.Categories)
	if err != nil {
		log.GetLog().Errorf("Unable to insert cafe. error: %v", err)
	}
	return err
}

func (c *CafesRepoImp) GetByID(ctx context.Context, id int32) (*models.Cafe, error) {
	var cafe models.Cafe
	err := c.postgres.QueryRow(ctx, "SELECT id, owner_id, name, description, opening_time, closing_time, capacity, phone_number, email, address, location, catagoires FROM cafes WHERE id = $1", id).Scan(&cafe.ID, &cafe.OwnerID, &cafe.Name, &cafe.Description, &cafe.OpeningTime, &cafe.ClosingTime, &cafe.Capacity, &cafe.ContactInfo.Phone, &cafe.ContactInfo.Email, &cafe.ContactInfo.Address, &cafe.ContactInfo.Location, &cafe.Categories)
	if err != nil {
		log.GetLog().Errorf("Unable to get cafe by id. error: %v", err)
	}
	return &cafe, err
}

func (c *CafesRepoImp) SearchCafe(ctx context.Context, name string, address string, location string, categories []string) ([]models.Cafe, error) {
	var cafes []models.Cafe
	rows, err := c.postgres.Query(ctx, "SELECT id, owner_id, name, description, opening_time, closing_time, capacity, phone_number, email, address, location, catagoires FROM cafes WHERE name LIKE $1 AND description LIKE $1 AND address LIKE $2 AND location LIKE $3 AND catagoires @> $4", name, address, location, categories)
	if err != nil {
		log.GetLog().Errorf("Unable to search cafe. error: %v", err)
	}
	for rows.Next() {
		var cafe models.Cafe
		err = rows.Scan(&cafe.ID, &cafe.OwnerID, &cafe.Name, &cafe.Description, &cafe.OpeningTime, &cafe.ClosingTime, &cafe.Capacity, &cafe.ContactInfo.Phone, &cafe.ContactInfo.Email, &cafe.ContactInfo.Address, &cafe.ContactInfo.Location, &cafe.Categories)
		if err != nil {
			log.GetLog().Errorf("Unable to scan cafe. error: %v", err)
		}
		cafes = append(cafes, cafe)
	}
	return cafes, err
}
