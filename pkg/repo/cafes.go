package repo

import (
	"barista/pkg/log"
	"barista/pkg/models"
	"context"
	"github.com/jackc/pgx/v5"
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
				location TEXT,
				province TEXT,
				city TEXT,
				address TEXT,
				categories TEXT,
				FOREIGN KEY (owner_id) REFERENCES users(id)
			);`)

	if err != nil {
		log.GetLog().WithError(err).WithField("table", "cafes").Fatal("Unable to create table")
	}

	_, err = postgres.Exec(context.Background(), `INSERT INTO cafes (id, owner_id, name, description, opening_time, closing_time, capacity, phone_number, email, location, province, city, address, categories)
			VALUES 
			(1, 12, 'Cafe One', 'A stylish cafe with a focus on specialty coffees and homemade desserts.', '2024-04-18 08:00:00', '2024-04-18 18:00:00', 25, '+12345678901', 'cafe_one@example.com', '123 Oak Street', 'California', 'Los Angeles', '123 Oak Street, Los Angeles, CA', 'Coffee, Desserts'),
			(2, 12, 'Sunrise Cafe', 'Start your day with our freshly brewed coffees and delicious breakfast options.', '2024-04-18 07:00:00', '2024-04-18 15:00:00', 20, '+12345678902', 'sunrise_cafe@example.com', '456 Maple Avenue', 'California', 'San Francisco', '456 Maple Avenue, San Francisco, CA', 'Coffee, Breakfast'),
			(3, 12, 'Mellow Brew', 'Relax and unwind with our selection of herbal teas and tranquil ambiance.', '2024-04-18 09:00:00', '2024-04-18 20:00:00', 15, '+12345678903', 'mellow_brew@example.com', '789 Pine Street', 'California', 'San Diego', '789 Pine Street, San Diego, CA', 'Tea, Relaxation'),
			(4, 12, 'City Buzz', 'A bustling cafe in the heart of downtown, serving up energizing coffees and light snacks.', '2024-04-18 07:30:00', '2024-04-18 17:30:00', 35, '+12345678904', 'city_buzz@example.com', '101 Market Street', 'California', 'San Jose', '101 Market Street, San Jose, CA', 'Coffee, Snacks'),
			(5, 12, 'Coastal Cafe', 'Enjoy ocean views and coastal vibes with our selection of coastal-inspired drinks and treats.', '2024-04-18 10:00:00', '2024-04-18 18:00:00', 40, '+12345678905', 'coastal_cafe@example.com', '234 Beach Boulevard', 'California', 'Santa Monica', '234 Beach Boulevard, Santa Monica, CA', 'Coffee, Coastal'),
			(6, 12, 'Green Haven', 'A green oasis in the city, offering organic coffees, teas, and vegetarian delights.', '2024-04-18 08:30:00', '2024-04-18 19:30:00', 25, '+12345678906', 'green_haven@example.com', '345 Park Avenue', 'California', 'Oakland', '345 Park Avenue, Oakland, CA', 'Coffee, Tea, Vegetarian'),
			(7, 12, 'Urban Roast', 'Fuel your day with our bold espresso drinks and hearty sandwiches.', '2024-04-18 06:30:00', '2024-04-18 16:30:00', 30, '+12345678907', 'urban_roast@example.com', '456 Broadway', 'California', 'Sacramento', '456 Broadway, Sacramento, CA', 'Coffee, Sandwiches'),
			(8, 12, 'Rustic Brew', 'Step into our cozy atmosphere and savor our rustic-inspired coffees and pastries.', '2024-04-18 08:00:00', '2024-04-18 18:00:00', 20, '+12345678908', 'rustic_brew@example.com', '567 Elm Street', 'California', 'Fresno', '567 Elm Street, Fresno, CA', 'Coffee, Pastries'),
			(9, 12, 'Lakeside Lounge', 'Relax by the lake with our refreshing iced drinks and lakeside bites.', '2024-04-18 09:00:00', '2024-04-18 17:00:00', 15, '+12345678909', 'lakeside_lounge@example.com', '678 Lakeview Drive', 'California', 'Lake Tahoe', '678 Lakeview Drive, Lake Tahoe, CA', 'Coffee, Snacks, Lakeside'),
			(10, 12, 'Tranquil Brew', 'Find peace and serenity with our selection of calming teas and quiet ambiance.', '2024-04-18 10:00:00', '2024-04-18 20:00:00', 25, '+12345678910', 'tranquil_brew@example.com', '789 Forest Avenue', 'California', 'Yosemite', '789 Forest Avenue, Yosemite, CA', 'Tea, Relaxation');
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

	rows, err := c.postgres.Query(ctx, "SELECT id, owner_id, name, description, opening_time, closing_time, capacity, phone_number, email, province, city, address, location, categories FROM cafes WHERE "+strings.Join(list, " AND "))
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
