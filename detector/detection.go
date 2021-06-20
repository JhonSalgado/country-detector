package detector

import (
	"strings"

	"github.com/JhonSalgado/country-detector/detector/locations"
)

// PlaceInfo is the struct with the info of the country/municipality detected
type PlaceInfo struct {
	Name         string
	Code         string
	Longitude    string
	Latitude     string
	Municipality locations.Place
}

// getInfo creates a PlaceInfo struct with the country information
func (detector CountryDetector) getInfo(countryInfo locations.Place, code string) PlaceInfo {
	place := PlaceInfo{
		Name:      countryInfo.Name,
		Code:      code,
		Longitude: countryInfo.Longitude,
		Latitude:  countryInfo.Latitude,
	}
	return place
}

// DetectFromText detects the country or municipality mentioned in a text and returns its information
func (detector CountryDetector) DetectFromText(text string) (PlaceInfo, bool) {
	found := false
	place := PlaceInfo{}
	lowercaseText := strings.ToLower(text)

	// detect country
	for country, code := range detector.translations {
		if strings.Contains(lowercaseText, country) {
			countryInfo, ok := detector.countries[code]
			if ok {
				place = detector.getInfo(countryInfo, code)
				found = true
			}
			break
		}
	}

	// check if we can detect the municipality
	if len(detector.municipalities) > 0 && (place.Code == detector.countryCode || !found) {
		for municipality, municipalityInfo := range detector.municipalities {
			if strings.Contains(lowercaseText, municipality) {
				// if we had not found the country we set it now
				if !found {
					countryInfo := detector.countries[detector.countryCode]
					place = detector.getInfo(countryInfo, detector.countryCode)
					found = true
				}
				place.Municipality = municipalityInfo
				break
			}
		}
	}
	return place, found
}
