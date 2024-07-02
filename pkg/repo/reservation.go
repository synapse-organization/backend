package repo

import (
	"barista/pkg/log"
	"barista/pkg/models"
	"context"
	"math/rand"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ReservationRepo interface {
	Create(ctx context.Context, reservation *models.Reservation) error
	GetByID(ctx context.Context, id int32) (*models.Reservation, error)
	GetByUserID(ctx context.Context, userID int32) (*[]models.Reservation, error)
	GetByCafeID(ctx context.Context, cafeID int32) ([]*models.Reservation, error)
	GetByDate(ctx context.Context, cafeID int32, startTime time.Time, endTime time.Time) (*[]models.Reservation, error)
	CountByTime(ctx context.Context, cafeID int32, startTime time.Time, endTime time.Time) (int32, error)
	GetFullyBookedDays(ctx context.Context, cafeID int32, startDate time.Time, openingTime int8, closingTime int8) ([]time.Time, error)
	GetAvailableTimeSlots(ctx context.Context, cafeID int32, day time.Time, cafeCapacity int32, openingTime int8, closingTime int8) ([]map[string]interface{}, error)
}

type ReservationRepoImp struct {
	postgres *pgxpool.Pool
}

func NewReservationRepoImp(postgres *pgxpool.Pool) *ReservationRepoImp {
	_, err := postgres.Exec(context.Background(),
		`CREATE TABLE IF NOT EXISTS reservations (
				id INTEGER PRIMARY KEY,
				cafe_id INTEGER,
				user_id INTEGER,
				transaction_id TEXT,
				start_time TIMESTAMP,
				end_time TIMESTAMP,
				people INTEGER,
				FOREIGN KEY (cafe_id) REFERENCES cafes(id),
				FOREIGN KEY (user_id) REFERENCES users(id),
				FOREIGN KEY (transaction_id) REFERENCES transactions(id)
			);`)
	if err != nil {
		log.GetLog().WithError(err).WithField("table", "reservations").Fatal("Unable to create table")
	}
	return &ReservationRepoImp{postgres: postgres}
}

func (r *ReservationRepoImp) Create(ctx context.Context, reservation *models.Reservation) error {
	reservation.ID = rand.Int31()
	_, err := r.postgres.Exec(ctx,
		`INSERT INTO reservations (id, cafe_id, user_id, transaction_id, start_time, end_time, people)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		reservation.ID, reservation.CafeID, reservation.UserID, reservation.TransactionID, reservation.StartTime, reservation.EndTime, reservation.People)
	if err != nil {
		log.GetLog().Errorf("Unable to insert reservation. error: %v", err)
	}
	return err
}

func (r *ReservationRepoImp) GetByID(ctx context.Context, id int32) (*models.Reservation, error) {
	var reservation models.Reservation
	err := r.postgres.QueryRow(ctx, "SELECT id, cafe_id, user_id, transaction_id, start_time, end_time, people FROM reservations WHERE id = $1", id).Scan(&reservation.ID, &reservation.CafeID, &reservation.UserID, &reservation.TransactionID, &reservation.StartTime, &reservation.EndTime, &reservation.People)
	if err != nil {
		log.GetLog().Errorf("Unable to get reservation by id. error: %v", err)
	}
	return &reservation, err
}

func (r *ReservationRepoImp) GetByUserID(ctx context.Context, userID int32) (*[]models.Reservation, error) {
	rows, err := r.postgres.Query(ctx, "SELECT id, cafe_id, user_id, transaction_id, start_time, end_time, people FROM reservations WHERE user_id = $1", userID)
	if err != nil {
		log.GetLog().Errorf("Unable to get reservations by user id. error: %v", err)
	}
	defer rows.Close()

	var reservations []models.Reservation
	for rows.Next() {
		var reservation models.Reservation
		err = rows.Scan(&reservation.ID, &reservation.CafeID, &reservation.UserID, &reservation.TransactionID, &reservation.StartTime, &reservation.EndTime, &reservation.People)
		if err != nil {
			log.GetLog().Errorf("Unable to scan reservation. error: %v", err)
			break
		}
		reservations = append(reservations, reservation)
	}
	return &reservations, err
}

func (r *ReservationRepoImp) GetByCafeID(ctx context.Context, cafeID int32) ([]*models.Reservation, error) {
	rows, err := r.postgres.Query(ctx, "SELECT id, cafe_id, user_id, transaction_id, start_time, end_time, people FROM reservations WHERE cafe_id = $1", cafeID)
	if err != nil {
		log.GetLog().Errorf("Unable to get reservations by cafe id. error: %v", err)
	}
	defer rows.Close()

	var reservations []*models.Reservation
	for rows.Next() {
		var reservation models.Reservation
		err = rows.Scan(&reservation.ID, &reservation.CafeID, &reservation.UserID, &reservation.TransactionID, &reservation.StartTime, &reservation.EndTime, &reservation.People)
		if err != nil {
			log.GetLog().Errorf("Unable to scan reservation. error: %v", err)
			break
		}
		reservations = append(reservations, &reservation)
	}
	return reservations, err
}

func (r *ReservationRepoImp) CountByTime(ctx context.Context, cafeID int32, startTime time.Time, endTime time.Time) (int32, error) {
	var totalPeople int32
	query := `
		SELECT COALESCE(SUM(people), 0)
		FROM reservations
		WHERE cafe_id = $1
		AND start_time <= $2
		AND end_time >= $3
	`
	err := r.postgres.QueryRow(ctx, query, cafeID, startTime, endTime).Scan(&totalPeople)
	if err != nil {
		log.GetLog().Errorf("Unable to count reservations. error: %v", err)
		return 0, err
	}
	return totalPeople, nil
}

func (r *ReservationRepoImp) GetFullyBookedDays(ctx context.Context, cafeID int32, startDate time.Time, openingTime int8, closingTime int8) ([]time.Time, error) {
    query := `
        SELECT date_trunc('day', start_time) AS day, COUNT(*)
        FROM reservations
        WHERE cafe_id = $1 AND start_time >= $2 AND
              EXTRACT(HOUR FROM start_time) >= $3 AND EXTRACT(HOUR FROM start_time) < $4
        GROUP BY day
        HAVING COUNT(*) >= ($4 - $3)
    `

    rows, err := r.postgres.Query(ctx, query, cafeID, startDate, openingTime, closingTime)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var fullyBookedDays []time.Time
    for rows.Next() {
        var day time.Time
        if err := rows.Scan(&day); err != nil {
            return nil, err
        }
        fullyBookedDays = append(fullyBookedDays, day)
    }

    return fullyBookedDays, nil
}


func (r *ReservationRepoImp) GetAvailableTimeSlots(ctx context.Context, cafeID int32, day time.Time, cafeCapacity int32, openingTime int8, closingTime int8) ([]map[string]interface{}, error) {
    dayStart := day.Add(time.Duration(openingTime) * time.Hour)
    dayEnd := day.Add(time.Duration(closingTime-1) * time.Hour)

    query := `
        WITH time_slots AS (
            SELECT generate_series($2::timestamp, $3::timestamp, '1 hour') AS slot_time
        )
        SELECT 
            time_slots.slot_time,
            ($4 - COALESCE(SUM(reservations.people), 0)) AS remaining_capacity
        FROM time_slots
        LEFT JOIN reservations ON time_slots.slot_time = reservations.start_time AND reservations.cafe_id = $1
        GROUP BY time_slots.slot_time
        HAVING ($4 - COALESCE(SUM(reservations.people), 0)) > 0
		ORDER BY time_slots.slot_time
    `

    rows, err := r.postgres.Query(ctx, query, cafeID, dayStart, dayEnd, cafeCapacity)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var timeSlots []map[string]interface{}
    for rows.Next() {
        var slotTime time.Time
        var remainingCapacity int32
        if err := rows.Scan(&slotTime, &remainingCapacity); err != nil {
            return nil, err
        }
        timeSlots = append(timeSlots, map[string]interface{}{
            "slot_time":         slotTime,
            "remaining_capacity": remainingCapacity,
        })
    }

    return timeSlots, nil
}


func (r *ReservationRepoImp) GetByDate(ctx context.Context, cafeID int32, startTime time.Time, endTime time.Time) (*[]models.Reservation, error) {
	rows, err := r.postgres.Query(ctx,
		`SELECT *
		FROM reservations
		WHERE cafe_id = $1
		AND start_time >= $2
		AND end_time <= $3
		ORDER BY start_time, end_time`,
		cafeID, startTime, endTime)
	if err != nil {
		log.GetLog().Errorf("Unable to get reservation by date. error: %v", err)
		return nil, err
	}

	var reservations []models.Reservation
	for rows.Next() {
		reservation := models.Reservation{}
		err = rows.Scan(&reservation.ID, &reservation.CafeID, &reservation.UserID, &reservation.TransactionID, &reservation.StartTime, &reservation.EndTime, &reservation.People)
		if err != nil {
			log.GetLog().Errorf("Unable to get reservation by date. error: %v", err)
			return nil, err
		}
		reservations = append(reservations, reservation)
	}

	return &reservations, nil
}