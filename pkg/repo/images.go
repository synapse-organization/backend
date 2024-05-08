package repo

import (
	"barista/pkg/log"
	"barista/pkg/models"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ImageRepo interface {
	Create(ctx context.Context, image *models.Image) error
	GetByReferenceID(ctx context.Context, id int32) ([]*models.Image, error)
	CheckExistence(ctx context.Context, imageID string) (bool, error)
	DeleteByID(ctx context.Context, id string) error
	DeleteByReferenceID(ctx context.Context, referenceID int32) error
	GetMainImage(ctx context.Context, referenceID int32) (string, error)
	UpdateByReferenceID(ctx context.Context, referenceID int32, id string) error
}

type ImageRepoImp struct {
	postgres *pgxpool.Pool
}

func NewImageRepoImp(postgres *pgxpool.Pool) *ImageRepoImp {
	_, err := postgres.Exec(context.Background(),
		`CREATE TABLE IF NOT EXISTS images (
				id TEXT PRIMARY KEY,
				reference_id INTEGER,
				create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			);`)
	if err != nil {
		log.GetLog().WithError(err).WithField("table", "images").Fatal("Unable to create table")
	}
	return &ImageRepoImp{postgres: postgres}
}

func (r *ImageRepoImp) Create(ctx context.Context, image *models.Image) error {
	_, err := r.postgres.Exec(ctx, "INSERT INTO images (id, reference_id) VALUES ($1, $2)", image.ID, image.Reference)
	if err != nil {
		log.GetLog().Errorf("Unable to insert image. error: %v", err)
	}
	return err
}

func (r *ImageRepoImp) GetByReferenceID(ctx context.Context, id int32) ([]*models.Image, error) {
	var images []*models.Image
	rows, err := r.postgres.Query(ctx, "SELECT id, reference_id FROM images WHERE reference_id = $1 ORDER BY create_at DESC", id)
	if err != nil {
		log.GetLog().Errorf("Unable to get images by cafe id. error: %v", err)
	}

	for rows.Next() {
		var image models.Image
		err := rows.Scan(&image.ID, &image.Reference)
		if err != nil {
			log.GetLog().Errorf("Unable to scan image. error: %v", err)
			continue
		}
		images = append(images, &image)
	}

	return images, err
}

func (r *ImageRepoImp) CheckExistence(ctx context.Context, imageID string) (bool, error) {
	var exists bool
	err := r.postgres.QueryRow(ctx,
		`SELECT EXISTS
		(SELECT 1 FROM images WHERE id = $1)`,
		imageID).Scan(&exists)
	if err != nil {
		log.GetLog().Errorf("Unable to check image existence. error: %v", err)
	}

	return exists, err
}

func (r *ImageRepoImp) DeleteByID(ctx context.Context, id string) error {
	_, err := r.postgres.Exec(ctx,
		`DELETE FROM images
		WHERE id = $1`, id)
	if err != nil {
		log.GetLog().Errorf("Unable to delete image by id. error: %v", err)
		return err
	}

	return err
}

func (r *ImageRepoImp) DeleteByReferenceID(ctx context.Context, referenceID int32) error {
	_, err := r.postgres.Exec(ctx,
		`DELETE FROM images
		WHERE reference_id = $1`, referenceID)
	if err != nil {
		log.GetLog().Errorf("Unable to delete image by reference id. error: %v", err)
		return err
	}

	return err
}

func (r *ImageRepoImp) GetMainImage(ctx context.Context, referenceID int32) (string, error) {
	imageID := ""
	err := r.postgres.QueryRow(ctx,
		`SELECT id
		FROM images
		WHERE reference_id = $1
		ORDER BY create_at
		LIMIT 1`, referenceID).Scan(&imageID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			imageID = ""
		}else {
			log.GetLog().Errorf("Unable to get main image. error: %v", err)
			return "", err
		}
	}

	return imageID, nil
}

func (r *ImageRepoImp) UpdateByReferenceID(ctx context.Context, referenceID int32, id string) error {
	_, err := r.postgres.Exec(ctx,
		`UPDATE images
		SET reference_id = $1
		WHERE id = $2`,
		referenceID, id)
	if err != nil {
		log.GetLog().Errorf("Unable to update image. error: %v", err)
		return err
	}

	return nil
}
