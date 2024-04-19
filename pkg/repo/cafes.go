package repo

import (
	"barista/pkg/log"
	"barista/pkg/models"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"math/rand"
	"strings"
)

func init() {
	log.GetLog().Info("Init CafesRepo")
}

type CafesRepo interface {
	Create(ctx context.Context, cafe *models.Cafe) error
	GetByID(ctx context.Context, id int32) (*models.Cafe, error)
	SearchCafe(ctx context.Context, name string, address string, location string, category string) ([]models.Cafe, error)
}

type CafesRepoImp struct {
	postgres *pgxpool.Pool
}

func NewCafeRepoImp(postgres *pgxpool.Pool) *CafesRepoImp {
	_, err := postgres.Exec(context.Background(),
		`CREATE TABLE IF NOT EXISTS cafes (
				id INTEGER PRIMARY KEY,
				owner_id INTEGER,
				name TEXT,
				description TEXT,
				opening_time INTEGER,
				closing_time INTEGER,
				capacity INTEGER,
				phone_number TEXT,
				email TEXT,
				location TEXT,
				province INTEGER,
				city INTEGER,
				address TEXT,
				categories TEXT,
				FOREIGN KEY (owner_id) REFERENCES users(id)
			);`)

	if err != nil {
		log.GetLog().WithError(err).WithField("table", "cafes").Fatal("Unable to create table")
	}

	_, err = postgres.Exec(context.Background(), `INSERT INTO cafes (id, owner_id, name, description, opening_time, closing_time, capacity, phone_number, email, location, province, city, address, categories)
			VALUES 
			(1, 12, 'Cafe One', 'A stylish cafe with a focus on specialty coffees and homemade desserts.', 8, 10, 25, '+12345678901', 'cafe_one@example.com', '123 Oak Street', 1, 2, '123 Oak Street, Los Angeles, CA', 'Coffee, Desserts');
	`)
	if err != nil {
		log.GetLog().Errorf("Unable to insert cafes. error: %v", err)
	}

	return &CafesRepoImp{postgres: postgres}
}

func (c *CafesRepoImp) Create(ctx context.Context, cafe *models.Cafe) error {
	cafe.ID = rand.Int31()
	categories := ""
	for i, category := range cafe.Categories {
		categories += string(category)
		if i != len(cafe.Categories)-1 {
			categories += ","
		}
	}
	_, err := c.postgres.Exec(ctx, "INSERT INTO cafes (id, owner_id, name, description, opening_time, closing_time, capacity, phone_number, email, province, city, address, location, categories) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)", cafe.ID, cafe.OwnerID, cafe.Name, cafe.Description, cafe.OpeningTime, cafe.ClosingTime, cafe.Capacity, cafe.ContactInfo.Phone, cafe.ContactInfo.Email, cafe.ContactInfo.Province, cafe.ContactInfo.City, cafe.ContactInfo.Address, cafe.ContactInfo.Location, categories)
	if err != nil {
		log.GetLog().Errorf("Unable to insert cafe. error: %v", err)
	}
	return err
}

func (c *CafesRepoImp) GetByID(ctx context.Context, id int32) (*models.Cafe, error) {
	var cafe models.Cafe
	catetories := ""
	err := c.postgres.QueryRow(ctx, "SELECT id, owner_id, name, description, opening_time, closing_time, capacity, phone_number, email, province, city, address, location, categories FROM cafes WHERE id = $1", id).Scan(&cafe.ID, &cafe.OwnerID, &cafe.Name, &cafe.Description, &cafe.OpeningTime, &cafe.ClosingTime, &cafe.Capacity, &cafe.ContactInfo.Phone, &cafe.ContactInfo.Email, &cafe.ContactInfo.Province, &cafe.ContactInfo.City, &cafe.ContactInfo.Address, &cafe.ContactInfo.Location, &catetories)

	for _, category := range strings.Split(catetories, ",") {
		cafe.Categories = append(cafe.Categories, models.CafeCategory(category))
	}

	if err != nil {
		log.GetLog().Errorf("Unable to get cafe by id. error: %v", err)
	}
	return &cafe, err
}

func (c *CafesRepoImp) SearchCafe(ctx context.Context, name string, province string, city string, category string) ([]models.Cafe, error) {
	var cafes []models.Cafe
	list := []string{}
	query := "SELECT id, owner_id, name, description, opening_time, closing_time, capacity, phone_number, email, province, city, address, location, categories FROM cafes"

	if name != "" {
		list = append(list, "name LIKE '%"+name+"%'")
	}
	if province != "" {
		list = append(list, "province = '"+province+"'")
	}
	if city != "" {
		list = append(list, "city = '"+city+"'")
	}
	if category != "" {
		list = append(list, "categories LIKE '%"+category+"%'")
	}

	if len(list) > 0 {
		query += " WHERE " + strings.Join(list, " AND ")

	}

	rows, err := c.postgres.Query(ctx, query)
	if err != nil {
		log.GetLog().Errorf("Unable to search cafe. error: %v", err)
	}

	for rows.Next() {
		var cafe models.Cafe
		catetories := ""
		err = rows.Scan(&cafe.ID, &cafe.OwnerID, &cafe.Name, &cafe.Description, &cafe.OpeningTime, &cafe.ClosingTime, &cafe.Capacity, &cafe.ContactInfo.Phone, &cafe.ContactInfo.Email, &cafe.ContactInfo.Province, &cafe.ContactInfo.City, &cafe.ContactInfo.Address, &cafe.ContactInfo.Location, &catetories)
		if err != nil {
			log.GetLog().Errorf("Unable to scan cafe. error: %v", err)
		}
		for _, category := range strings.Split(catetories, ",") {
			cafe.Categories = append(cafe.Categories, models.CafeCategory(strings.TrimSpace(category)))
		}
		cafes = append(cafes, cafe)
	}

	return cafes, err
}
