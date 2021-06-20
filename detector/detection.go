package detector

import (
	"strings"

	"github.com/JhonSalgado/country-detector/detector/locations"
)

type PlaceInfo struct {
	Name         string
	Code         string
	Longitude    string
	Latitude     string
	Municipality locations.Place
}

func (detector CountryDetector) getInfo(countryInfo locations.Place, code string) PlaceInfo {
	place := PlaceInfo{
		Name:      countryInfo.Name,
		Code:      code,
		Longitude: countryInfo.Longitude,
		Latitude:  countryInfo.Latitude,
	}
	return place
}

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
