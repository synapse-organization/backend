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

type UpdateCafeType string

const (
	UpdateCafeName             UpdateCafeType = "name"
	UpdateCafeDescription      UpdateCafeType = "description"
	UpdateCafeOpeningTime      UpdateCafeType = "opening_time"
	UpdateCafeClosingTime      UpdateCafeType = "closing_time"
	UpdateCafeCapacity         UpdateCafeType = "capacity"
	UpdateCafePhoneNumber      UpdateCafeType = "phone_number"
	UpdateCafeEmail            UpdateCafeType = "email"
	UpdateCafeLocation         UpdateCafeType = "location"
	UpdateCafeProvince         UpdateCafeType = "province"
	UpdateCafeCity             UpdateCafeType = "city"
	UpdateCafeAddress          UpdateCafeType = "address"
	UpdateCafeCategories       UpdateCafeType = "categories"
	UpdateCafeAmenities        UpdateCafeType = "amenities"
	UpdateCafeReservationPrice UpdateCafeType = "reservation_price"
)

type CafesRepo interface {
	Create(ctx context.Context, cafe *models.Cafe) (int32, error)
	GetByID(ctx context.Context, id int32) (*models.Cafe, error)
	SearchCafe(ctx context.Context, name string, address string, location string, category string) ([]models.Cafe, error)
	GetByCafeIDs(ctx context.Context, ids []int32) ([]models.Cafe, error)
	GetByOwnerID(ctx context.Context, id int32) (*models.Cafe, error)
	Update(ctx context.Context, id int32, updateCafeType UpdateCafeType, value interface{}) error
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
				reservation_price FLOAT,
				FOREIGN KEY (owner_id) REFERENCES users(id)
			);`)

	if err != nil {
		log.GetLog().WithError(err).WithField("table", "cafes").Fatal("Unable to create table")
	}

	_, err = postgres.Exec(context.Background(), `INSERT INTO cafes (id, owner_id, name, description, opening_time, closing_time, capacity, phone_number, email, location, province, city, address, categories, amenities, reservation_price)
			VALUES
			(1, 31, 'Bean & Blossom', 'A stylish cafe with a focus on specialty coffees and homemade desserts.', 8, 23, 25, 12345678901, 'beanandblossom@gmail.com', '', 1, 2, '123 Oak Street, Los Angeles, CA', 'Coffee,Desserts', 'وای فای', 15.0),
			(2, 32, 'Urban Brew Café', 'A cozy cafe offering a variety of teas and snacks.', 9, 23, 30, 23456789012, 'urbanbrewcafe@gmail.com', '', 1, 2, '456 Maple Avenue, Los Angeles, CA', 'Tea,Snacks', 'تلویزیون', 10.0),
			(3, 33, 'The Coffee Corner', 'An elegant cafe known for its artisan pastries.', 7, 21, 20, 34567890123, 'thecoffeecorner@gmail.com', '', 1, 2, '789 Pine Street, Los Angeles, CA', 'Pastries,Coffee', 'پارکینگ', 20.0),
			(4, 34, 'Whisk & Sip', 'A modern cafe with a wide selection of vegan options.', 8, 20, 35, 45678901234, 'whiskandsip@gmail.com', '', 1, 2, '321 Birch Avenue, Los Angeles, CA', 'Vegan,Snacks', 'موسیقی زنده', 18.0),
			(5, 35, 'Café Serenity', 'A family-friendly cafe with a play area for children.', 9, 19, 40, 56789012345, 'cafeserenity@gmail.com', '', 1, 2, '654 Cedar Street, Los Angeles, CA', 'Family,Desserts', 'بردگیم', 12.0),
			(6, 36, 'Roast & Toast', 'A cafe with a great selection of international cuisines.', 10, 21, 45, 67890123456, 'roastandtoast@gmail.com', '', 1, 2, '987 Willow Lane, Los Angeles, CA', 'International,Main Dish', 'اجازه ورود حیوانات خانگی', 22.0),
			(7, 37, 'Mocha Mystique', 'A cozy cafe perfect for working or studying.', 8, 22, 30, 78901234567, 'mochamystique@gmail.com', '', 1, 2, '321 Elm Street, Los Angeles, CA', 'Coffee,Snacks', 'قلیان', 17.0),
			(8, 38, 'Brewed Awakening', 'A rustic cafe with a wide selection of herbal teas.', 9, 23, 25, 89012345678, 'brewedawakening@gmail.com', '', 1, 2, '654 Spruce Avenue, Los Angeles, CA', 'Herbal Tea,Desserts', 'دود آزاد', 16.0),
			(9, 39, 'The Java Lounge', 'A lively cafe with live music performances.', 7, 24, 50, 90123456789, 'thejavalounge@gmail.com', '', 1, 2, '987 Ash Lane, Los Angeles, CA', 'Coffee,Live Music', 'غذای گیاهی', 25.0),
			(10, 40, 'Espresso Escape', 'A minimalist cafe with a serene atmosphere.', 8, 21, 15, 12309845678, 'espressoescape@gmail.com', '', 1, 2, '123 Fir Street, Los Angeles, CA', 'Tea,Snacks', 'غذای وگان', 14.0),
			(11, 41, 'Crème & Beans', 'A hip cafe with unique and creative drinks.', 9, 22, 35, 23410956789, 'cremeandbeans@gmail.com', '', 1, 2, '456 Cedar Avenue, Los Angeles, CA', 'Drinks,Snacks', 'دسترسی برای معلولان', 19.0),
			(12, 42, 'The Daily Grind', 'A cafe offering a great selection of board games.', 10, 20, 30, 34521067890, 'thedailygrind@gmail.com', '', 1, 2, '789 Pine Street, Los Angeles, CA', 'Coffee,Games', 'اتاق جلسه', 15.0),
			(13, 43, 'Mellow Brew', 'A cafe with a large outdoor seating area.', 8, 22, 45, 45632178901, 'mellowbrew@gmail.com', '', 1, 2, '321 Birch Avenue, Los Angeles, CA', 'Tea,Snacks', 'فضای نشستن بیرون', 18.0),
			(14, 44, 'Café Mirage', 'A pet-friendly cafe with great coffee.', 9, 23, 40, 56743289012, 'cafemirage@gmail.com', '', 1, 2, '654 Cedar Street, Los Angeles, CA', 'Coffee,Snacks', 'اجازه ورود حیوانات خانگی', 17.0),
			(15, 45, 'Harvest Brews', 'A cafe with excellent desserts and snacks.', 10, 24, 50, 67854390123, 'harvestbrews@gmail.com', '', 1, 2, '987 Willow Lane, Los Angeles, CA', 'Desserts,Snacks', 'وای فای', 21.0),
			(16, 46, 'Java Junction', 'A modern cafe with great food and drinks.', 8, 20, 35, 78965401234, 'javajunction@gmail.com', '', 1, 2, '321 Elm Street, Los Angeles, CA', 'Main Dish,Drinks', 'تلویزیون', 22.0),
			(17, 47, 'The Cozy Nook', 'A small cafe with a selection of vegan foods.', 9, 21, 20, 89076512345, 'thecozynook@gmail.com', '', 1, 2, '654 Spruce Avenue, Los Angeles, CA', 'Vegan,Snacks', 'پارکینگ', 18.0),
			(18, 48, 'Velvet Roast', 'A cafe with delicious pastries and coffee.', 10, 23, 25, 90187623456, 'velvetroast@gmail.com', '', 1, 2, '987 Ash Lane, Los Angeles, CA', 'Pastries,Coffee', 'موسیقی زنده', 16.0),
			(19, 49, 'Caffeine Haven', 'A lively cafe with various entertainment options.', 8, 22, 50, 12398734567, 'caffeinehaven@gmail.com', '', 1, 2, '123 Fir Street, Los Angeles, CA', 'Coffee,Entertainment', 'بردگیم', 23.0),
			(20, 50, 'Artisan Aroma', 'A cafe offering a variety of teas and herbal drinks.', 9, 21, 30, 23409845678, 'artisanaroma@gmail.com', '', 1, 2, '456 Cedar Avenue, Los Angeles, CA', 'Herbal Tea,Desserts', 'وای فای', 14.0),
			(21, 51, 'Brewed Bliss', 'A rustic cafe with a wide selection of coffees.', 10, 23, 25, 34510956789, 'brewedbliss@gmail.com', '', 1, 2, '789 Pine Street, Los Angeles, CA', 'Coffee,Desserts', 'تلویزیون', 20.0),
			(22, 52, 'The Tea Emporium', 'A family-friendly cafe with a variety of snacks.', 8, 22, 40, 45621067890, 'theteaemporium@gmail.com', '', 1, 2, '321 Birch Avenue, Los Angeles, CA', 'Family,Snacks', 'پارکینگ', 18.0),
			(23, 53, 'Café Delights', 'A cafe with a great selection of international dishes.', 9, 24, 45, 56732178901, 'cafedelights@gmail.com', '', 1, 2, '987 Willow Lane, Los Angeles, CA', 'International,Main Dish', 'موسیقی زنده', 24.0),
			(24, 54, 'Sip & Savor', 'A cozy cafe with excellent coffee and pastries.', 8, 20, 20, 67843289012, 'sipandsavor@gmail.com', '', 1, 2, '321 Elm Street, Los Angeles, CA', 'Coffee,Pastries', 'وای فای', 19.0),
			(25, 55, 'The Coffee House', 'A modern cafe with a wide variety of drinks.', 9, 21, 35, 78954390123, 'thecoffeehouse@gmail.com', '', 1, 2, '654 Spruce Avenue, Los Angeles, CA', 'Drinks,Snacks', 'تلویزیون', 22.0),
			(26, 56, 'Sweet Aroma Café', 'A family-friendly cafe with a variety of snacks.', 10, 22, 40, 89065401234, 'sweetaromacafe@gmail.com', '', 1, 2, '987 Ash Lane, Los Angeles, CA', 'Family,Snacks', 'اجازه ورود حیوانات خانگی', 20.0),
			(27, 57, 'The Rustic Bean', 'A rustic cafe with a wide selection of herbal teas.', 8, 20, 30, 90176512345, 'therusticbean@gmail.com', '', 1, 2, '123 Fir Street, Los Angeles, CA', 'Herbal Tea,Desserts', 'دسترسی برای معلولان', 18.0),
			(28, 58, 'Brew & Mingle', 'A lively cafe with various entertainment options.', 9, 23, 50, 12387623456, 'brewandmingle@gmail.com', '', 1, 2, '456 Cedar Avenue, Los Angeles, CA', 'Coffee,Entertainment', 'بردگیم', 23.0),
			(29, 59, 'Whispering Pines Café', 'A serene cafe perfect for relaxation.', 10, 21, 25, 23498734567, 'whisperingpinescafe@gmail.com', '', 1, 2, '789 Pine Street, Los Angeles, CA', 'Tea,Snacks', 'موسیقی زنده', 17.0),
			(30, 60, 'The Steaming Cup', 'A cozy cafe offering a variety of teas and snacks.', 8, 22, 30, 34509845678, 'thesteamingcup@gmail.com', '', 1, 2, '321 Birch Avenue, Los Angeles, CA', 'Tea,Snacks', 'وای فای', 15.0)
			ON CONFLICT DO NOTHING;
	`)
	if err != nil {
		log.GetLog().Errorf("Unable to insert cafes. error: %v", err)
	}

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

	_, err := c.postgres.Exec(ctx, "INSERT INTO cafes (id, owner_id, name, description, opening_time, closing_time, capacity, phone_number, email, province, city, address, location, categories, amenities, reservation_price) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)", cafe.ID, cafe.OwnerID, cafe.Name, cafe.Description, cafe.OpeningTime, cafe.ClosingTime, cafe.Capacity, cafe.ContactInfo.Phone, cafe.ContactInfo.Email, cafe.ContactInfo.Province, cafe.ContactInfo.City, cafe.ContactInfo.Address, cafe.ContactInfo.Location, categories, amenities, cafe.ReservationPrice)
	if err != nil {
		log.GetLog().Errorf("Unable to insert cafe. error: %v", err)
	}
	return cafe.ID, err
}

func (c *CafesRepoImp) GetByID(ctx context.Context, id int32) (*models.Cafe, error) {
	var cafe models.Cafe
	var categories, amenities string
	err := c.postgres.QueryRow(ctx, "SELECT id, owner_id, name, description, opening_time, closing_time, capacity, phone_number, email, province, city, address, location, categories, amenities, reservation_price FROM cafes WHERE id = $1", id).Scan(&cafe.ID, &cafe.OwnerID, &cafe.Name, &cafe.Description, &cafe.OpeningTime, &cafe.ClosingTime, &cafe.Capacity, &cafe.ContactInfo.Phone, &cafe.ContactInfo.Email, &cafe.ContactInfo.Province, &cafe.ContactInfo.City, &cafe.ContactInfo.Address, &cafe.ContactInfo.Location, &categories, &amenities, &cafe.ReservationPrice)
	if err != nil {
		log.GetLog().Errorf("Unable to get cafe by id. error: %v", err)
	}

	if categories == "" {
		cafe.Categories = []models.CafeCategory{}
	} else {
		for _, category := range strings.Split(categories, ",") {
			cafe.Categories = append(cafe.Categories, models.CafeCategory(category))
		}
	}

	if amenities == "" {
		cafe.Amenities = []models.AmenityCategory{}
	} else {
		for _, amenity := range strings.Split(amenities, ",") {
			cafe.Amenities = append(cafe.Amenities, models.AmenityCategory(amenity))
		}
	}
	
	return &cafe, err
}

func (c *CafesRepoImp) SearchCafe(ctx context.Context, name string, province string, city string, category string) ([]models.Cafe, error) {
	var cafes []models.Cafe
	list := []string{}
	query := "SELECT id, owner_id, name, description, opening_time, closing_time, capacity, phone_number, email, province, city, address, location, categories, amenities, reservation_price FROM cafes"

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
		err = rows.Scan(&cafe.ID, &cafe.OwnerID, &cafe.Name, &cafe.Description, &cafe.OpeningTime, &cafe.ClosingTime, &cafe.Capacity, &cafe.ContactInfo.Phone, &cafe.ContactInfo.Email, &cafe.ContactInfo.Province, &cafe.ContactInfo.City, &cafe.ContactInfo.Address, &cafe.ContactInfo.Location, &categories, &amenities, &cafe.ReservationPrice)
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

	query := "SELECT id, owner_id, name, description, opening_time, closing_time, capacity, phone_number, email, province, city, address, location, categories, amenities, reservation_price FROM cafes WHERE id IN ("
	query += strings.Join(listIds, ",")
	query += ")"

	rows, err := c.postgres.Query(ctx, query)
	if err != nil {
		log.GetLog().Errorf("Unable to get cafes by ids. error: %v", err)
	}

	for rows.Next() {
		var cafe models.Cafe
		catetories := ""
		amenities := ""
		err = rows.Scan(&cafe.ID, &cafe.OwnerID, &cafe.Name, &cafe.Description, &cafe.OpeningTime, &cafe.ClosingTime, &cafe.Capacity, &cafe.ContactInfo.Phone, &cafe.ContactInfo.Email, &cafe.ContactInfo.Province, &cafe.ContactInfo.City, &cafe.ContactInfo.Address, &cafe.ContactInfo.Location, &catetories, &amenities, &cafe.ReservationPrice)
		if err != nil {
			log.GetLog().Errorf("Unable to scan cafe. error: %v", err)
		}

		for _, category := range strings.Split(catetories, ",") {
			cafe.Categories = append(cafe.Categories, models.CafeCategory(strings.TrimSpace(category)))
		}

		for _, amenity := range strings.Split(amenities, ",") {
			cafe.Amenities = append(cafe.Amenities, models.AmenityCategory(strings.TrimSpace(amenity)))
		}

		cafes = append(cafes, cafe)
	}

	return cafes, err
}

func (c *CafesRepoImp) GetByOwnerID(ctx context.Context, id int32) (*models.Cafe, error) {
	var cafe models.Cafe
	categories := ""
	amenities := ""
	err := c.postgres.QueryRow(ctx, `SELECT id, owner_id, name, description, opening_time, closing_time, capacity, phone_number, email, province, city, address, location, categories, amenities, reservation_price FROM cafes WHERE owner_id = $1`, id).Scan(&cafe.ID, &cafe.OwnerID, &cafe.Name, &cafe.Description, &cafe.OpeningTime, &cafe.ClosingTime, &cafe.Capacity, &cafe.ContactInfo.Phone, &cafe.ContactInfo.Email, &cafe.ContactInfo.Province, &cafe.ContactInfo.City, &cafe.ContactInfo.Address, &cafe.ContactInfo.Location, &categories, &amenities, &cafe.ReservationPrice)
	if err != nil {
		log.GetLog().Errorf("Unable to get cafe by owner id. error: %v", err)
	}

	for _, category := range strings.Split(categories, ",") {
		cafe.Categories = append(cafe.Categories, models.CafeCategory(category))
	}

	for _, amenity := range strings.Split(amenities, ",") {
		cafe.Amenities = append(cafe.Amenities, models.AmenityCategory(amenity))
	}

	return &cafe, err
}

func (c *CafesRepoImp) Update(ctx context.Context, id int32, updateCafeType UpdateCafeType, value interface{}) error {
	columnName := string(updateCafeType)
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
		`DELETE FROM cafes
		WHERE id = $1`, id)
	if err != nil {
		log.GetLog().Errorf("Unable to delete cafe. error: %v", err)
		return err
	}

	return err
}
