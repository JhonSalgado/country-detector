package detector

import (
	"testing"

	"github.com/JhonSalgado/country-detector/detector/locations"
)

var detector CountryDetector
var detectorChile CountryDetector

func init() {
	detector = GetDetector()
	detectorChile = GetDetectorChile()
}

func comparePlacesInfo(t *testing.T, want PlaceInfo, got PlaceInfo) {
	if want.Name != got.Name {
		t.Fatalf("Different countries.\n Want: '%s'\n Got: '%s'", want.Name, got.Name)
	}
	if want.Code != got.Code {
		t.Fatalf("Different Codes.\n Want: '%s'\n Got: '%s'", want.Code, got.Code)
	}
	if want.Latitude != got.Latitude {
		t.Fatalf("Different latitude.\n Want: %v\n Got: %v", want.Latitude, got.Latitude)
	}
	if want.Longitude != got.Longitude {
		t.Fatalf("Different Longitude.\n Want: %v\n Got: %v", want.Longitude, got.Longitude)
	}
	if want.Municipality != got.Municipality {
		t.Fatalf("Different Communes.\n Want: %v\n Got: %v", want.Municipality, got.Municipality)
	}
}

func TestDetectFromText(t *testing.T) {
	want := PlaceInfo{
		Name:      "united states",
		Code:      "us",
		Latitude:  "37.09024",
		Longitude: "-95.712891",
	}
	got, found := detector.DetectFromText("temblor en Estados Unidos?")
	if !found {
		t.Fatalf("Expected to find '%s'. Got: nothing found", want.Name)
	}
	comparePlacesInfo(t, want, got)
}

func TestDetectFromTextWithDifferentSpaces(t *testing.T) {
	want := PlaceInfo{
		Name:      "united states",
		Code:      "us",
		Latitude:  "37.09024",
		Longitude: "-95.712891",
	}
	gotLeadingSpace, foundLeadingSpace := detector.DetectFromText("I live in usa")
	gotTrailingSpace, foundTrailingSpace := detector.DetectFromText("usa is my country")
	gotNoCountry, foundNoCountry := detector.DetectFromText("usability is important")

	if !foundLeadingSpace {
		t.Fatalf("Expected to find '%s'. Got: nothing found", want.Name)
	}
	comparePlacesInfo(t, want, gotLeadingSpace)

	if !foundTrailingSpace {
		t.Fatalf("Expected to find '%s'. Got: nothing found", want.Name)
	}
	comparePlacesInfo(t, want, gotTrailingSpace)

	if foundNoCountry {
		t.Fatalf("Expected to not find any Country. Got: %s", gotNoCountry.Name)
	}
}

func TestDetectFromTextUnicode(t *testing.T) {
	want := PlaceInfo{
		Name:      "japan",
		Code:      "jp",
		Latitude:  "36.204824",
		Longitude: "138.252924",
	}
	got, found := detector.DetectFromText("日本の地震、津波の危険性があります")
	if !found {
		t.Fatalf("Expected to find '%s'. Got: nothing found", want.Name)
	}
	comparePlacesInfo(t, want, got)
}

func TestDetectFromTextNoCountry(t *testing.T) {
	got, found := detector.DetectFromText("the shaking felt very strong in my city")
	if found {
		t.Fatalf("Expected to not find any Country. Got: %s", got.Name)
	}
}

func TestDetectFromTextChile(t *testing.T) {
	want := PlaceInfo{
		Name:         "chile",
		Code:         "cl",
		Latitude:     "-35.675147",
		Longitude:    "-71.542969",
		Municipality: locations.Place{Name: "santiago", Latitude: "-33.4489", Longitude: "-70.6693"},
	}
	got, found := detectorChile.DetectFromText("estudiando en santiago, chile")
	if !found {
		t.Fatalf("Expected to find %s. Got: nothing found", want.Municipality.Name)
	}
	comparePlacesInfo(t, want, got)
}

func TestDetectFromTextChileCommuneOnly(t *testing.T) {
	want := PlaceInfo{
		Name:         "chile",
		Code:         "cl",
		Latitude:     "-35.675147",
		Longitude:    "-71.542969",
		Municipality: locations.Place{Name: "temuco", Latitude: "-38.7359", Longitude: "-72.5904"},
	}
	got, found := detectorChile.DetectFromText("TEMBLOR EN TEMUCO!!")
	if !found {
		t.Fatalf("Expected to find %s. Got: nothing found", want.Municipality.Name)
	}
	comparePlacesInfo(t, want, got)
}

func TestDetectFromTextChileLongCommune(t *testing.T) {
	want := PlaceInfo{
		Name:         "chile",
		Code:         "cl",
		Latitude:     "-35.675147",
		Longitude:    "-71.542969",
		Municipality: locations.Place{Name: "san vicente de tagua tagua", Latitude: "-34.2812", Longitude: "-71.8571"},
	}
	got, found := detectorChile.DetectFromText("pasando la tarde en San Vicente de Tagua Tagua")
	if !found {
		t.Fatalf("Expected to find %s. Got: nothing found", want.Municipality.Name)
	}
	comparePlacesInfo(t, want, got)
}

func TestDetectFromTextChileNoCommune(t *testing.T) {
	want := PlaceInfo{
		Name:      "chile",
		Code:      "cl",
		Latitude:  "-35.675147",
		Longitude: "-71.542969",
	}
	got, found := detectorChile.DetectFromText("TEMBLOR EN CHILE!!")
	if !found {
		t.Fatalf("Expected to find %s. Got: nothing found", want.Name)
	}
	comparePlacesInfo(t, want, got)
}
