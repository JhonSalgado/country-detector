package detector

import (
	"fmt"
	"regexp"

	"github.com/JhonSalgado/country-detector/detector/locations"
	"github.com/JhonSalgado/text-processor/processor"
)

var textProcessor = processor.GetTextProcessor()

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

// ContainsSentence checks if a word or sentence is present in a text
func (detector CountryDetector) ContainsSentence(text string, sentence string) bool {

	// we will check if a prefix of the substring is contained in the text before making the full comparison
	prefixSize := 4
	sentenceLength := len(sentence)

	// last index of the text to verify, beyond that index the sentence does not fit
	lastIndex := len(text) - sentenceLength

	// the sentence cannot be contained in a text shorter than it
	if lastIndex < 0 {
		return false
	}

	// if the sentence size is lower than the prefix size we use the whole sentence to search
	if prefixSize > sentenceLength {
		prefixSize = sentenceLength
	}

	prefix := sentence[0:prefixSize]
	// if the prefix is contained in the text we will compare the entire sentence
	if match, _ := regexp.MatchString(fmt.Sprintf(`(\s|^)%s`, prefix), text[0:lastIndex+prefixSize]); match {
		// this regular expression indicates that the sentence must be at the beginning of the text
		// or after a space, it cannot start or end in the middle of a word
		match, _ = regexp.MatchString(fmt.Sprintf(`(\s|^)%s(\s|$)`, sentence), text)
		return match
	} else {
		return false
	}
}

// detectInSingleLang detects the country or municipality mentioned in a text in just one lang and returns its information
func (detector CountryDetector) detectInSingleLang(text string, lang string) (PlaceInfo, bool) {

	found := false
	place := PlaceInfo{}

	// detect country
	for country, code := range detector.translations[lang] {
		if detector.ContainsSentence(text, country) {
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
			if detector.ContainsSentence(text, municipality) {
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

// Detect detects the country or municipality mentioned in a text in a list of possible langs and returns its information
func (detector CountryDetector) Detect(text string, langs []string) (PlaceInfo, bool, error) {

	found := false
	place := PlaceInfo{}
	cleanText := textProcessor.CleanText(text)
	for _, lang := range langs {
		if _, ok := detector.translations[lang]; ok {
			place, found = detector.detectInSingleLang(cleanText, lang)
			if found {
				return place, found, nil
			}
		} else {
			err := fmt.Errorf("lang code \"%s\" not supported", lang)
			return place, found, err
		}
	}
	return place, found, nil
}

// DetectInAnyLang detects the country or municipality mentioned in a text in any language and returns its information
func (detector CountryDetector) DetectInAnyLang(text string) (PlaceInfo, bool) {

	found := false
	place := PlaceInfo{}
	cleanText := textProcessor.CleanText(text)
	for lang := range detector.translations {
		place, found = detector.detectInSingleLang(cleanText, lang)
		if found {
			return place, found
		}
	}
	return place, found
}
