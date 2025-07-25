package partner

import (
	"DistanceTrackerServer/models"
	"DistanceTrackerServer/utils"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"time"
)

const (
	savedLocationsToRetrieve = 3
)

var (
	emailFromContext  = utils.EmailFromContext
	dbConnFromContext = utils.DBConnFromContext
	sin               = math.Sin
	cos               = math.Cos
)

func retrievePartnerLocation(dbConn *sql.DB, userId int) (models.LocationFromDB, error) {
	partnerId, err := utils.GetPartnerIdByUserId(dbConn, userId)
	if err != nil {
		return models.LocationFromDB{}, fmt.Errorf("failed to retrieve partner ID: %w", err)
	}

	if partnerId == 0 {
		return models.LocationFromDB{}, fmt.Errorf("no partner linked for user ID %d", userId)
	}

	query := `
		SELECT latitude, longitude, created_at FROM locations 
		WHERE user_id = ? AND is_valid = TRUE
		ORDER BY created_at DESC 
		LIMIT 1`
	row := dbConn.QueryRow(query, partnerId)

	var partnerLocation models.LocationFromDB
	err = row.Scan(&partnerLocation.Latitude, &partnerLocation.Longitude, &partnerLocation.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.LocationFromDB{}, fmt.Errorf("no valid location found for partner ID %d", partnerId)
		}
		return models.LocationFromDB{}, fmt.Errorf("failed to scan partner location: %w", err)
	}

	return partnerLocation, nil
}

func insertLocationToDB(location models.Location, dbConn *sql.DB, userId int, isValid bool) error {
	query := `
		INSERT INTO locations (user_id, latitude, longitude, created_at, is_valid) 
		VALUES (?, ?, ?, ?, ?)`
	_, err := dbConn.Exec(query, userId, location.Latitude, location.Longitude, time.Now(), isValid)
	if err != nil {
		return fmt.Errorf("failed to insert location into database: %w", err)
	}
	return nil
}

func validateDistanceRequest(currentLocation models.Location, dbConn *sql.DB, userId int) error {
	if currentLocation.Latitude < -90 || currentLocation.Latitude > 90 {
		return fmt.Errorf("latitude must be between -90 and 90")
	}
	if currentLocation.Longitude < -180 || currentLocation.Longitude > 180 {
		return fmt.Errorf("longitude must be between -180 and 180")
	}

	locations, err := getLastNLocations(userId, dbConn, savedLocationsToRetrieve)
	if err != nil {
		return fmt.Errorf("failed to retrieve last locations: %w", err)
	}

	// If we don't have enough locations, we can't validate the distance request, so we allow it to proceed
	if len(locations) < (savedLocationsToRetrieve) {
		return nil
	}

	const maxSpeedKmh = 150.0           // Standard hard limit
	const highDistanceThreshold = 300.0 // km, allow high speed for longer distances
	const insaneSpeedThreshold = 1000.0 // km/h, probably spoofed or plane teleport

	newDistance := calculateDistance(locations[0].ToLocation(), currentLocation)
	newTimeDiff := time.Since(locations[0].CreatedAt).Hours()

	if newTimeDiff <= 0 {
		return fmt.Errorf("invalid time difference for new location")
	}

	newSpeed := newDistance / newTimeDiff

	// Reject speeds that are implausibly fast under any circumstance
	if newSpeed > insaneSpeedThreshold {
		return fmt.Errorf("speed %.2f km/h is unreasonably high", newSpeed)
	}

	// Allow high-speed travel if distance is very large
	if newDistance > highDistanceThreshold {
		// Only allow if the time justifies it
		if newSpeed > 900.0 { // Max reasonable commercial plane speed
			return fmt.Errorf("high distance (%.2f km) but speed %.2f km/h exceeds realistic air travel", newDistance, newSpeed)
		}
		// Otherwise allow — user may be on a train or flight
		return nil
	}

	// Normal validation for everyday movement
	if newSpeed > maxSpeedKmh {
		return fmt.Errorf("new speed %.2f km/h exceeds maximum allowed speed of %.2f km/h", newSpeed, maxSpeedKmh)
	}

	return nil
}

func getLastNLocations(userID int, dbConn *sql.DB, n int) ([]models.LocationFromDB, error) {
	query := `
		SELECT latitude, longitude, created_at FROM locations 
        WHERE user_id = ? AND is_valid = TRUE
        ORDER BY created_at DESC 
        LIMIT ?`
	rows, err := dbConn.Query(query, userID, n)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve last %d locations for user %d: %w", n, userID, err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Printf("Error closing rows: %v\n", err)
		}
	}(rows)

	var locations []models.LocationFromDB
	for rows.Next() {
		var loc models.LocationFromDB
		if err := rows.Scan(&loc.Latitude, &loc.Longitude, &loc.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan location: %w", err)
		}
		loc.UserID = userID
		locations = append(locations, loc)
	}

	return locations, nil
}

func calculateDistance(loc1, loc2 models.Location) float64 {
	// Haversine formula to calculate the distance between two points on the Earth
	const R = 6371 // Radius of the Earth in kilometers
	lat1 := degreesToRadians(loc1.Latitude)
	lon1 := degreesToRadians(loc1.Longitude)
	lat2 := degreesToRadians(loc2.Latitude)
	lon2 := degreesToRadians(loc2.Longitude)

	dlat := lat2 - lat1
	dlon := lon2 - lon1

	a := sin(dlat/2)*sin(dlat/2) +
		cos(lat1)*cos(lat2)*
			sin(dlon/2)*sin(dlon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return math.Abs(R * c) // Distance in kilometers
}

func degreesToRadians(degrees float64) float64 {
	return degrees * (3.141592653589793 / 180)
}
