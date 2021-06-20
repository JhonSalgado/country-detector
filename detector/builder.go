package detector

import (
	"github.com/JhonSalgado/country-detector/detector/locations"
)

type CountryDetector struct {
	translations map[string]string
	countries    map[string]locations.Place
	communes     map[string]locations.Place
	countryCode  string
}

// GetDetector returns a country detector
func GetDetector() CountryDetector {
	detector := CountryDetector{}
	detector.translations = locations.Translations
	detector.countries = locations.Countries
	return detector
}

// GetDetectorWithCommunes returns a detector that have information about communes of a certain country
func GetDetectorWithCommunes(countryCode string) CountryDetector {
	detector := GetDetector()
	if communes, ok := locations.Communes[countryCode]; ok {
		detector.communes = communes
	}
	detector.countryCode = countryCode
	return detector
}

// GetDetectorChile returns a country detector that have information about chilean communes
func GetDetectorChile() CountryDetector {
	detector := GetDetectorWithCommunes("cl")
	return detector
}
