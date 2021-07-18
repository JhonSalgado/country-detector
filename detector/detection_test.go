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

func TestDetect(t *testing.T) {
	want := PlaceInfo{
		Name:      "united states",
		Code:      "us",
		Latitude:  "37.09024",
		Longitude: "-95.712891",
	}
	got, found, _ := detector.Detect("temblor en Estados Unidos?", []string{"es"})
	if !found {
		t.Fatalf("Expected to find '%s'. Got: nothing found", want.Name)
	}
	comparePlacesInfo(t, want, got)
}

func TestDetectWithDifferentSpaces(t *testing.T) {
	want := PlaceInfo{
		Name:      "united states",
		Code:      "us",
		Latitude:  "37.09024",
		Longitude: "-95.712891",
	}
	gotLeadingSpace, foundLeadingSpace, _ := detector.Detect("I live in usa", []string{"en"})
	gotTrailingSpace, foundTrailingSpace, _ := detector.Detect("usa is my country", []string{"en"})
	gotNoCountry, foundNoCountry, _ := detector.Detect("usability is important", []string{"en"})

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

func TestDetectUnicode(t *testing.T) {
	want := PlaceInfo{Name: "russia", Code: "ru", Latitude: "61.52401", Longitude: "105.318756"}
	got, found, _ := detector.Detect("в россия очень холодно", []string{"ru"})
	if !found {
		t.Fatalf("Expected to find '%s'. Got: nothing found", want.Name)
	}
	comparePlacesInfo(t, want, got)
}

func TestDetectNoCountry(t *testing.T) {
	got, found, _ := detector.Detect("the shaking felt very strong in my city", []string{"en"})
	if found {
		t.Fatalf("Expected to not find any Country. Got: %s", got.Name)
	}
}

func TestDetectChile(t *testing.T) {
	want := PlaceInfo{
		Name:         "chile",
		Code:         "cl",
		Latitude:     "-35.675147",
		Longitude:    "-71.542969",
		Municipality: locations.Place{Name: "santiago", Latitude: "-33.4489", Longitude: "-70.6693"},
	}
	got, found, _ := detectorChile.Detect("estudiando en santiago, chile", []string{"es"})
	if !found {
		t.Fatalf("Expected to find %s. Got: nothing found", want.Municipality.Name)
	}
	comparePlacesInfo(t, want, got)
}

func TestDetectChileCommuneOnly(t *testing.T) {
	want := PlaceInfo{
		Name:         "chile",
		Code:         "cl",
		Latitude:     "-35.675147",
		Longitude:    "-71.542969",
		Municipality: locations.Place{Name: "temuco", Latitude: "-38.7359", Longitude: "-72.5904"},
	}
	got, found, _ := detectorChile.Detect("TEMBLOR EN TEMUCO!!", []string{"es"})
	if !found {
		t.Fatalf("Expected to find %s. Got: nothing found", want.Municipality.Name)
	}
	comparePlacesInfo(t, want, got)
}

func TestDetectChileLongCommune(t *testing.T) {
	want := PlaceInfo{
		Name:         "chile",
		Code:         "cl",
		Latitude:     "-35.675147",
		Longitude:    "-71.542969",
		Municipality: locations.Place{Name: "san vicente de tagua tagua", Latitude: "-34.2812", Longitude: "-71.8571"},
	}
	got, found, _ := detectorChile.Detect("pasando la tarde en San Vicente de Tagua Tagua", []string{"es"})
	if !found {
		t.Fatalf("Expected to find %s. Got: nothing found", want.Municipality.Name)
	}
	comparePlacesInfo(t, want, got)
}

func TestDetectChileNoCommune(t *testing.T) {
	want := PlaceInfo{
		Name:      "chile",
		Code:      "cl",
		Latitude:  "-35.675147",
		Longitude: "-71.542969",
	}
	got, found, _ := detectorChile.Detect("TEMBLOR EN CHILE!!", []string{"es"})
	if !found {
		t.Fatalf("Expected to find %s. Got: nothing found", want.Name)
	}
	comparePlacesInfo(t, want, got)
}

func TestDetectInvalidLang(t *testing.T) {
	got, found, err := detectorChile.Detect("chile!!", []string{"invalidCode"})
	if found || err == nil {
		t.Fatalf("Expected an error because of the invalid code. Got a detection: %s", got.Name)
	}
}

func TestDetectInAnyLang(t *testing.T) {
	want := PlaceInfo{
		Name:      "united states",
		Code:      "us",
		Latitude:  "37.09024",
		Longitude: "-95.712891",
	}
	got, found := detector.DetectInAnyLang("temblor en Estados Unidos?")
	if !found {
		t.Fatalf("Expected to find '%s'. Got: nothing found", want.Name)
	}
	comparePlacesInfo(t, want, got)
}

func TestDetectInAnyLangWithDifferentSpaces(t *testing.T) {
	want := PlaceInfo{
		Name:      "united states",
		Code:      "us",
		Latitude:  "37.09024",
		Longitude: "-95.712891",
	}
	gotLeadingSpace, foundLeadingSpace := detector.DetectInAnyLang("I live in usa")
	gotTrailingSpace, foundTrailingSpace := detector.DetectInAnyLang("usa is my country")
	gotNoCountry, foundNoCountry := detector.DetectInAnyLang("usability is important")

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

func TestDetectInAnyLangUnicode(t *testing.T) {
	want := PlaceInfo{Name: "russia", Code: "ru", Latitude: "61.52401", Longitude: "105.318756"}
	got, found := detector.DetectInAnyLang("в россия очень холодно")
	if !found {
		t.Fatalf("Expected to find '%s'. Got: nothing found", want.Name)
	}
	comparePlacesInfo(t, want, got)
}

func TestDetectInAnyLangNoCountry(t *testing.T) {
	got, found := detector.DetectInAnyLang("the shaking felt very strong in my city")
	if found {
		t.Fatalf("Expected to not find any Country. Got: %s", got.Name)
	}
}

func TestDetectInAnyLangChile(t *testing.T) {
	want := PlaceInfo{
		Name:         "chile",
		Code:         "cl",
		Latitude:     "-35.675147",
		Longitude:    "-71.542969",
		Municipality: locations.Place{Name: "santiago", Latitude: "-33.4489", Longitude: "-70.6693"},
	}
	got, found := detectorChile.DetectInAnyLang("estudiando en santiago, chile")
	if !found {
		t.Fatalf("Expected to find %s. Got: nothing found", want.Municipality.Name)
	}
	comparePlacesInfo(t, want, got)
}

func TestDetectInAnyLangChileCommuneOnly(t *testing.T) {
	want := PlaceInfo{
		Name:         "chile",
		Code:         "cl",
		Latitude:     "-35.675147",
		Longitude:    "-71.542969",
		Municipality: locations.Place{Name: "temuco", Latitude: "-38.7359", Longitude: "-72.5904"},
	}
	got, found := detectorChile.DetectInAnyLang("TEMBLOR EN TEMUCO!!")
	if !found {
		t.Fatalf("Expected to find %s. Got: nothing found", want.Municipality.Name)
	}
	comparePlacesInfo(t, want, got)
}

func TestDetectInAnyLangChileLongCommune(t *testing.T) {
	want := PlaceInfo{
		Name:         "chile",
		Code:         "cl",
		Latitude:     "-35.675147",
		Longitude:    "-71.542969",
		Municipality: locations.Place{Name: "san vicente de tagua tagua", Latitude: "-34.2812", Longitude: "-71.8571"},
	}
	got, found := detectorChile.DetectInAnyLang("pasando la tarde en San Vicente de Tagua Tagua")
	if !found {
		t.Fatalf("Expected to find %s. Got: nothing found", want.Municipality.Name)
	}
	comparePlacesInfo(t, want, got)
}

func TestDetectInAnyLangChileNoCommune(t *testing.T) {
	want := PlaceInfo{
		Name:      "chile",
		Code:      "cl",
		Latitude:  "-35.675147",
		Longitude: "-71.542969",
	}
	got, found := detectorChile.DetectInAnyLang("TEMBLOR EN CHILE!!")
	if !found {
		t.Fatalf("Expected to find %s. Got: nothing found", want.Name)
	}
	comparePlacesInfo(t, want, got)
}
