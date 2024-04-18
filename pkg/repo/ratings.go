package repo

import (
	"barista/pkg/log"
	"barista/pkg/models"
	"context"
	"github.com/jackc/pgx/v5"
)

type RatingsRepo interface {
	Create(ctx context.Context, rating *models.Rating) error
	GetByID(ctx context.Context, id int32) (*models.Rating, error)
	GetByCafeID(ctx context.Context, cafeID int32) ([]*models.Rating, error)
	GetCafesRating(ctx context.Context, cafeID int32) (float64, error)
	GetByUserID(ctx context.Context, userID int32) ([]*models.Rating, error)
}

type RatingsRepoImp struct {
	postgres *pgx.Conn
}

func NewRatingsRepoImp(postgres *pgx.Conn) *RatingsRepoImp {
	_, err := postgres.Exec(context.Background(),
		`CREATE TABLE IF NOT EXISTS ratings (
				id INTEGER PRIMARY KEY,
				cafe_id INTEGER,
				user_id INTEGER,
				rating INTEGER,
				FOREIGN KEY (cafe_id) REFERENCES cafes(id),
				FOREIGN KEY (user_id) REFERENCES users(id),
    			Unique(cafe_id, user_id)
			);`)
	if err != nil {
		log.GetLog().WithError(err).WithField("table", "ratings").Fatal("Unable to create table")
	}
	return &RatingsRepoImp{postgres: postgres}
}

func (r *RatingsRepoImp) Create(ctx context.Context, rating *models.Rating) error {
	_, err := r.postgres.Exec(ctx, "INSERT INTO ratings (id, cafe_id, user_id, rating) VALUES ($1, $2, $3, $4)", rating.ID, rating.CafeID, rating.UserID, rating.Rating)
	if err != nil {
		log.GetLog().Errorf("Unable to insert rating. error: %v", err)
	}
	return err
}

func (r *RatingsRepoImp) GetByID(ctx context.Context, id int32) (*models.Rating, error) {
	var rating models.Rating
	err := r.postgres.QueryRow(ctx, "SELECT id, cafe_id, user_id, rating FROM ratings WHERE id = $1", id).Scan(&rating.ID, &rating.CafeID, &rating.UserID, &rating.Rating)
	if err != nil {
		log.GetLog().Errorf("Unable to get rating by id. error: %v", err)
	}
	return &rating, err
}

func (r *RatingsRepoImp) GetByCafeID(ctx context.Context, cafeID int32) ([]*models.Rating, error) {
	rows, err := r.postgres.Query(ctx, "SELECT id, cafe_id, user_id, rating FROM ratings WHERE cafe_id = $1", cafeID)
	if err != nil {
		log.GetLog().Errorf("Unable to get ratings by cafe id. error: %v", err)
	}
	defer rows.Close()

	var ratings []*models.Rating
	for rows.Next() {
		var rating models.Rating
		err = rows.Scan(&rating.ID, &rating.CafeID, &rating.UserID, &rating.Rating)
		if err != nil {
			log.GetLog().Errorf("Unable to scan rating. error: %v", err)
			return nil, err
		}
		ratings = append(ratings, &rating)
	}
	return ratings, nil
}

func (r *RatingsRepoImp) GetCafesRating(ctx context.Context, cafeID int32) (float64, error) {
	rows, err := r.postgres.Query(ctx, "SELECT rating FROM ratings WHERE cafe_id = $1", cafeID)
	if err != nil {
		log.GetLog().Errorf("Unable to get ratings by cafe id. error: %v", err)
	}
	defer rows.Close()

	var sum float64
	var count int
	for rows.Next() {
		var rating int
		err = rows.Scan(&rating)
		if err != nil {
			log.GetLog().Errorf("Unable to scan rating. error: %v", err)
			return 0, err
		}
		sum += float64(rating)
		count++
	}
	if count == 0 {
		return 0, nil
	}
	return sum / float64(count), nil
}

func (r *RatingsRepoImp) GetByUserID(ctx context.Context, userID int32) ([]*models.Rating, error) {
	rows, err := r.postgres.Query(ctx, "SELECT id, cafe_id, user_id, rating FROM ratings WHERE user_id = $1", userID)
	if err != nil {
		log.GetLog().Errorf("Unable to get ratings by user id. error: %v", err)
	}
	defer rows.Close()

	var ratings []*models.Rating
	for rows.Next() {
		var rating models.Rating
		err = rows.Scan(&rating.ID, &rating.CafeID, &rating.UserID, &rating.Rating)
		if err != nil {
			log.GetLog().Errorf("Unable to scan rating. error: %v", err)
			return nil, err
		}
		ratings = append(ratings, &rating)
	}
	return ratings, nil
}
