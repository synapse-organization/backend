package repo

import (
	"barista/pkg/log"
	"barista/pkg/models"
	"context"
	"math/rand"

	"github.com/jackc/pgx/v5/pgxpool"
)

type FavoritesRepo interface {
	Create(ctx context.Context, favorite *models.Favorite) (int32, error)
	GetByID(ctx context.Context, id int32) (*models.Favorite, error)
	CheckExists(ctx context.Context, userID int32, cafeID int32) (bool, error)
	DeleteByID(ctx context.Context, id int32) error
	DeleteByIDs(ctx context.Context, userID int32, cafeID int32) error
	GetFavoritesByUserID(ctx context.Context, userID int32) ([]*models.Favorite, error)
}

type FavoritesRepoImp struct {
	postgres *pgxpool.Pool
}

func NewFavoritesRepoImp(postgres *pgxpool.Pool) *FavoritesRepoImp {
	_, err := postgres.Exec(context.Background(),
		`CREATE TABLE IF NOT EXISTS favorites (
				id INTEGER PRIMARY KEY,
				user_id INTEGER,
				cafe_id INTEGER,
				FOREIGN KEY (user_id) REFERENCES users(id),
				FOREIGN KEY (cafe_id) REFERENCES cafes(id)
			);`)

	if err != nil {
		log.GetLog().WithError(err).WithField("table", "favorites").Fatal("Unable to create table")
	}

	return &FavoritesRepoImp{postgres: postgres}
}

func (r *FavoritesRepoImp) Create(ctx context.Context, favorite *models.Favorite) (int32, error) {
	id := rand.Int31()
	_, err := r.postgres.Exec(ctx, "INSERT INTO favorites (id, user_id, cafe_id) VALUES ($1, $2, $3)", id, favorite.UserID, favorite.CafeID)
	if err != nil {
		log.GetLog().Errorf("Unable to insert favorite. error: %v", err)
	}
	return id, err
}

func (r *FavoritesRepoImp) GetByID(ctx context.Context, id int32) (*models.Favorite, error) {
	var favorite models.Favorite
	err := r.postgres.QueryRow(ctx, "SELECT id, user_id, cafe_id FROM favorites WHERE id = $1", id).Scan(&favorite.ID, &favorite.UserID, &favorite.CafeID)
	if err != nil {
		log.GetLog().Errorf("Unable to get favorite by id. error: %v", err)
	}

	return &favorite, err
}

func (r *FavoritesRepoImp) CheckExists(ctx context.Context, userID int32, cafeID int32) (bool, error) {
	var exists bool
	err := r.postgres.QueryRow(ctx,
		`SELECT EXISTS(
			SELECT 1
			FROM favorites
			WHERE user_id = $1
			AND cafe_id = $2)`,
		userID, cafeID).Scan(&exists)
	if err != nil {
		log.GetLog().Errorf("Unable to check favorite existence. error: %v", err)
	}

	return exists, nil
}

func (r *FavoritesRepoImp) DeleteByID(ctx context.Context, id int32) error {
	_, err := r.postgres.Exec(ctx,
		`DELETE FROM favorites
		WHERE id = $1`, id)
	if err != nil {
		log.GetLog().Errorf("Unable to delete favorite by id. error: %v", err)
		return err
	}

	return err
}

func (r *FavoritesRepoImp) DeleteByIDs(ctx context.Context, userID int32, cafeID int32) error {
	_, err := r.postgres.Exec(ctx,
		`DELETE FROM favorites
		WHERE user_id = $1
		AND cafe_id = $2`,
		userID, cafeID)
	if err != nil {
		log.GetLog().Errorf("Unable to delete favorite by ids. error: %v", err)
		return err
	}

	return err
}

func (r *FavoritesRepoImp) GetFavoritesByUserID(ctx context.Context, userID int32) ([]*models.Favorite, error) {
	var favorites []*models.Favorite
	rows, err := r.postgres.Query(ctx, "SELECT id, user_id, cafe_id FROM favorites WHERE user_id = $1", userID)
	if err != nil {
		log.GetLog().Errorf("Unable to get favorites by user id. error: %v", err)
	}

	for rows.Next() {
		var favorite models.Favorite
		err := rows.Scan(&favorite.ID, &favorite.UserID, &favorite.CafeID)
		if err != nil {
			log.GetLog().Errorf("Unable to scan favorite. error: %v", err)
			continue
		}
		favorites = append(favorites, &favorite)
	}

	return favorites, err
}
