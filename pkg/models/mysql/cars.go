package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"showserenity.net/car-rental-system/pkg/models"
)

type CarModel struct {
	DB *sql.DB
}

func (m *CarModel) GetCar(id int) (*models.Car, error) {
	stmt := `SELECT id, model, carType, seats, color, location, image_url, age_requirement, description, cost FROM cars
WHERE id = ?`

	row := m.DB.QueryRow(stmt, id)

	s := &models.Car{}

	err := row.Scan(&s.ID, &s.Model, &s.CarType, &s.Seats, &s.Color, &s.Location, &s.Image, &s.AgeRequirement, &s.Description, &s.Cost)
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
	stmt := `SELECT id, model, carType, seats, color, location, image_url, age_requirement, description, cost FROM cars
WHERE id = ?`

	rows, err := m.DB.Query(stmt, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	carsList := []*models.Car{}

	for rows.Next() {
		s := &models.Car{}
		err := rows.Scan(&s.ID, &s.Model, &s.CarType, &s.Seats, &s.Color, &s.Location, &s.Image, &s.AgeRequirement, &s.Description, &s.Cost)
		if err != nil {
			return nil, err
		}
		carsList = append(carsList, s)
	}

	return carsList, nil
}

func (m *CarModel) GetCarsByType(carsType string) ([]*models.Car, error) {
	stmt := `SELECT id, model, carType, seats, color, location, image_url,
       age_requirement, description, cost
	FROM cars
	WHERE carType = ?`

	rows, err := m.DB.Query(stmt, carsType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	carsList := []*models.Car{}

	for rows.Next() {
		s := &models.Car{}
		err := rows.Scan(&s.ID, &s.Model, &s.CarType, &s.Seats, &s.Color, &s.Location, &s.Image, &s.AgeRequirement, &s.Description, &s.Cost)
		if err != nil {
			return nil, err
		}
		carsList = append(carsList, s)
	}

	return carsList, nil
}

func (m *CarModel) InsertCar(seats, ageRequirement, cost int, model, carType, color, location, imageUrl, description string) (int, error) {
	stmt := `INSERT INTO cars (seats, age_requirement, cost, model, carType, color, location, image_url, description)
	VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := m.DB.Exec(stmt, seats, ageRequirement, cost, model, carType, color, location, imageUrl, description)
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

func (m *CarModel) LatestCars() ([]*models.Car, error) {
	stmt := `SELECT id, model, carType, seats, color, location, image_url,
       age_requirement, description, cost
	FROM cars
	ORDER BY id DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	cars := []*models.Car{}

	for rows.Next() {
		s := &models.Car{}

		err := rows.Scan(&s.ID, &s.Model, &s.CarType, &s.Seats, &s.Color, &s.Location, &s.Image, &s.AgeRequirement, &s.Description, &s.Cost)
		if err != nil {
			return nil, err
		}

		cars = append(cars, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return cars, nil
}
