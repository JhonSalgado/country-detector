package detector

import (
	"github.com/JhonSalgado/country-detector/detector/locations"
)

type CountryDetector struct {
	translations   map[string]string
	countries      map[string]locations.Place
	municipalities map[string]locations.Place
	countryCode    string
}

// GetDetector returns a country detector
func GetDetector() CountryDetector {
	detector := CountryDetector{}
	detector.translations = locations.Translations
	detector.countries = locations.Countries
	return detector
}

// GetDetectorWithMunicipalities returns a detector that have information about municipalities of a certain country
func GetDetectorWithMunicipalities(countryCode string) CountryDetector {
	detector := GetDetector()
	if municipalities, ok := locations.Municipalities[countryCode]; ok {
		detector.municipalities = municipalities
	}
	detector.countryCode = countryCode
	return detector
}

// GetDetectorChile returns a country detector that have information about chilean municipalities
func GetDetectorChile() CountryDetector {
	detector := GetDetectorWithMunicipalities("cl")
	return detector
}
