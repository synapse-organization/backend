package repo

import (
	"barista/pkg/log"
	"barista/pkg/models"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CommentsRepo interface {
	Create(ctx context.Context, comment *models.Comment) error
	GetByID(ctx context.Context, id int32) (*models.Comment, error)
	GetByCafeID(ctx context.Context, cafeID int32) ([]*models.Comment, error)
	// GetLast(ctx context.Context, cafeID int32, limit int32, offset int) ([]*models.Comment, error)
	GetAllByCafeID(ctx context.Context, cafeID int32) ([]*models.Comment, error)
}

type CommentsRepoImp struct {
	postgres *pgxpool.Pool
}

func NewCommentsRepoImp(postgres *pgxpool.Pool) *CommentsRepoImp {
	_, err := postgres.Exec(context.Background(),
		`CREATE TABLE IF NOT EXISTS comments (
				id INTEGER PRIMARY KEY,
				user_id INTEGER,
				cafe_id INTEGER,
				comment TEXT,
				date	TIMESTAMP,
				FOREIGN KEY (cafe_id) REFERENCES cafes(id),
				FOREIGN KEY (user_id) REFERENCES users(id)
			);`)
	if err != nil {
		log.GetLog().WithError(err).WithField("table", "comments").Fatal("Unable to create table")
	}
	return &CommentsRepoImp{postgres: postgres}
}

func (c *CommentsRepoImp) Create(ctx context.Context, comment *models.Comment) error {
	_, err := c.postgres.Exec(ctx, "INSERT INTO comments (id, user_id, cafe_id, comment, date) VALUES ($1, $2, $3, $4, $5)", comment.ID, comment.UserID, comment.CafeID, comment.Comment, comment.Date)
	if err != nil {
		log.GetLog().Errorf("Unable to insert comment. error: %v", err)
	}
	return err
}

func (c *CommentsRepoImp) GetByID(ctx context.Context, id int32) (*models.Comment, error) {
	var comment models.Comment
	err := c.postgres.QueryRow(ctx, "SELECT id, cafe_id, user_id, comment FROM comments WHERE id = $1", id).Scan(&comment.ID, &comment.CafeID, &comment.UserID, &comment.Comment)
	if err != nil {
		log.GetLog().Errorf("Unable to get comment by id. error: %v", err)
	}
	return &comment, err
}

func (c *CommentsRepoImp) GetByCafeID(ctx context.Context, cafeID int32) ([]*models.Comment, error) {
	rows, err := c.postgres.Query(ctx, "SELECT id, cafe_id, user_id, comment FROM comments WHERE cafe_id = $1", cafeID)
	if err != nil {
		log.GetLog().Errorf("Unable to get comments by cafe id. error: %v", err)
	}
	defer rows.Close()

	var comments []*models.Comment
	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(&comment.ID, &comment.CafeID, &comment.UserID, &comment.Comment)
		if err != nil {
			log.GetLog().Errorf("Unable to scan comments. error: %v", err)
		}
		comments = append(comments, &comment)
	}
	return comments, err
}

// func (c *CommentsRepoImp) GetLast(ctx context.Context, cafeID int32, limit int32, offset int) ([]*models.Comment, error) {
// 	rows, err := c.postgres.Query(ctx,
// 		`SELECT id, user_id, cafe_id, comment, date 
// 		FROM comments
// 		WHERE cafe_id = $1
// 		ORDER BY date DESC
// 		LIMIT $2
// 		OFFSET $3`, cafeID, limit, offset)
// 	if err != nil {
// 		log.GetLog().Errorf("Unable to get last %v comments. error: %v", limit, err)
// 	}
// 	defer rows.Close()

// 	var comments []*models.Comment
// 	for rows.Next() {
// 		var comment models.Comment
// 		err := rows.Scan(&comment.ID, &comment.UserID, &comment.CafeID, &comment.Comment, &comment.Date)
// 		if err != nil {
// 			log.GetLog().Errorf("Unable to scan comments. error: %v", err)
// 		}
// 		comments = append(comments, &comment)
// 	}

// 	return comments, err
// }

func (c *CommentsRepoImp) GetAllByCafeID(ctx context.Context, cafeID int32) ([]*models.Comment, error) {
	rows, err := c.postgres.Query(ctx,
		`SELECT id, user_id, cafe_id, comment, date 
		FROM comments
		WHERE cafe_id = $1
		ORDER BY date DESC`, cafeID)
	if err != nil {
		log.GetLog().Errorf("Unable to get all comments. error: %v", err)
	}
	defer rows.Close()

	var comments []*models.Comment
	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(&comment.ID, &comment.UserID, &comment.CafeID, &comment.Comment, &comment.Date)
		if err != nil {
			log.GetLog().Errorf("Unable to scan comments. error: %v", err)
		}
		comments = append(comments, &comment)
	}

	return comments, err
}
