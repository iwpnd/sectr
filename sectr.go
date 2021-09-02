package sectr

import (
	"math"
)

// Point ...
type Point struct {
	lng, lat float64
}

// Sector ...
type Sector struct {
	Type        string        `json:"type"`
	Coordinates [][][]float64 `json:"coordinates"`
}

// NewPoint ...
func NewPoint(lng, lat float64) Point {
	return Point{lng, lat}
}

func radToDegree(rad float64) float64 {
	return rad * 180 / math.Pi
}

func degreeToRad(degree float64) float64 {
	return degree * math.Pi / 180
}

func distanceToRadians(distance float64) float64 {
	const factor = 6371008.8 // earth radius

	return distance / factor
}

func bearingToAngle(bearing float64) float64 {
	angle := math.Mod(bearing, 360)

	if angle < 0 {
		angle = angle + 360
	}

	return angle
}

func destination(start Point, distance, bearing float64) Point {
	φ1 := degreeToRad(start.lat)
	λ1 := degreeToRad(start.lng)
	bearingRad := degreeToRad(bearing)
	distanceRad := distanceToRadians(distance)

	φ2 := math.Sin(
		math.Sin(φ1)*math.Cos(distanceRad) +
			math.Cos(φ1)*math.Sin(distanceRad)*math.Cos(bearingRad))

	λ2 := λ1 + math.Atan2(
		math.Sin(bearingRad)*
			math.Sin(distanceRad)*
			math.Cos(φ1),
		math.Cos(distanceRad)*
			math.Sin(φ1)*
			math.Sin(φ2))

	lng := radToDegree(λ2)
	lat := radToDegree(φ2)

	return Point{lng: lng, lat: lat}
}

// CreateSector ...
func CreateSector(center Point, radius, bearing1, bearing2 float64) []Point {
	const steps = 64 // fix for now
	var endDegree float64

	angle1 := bearingToAngle(bearing1)
	angle2 := bearingToAngle(bearing2)

	startDegree := angle1

	if angle1 < angle2 {
		endDegree = angle2
	} else {
		endDegree = angle2 + 360
	}

	α := startDegree
	var sector []Point

	sector = append(sector, center)

	for i := 1; ; i++ {
		if α < endDegree {
			sector = append(sector, destination(center, radius, α))
			α = startDegree + float64((i*360)/steps)
		}

		if α >= endDegree {
			sector = append(sector, destination(center, radius, endDegree))
			sector = append(sector, center)
			return sector
		}
	}
}
