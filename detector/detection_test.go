package detector

import (
	"testing"
)

var detector countryDetector
var detectorChile countryDetector

func init() {
	detector = GetDetector()
	detectorChile = GetDetectorChile()
}

func comparePlacesInfo(t *testing.T, want placeInfo, got placeInfo) {
	if want.country != got.country {
		t.Fatalf("Different countries.\n Want: '%s'\n Got: '%s'", want.country, got.country)
	}
	if want.code != got.code {
		t.Fatalf("Different codes.\n Want: '%s'\n Got: '%s'", want.code, got.code)
	}
	if want.location != got.location {
		t.Fatalf("Different codes.\n Want: %v\n Got: %v", want.location, got.location)
	}
	if want.commune != got.commune {
		t.Fatalf("Different codes.\n Want: %v\n Got: %v", want.commune, got.commune)
	}
}

func TestDetectFromText(t *testing.T) {
	want := placeInfo{
		country:  "United States",
		code:     "us",
		location: location{lat: "37.09024", lon: "-95.712891"},
	}
	got, found := detector.DetectFromText("temblor en Estados Unidos?")
	if !found {
		t.Fatalf("Expected to find %s. Got: nothing found", want.country)
	}
	comparePlacesInfo(t, want, got)
}

func TestDetectFromTextChile(t *testing.T) {
	want := placeInfo{
		country:  "chile",
		code:     "cl",
		location: location{lat: "-35.675147", lon: "-71.542969"},
		commune:  place{name: "temuco", location: location{lat: "-38.7359", lon: "-72.5904"}},
	}
	got, found := detectorChile.DetectFromText("TEMBLOR EN TEMUCO!!")
	if !found {
		t.Fatalf("Expected to find %s. Got: nothing found", want.country)
	}
	comparePlacesInfo(t, want, got)
}
