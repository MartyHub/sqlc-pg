// Code generated by github.com/MartyHub/sqlc-pg version dev
// DO NOT EDIT

package sqlc

import (
	"context"
	"github.com/jackc/pgx/v5"
	"time"
)

type VenueRepository interface {
	CreateVenue(ctx context.Context, params *CreateVenueParams) (int32, error)
	DeleteVenue(ctx context.Context, slug string) error
	GetVenue(ctx context.Context, params *GetVenueParams) (*ListVenues, error)
	ListVenues(ctx context.Context, city string) ([]*ListVenues, error)
	UpdateVenueName(ctx context.Context, params *ListCities) (int32, error)
	VenueCountByCity(ctx context.Context) ([]*VenueCountByCity, error)
}

type venueRepository struct {
	db Database
}

func NewVenueRepository(db Database) VenueRepository {
	return venueRepository{db: db}
}

const createVenueStmt = `INSERT INTO venue (
    slug,
    name,
    city,
    created_at,
    spotify_playlist,
    status,
    statuses,
    tags
) VALUES (
    $1,
    $2,
    $3,
    NOW(),
    $4,
    $5,
    $6,
    $7
) RETURNING id`

type CreateVenueParams struct {
	Slug            string   `db:"slug"`
	Name            string   `db:"name"`
	City            string   `db:"city"`
	SpotifyPlaylist string   `db:"spotify_playlist"`
	Status          Status   `db:"status"`
	Statuses        []Status `db:"statuses"`
	Tags            []string `db:"tags"`
}

func ScanCreateVenue(row pgx.CollectableRow) (int32, error) {
	var result int32

	err := row.Scan(&result)

	return result, err
}

func (repo venueRepository) CreateVenue(ctx context.Context, params *CreateVenueParams) (int32, error) {
	rows, err := repo.db.Query(ctx, createVenueStmt,
		params.Slug,
		params.Name,
		params.City,
		params.SpotifyPlaylist,
		params.Status,
		params.Statuses,
		params.Tags,
	)
	if err != nil {
		var result int32

		return result, err
	}

	return CollectExactlyOneRow(rows, ScanCreateVenue)
}

const deleteVenueStmt = `DELETE FROM venue
WHERE slug = $1 AND slug = $1`

func (repo venueRepository) DeleteVenue(ctx context.Context, slug string) error {
	_, err := repo.db.Exec(ctx, deleteVenueStmt, slug)

	return err
}

const getVenueStmt = `SELECT id, status, statuses, slug, name, city, spotify_playlist, songkick_id, tags, created_at
FROM venue
WHERE slug = $1 AND city = $2`

type GetVenueParams struct {
	Slug string `db:"slug"`
	City string `db:"city"`
}

func (repo venueRepository) GetVenue(ctx context.Context, params *GetVenueParams) (*ListVenues, error) {
	rows, err := repo.db.Query(ctx, getVenueStmt,
		params.Slug,
		params.City,
	)
	if err != nil {
		return nil, err
	}

	return CollectExactlyOneRow(rows, ScanListVenues)
}

const listVenuesStmt = `SELECT id, status, statuses, slug, name, city, spotify_playlist, songkick_id, tags, created_at
FROM venue
WHERE city = $1
ORDER BY name`

type ListVenues struct {
	ID              int32     `db:"id"`
	Status          Status    `db:"status"`
	Statuses        []Status  `db:"statuses"`
	Slug            string    `db:"slug"`
	Name            string    `db:"name"`
	City            string    `db:"city"`
	SpotifyPlaylist string    `db:"spotify_playlist"`
	SongkickID      *string   `db:"songkick_id"`
	Tags            []string  `db:"tags"`
	CreatedAt       time.Time `db:"created_at"`
}

func ScanListVenues(row pgx.CollectableRow) (*ListVenues, error) {
	result := new(ListVenues)

	if err := row.Scan(
		&result.ID,
		&result.Status,
		&result.Statuses,
		&result.Slug,
		&result.Name,
		&result.City,
		&result.SpotifyPlaylist,
		&result.SongkickID,
		&result.Tags,
		&result.CreatedAt,
	); err != nil {
		return nil, err
	}

	return result, nil
}

func (repo venueRepository) ListVenues(ctx context.Context, city string) ([]*ListVenues, error) {
	rows, err := repo.db.Query(ctx, listVenuesStmt, city)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, ScanListVenues)
}

const updateVenueNameStmt = `UPDATE venue
SET name = $2
WHERE slug = $1
RETURNING id`

func (repo venueRepository) UpdateVenueName(ctx context.Context, params *ListCities) (int32, error) {
	rows, err := repo.db.Query(ctx, updateVenueNameStmt,
		params.Slug,
		params.Name,
	)
	if err != nil {
		var result int32

		return result, err
	}

	return CollectExactlyOneRow(rows, ScanCreateVenue)
}

const venueCountByCityStmt = `SELECT
    city,
    count(*)
FROM venue
GROUP BY 1
ORDER BY 1`

type VenueCountByCity struct {
	City  string `db:"city"`
	Count int    `db:"count"`
}

func ScanVenueCountByCity(row pgx.CollectableRow) (*VenueCountByCity, error) {
	result := new(VenueCountByCity)

	if err := row.Scan(
		&result.City,
		&result.Count,
	); err != nil {
		return nil, err
	}

	return result, nil
}

func (repo venueRepository) VenueCountByCity(ctx context.Context) ([]*VenueCountByCity, error) {
	rows, err := repo.db.Query(ctx, venueCountByCityStmt)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, ScanVenueCountByCity)
}