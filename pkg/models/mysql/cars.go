package mysql

import (
	"database/sql"
	"errors"
	"showserenity.net/car-rental-system/pkg/models"
)

type CarModel struct {
	DB *sql.DB
}

func (m *CarModel) GetCar(id int) (*models.Car, error) {
	stmt := `SELECT id, model, carType, seats, color, location FROM cars
WHERE location = ?`

	row := m.DB.QueryRow(stmt, id)

	s := &models.Car{}

	err := row.Scan(&s.ID, &s.Model, &s.CarType, &s.Seats, &s.Color, &s.Location)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

func (m *CarModel) GetCars(id int) ([]*models.Car, error) {
	stmt := `SELECT id, model, carType, seats, color, location FROM cars
WHERE location = ?`

	rows, err := m.DB.Query(stmt, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	carsList := []*models.Car{}

	for rows.Next() {
		s := &models.Car{}
		err := rows.Scan(&s.ID, &s.Model, &s.CarType, &s.Seats, &s.Color, &s.Location)
		if err != nil {
			return nil, err
		}
		carsList = append(carsList, s)
	}

	return carsList, nil
}
