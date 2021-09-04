package sectr

import (
	"bytes"
	"math"
	"testing"
)

// test helper to appoximate coordinate equality
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
	NewSector(p, 1000, 0, 359)
}

func TestTerminal(t *testing.T) {
	test := []struct {
		origin, expected  Point
		bearing, distance float64
	}{
		{
			origin:   Point{lng: 13.35, lat: 52.45},
			distance: 1112.758,
			bearing:  90,
			expected: Point{lng: 13.3664, lat: 52.45},
		},
		{
			origin:   Point{lng: 0.0, lat: 0.0},
			distance: 10000,
			bearing:  180,
			expected: Point{lng: 0.0, lat: -0.089932},
		},
		{
			origin:   Point{lng: 13.35, lat: -52.45},
			distance: 10000,
			bearing:  180,
			expected: Point{lng: 13.35, lat: -52.539932},
		},
	}

	for _, test := range test {
		got := terminal(test.origin, test.distance, test.bearing)

		if !approxEqual(test.expected.lat, got.lat, 0.00001) {
			t.Errorf("Expected %+v, got: %+v", test.expected.lat, got.lat)
		}

		if !approxEqual(test.expected.lng, got.lng, 0.00001) {
			t.Errorf("Expected %+v, got: %+v", test.expected.lng, got.lng)
		}
	}
}

func TestSector(t *testing.T) {
	origin := Point{lat: 52.25, lng: 13.37}
	s := NewSector(origin, 100, 0, 90)

	sj := s.JSON()

	if !bytes.Contains(sj, []byte(`"type":"Polygon"`)) {
		t.Errorf("Sector geometry should have type Polygon")
	}

	if !bytes.Contains(sj, []byte(`"coordinates":[[[13.37,52.25],[13.37,52.25089932],[13.37012803,52.2508959],[13.3702803,52.2508828],[13.37040491,52.25086448],[13.37055029,52.25083383],[13.37068965,52.25079405],[13.37080006,52.25075423],[13.37092446,52.2506989],[13.37103872,52.25063591],[13.3711253,52.25057807],[13.37121783,52.25050289],[13.37128479,52.25043599],[13.37135219,52.25035138],[13.37140478,52.25026293],[13.37143686,52.25018697],[13.37146091,52.250094],[13.37146896,52.24999999],[13.37,52.25]]]`)) {
		t.Errorf("Sector geometry coordinates should match expected")
	}
}
func TestSectorCircle(t *testing.T) {
	origin := Point{lat: 52.25, lng: 13.37}
	s := NewSector(origin, 100, 0, 0)

	sj := s.JSON()

	if !bytes.Contains(sj, []byte(`"type":"Polygon"`)) {
		t.Errorf("Sector geometry should have type Polygon")
	}

	if !bytes.Contains(sj, []byte(`"coordinates":[[[13.36987197,52.2508959],[13.3697197,52.2508828],[13.36959509,52.25086448],[13.36944971,52.25083383],[13.36931035,52.25079405],[13.36919994,52.25075423],[13.36907554,52.2506989],[13.36896128,52.25063591],[13.3688747,52.25057807],[13.36878217,52.25050289],[13.36871521,52.25043599],[13.36864781,52.25035138],[13.36859522,52.25026293],[13.36856314,52.25018697],[13.36853909,52.250094],[13.36853104,52.24999999],[13.36853664,52.24992161],[13.36855804,52.24982839],[13.36858796,52.24975211],[13.36863802,52.2496631],[13.368703,52.24957779],[13.36876804,52.24951019],[13.36885842,52.24943403],[13.36896131,52.24936408],[13.36905579,52.24931108],[13.36917858,52.24925443],[13.36928785,52.24921343],[13.36942604,52.24917217],[13.36957053,52.24913997],[13.36969459,52.24912033],[13.36984646,52.24910561],[13.37,52.24910068],[13.37012803,52.2491041],[13.37028028,52.2491172],[13.37040489,52.24913552],[13.37055027,52.24916616],[13.37068962,52.24920595],[13.37080004,52.24924576],[13.37092443,52.24930109],[13.37103869,52.24936408],[13.37112527,52.24942192],[13.37121781,52.2494971],[13.37128477,52.24956399],[13.37135217,52.2496486],[13.37140476,52.24973706],[13.37143685,52.24981301],[13.37146091,52.24990599],[13.37146896,52.24999999],[13.37146337,52.25007837],[13.37144197,52.25017159],[13.37141206,52.25024788],[13.371362,52.25033688],[13.37129702,52.2504222],[13.37123198,52.2504898],[13.37114161,52.25056596],[13.37103872,52.25063591],[13.37094424,52.25068892],[13.37082144,52.25074557],[13.37071218,52.25078656],[13.37057398,52.25082783],[13.37042949,52.25086002],[13.37030542,52.25087967],[13.37015355,52.25089439],[13.36987197,52.2508959]]]`)) {
		t.Errorf("Sector geometry coordinates should match expected")
	}
}
