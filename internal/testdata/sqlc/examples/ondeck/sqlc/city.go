// Code generated by github.com/MartyHub/sqlc-pg version dev
// DO NOT EDIT

package sqlc

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type CityRepository interface {
	CreateCity(ctx context.Context, params *CreateCityParams) (*ListCities, error)
	GetCity(ctx context.Context, slug string) (*ListCities, error)
	ListCities(ctx context.Context) ([]*ListCities, error)
	UpdateCityName(ctx context.Context, params *ListCities) error
}

type cityRepository struct {
	db Database
}

func NewCityRepository(db Database) CityRepository {
	return cityRepository{db: db}
}

const createCityStmt = `INSERT INTO city (
    name,
    slug
) VALUES (
    $1,
    $2
) RETURNING slug, name`

type CreateCityParams struct {
	Name string `db:"name"`
	Slug string `db:"slug"`
}

func (repo cityRepository) CreateCity(ctx context.Context, params *CreateCityParams) (*ListCities, error) {
	rows, err := repo.db.Query(ctx, createCityStmt,
		params.Name,
		params.Slug,
	)
	if err != nil {
		return nil, err
	}

	return CollectExactlyOneRow(rows, ScanListCities)
}

const getCityStmt = `SELECT slug, name
FROM city
WHERE slug = $1`

func (repo cityRepository) GetCity(ctx context.Context, slug string) (*ListCities, error) {
	rows, err := repo.db.Query(ctx, getCityStmt, slug)
	if err != nil {
		return nil, err
	}

	return CollectExactlyOneRow(rows, ScanListCities)
}

const listCitiesStmt = `SELECT slug, name
FROM city
ORDER BY name`

type ListCities struct {
	Slug string `db:"slug"`
	Name string `db:"name"`
}

func ScanListCities(row pgx.CollectableRow) (*ListCities, error) {
	result := new(ListCities)

	if err := row.Scan(
		&result.Slug,
		&result.Name,
	); err != nil {
		return nil, err
	}

	return result, nil
}

func (repo cityRepository) ListCities(ctx context.Context) ([]*ListCities, error) {
	rows, err := repo.db.Query(ctx, listCitiesStmt)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, ScanListCities)
}

const updateCityNameStmt = `UPDATE city
SET name = $2
WHERE slug = $1`

func (repo cityRepository) UpdateCityName(ctx context.Context, params *ListCities) error {
	_, err := repo.db.Exec(ctx, updateCityNameStmt,
		params.Slug,
		params.Name,
	)

	return err
}