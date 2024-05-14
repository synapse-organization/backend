package repo

import (
	"barista/pkg/log"
	"barista/pkg/models"
	"context"
	"math/rand"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/cast"
)

func init() {
	log.GetLog().Info("Init CafesRepo")
}

type UpdateType string

const (
	UpdateName        UpdateType = "name"
	UpdateDescription UpdateType = "description"
	UpdateOpeningTime UpdateType = "opening_time"
	UpdateClosingTime UpdateType = "closing_time"
	UpdateCapacity    UpdateType = "capacity"
	UpdatePhoneNumber UpdateType = "phone_number"
	UpdateEmail       UpdateType = "email"
	UpdateLocation    UpdateType = "location"
	UpdateProvince    UpdateType = "province"
	UpdateCity        UpdateType = "city"
	UpdateAddress     UpdateType = "address"
	UpdateCategories  UpdateType = "categories"
	UpdateAmenities   UpdateType = "amenities"
)

type CafesRepo interface {
	Create(ctx context.Context, cafe *models.Cafe) (int32, error)
	GetByID(ctx context.Context, id int32) (*models.Cafe, error)
	SearchCafe(ctx context.Context, name string, address string, location string, category string) ([]models.Cafe, error)
	GetByCafeIDs(ctx context.Context, ids []int32) ([]models.Cafe, error)
	Update(ctx context.Context, id int32, updateType UpdateType, value interface{}) error
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
				phone_number BIGINT,
				email TEXT,
				location TEXT,
				province INTEGER,
				city INTEGER,
				address TEXT,
				categories TEXT,
				amenities TEXT,
				FOREIGN KEY (owner_id) REFERENCES users(id)
			);`)

	if err != nil {
		log.GetLog().WithError(err).WithField("table", "cafes").Fatal("Unable to create table")
	}

	// _, err = postgres.Exec(context.Background(), `INSERT INTO cafes (id, owner_id, name, description, opening_time, closing_time, capacity, phone_number, email, location, province, city, address, categories, amenities)
	// 		VALUES
	// 		(1, 12, 'Cafe One', 'A stylish cafe with a focus on specialty coffees and homemade desserts.', 8, 10, 25, '+12345678901', 'cafe_one@example.com', '112', 1, 2, '123 Oak Street, Los Angeles, CA', 'Coffee,Desserts', 'وای فای');
	// `)
	// if err != nil {
	// 	log.GetLog().Errorf("Unable to insert cafes. error: %v", err)
	// }

	return &CafesRepoImp{postgres: postgres}
}

func (c *CafesRepoImp) Create(ctx context.Context, cafe *models.Cafe) (int32, error) {
	cafe.ID = rand.Int31()
	categories := ""
	for i, category := range cafe.Categories {
		categories += string(category)
		if i != len(cafe.Categories)-1 {
			categories += ","
		}
	}

	amenities := ""
	for i, amenity := range cafe.Amenities {
		amenities += string(amenity)
		if i != len(cafe.Amenities)-1 {
			amenities += ","
		}
	}

	_, err := c.postgres.Exec(ctx, "INSERT INTO cafes (id, owner_id, name, description, opening_time, closing_time, capacity, phone_number, email, province, city, address, location, categories, amenities) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)", cafe.ID, cafe.OwnerID, cafe.Name, cafe.Description, cafe.OpeningTime, cafe.ClosingTime, cafe.Capacity, cafe.ContactInfo.Phone, cafe.ContactInfo.Email, cafe.ContactInfo.Province, cafe.ContactInfo.City, cafe.ContactInfo.Address, cafe.ContactInfo.Location, categories, amenities)
	if err != nil {
		log.GetLog().Errorf("Unable to insert cafe. error: %v", err)
	}
	return cafe.ID, err
}

func (c *CafesRepoImp) GetByID(ctx context.Context, id int32) (*models.Cafe, error) {
	var cafe models.Cafe
	categories := ""
	amenities := ""
	err := c.postgres.QueryRow(ctx, "SELECT id, owner_id, name, description, opening_time, closing_time, capacity, phone_number, email, province, city, address, location, categories, amenities FROM cafes WHERE id = $1", id).Scan(&cafe.ID, &cafe.OwnerID, &cafe.Name, &cafe.Description, &cafe.OpeningTime, &cafe.ClosingTime, &cafe.Capacity, &cafe.ContactInfo.Phone, &cafe.ContactInfo.Email, &cafe.ContactInfo.Province, &cafe.ContactInfo.City, &cafe.ContactInfo.Address, &cafe.ContactInfo.Location, &categories, &amenities)
	if err != nil {
		log.GetLog().Errorf("Unable to get cafe by id. error: %v", err)
	}

	for _, category := range strings.Split(categories, ",") {
		cafe.Categories = append(cafe.Categories, models.CafeCategory(category))
	}

	for _, amenity := range strings.Split(amenities, ",") {
		cafe.Amenities = append(cafe.Amenities, models.AmenityCategory(amenity))
	}

	return &cafe, err
}

func (c *CafesRepoImp) SearchCafe(ctx context.Context, name string, province string, city string, category string) ([]models.Cafe, error) {
	var cafes []models.Cafe
	list := []string{}
	query := "SELECT id, owner_id, name, description, opening_time, closing_time, capacity, phone_number, email, province, city, address, location, categories, amenities FROM cafes"

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
		categories := ""
		amenities := ""
		err = rows.Scan(&cafe.ID, &cafe.OwnerID, &cafe.Name, &cafe.Description, &cafe.OpeningTime, &cafe.ClosingTime, &cafe.Capacity, &cafe.ContactInfo.Phone, &cafe.ContactInfo.Email, &cafe.ContactInfo.Province, &cafe.ContactInfo.City, &cafe.ContactInfo.Address, &cafe.ContactInfo.Location, &categories, &amenities)
		if err != nil {
			log.GetLog().Errorf("Unable to scan cafe. error: %v", err)
		}

		for _, category := range strings.Split(categories, ",") {
			cafe.Categories = append(cafe.Categories, models.CafeCategory(strings.TrimSpace(category)))
		}

		for _, amenity := range strings.Split(amenities, ",") {
			cafe.Amenities = append(cafe.Amenities, models.AmenityCategory(strings.TrimSpace(amenity)))
		}

		cafes = append(cafes, cafe)
	}

	return cafes, err
}

func (c *CafesRepoImp) GetByCafeIDs(ctx context.Context, ids []int32) ([]models.Cafe, error) {
	var cafes []models.Cafe
	listIds := []string{}
	if len(ids) == 0 {
		return cafes, nil
	}

	for _, id := range ids {
		listIds = append(listIds, cast.ToString(id))
	}

	query := "SELECT id, owner_id, name, description, opening_time, closing_time, capacity, phone_number, email, province, city, address, location, categories FROM cafes WHERE id IN ("
	query += strings.Join(listIds, ",")
	query += ")"

	rows, err := c.postgres.Query(ctx, query)
	if err != nil {
		log.GetLog().Errorf("Unable to get cafes by ids. error: %v", err)
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

func (c *CafesRepoImp) Update(ctx context.Context, id int32, updateType UpdateType, value interface{}) error {
	columnName := string(updateType)
	if columnName == "categories" {
		categoriesVal := value.([]models.CafeCategory)
		categories := ""
		for i, category := range categoriesVal {
			categories += string(category)
			if i != len(categoriesVal)-1 {
				categories += ","
			}
		}
		value = categories
	} else if columnName == "amenities" {
		amenitiesVal := value.([]models.AmenityCategory)
		amenities := ""
		for i, category := range amenitiesVal {
			amenities += string(category)
			if i != len(amenitiesVal)-1 {
				amenities += ","
			}
		}
		value = amenities
	}

	query := "UPDATE cafes SET " + columnName + " = $1 WHERE id = $2"
	_, err := c.postgres.Exec(ctx, query, value, id)
	if err != nil {
		log.GetLog().Errorf("Unable to update cafe. error: %v", err)
		return err
	}

	return nil
}

func (c *CafesRepoImp) DeleteByID(ctx context.Context, id int32) error {
	_, err := c.postgres.Exec(ctx,
		`DELETE FROM menu_items
		WHERE id = $1`, id)
	if err != nil {
		log.GetLog().Errorf("Unable to delete menu item. error: %v", err)
		return err
	}

	return err
}
