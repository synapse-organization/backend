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
	Create(ctx context.Context, menuItem *models.MenuItem) error
	GetItemsByCafeID(ctx context.Context, cafeID int32) ([]*models.MenuItem, error)
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
	return &MenuItemsRepoImp{postgres: postgres}
}

func (c *MenuItemsRepoImp) Create(ctx context.Context, menuItem *models.MenuItem) error {
	menuItem.ID = rand.Int31()
	ingredients := ""
	for i, ingredient := range menuItem.Ingredients {
		ingredients += ingredient
		if i != len(menuItem.Ingredients)-1 {
			ingredients += ","
		}
	}

	_, err := c.postgres.Exec(ctx,
		`INSERT INTO menu_items (id, cafe_id, name, price, category, ingredients)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		menuItem.ID, menuItem.CafeID, menuItem.Name, menuItem.Price, menuItem.Category, ingredients)
	if err != nil {
		log.GetLog().Errorf("Unable to insert menu item. error: %v", err)
	}
	return err
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
