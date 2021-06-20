package detector

import (
	"strings"

	"github.com/JhonSalgado/country-detector/detector/locations"
)

type PlaceInfo struct {
	Name      string
	Code      string
	Longitude string
	Latitude  string
	Commune   locations.Place
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

	// check if we can detect the commune
	if len(detector.communes) > 0 && (place.Code == detector.countryCode || !found) {
		for commune, communeInfo := range detector.communes {
			if strings.Contains(lowercaseText, commune) {
				if !found {
					countryInfo := detector.countries[detector.countryCode]
					place = detector.getInfo(countryInfo, detector.countryCode)
					found = true
				}
				place.Commune = communeInfo
				break
			}
		}
	}
	return place, found
}
