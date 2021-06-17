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
	if want.Commune != got.Commune {
		t.Fatalf("Different Communes.\n Want: %v\n Got: %v", want.Commune, got.Commune)
	}
}

// Detect from text
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
		Name:      "chile",
		Code:      "cl",
		Latitude:  "-35.675147",
		Longitude: "-71.542969",
		Commune:   locations.Place{Name: "temuco", Latitude: "-38.7359", Longitude: "-72.5904"},
	}
	got, found := detectorChile.DetectFromText("TEMBLOR EN TEMUCO!!")
	if !found {
		t.Fatalf("Expected to find %s. Got: nothing found", want.Commune.Name)
	}
	comparePlacesInfo(t, want, got)
}

func TestDetectFromTextChileLongCommune(t *testing.T) {
	want := PlaceInfo{
		Name:      "chile",
		Code:      "cl",
		Latitude:  "-35.675147",
		Longitude: "-71.542969",
		Commune:   locations.Place{Name: "san vicente de tagua tagua", Latitude: "-34.2812", Longitude: "-71.8571"},
	}
	got, found := detectorChile.DetectFromText("pasando la tarde en San Vicente de Tagua Tagua")
	if !found {
		t.Fatalf("Expected to find %s. Got: nothing found", want.Commune.Name)
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

// // Detect from List
// func TestDetectFromList(t *testing.T) {
// 	want := PlaceInfo{
// 		Name:      "argentina",
// 		Code:      "ar",
// 		Latitude:  "-38.416097",
// 		Longitude: "-63.616672",
// 	}

// 	got, found := detector.DetectFromList([]string{"argentina", "is", "beautiful"})
// 	if !found {
// 		t.Fatalf("Expected to find '%s'. Got: nothing found", want.Name)
// 	}
// 	comparePlacesInfo(t, want, got)
// }

// func TestDetectFromListNoCountry(t *testing.T) {
// 	got, found := detector.DetectFromList([]string{"I", "love", "charmander"})
// 	if found {
// 		t.Fatalf("Expected to not find any Country. Got: %s", got.Name)
// 	}
// }

// func TestDetectFromListChile(t *testing.T) {
// 	want := PlaceInfo{
// 		Name:     "chile",
// 		Code:     "cl",
// 		Latitude: "-35.675147", Longitude: "-71.542969",
// 		Commune: locations.Place{Name: "alto hospicio", Latitude: "-20.2687", Longitude: "-70.1049"},
// 	}
// 	got, found := detectorChile.DetectFromList([]string{"Please", "come", "to", "alto", "hospicio"})
// 	if !found {
// 		t.Fatalf("Expected to find %s. Got: nothing found", want.Commune.Name)
// 	}
// 	comparePlacesInfo(t, want, got)
// }

// func TestDetectFromListChileNoCommune(t *testing.T) {
// 	want := PlaceInfo{
// 		Name:      "chile",
// 		Code:      "cl",
// 		Latitude:  "-35.675147",
// 		Longitude: "-71.542969",
// 	}
// 	got, found := detectorChile.DetectFromList([]string{"chile", "has", "great", "beaches"})
// 	if !found {
// 		t.Fatalf("Expected to find %s. Got: nothing found", want.Name)
// 	}
// 	comparePlacesInfo(t, want, got)
// }

// // Detect from Map
// func TestDetectFromMap(t *testing.T) {
// 	want := PlaceInfo{
// 		Name:      "afghanistan",
// 		Code:      "af",
// 		Latitude:  "33.93911",
// 		Longitude: "67.709953",
// 	}

// 	got, found := detector.DetectFromMap(map[string]int{"Αφγανιστάν": 1})
// 	if !found {
// 		t.Fatalf("Expected to find '%s'. Got: nothing found", want.Name)
// 	}
// 	comparePlacesInfo(t, want, got)
// }

// func TestDetectFromMapNoCountry(t *testing.T) {
// 	got, found := detector.DetectFromMap(map[string]int{"I": 1, "love": 3, "charizard": 1})
// 	if found {
// 		t.Fatalf("Expected to not find any Country. Got: %s", got.Name)
// 	}
// }

// func TestDetectFromMapChile(t *testing.T) {
// 	want := PlaceInfo{
// 		Name:     "chile",
// 		Code:     "cl",
// 		Latitude: "-35.675147", Longitude: "-71.542969",
// 		Commune: locations.Place{Name: "los vilos", Latitude: "-31.9122", Longitude: "-71.5112"},
// 	}
// 	got, found := detectorChile.DetectFromMap(map[string]int{"I": 1, "live": 1, "in": 1, "los": 1, "vilos": 1})
// 	if !found {
// 		t.Fatalf("Expected to find %s. Got: nothing found", want.Commune.Name)
// 	}
// 	comparePlacesInfo(t, want, got)
// }

// func TestDetectFromMapChileNoCommune(t *testing.T) {
// 	want := PlaceInfo{
// 		Name:      "chile",
// 		Code:      "cl",
// 		Latitude:  "-35.675147",
// 		Longitude: "-71.542969",
// 	}
// 	got, found := detectorChile.DetectFromMap(map[string]int{"chile": 1, "has": 1, "great": 2, "beaches": 1})
// 	if !found {
// 		t.Fatalf("Expected to find %s. Got: nothing found", want.Name)
// 	}
// 	comparePlacesInfo(t, want, got)
// }
