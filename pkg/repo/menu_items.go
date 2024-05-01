package repo

import (
	"barista/pkg/log"
	"barista/pkg/models"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func init() {
	log.GetLog().Info("Init MenuItemsRepo")
}

type MenuItemsRepo interface {
	Create(ctx context.Context, menuItem *models.MenuItem) error
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
			image_id TEXT,
			FOREIGN KEY (cafe_id) REFERENCES cafes(id),
			FOREIGN KEY (image_id) REFERENCES images(id)
		);`)
	if err != nil {
		log.GetLog().WithError(err).WithField("table", "menu_items").Fatal("Unable to create table")
	}
	return &MenuItemsRepoImp{postgres: postgres}
}

func (c *MenuItemsRepoImp) Create(ctx context.Context, menuItem *models.MenuItem) error {
	_, err := c.postgres.Exec(ctx,
		`INSERT INTO menu_items (id, cafe_id, name, price, category, image_id)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		menuItem.ID, menuItem.CafeID, menuItem.Name, menuItem.Price, menuItem.Category, menuItem.ImageID)
	if err != nil {
		log.GetLog().Errorf("Unable to insert menu item. error: %v", err)
	}
	return err
}
