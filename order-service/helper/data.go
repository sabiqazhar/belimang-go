package helper

import (
	"fmt"
	"math"

	"github.com/jackc/pgx/v5/pgtype"
)

func FloatToNumeric(f float64) pgtype.Numeric {
	var n pgtype.Numeric
	// Use string formatting to avoid float precision issues
	err := n.Scan(fmt.Sprintf("%.2f", f))
	if err != nil {
		return pgtype.Numeric{}
	} // .2f matches your DECIMAL(10, 2)
	return n
}

// toRadians converts a degree value to radians.
func toRadians(deg float64) float64 {
	return deg * math.Pi / 180
}

// CalculateHaversineDistance calculates the distance between two points on Earth.
func CalculateHaversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	// Earth's radius in meters.
	const R = 6371000

	// Convert latitude and longitude from degrees to radians.
	radLat1 := toRadians(lat1)
	radLat2 := toRadians(lat2)
	dLat := toRadians(lat2 - lat1)
	dLon := toRadians(lon2 - lon1)

	// Apply the Haversine formula.
	a := math.Pow(math.Sin(dLat/2), 2) + math.Cos(radLat1)*math.Cos(radLat2)*math.Pow(math.Sin(dLon/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	// Calculate the final distance.
	distance := R * c

	return distance // Distance in meters
}

func CalculateDeliveryTime(distanceMeters float64) float64 {
	// The formula is: ((distanceMeters / 1000) / 40) * 60
	// This simplifies to: distanceMeters * 0.0015
	const minutesPerMeter = 0.0015
	return distanceMeters * minutesPerMeter
}
