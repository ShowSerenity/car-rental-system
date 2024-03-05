package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"showserenity.net/car-rental-system/pkg/models"
)

type RentModel struct {
	DB *sql.DB
}

func (m *RentModel) InsertRent(renterID, carsID, expires int, bill float64) (int, error) {
	stmt := `INSERT INTO rentbook (renter_id, cars_id, bill, rent_start, rent_end)
	VALUES(?, ?, ?, CURRENT_TIMESTAMP, DATE_ADD(CURRENT_TIMESTAMP, INTERVAL ? MINUTE))`
	result, err := m.DB.Exec(stmt, renterID, carsID, bill, expires)
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

func (m *RentModel) GetRent(id int) (*models.Rent, error) {
	stmt := `SELECT 
    users.name,
    cars.model,
    cars.carType,
    cars.seats,
    cars.color,
    cars.location,
    rentbook.id,
    rentbook.rent_start,
    rentbook.rent_end,
	rentbook.bill
FROM
    rentbook
JOIN
    users ON rentbook.renter_id = users.id
JOIN
    cars ON rentbook.cars_id = cars.id
WHERE rentbook.id = ?`

	row := m.DB.QueryRow(stmt, id)

	s := &models.Rent{}

	err := row.Scan(&s.RenterName, &s.Model, &s.CarType, &s.Seats, &s.Color, &s.Location, &s.ID, &s.RentStart, &s.RentEnd, &s.Bill)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

func (m *RentModel) LatestRents(id int) ([]*models.Rent, error) {
	stmt := `SELECT 
    users.name,
    cars.model,
    cars.carType,
    cars.seats,
    cars.color,
    cars.location,
    rentbook.rent_start,
    rentbook.rent_end,
	rentbook.bill
FROM
    rentbook
JOIN
    users ON rentbook.renter_id = users.id
JOIN
    cars ON rentbook.cars_id = cars.id
WHERE rentbook.renter_id = ? ORDER BY rentbook.rent_start DESC LIMIT 10`

	rows, err := m.DB.Query(stmt, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	rents := []*models.Rent{}

	for rows.Next() {
		s := &models.Rent{}

		err := rows.Scan(&s.RenterName, &s.Model, &s.CarType, &s.Seats, &s.Color, &s.Location, &s.RentStart, &s.RentEnd, &s.Bill)
		if err != nil {
			return nil, err
		}

		rents = append(rents, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return rents, nil
}
