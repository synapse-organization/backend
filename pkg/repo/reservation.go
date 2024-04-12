package repo

import (
	"barista/pkg/log"
	"barista/pkg/models"
	"context"
	"github.com/jackc/pgx/v5"
)

type ReservationRepo interface {
	Create(ctx context.Context, reservation *models.Reservation) error
	GetByID(ctx context.Context, id int32) (*models.Reservation, error)
	GetByUserID(ctx context.Context, userID int32) ([]*models.Reservation, error)
	GetByCafeID(ctx context.Context, cafeID int32) ([]*models.Reservation, error)
}

type ReservationRepoImp struct {
	postgres *pgx.Conn
}

func NewReservationRepoImp(postgres *pgx.Conn) *ReservationRepoImp {
	_, err := postgres.Exec(context.Background(),
		`CREATE TABLE IF NOT EXISTS reservations (
				id INTEGER PRIMARY KEY,
				cafe_id INTEGER,
				user_id INTEGER,
				start_time TIMESTAMP,
				end_time TIMESTAMP,
				people INTEGER,
				FOREIGN KEY (cafe_id) REFERENCES cafes(id),
				FOREIGN KEY (user_id) REFERENCES users(id)
			);`)
	if err != nil {
		log.GetLog().WithError(err).WithField("table", "reservations").Fatal("Unable to create table")
	}
	return &ReservationRepoImp{postgres: postgres}
}

func (r *ReservationRepoImp) Create(ctx context.Context, reservation *models.Reservation) error {
	_, err := r.postgres.Exec(ctx, "INSERT INTO reservations (id, cafe_id, user_id, start_time, end_time, people) VALUES ($1, $2, $3, $4, $5, $6)", reservation.ID, reservation.CafeID, reservation.UserID, reservation.StartTime, reservation.EndTime, reservation.People)
	if err != nil {
		log.GetLog().Errorf("Unable to insert reservation. error: %v", err)
	}
	return err
}

func (r *ReservationRepoImp) GetByID(ctx context.Context, id int32) (*models.Reservation, error) {
	var reservation models.Reservation
	err := r.postgres.QueryRow(ctx, "SELECT id, cafe_id, user_id, start_time, end_time, people FROM reservations WHERE id = $1", id).Scan(&reservation.ID, &reservation.CafeID, &reservation.UserID, &reservation.StartTime, &reservation.EndTime, &reservation.People)
	if err != nil {
		log.GetLog().Errorf("Unable to get reservation by id. error: %v", err)
	}
	return &reservation, err
}

func (r *ReservationRepoImp) GetByUserID(ctx context.Context, userID int32) ([]*models.Reservation, error) {
	rows, err := r.postgres.Query(ctx, "SELECT id, cafe_id, user_id, start_time, end_time, people FROM reservations WHERE user_id = $1", userID)
	if err != nil {
		log.GetLog().Errorf("Unable to get reservations by user id. error: %v", err)
	}
	defer rows.Close()

	var reservations []*models.Reservation
	for rows.Next() {
		var reservation models.Reservation
		err = rows.Scan(&reservation.ID, &reservation.CafeID, &reservation.UserID, &reservation.StartTime, &reservation.EndTime, &reservation.People)
		if err != nil {
			log.GetLog().Errorf("Unable to scan reservation. error: %v", err)
			break
		}
		reservations = append(reservations, &reservation)
	}
	return reservations, err
}

func (r *ReservationRepoImp) GetByCafeID(ctx context.Context, cafeID int32) ([]*models.Reservation, error) {
	rows, err := r.postgres.Query(ctx, "SELECT id, cafe_id, user_id, start_time, end_time, people FROM reservations WHERE cafe_id = $1", cafeID)
	if err != nil {
		log.GetLog().Errorf("Unable to get reservations by cafe id. error: %v", err)
	}
	defer rows.Close()

	var reservations []*models.Reservation
	for rows.Next() {
		var reservation models.Reservation
		err = rows.Scan(&reservation.ID, &reservation.CafeID, &reservation.UserID, &reservation.StartTime, &reservation.EndTime, &reservation.People)
		if err != nil {
			log.GetLog().Errorf("Unable to scan reservation. error: %v", err)
			break
		}
		reservations = append(reservations, &reservation)
	}
	return reservations, err
}
