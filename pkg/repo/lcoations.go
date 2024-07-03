package repo

import (
	"barista/pkg/log"
	"barista/pkg/models"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LocationsRepo interface {
	SetLocation(ctx context.Context, location *models.Location) error
	FindAll(ctx context.Context) ([]*models.Location, error)
	GetCafeLocation(ctx context.Context, id int32) (models.Location, error)
}

type LocationsRepoImp struct {
	postgres *pgxpool.Pool
}

func NewLocationsRepoImp(postgres *pgxpool.Pool) *LocationsRepoImp {
	_, err := postgres.Exec(context.Background(),
		`CREATE TABLE IF NOT EXISTS locations (
				id BIGINT PRIMARY KEY,
				latitude FLOAT,
				longitude FLOAT
			);`)
	if err != nil {
		log.GetLog().WithError(err).WithField("table", "locations").Fatal("Unable to create table")
	}
	return &LocationsRepoImp{postgres: postgres}
}

func (r *LocationsRepoImp) SetLocation(ctx context.Context, location *models.Location) error {
	_, err := r.postgres.Exec(ctx, "INSERT INTO locations (id, latitude, longitude) VALUES ($1, $2, $3) ON CONFLICT (id) DO UPDATE SET latitude = $2, longitude = $3", location.CafeID, location.Lat, location.Lng)
	if err != nil {
		log.GetLog().Errorf("Unable to insert location. error: %v", err)
	}
	return err
}

func (r *LocationsRepoImp) FindAll(ctx context.Context) ([]*models.Location, error) {
	var locations []*models.Location
	rows, err := r.postgres.Query(ctx, "SELECT id, latitude, longitude FROM locations")
	if err != nil {
		log.GetLog().Errorf("Unable to get locations. error: %v", err)
	}
	for rows.Next() {
		var location models.Location
		err := rows.Scan(&location.CafeID, &location.Lat, &location.Lng)
		if err != nil {
			log.GetLog().Errorf("Unable to scan location. error: %v", err)
			continue
		}
		locations = append(locations, &location)
	}
	return locations, err
}

func (r *LocationsRepoImp) GetCafeLocation(ctx context.Context, id int32) (models.Location, error) {
	var location models.Location
	err := r.postgres.QueryRow(ctx, "SELECT id, latitude, longitude FROM locations where id=$1", id).Scan(&location.CafeID, &location.Lat, &location.Lng)
	if err != nil && err.Error() != "no rows in result set" {
		log.GetLog().Errorf("Unable to get locations. error: %v", err)
		return location, err
	}
	return location, nil
}
