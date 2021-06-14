package detector

import (
	"strings"

	"github.com/JhonSalgado/country-detector/detector/locations"
	"github.com/JhonSalgado/text-processor/processor"
)

type PlaceInfo struct {
	Country   string
	Code      string
	Longitude string
	Latitude  string
	Commune   locations.Place
}

type codeCompleteness struct {
	got    int
	needed int
}

func (detector CountryDetector) getInfo(countryInfo locations.Place, code string) PlaceInfo {
	place := PlaceInfo{
		Country:   countryInfo.Name,
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
	for country, code := range detector.heavyTranslations {
		if strings.Contains(lowercaseText, country) {
			countryInfo, ok := detector.countries[code]
			if ok {
				place = detector.getInfo(countryInfo, code)
				found = true
			}
			break
		}
	}
	return place, found
}

func (detector CountryDetector) DetectFromTextLight(text string) (PlaceInfo, bool) {
	textProcessor := processor.GetTextProcessor()
	wordList := textProcessor.GetWordsSet(text)
	return detector.DetectFromList(wordList)
}

func (detector CountryDetector) checkCompleteness(codeCompleteness map[string]codeCompleteness) (string, bool) {
	found := false
	maxCompletition := 0
	var maxCompletitionCode string
	for code, completeness := range codeCompleteness {
		if completeness.got >= completeness.needed && completeness.needed > maxCompletition {
			found = true
			maxCompletition = completeness.needed
			maxCompletitionCode = code
		}
	}
	return maxCompletitionCode, found
}

func (detector CountryDetector) DetectFromMap(tokens map[string]int) (PlaceInfo, bool) {
	tokenList := make([]string, 0, len(tokens))
	for token := range tokens {
		tokenList = append(tokenList, token)
	}
	return detector.DetectFromList(tokenList)
}

func (detector CountryDetector) DetectFromList(tokens []string) (PlaceInfo, bool) {
	place := PlaceInfo{}
	codesGot := make(map[string]codeCompleteness)
	for _, token := range tokens {
		token = strings.ToLower(token)
		if codes, ok := detector.translations[token]; ok {
			for code, expected := range codes {
				if completition, ok := codesGot[code]; ok {
					completition.got += 1
					if completition.needed > expected {
						completition.needed = expected
					}
					codesGot[code] = completition
				} else {
					codesGot[code] = codeCompleteness{
						got:    1,
						needed: expected,
					}
				}
			}
		}
	}
	code, found := detector.checkCompleteness(codesGot)
	if !found {
		return place, found
	}
	countryInfo, codeFound := detector.countries[code]
	if !codeFound {
		return place, codeFound
	}
	place = detector.getInfo(countryInfo, code)
	return place, codeFound
}
