package detector

import (
	"bufio"
	"os"
	"strings"
)

type location struct {
	lat string
	lon string
}

type place struct {
	location
	name string
}

type countryDetector struct {
	translations map[string]string
	countries    map[string]place
	communes     map[string]location
}

var translationsPath = "./static/translations.txt"
var countriesPath = "./static/countries.txt"
var communesPath = "./static/comunas.txt"

// GetDetector returns a country detector
func GetDetector() countryDetector {
	detector := countryDetector{}
	detector.translations = make(map[string]string)
	detector.loadTranslations()
	detector.countries = make(map[string]place)
	detector.loadCountries()
	return detector
}

// GetDetectorChile returns a country detector that have information about chilean communes
func GetDetectorChile() countryDetector {
	detector := GetDetector()
	detector.communes = make(map[string]location)
	detector.loadCommunes()
	return detector
}

func (detector countryDetector) loadTranslations() error {
	file, err := os.Open(translationsPath)
	if err != nil {
		return err
	}
	defer file.Close()
	fileScanner := bufio.NewScanner(file)
	for fileScanner.Scan() {
		line := fileScanner.Text()
		transaltionInfo := strings.Split(line, "\t")
		detector.translations[transaltionInfo[2]] = transaltionInfo[0]
	}
	return nil
}

func (detector countryDetector) loadCountries() error {
	file, err := os.Open(countriesPath)
	if err != nil {
		return err
	}
	defer file.Close()
	fileScanner := bufio.NewScanner(file)
	for fileScanner.Scan() {
		line := fileScanner.Text()
		countryInfo := strings.Split(line, "\t")
		placeInfo := place{
			name:     countryInfo[3],
			location: location{lon: countryInfo[2], lat: countryInfo[1]},
		}
		detector.countries[strings.ToLower(countryInfo[0])] = placeInfo
	}
	return nil
}

func (detector countryDetector) loadCommunes() error {
	file, err := os.Open(communesPath)
	if err != nil {
		return err
	}
	defer file.Close()
	fileScanner := bufio.NewScanner(file)
	for fileScanner.Scan() {
		line := fileScanner.Text()
		communeInfo := strings.Split(line, "\t")
		locationInfo := location{
			lat: communeInfo[1],
			lon: communeInfo[2],
		}
		detector.communes[communeInfo[0]] = locationInfo
	}
	return nil
}
