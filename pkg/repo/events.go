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
			(61, 1, 'Live Jazz Night', 'Enjoy an evening of live jazz music.', '2024-07-10 19:00:00', '2024-07-10 21:00:00', 20.0, 50, 30, true),
			(62, 1, 'Art Exhibition', 'Showcasing local artists.', '2024-07-11 10:00:00', '2024-07-11 17:00:00', 10.0, 100, 60, true),
			(63, 1, 'Coffee Tasting', 'Sample different coffee blends.', '2024-07-12 14:00:00', '2024-07-12 16:00:00', 15.0, 30, 25, true),
			(64, 4, 'Vegan Cooking Workshop', 'Learn to cook vegan dishes.', '2024-07-13 11:00:00', '2024-07-13 14:00:00', 25.0, 20, 18, true),
			(65, 5, 'Kids Storytelling', 'Fun storytelling session for kids.', '2024-07-14 10:00:00', '2024-07-14 11:00:00', 5.0, 40, 35, true),
			(66, 6, 'International Cuisine Night', 'Experience dishes from around the world.', '2024-07-15 18:00:00', '2024-07-15 21:00:00', 30.0, 60, 50, true),
			(67, 7, 'Poetry Reading', 'Enjoy an evening of poetry readings.', '2024-07-16 19:00:00', '2024-07-16 20:30:00', 10.0, 40, 35, true),
			(68, 8, 'Herbal Tea Workshop', 'Learn about and taste herbal teas.', '2024-07-17 15:00:00', '2024-07-17 17:00:00', 12.0, 30, 25, true),
			(69, 9, 'Rock Band Performance', 'Live performance by a local rock band.', '2024-07-18 20:00:00', '2024-07-18 22:00:00', 25.0, 70, 60, true),
			(70, 10, 'Mindfulness Meditation', 'Guided mindfulness meditation session.', '2024-07-19 09:00:00', '2024-07-19 10:00:00', 8.0, 20, 15, true),
			(71, 11, 'Creative Writing Workshop', 'Improve your creative writing skills.', '2024-07-20 14:00:00', '2024-07-20 17:00:00', 20.0, 25, 20, true),
			(72, 12, 'Board Games Night', 'Enjoy an evening of board games.', '2024-07-21 18:00:00', '2024-07-21 21:00:00', 10.0, 40, 35, true),
			(73, 13, 'Outdoor Movie Night', 'Watch a movie under the stars.', '2024-07-22 20:00:00', '2024-07-22 22:00:00', 15.0, 60, 50, true),
			(74, 14, 'Pet Adoption Day', 'Meet and adopt pets in need of homes.', '2024-07-23 10:00:00', '2024-07-23 14:00:00', 0.0, 100, 80, true),
			(75, 15, 'Dessert Baking Class', 'Learn to bake delicious desserts.', '2024-07-24 13:00:00', '2024-07-24 15:00:00', 20.0, 20, 18, true),
			(76, 16, 'Trivia Night', 'Test your knowledge at our trivia night.', '2024-07-25 19:00:00', '2024-07-25 21:00:00', 10.0, 40, 35, true),
			(77, 17, 'Vegan Food Fair', 'Sample and buy various vegan foods.', '2024-07-26 10:00:00', '2024-07-26 14:00:00', 5.0, 50, 45, true),
			(78, 18, 'Live Acoustic Music', 'Enjoy an evening of acoustic music.', '2024-07-27 19:00:00', '2024-07-27 21:00:00', 15.0, 50, 45, true),
			(79, 19, 'Board Game Tournament', 'Compete in our board game tournament.', '2024-07-28 16:00:00', '2024-07-28 19:00:00', 10.0, 30, 25, true),
			(80, 20, 'Herbal Tea Tasting', 'Taste and learn about herbal teas.', '2024-07-29 14:00:00', '2024-07-29 16:00:00', 12.0, 20, 18, true),
			(81, 21, 'Children''s Storytime', 'Fun storytime session for children.', '2024-07-30 10:00:00', '2024-07-30 11:00:00', 5.0, 30, 25, true),
			(82, 22, 'Creative Writing Session', 'Enhance your creative writing skills.', '2024-07-31 15:00:00', '2024-07-31 17:00:00', 20.0, 25, 20, true),
			(83, 23, 'Outdoor Art Class', 'Join us for an outdoor art class.', '2024-08-01 11:00:00', '2024-08-01 13:00:00', 15.0, 20, 18, true),
			(84, 24, 'Pet Adoption Event', 'Adopt pets from local shelters.', '2024-08-02 10:00:00', '2024-08-02 14:00:00', 0.0, 100, 80, true),
			(85, 25, 'Dessert Tasting', 'Taste a variety of desserts.', '2024-08-03 14:00:00', '2024-08-03 16:00:00', 10.0, 30, 25, true),
			(86, 26, 'Live Music Night', 'Enjoy live music performances.', '2024-08-04 19:00:00', '2024-08-04 21:00:00', 20.0, 50, 45, true),
			(87, 27, 'Vegan Cooking Demo', 'Watch a live vegan cooking demo.', '2024-08-05 11:00:00', '2024-08-05 13:00:00', 15.0, 20, 18, true),
			(88, 28, 'Poetry Open Mic', 'Share your poetry at our open mic.', '2024-08-06 18:00:00', '2024-08-06 20:00:00', 10.0, 40, 35, true),
			(89, 29, 'Rock Music Night', 'Enjoy an evening of rock music.', '2024-08-07 20:00:00', '2024-08-07 22:00:00', 25.0, 60, 50, true),
			(90, 30, 'Mindfulness Workshop', 'Learn mindfulness techniques.', '2024-08-08 09:00:00', '2024-08-08 11:00:00', 10.0, 20, 18, true);
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
