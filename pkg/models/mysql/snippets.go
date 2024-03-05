package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"showserenity.net/car-rental-system/pkg/models"
)

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(carsId int, title, image, content, expires string) (int, error) {
	stmt := `INSERT INTO snippets (title, image_url, cars_id, content, created, expires)
	VALUES(?, ?, ?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`
	result, err := m.DB.Exec(stmt, title, image, carsId, content, expires)
	if err != nil {
		return 0, nil
	}
	fmt.Print("successfully inserted: ", result)

	id, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}

	return int(id), nil
}

func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	stmt := `SELECT
    snippets.id,
    snippets.title,
    cars.image_url,
    snippets.content,
    snippets.created,
    snippets.expires,
    cars.id,
    cars.model,
    cars.carType,
    cars.seats,
    cars.color,
    cars.location
FROM snippets
JOIN cars ON snippets.cars_id = cars.id
WHERE snippets.id = ?`

	row := m.DB.QueryRow(stmt, id)

	s := &models.Snippet{}

	err := row.Scan(&s.ID, &s.Title, &s.Image, &s.Content, &s.Created, &s.Expires, &s.CarsID, &s.Model, &s.CarType, &s.Seats, &s.Color, &s.Location)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

func (m *SnippetModel) LatestSnippets() ([]*models.Snippet, error) {
	stmt := `SELECT
    snippets.id,
    snippets.title,
    cars.image_url,
    snippets.content,
    snippets.created,
    snippets.expires,
    cars.model,
    cars.carType,
    cars.seats,
    cars.color,
    cars.location
FROM snippets
JOIN cars ON snippets.cars_id = cars.id
ORDER BY snippets.created DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	snippets := []*models.Snippet{}

	for rows.Next() {
		s := &models.Snippet{}

		err := rows.Scan(&s.ID, &s.Title, &s.Image, &s.Content, &s.Created, &s.Expires, &s.Model, &s.CarType, &s.Seats, &s.Color, &s.Location)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}

func (m *SnippetModel) GetByType(carsType string) ([]*models.Snippet, error) {
	stmt := `SELECT
    snippets.id,
    snippets.title,
    cars.image_url,
    snippets.content,
    snippets.created,
    snippets.expires,
    cars.model,
    cars.carType,
    cars.seats,
    cars.color,
    cars.location
FROM snippets
JOIN cars ON snippets.cars_id = cars.id
WHERE cars.carType = ?;
`

	rows, err := m.DB.Query(stmt, carsType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	carsList := []*models.Snippet{}

	for rows.Next() {
		s := &models.Snippet{}
		err := rows.Scan(&s.ID, &s.Title, &s.Image, &s.Content, &s.Created, &s.Expires, &s.Model, &s.CarType, &s.Seats, &s.Color, &s.Location)
		if err != nil {
			return nil, err
		}
		carsList = append(carsList, s)
	}

	return carsList, nil
}
