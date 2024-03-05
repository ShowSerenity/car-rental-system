package main

import (
	"fmt"
	"math/rand"
	"net/url"
	"time"
)

// Coordinate represents a pair of latitude and longitude.
type Coordinate struct {
	Latitude  float64
	Longitude float64
}

func (app *application) generateRandomMap() string {
	astanaCenter := Coordinate{Latitude: 51.1694, Longitude: 71.4491}
	radiusMeters := 5000

	return generateLocationCode(astanaCenter, radiusMeters)
}

func generateLocationCode(center Coordinate, radius int) string {
	rand.Seed(time.Now().UnixNano())

	deltaX := (rand.Float64() * float64(2*radius)) - float64(radius)
	deltaY := (rand.Float64() * float64(2*radius)) - float64(radius)

	// Calculate the new coordinates based on the center
	newLatitude := center.Latitude + (deltaX / 111000) // Approximate conversion from meters to degrees
	newLongitude := center.Longitude + (deltaY / (111000 * 1.5))
	// Convert latitude and longitude to string format
	latLngStr := fmt.Sprintf("%.6f,%.6f", newLatitude, newLongitude)

	// Encode the location string
	encodedLocation := url.QueryEscape(latLngStr)
	fmt.Print("encodedLocation: ", encodedLocation)

	return constructGoogleMapURL(encodedLocation)
}

func constructGoogleMapURL(coordinate string) string {
	return fmt.Sprintf("https://www.google.com/maps/embed/v1/place?key=AIzaSyCZ8HNo3uR57W6cRFPdwBfhq-HhTiZhsFU&q=%s&zoom=16",
		coordinate)
}
