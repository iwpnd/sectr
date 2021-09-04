package sectr

import (
	"encoding/json"
	"math"
)

const earthRadius = 6371008.8 // earth radius

// Point ...
type Point struct {
	lng, lat float64
}

// Sector ...
type Sector struct {
	Type        string        `json:"type"`
	Coordinates [][][]float64 `json:"coordinates"`
	origin      Point
	radius      float64
	bearing1    float64
	bearing2    float64
}

func radToDegree(rad float64) float64 {
	return rad * 180 / math.Pi
}

func degreeToRad(degree float64) float64 {
	return degree * math.Pi / 180
}

func distanceToRadians(distance float64) float64 {
	const r = earthRadius

	return distance / r
}

// terminal calculates the terminal position travelling a distance
// from a given origin
// see https://www.movable-type.co.uk/scripts/latlong.html
func terminal(start Point, distance, bearing float64) Point {
	φ1 := degreeToRad(start.lat)
	λ1 := degreeToRad(start.lng)
	bearingRad := degreeToRad(bearing)
	distanceRad := distanceToRadians(distance)

	φ2 := math.Asin(
		math.Sin(φ1)*
			math.Cos(distanceRad) +
			math.Cos(φ1)*
				math.Sin(distanceRad)*
				math.Cos(bearingRad))

	λ2 := λ1 + math.Atan2(
		math.Sin(bearingRad)*
			math.Sin(distanceRad)*
			math.Cos(φ1),
		math.Cos(distanceRad)-
			math.Sin(φ1)*
				math.Sin(φ2))

	// cap decimals at .00000001 degree ~= 1.11mm
	lng := math.Round(radToDegree(λ2)*100000000) / 100000000
	lat := math.Round(radToDegree(φ2)*100000000) / 100000000

	return Point{lng: lng, lat: lat}
}

func bearingToAngle(bearing float64) float64 {
	angle := math.Mod(bearing, 360)

	if angle < 0 {
		angle = angle + 360
	}

	return angle
}

// NewSector creates a sector from a given origin point, a radius and two bearings
func NewSector(origin Point, radius, bearing1, bearing2 float64) *Sector {
	// to cap the maximum positions in a sector to 64 terminal + 2* origin
	// the higher the smoother, yet the bigger the coordinate array
	const steps = 64
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

	s := &Sector{
		Type:     "Polygon",
		origin:   origin,
		bearing1: bearing1,
		bearing2: bearing2,
		radius:   radius,
	}

	s.addPoint(origin)

	for i := 1; ; i++ {
		if α < endDegree {
			t := terminal(origin, radius, α)
			s.addPoint(t)

			α = startDegree + float64((i*360)/steps)
		}

		if α >= endDegree {
			t := terminal(origin, radius, endDegree)
			s.addPoint(t)
			s.addPoint(origin)

			return s
		}
	}
}

func (s *Sector) addPoint(p Point) {
	if len(s.Coordinates) == 0 {
		s.Coordinates = append(s.Coordinates, [][]float64{{p.lng, p.lat}})
		return
	}

	s.Coordinates[0] = append(s.Coordinates[0], []float64{p.lng, p.lat})
}

// JSON exports the Sector as json string
func (s Sector) JSON() string {
	j, _ := json.Marshal(s)
	return string(j)
}
