package sectr

import (
	"fmt"
	"math"
	"testing"
)

func approxEqual(want, got, tolerance float64) bool {
	diff := math.Abs(want - got)
	mean := math.Abs(want+got) / 2

	if math.IsNaN(diff / mean) {
		return true
	}

	return (diff / mean) < tolerance
}

func BenchmarkSectr(b *testing.B) {
	p := Point{lng: 13.37, lat: 52.25}
	sectr := NewSector(p, 1000, 0, 359)
	sectr.JSON()
}

func TestDestination(t *testing.T) {
	test := []struct {
		start, expected   Point
		bearing, distance float64
	}{
		{
			start:    Point{lng: 13.35, lat: 52.45},
			distance: 1112.758,
			bearing:  90,
			expected: Point{lng: 13.3664, lat: 52.45},
		},
		{
			start:    Point{lng: 0.0, lat: 0.0},
			distance: 10000,
			bearing:  180,
			expected: Point{lng: 0.0, lat: -0.089932},
		},
		{
			start:    Point{lng: 13.35, lat: -52.45},
			distance: 10000,
			bearing:  180,
			expected: Point{lng: 13.35, lat: -52.539932},
		},
	}

	for _, test := range test {
		got := destination(test.start, test.distance, test.bearing)
		fmt.Println("got: ", got)

		if !approxEqual(test.expected.lat, got.lat, 0.00001) {
			t.Errorf("Expected %+v, got: %+v", test.expected.lat, got.lat)
		}

		if !approxEqual(test.expected.lng, got.lng, 0.00001) {
			t.Errorf("Expected %+v, got: %+v", test.expected.lng, got.lng)
		}
	}
}
