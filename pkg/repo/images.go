package repo

import (
	"barista/pkg/log"
	"barista/pkg/models"
	"context"
	"github.com/jackc/pgx/v5"
)

type ImageRepo interface {
	Create(ctx context.Context, image *models.Image) error
	GetByCafeID(ctx context.Context, id int32) ([]*models.Image, error)
}

type ImageRepoImp struct {
	postgres *pgx.Conn
}

func NewImageRepoImp(postgres *pgx.Conn) *ImageRepoImp {
	_, err := postgres.Exec(context.Background(),
		`CREATE TABLE IF NOT EXISTS images (
				id TEXT PRIMARY KEY,
				cafe_id INTEGER,
				create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				FOREIGN KEY (cafe_id) REFERENCES cafes(id)
			);`)
	if err != nil {
		log.GetLog().WithError(err).WithField("table", "images").Fatal("Unable to create table")
	}
	return &ImageRepoImp{postgres: postgres}
}

func (r *ImageRepoImp) Create(ctx context.Context, image *models.Image) error {
	_, err := r.postgres.Exec(ctx, "INSERT INTO images (id, cafe_id) VALUES ($1, $2)", image.ID, image.CafeID)
	if err != nil {
		log.GetLog().Errorf("Unable to insert image. error: %v", err)
	}
	return err
}

func (r *ImageRepoImp) GetByCafeID(ctx context.Context, id int32) ([]*models.Image, error) {
	var images []*models.Image
	rows, err := r.postgres.Query(ctx, "SELECT id, cafe_id FROM images WHERE cafe_id = $1 ORDER BY create_at DESC", id)
	if err != nil {
		log.GetLog().Errorf("Unable to get images by cafe id. error: %v", err)
	}

	for rows.Next() {
		var image models.Image
		err := rows.Scan(&image.ID, &image.CafeID)
		if err != nil {
			log.GetLog().Errorf("Unable to scan image. error: %v", err)
			continue
		}
		images = append(images, &image)
	}

	return images, err
}