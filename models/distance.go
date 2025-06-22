package models

import "time"

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
