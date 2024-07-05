package repo

import (
	"barista/pkg/log"
	"barista/pkg/models"
	"context"
	"math/rand"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

func init() {
	log.GetLog().Info("Init MenuItemsRepo")
}

type MenuItemsRepo interface {
	Create(ctx context.Context, menuItem *models.MenuItem) (int32, error)
	GetItemsByCafeID(ctx context.Context, cafeID int32) ([]*models.MenuItem, error)
	GetByID(ctx context.Context, id int32) (*models.MenuItem, error)
	UpdateName(ctx context.Context, id int32, newName string) error
	UpdatePrice(ctx context.Context, id int32, newPrice float64) error
	UpdateIngredients(ctx context.Context, id int32, newIngredients []string) error
	// UpdateImageID(ctx context.Context, id int32, newImage string) error
	DeleteByID(ctx context.Context, id int32) error
}

type MenuItemsRepoImp struct {
	postgres *pgxpool.Pool
}

func NewMenuItemRepoImp(postgres *pgxpool.Pool) *MenuItemsRepoImp {
	_, err := postgres.Exec(context.Background(),
		`CREATE TABLE IF NOT EXISTS menu_items (
			id INT PRIMARY KEY,
			cafe_id INT,
			name TEXT,
			price FLOAT,
			category TEXT,
			ingredients TEXT,
			FOREIGN KEY (cafe_id) REFERENCES cafes(id)
		);`)
	if err != nil {
		log.GetLog().WithError(err).WithField("table", "menu_items").Fatal("Unable to create table")
	}

	_, err = postgres.Exec(context.Background(), `INSERT INTO menu_items (id, cafe_id, name, price, category, ingredients)
		VALUES
		(91, 1, 'نسکافه', 50000.0, 'coffee', 'آب, دانه های قهوه'),
		(92, 1, 'کاپوچینو', 70000.0, 'coffee', 'اسپرسو, شیر, فوم'),
		(93, 1, 'لاته', 70000.0, 'coffee', 'اسپرسو, شیر, فوم'),
		(94, 1, 'چای سبز', 60000.0, 'tea', 'برگ چای سبز, آب'),
		(95, 1, 'چای لاته', 65000.0, 'tea', 'چای سیاه, ادویه, شیر, آب'),
		(96, 1, 'مافین بلوبری', 60000.0, 'dessert', 'آرد, شکر, بلوبری, تخم مرغ, کره, بکینگ پودر'),
		(97, 1, 'کیک شکلاتی', 55000.0, 'dessert', 'آرد, شکر, پودر کاکائو, تخم مرغ, کره, بکینگ پودر'),
		(98, 1, 'سالاد سزار', 130000.0, 'appetizer', 'کاهو, کروتون, پنیز پارمزان, چاشنی سزار'),
		(99, 1, 'ساندویچ مرغ کبابی', 160000.0, 'main_dish', 'نان, پنیر, کره, مرغ'),
		(100, 1, 'پیتزا مارگاریتا', 220000.0, 'main_dish', 'موزارلا، گوجه گیلاسی، ريحان ایتالیایی، سس مارينارا'),
		(101, 1, 'لمون بری', 105000.0, 'drink', 'لیمو, پوره میوه های قرمز, سودا')
		ON CONFLICT DO NOTHING;
	`)
	if err != nil {
		log.GetLog().Errorf("Unable to insert menu items. error: %v", err)
	}

	return &MenuItemsRepoImp{postgres: postgres}
}

func (c *MenuItemsRepoImp) Create(ctx context.Context, menuItem *models.MenuItem) (int32, error) {
	menuItem.ID = rand.Int31()
	ingredients := strings.Join(menuItem.Ingredients, ",")

	_, err := c.postgres.Exec(ctx,
		`INSERT INTO menu_items (id, cafe_id, name, price, category, ingredients)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		menuItem.ID, menuItem.CafeID, menuItem.Name, menuItem.Price, menuItem.Category, ingredients)
	if err != nil {
		log.GetLog().Errorf("Unable to insert menu item. error: %v", err)
	}
	return menuItem.ID, err
}

func (c *MenuItemsRepoImp) GetItemsByCafeID(ctx context.Context, cafeID int32) ([]*models.MenuItem, error) {
	rows, err := c.postgres.Query(ctx,
		`SELECT id, cafe_id, name, price, category, ingredients
		FROM menu_items
		WHERE cafe_id = $1`, cafeID)
	if err != nil {
		log.GetLog().Errorf("Unable to get menu items by cafe id. error: %v", err)
		return nil, err
	}
	defer rows.Close()

	var menu []*models.MenuItem
	for rows.Next() {
		var item models.MenuItem
		ingredients := ""
		err := rows.Scan(&item.ID, &item.CafeID, &item.Name, &item.Price, &item.Category, &ingredients)
		if err != nil {
			log.GetLog().Errorf("Unable to scan menu item. error: %v", err)
			return nil, err
		}

		item.Ingredients = append(item.Ingredients, strings.Split(ingredients, ",")...)

		menu = append(menu, &item)
	}

	return menu, err
}

func (c *MenuItemsRepoImp) GetByID(ctx context.Context, id int32) (*models.MenuItem, error) {
	var item models.MenuItem
	ingredients := ""
	err := c.postgres.QueryRow(ctx,
		`SELECT id, cafe_id, name, price, category, ingredients
		FROM menu_items
		WHERE id = $1`, id).Scan(&item.ID, &item.CafeID, &item.Name, &item.Price, &item.Category, &ingredients)
	if err != nil {
		log.GetLog().Errorf("Unable to get menu item by id. error: %v", err)
		return nil, err
	}

	item.Ingredients = append(item.Ingredients, strings.Split(ingredients, ",")...)

	return &item, err
}

func (c *MenuItemsRepoImp) UpdateName(ctx context.Context, id int32, newName string) error {
	_, err := c.postgres.Exec(ctx,
		`UPDATE menu_items
		SET name = $1
		WHERE id = $2`,
		newName, id)
	if err != nil {
		log.GetLog().Errorf("Unable to update menu items name. error: %v", err)
		return err
	}

	return err
}

func (c *MenuItemsRepoImp) UpdatePrice(ctx context.Context, id int32, newPrice float64) error {
	_, err := c.postgres.Exec(ctx,
		`UPDATE menu_items
		SET price = $1
		WHERE id = $2`,
		newPrice, id)
	if err != nil {
		log.GetLog().Errorf("Unable to update menu items price. error: %v", err)
		return err
	}

	return err
}

func (c *MenuItemsRepoImp) UpdateIngredients(ctx context.Context, id int32, newIngredients []string) error {
	ingreds := strings.Join(newIngredients, ",")
	_, err := c.postgres.Exec(ctx,
		`UPDATE menu_items
		SET ingredients = $1
		WHERE id = $2`,
		ingreds, id)
	if err != nil {
		log.GetLog().Errorf("Unable to update menu items ingredients. error: %v", err)
		return err
	}

	return err
}

// func (c *MenuItemsRepoImp) UpdateImageID(ctx context.Context, id int32, newImage string) error {
// 	_, err := c.postgres.Exec(ctx,
// 		`UPDATE menu_items
// 		SET image_id = $1
// 		WHERE id = $2`,
// 		newImage, id)
// 	if err != nil {
// 		log.GetLog().Errorf("Unable to update menu items image in menu items. error: %v", err)
// 		return err
// 	}

// 	return err
// }

func (c *MenuItemsRepoImp) DeleteByID(ctx context.Context, id int32) error {
	_, err := c.postgres.Exec(ctx,
		`DELETE FROM menu_items
		WHERE id = $1`, id)
	if err != nil {
		log.GetLog().Errorf("Unable to delete menu item. error: %v", err)
		return err
	}

	return err
}
