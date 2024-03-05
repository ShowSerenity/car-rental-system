package models

import (
	"errors"
	"time"
)

var (
	ErrNoRecord           = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail     = errors.New("models: duplicate email")
)

type Snippet struct {
	ID       int
	Title    string
	Image    string
	Content  string
	Created  time.Time
	Expires  time.Time
	CarsID   int
	Model    string
	CarType  string
	Seats    string
	Color    string
	Location string
}

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
	Active         bool
	Avatar         string
	Role           string
}

type Car struct {
	ID             int
	Model          string
	CarType        string
	Seats          int
	Color          string
	Location       string
	Image          string
	AgeRequirement int
	Description    string
	Cost           int
}

type Rent struct {
	ID         int
	RenterID   int
	RenterName string
	CarsID     int
	Model      string
	CarType    string
	Seats      int
	Color      string
	Location   string
	RentStart  time.Time
	RentEnd    time.Time
	Bill       float64
}
