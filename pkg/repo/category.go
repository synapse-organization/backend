package repo

import (
	"barista/pkg/log"
	"barista/pkg/models"
	"context"
	"github.com/jackc/pgx/v5"
)

type CategoryRepo interface {
	Create(ctx context.Context, category *models.Category) error
	GetByID(ctx context.Context, id int32) (*models.Category, error)
}

type CategoryRepoImp struct {
	postgres *pgx.Conn
}

func NewCategoryRepoImp(postgres *pgx.Conn) *CategoryRepoImp {
	_, err := postgres.Exec(context.Background(),
		`CREATE TABLE IF NOT EXISTS categories (
				id INTEGER PRIMARY KEY,
				name TEXT,
				category_type TEXT,
				UNIQUE(name)
			);`)
	if err != nil {
		log.GetLog().WithError(err).WithField("table", "categories").Fatal("Unable to create table")
	}
	return &CategoryRepoImp{postgres: postgres}
}

func (c *CategoryRepoImp) Create(ctx context.Context, category *models.Category) error {
	_, err := c.postgres.Exec(ctx, "INSERT INTO categories (id, name, category_type) VALUES ($1, $2, $3)", category.ID, category.Name, category.Type)
	if err != nil {
		log.GetLog().Errorf("Unable to insert category. error: %v", err)
	}
	return err
}

func (c *CategoryRepoImp) GetByID(ctx context.Context, id int32) (*models.Category, error) {
	var category models.Category
	err := c.postgres.QueryRow(ctx, "SELECT id, name, category_type FROM categories WHERE id = $1", id).Scan(&category.ID, &category.Name, category.Type)
	if err != nil {
		log.GetLog().Errorf("Unable to get category by id. error: %v", err)
	}
	return &category, err
}
