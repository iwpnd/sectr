package sectr

import (
	"testing"
)

func BenchmarkSectr(b *testing.B) {
	p := Point{lng: 13.37, lat: 52.25}
	sectr := NewSector(p, 1000, 0, 359)
	sectr.JSON()
}
