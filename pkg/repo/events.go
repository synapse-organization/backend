package repo

import (
	"barista/pkg/log"
	"barista/pkg/models"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UpdateEventType string

const (
	UpdateEventName          UpdateEventType = "name"
	UpdateEventDescription   UpdateEventType = "description"
	UpdateEventStartTime     UpdateEventType = "start_time"
	UpdateEventEndTime       UpdateEventType = "end_time"
	UpdateEventPrice         UpdateEventType = "price"
	UpdateEventCapacity      UpdateEventType = "capacity"
	UpdateEventAttendees     UpdateEventType = "current_attendees"
	UpdateEventReservability UpdateEventType = "reservable"
)

type EventRepo interface {
	CreateEventForCafe(ctx context.Context, event *models.Event) error
	CreateEventForUser(ctx context.Context, userID int32, eventID int32, transactionID string) error
	GetEventByID(ctx context.Context, id int32) (*models.Event, error)
	GetEventsByCafeID(ctx context.Context, cafeID int32) ([]*models.Event, error)
	GetEventsByUserID(ctx context.Context, userID int32) ([]*models.Event, error)
	GetAllEventsNearestStartTime(ctx context.Context, limit int32) ([]*models.Event, error)
	UpdateEvent(ctx context.Context, id int32, updateEventType UpdateEventType, value interface{}) error
	DeleteByID(ctx context.Context, id int32) error
}

type EventRepoImp struct {
	postgres *pgxpool.Pool
}

func NewEventRepoImp(postgres *pgxpool.Pool) *EventRepoImp {
	_, err := postgres.Exec(context.Background(),
		`CREATE TABLE IF NOT EXISTS events (
				id INTEGER PRIMARY KEY,
				cafe_id INTEGER,
				name TEXT,
				description TEXT,
				start_time TIMESTAMP,
				end_time TIMESTAMP,
				price FLOAT,
				capacity INTEGER,
				current_attendees INTEGER,
				reservable BOOLEAN,
				FOREIGN KEY (cafe_id) REFERENCES cafes(id)
			);`)
	if err != nil {
		log.GetLog().WithError(err).WithField("table", "events").Fatal("Unable to create table")
	}

	_, err = postgres.Exec(context.Background(), `INSERT INTO events (id, cafe_id, name, description, start_time, end_time, price, capacity, current_attendees, reservable)
			VALUES
			(61, 1, 'بازی مافیا', 'اگه پایه یه بازی جذاب مافیا هستی رویداد رو ثبت نام کن و بیا پیشمون.', '2024-07-10 19:00:00', '2024-07-10 21:00:00', 10000.0, 10, 5, true),
			(62, 1, 'مسابقه فوتبال', 'فوتبال دیدن با ما بیشتر خوش میگذره', '2024-07-11 16:00:00', '2024-07-11 17:00:00', 100000.0, 30, 15, true),
			(63, 1, 'شعرخوانی', 'یک روز شعرخوانی کنار حوض زیبای باکارا', '2024-07-12 18:00:00', '2024-07-12 20:00:00', 15000.0, 20, 7, true),
			(64, 2, 'مسابقه باریستا', 'اگه میخوای مهارت های خودت رو به همه نشون بدی تو این مسابقه شرکت کن', '2024-07-13 11:00:00', '2024-07-13 14:00:00', 25000.0, 20, 18, true),
			(65, 2, 'استنداپ کمدی', 'کمدین این برنامه یک سورپرایزه', '2024-07-14 10:00:00', '2024-07-14 11:00:00', 5000.0, 20, 2, true),
			(66, 3, 'تخفیف دانشجویی', 'در این رویداد برای دانشجو های عزیز 20 درصد تخفیف در نظر گرفتیم', '2024-07-15 18:00:00', '2024-07-15 21:00:00', 30.0, 25, 10, true),
			(67, 5, 'جز نایت', '1 ساعت اجرای خواننده های بی نظیر جز', '2024-07-16 19:00:00', '2024-07-16 20:30:00', 10.0, 20, 11, true)
	`)
	if err != nil {
		log.GetLog().Errorf("Unable to insert events. error: %v", err)
	}

	_, err = postgres.Exec(context.Background(),
		`CREATE TABLE IF NOT EXISTS event_reservations (
    				event_id INTEGER,
    				user_id INTEGER,
					transaction_id TEXT,
    				FOREIGN KEY (event_id) REFERENCES events(id),
    				FOREIGN KEY (user_id) REFERENCES users(id),
					FOREIGN KEY (transaction_id) REFERENCES transactions(id)
    			);`)
	if err != nil {
		log.GetLog().WithError(err).WithField("table", "event_reservations").Fatal("Unable to create table")
	}

	return &EventRepoImp{postgres: postgres}
}

func (e *EventRepoImp) CreateEventForCafe(ctx context.Context, event *models.Event) error {
	_, err := e.postgres.Exec(ctx, "INSERT INTO events (id, cafe_id, name, description, start_time, end_time, price, capacity, current_attendees, reservable) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)", event.ID, event.CafeID, event.Name, event.Description, event.StartTime, event.EndTime, event.Price, event.Capacity, event.CurrentAttendees, event.Reservable)
	if err != nil {
		log.GetLog().Errorf("Unable to insert event. error: %v", err)
	}
	return err
}

func (e *EventRepoImp) CreateEventForUser(ctx context.Context, userID int32, eventID int32, transactionID string) error {
	_, err := e.postgres.Exec(ctx, "INSERT INTO event_reservations (event_id, user_id, transaction_id) VALUES ($1, $2, $3)", eventID, userID, transactionID)
	if err != nil {
		log.GetLog().Errorf("Unable to insert event participant. error: %v", err)
	}
	return err
}

func (e *EventRepoImp) GetEventByID(ctx context.Context, id int32) (*models.Event, error) {
	var event models.Event
	err := e.postgres.QueryRow(ctx, "SELECT id, cafe_id, name, description, start_time, end_time, price, capacity, current_attendees, reservable FROM events WHERE id = $1", id).Scan(&event.ID, &event.CafeID, &event.Name, &event.Description, &event.StartTime, &event.EndTime, &event.Price, &event.Capacity, &event.CurrentAttendees, &event.Reservable)
	if err != nil {
		log.GetLog().Errorf("Unable to get event by id. error: %v", err)
	}
	return &event, err
}

func (e *EventRepoImp) GetEventsByCafeID(ctx context.Context, cafeID int32) ([]*models.Event, error) {
	rows, err := e.postgres.Query(ctx, "SELECT id, cafe_id, name, description, start_time, end_time, price, capacity, current_attendees, reservable FROM events WHERE cafe_id = $1", cafeID)
	if err != nil {
		log.GetLog().Errorf("Unable to get events by cafe id. error: %v", err)
	}
	defer rows.Close()

	var events []*models.Event
	for rows.Next() {
		var event models.Event
		err = rows.Scan(&event.ID, &event.CafeID, &event.Name, &event.Description, &event.StartTime, &event.EndTime, &event.Price, &event.Capacity, &event.CurrentAttendees, &event.Reservable)
		if err != nil {
			log.GetLog().Errorf("Unable to scan event. error: %v", err)
			return nil, err
		}
		events = append(events, &event)
	}
	return events, nil
}

func (e *EventRepoImp) GetEventsByUserID(ctx context.Context, userID int32) ([]*models.Event, error) {
	rows, err := e.postgres.Query(ctx, "SELECT e.id, e.cafe_id, e.name, e.description, e.start_time, e.end_time, e.price, e.capacity, e.current_attendees, e.reservable FROM events e JOIN event_reservations ep ON e.id = ep.event_id WHERE ep.user_id = $1", userID)
	if err != nil {
		log.GetLog().Errorf("Unable to get events by user id. error: %v", err)
	}
	defer rows.Close()

	var events []*models.Event
	for rows.Next() {
		var event models.Event
		err = rows.Scan(&event.ID, &event.CafeID, &event.Name, &event.Description, &event.StartTime, &event.EndTime, &event.Price, &event.Capacity, &event.CurrentAttendees, &event.Reservable)
		if err != nil {
			log.GetLog().Errorf("Unable to scan event. error: %v", err)
			return nil, err
		}
		events = append(events, &event)
	}
	return events, nil
}

func (c *EventRepoImp) GetAllEventsNearestStartTime(ctx context.Context, limit int32) ([]*models.Event, error) {
	rows, err := c.postgres.Query(ctx, "SELECT id, cafe_id, name, description, start_time, end_time, price, capacity, current_attendees, reservable FROM events ORDER BY start_time ASC LIMIT $1", limit)
	if err != nil {
		log.GetLog().Errorf("Unable to get all events. error: %v", err)
	}
	defer rows.Close()

	var events []*models.Event
	for rows.Next() {
		var event models.Event
		err = rows.Scan(&event.ID, &event.CafeID, &event.Name, &event.Description, &event.StartTime, &event.EndTime, &event.Price, &event.Capacity, &event.CurrentAttendees, &event.Reservable)
		if err != nil {
			log.GetLog().Errorf("Unable to scan event. error: %v", err)
			return nil, err
		}
		events = append(events, &event)
	}
	return events, nil
}

func (c *EventRepoImp) UpdateEvent(ctx context.Context, id int32, updateEventType UpdateEventType, value interface{}) error {
	columnName := string(updateEventType)

	query := "UPDATE events SET " + columnName + " = $1 WHERE id = $2"
	_, err := c.postgres.Exec(ctx, query, value, id)
	if err != nil {
		log.GetLog().Errorf("Unable to update event. error: %v", err)
		return err
	}

	return nil
}

func (c *EventRepoImp) DeleteByID(ctx context.Context, id int32) error {
	_, err := c.postgres.Exec(ctx,
		`DELETE FROM events
		WHERE id = $1`, id)
	if err != nil {
		log.GetLog().Errorf("Unable to delete event. error: %v", err)
		return err
	}

	return err
}
