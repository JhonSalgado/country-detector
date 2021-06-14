package detector

import (
	"github.com/JhonSalgado/country-detector/detector/locations"
)

type CountryDetector struct {
	translations      map[string]locations.Codes
	heavyTranslations map[string]string
	countries         map[string]locations.Place
	communes          map[string]locations.Place
}

// GetDetector returns a country detector
func GetDetector() CountryDetector {
	detector := CountryDetector{}
	detector.translations = locations.Translations
	detector.heavyTranslations = locations.HeavyTranslations
	detector.countries = locations.Countries
	return detector
}

// GetDetectorChile returns a country detector that have information about chilean communes
func GetDetectorChile() CountryDetector {
	detector := GetDetector()
	detector.communes = locations.Communes
	return detector
}
