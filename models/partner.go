package models

import (
	"fmt"
	"time"
)

type Location struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

type LocationFromDB struct {
	ID               int
	UserID           int
	Latitude         float64
	Longitude        float64
	CreatedAt        time.Time
	IsValid          bool
	ValidationReason string // optional
}

func (l *LocationFromDB) ToLocation() Location {
	return Location{
		Longitude: l.Longitude,
		Latitude:  l.Latitude,
	}
}

func (l *LocationFromDB) ToString() string {
	return "{latitude: " + fmt.Sprintf("%f", l.Latitude) +
		", longitude: " + fmt.Sprintf("%f", l.Longitude) +
		", created_at: " + l.CreatedAt.String() +
		", is_valid: " + fmt.Sprintf("%t", l.IsValid) +
		", validation_reason: " + l.ValidationReason + "}"
}

type UserInformation struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
}
