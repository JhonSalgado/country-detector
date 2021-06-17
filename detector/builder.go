package detector

import (
	"github.com/JhonSalgado/country-detector/detector/locations"
)

type CountryDetector struct {
	hashedTranslations map[string]locations.Codes
	translations       map[string]string
	countries          map[string]locations.Place
	communes           map[string]locations.Place
	countryCode        string
}

// GetDetector returns a country detector
func GetDetector() CountryDetector {
	detector := CountryDetector{}
	detector.hashedTranslations = locations.HashedTranslations
	detector.translations = locations.Translations
	detector.countries = locations.Countries
	return detector
}

// GetDetectorChile returns a country detector that have information about chilean communes
func GetDetectorChile() CountryDetector {
	detector := GetDetector()
	detector.communes = locations.CommunesChile
	detector.countryCode = "cl"
	return detector
}
