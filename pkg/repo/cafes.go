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
			(1, 31, 'باکارا', 'یک کافه شیک با تمرکز بر قهوه‌های تخصصی و دسرهای خوشمزه.', 8, 23, 40, 09373564072, 'bakara@gmail.com', '', 8, 329, 'خیابان انقلاب, بعد از تئاتر شهر', 'coffee_shop,restaurant', 'وای فای,موسیقی زنده,غذای وگان,تلویزیون', 40000.0),
			(2, 32, 'وی کافه', 'یک کافه دنج با ارائه مجموعه‌ای از چای‌ها و تنقلات.', 9, 23, 30, 02188803714, 'vcafe@gmail.com', '', 8, 329, 'خیابان فلسطین, پایین تر از بزرگمهر', 'coffee_shop', 'دود آزاد,فضای کار', 20000.0),
			(3, 33, 'ویونا', 'یک کافه زیبا که به شیرینی‌های دست‌ساز خود معروف است.', 7, 21, 25, 34567890123, 'viona@gmail.com', '', 8, 329, 'خیابان کریمخان زند', 'tea_house', 'پارکینگ', 20000.0),
			(4, 34, 'نادری', 'یک کافه مدرن با مجموعه‌ای گسترده از گزینه‌های خوشمزه.', 8, 20, 35, 45678901234, 'naderi@gmail.com', '', 8, 329, '321 Birch Avenue, Los Angeles, CA', 'food_court', 'موسیقی زنده', 18000.0),
			(5, 35, 'جز', 'یک کافه خانوادگی با یک منطقه بازی برای کودکان.', 9, 19, 40, 56789012345, 'jazcafe@gmail.com', '', 4, 101, '654 Cedar Street, Los Angeles, CA', 'ice_cream', 'بردگیم', 12000.0),
			(6, 36, 'راک', 'یک کافه با مجموعه‌ای عالی از غذاهای بین‌المللی.', 10, 21, 45, 67890123456, 'rockcafeiran@gmail.com', '', 4, 101, '987 Willow Lane, Los Angeles, CA', 'dessert_shop', 'اجازه ورود حیوانات خانگی', 22000.0),
			(7, 37, 'آبی', 'یک کافه دنج مناسب برای کار یا مطالعه.', 8, 22, 30, 78901234567, 'abicafe@gmail.com', '', 11, 1153, '321 Elm Street, Los Angeles, CA', 'ice_cream', 'قلیان', 17000.0),
			(8, 38, 'ریک', 'یک کافه روستیک با مجموعه‌ای گسترده از چای‌های گیاهی.', 9, 23, 25, 89012345678, 'brewedawakening@gmail.com', '', 11, 1153, '654 Spruce Avenue, Los Angeles, CA', 'dessert_shop', 'دود آزاد', 16000.0),
			(9, 39, 'کندوک', 'یک کافه پرجنب و جوش با اجرای موسیقی زنده.', 7, 24, 50, 90123456789, 'thejavalounge@gmail.com', '', 11, 1153, '987 Ash Lane, Los Angeles, CA', 'restaurant', 'غذای گیاهی', 25000.0),
			(10, 40, 'شیلا', 'یک کافه مینیمالیستی با فضای آرام.', 8, 21, 15, 12309845678, 'espressoescape@gmail.com', '', 4, 101, '123 Fir Street, Los Angeles, CA', 'restaurant', 'غذای وگان', 14000.0),
			(11, 41, 'کندیک', 'یک کافه هیپ با نوشیدنی‌های خلاقانه و منحصر به فرد.', 9, 22, 35, 23410956789, 'cremeandbeans@gmail.com', '', 4, 101, '456 Cedar Avenue, Los Angeles, CA', 'restaurant', 'دسترسی برای معلولان', 19000.0),
			(12, 42, 'ارکیده', 'یک کافه با مجموعه‌ای عالی از بازی‌های رومیزی.', 10, 20, 30, 34521067890, 'thedailygrind@gmail.com', '', 17, 790, '789 Pine Street, Los Angeles, CA', 'food_court', 'اتاق جلسه', 15000.0),
			(13, 43, 'زیتون', 'یک کافه با یک منطقه بزرگ نشستن در فضای باز.', 8, 22, 45, 45632178901, 'mellowbrew@gmail.com', '', 20, 705, '321 Birch Avenue, Los Angeles, CA', 'tea_house', 'فضای نشستن بیرون', 18000.0),
			(14, 44, 'گودو', 'یک کافه دوستانه با قهوه عالی.', 9, 23, 40, 56743289012, 'cafemirage@gmail.com', '', 20, 705, '654 Cedar Street, Los Angeles, CA', 'ice_cream', 'اجازه ورود حیوانات خانگی', 17000.0),
			(15, 45, 'درسا', 'یک کافه با دسرها و تنقلات عالی.', 10, 24, 50, 67854390123, 'harvestbrews@gmail.com', '', 20, 705, '987 Willow Lane, Los Angeles, CA', 'dessert_shop', 'وای فای', 21000.0)
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
