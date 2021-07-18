package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/JhonSalgado/text-processor/processor"
)

var textProcessor = processor.GetTextProcessor()

type Translations map[string]string

type Langs map[string]Translations

var (
	originFolder         = "./static/"
	municipalitiesFolder = "./static/municipalities/"
	destinationFolder    = "./detector/locations/"
)

// groupTranslationsByLang groups in a map the translations by language
func groupTranslationsByLang(filename string) Langs {
	// open origin file
	originFile, err := os.Open(originFolder + filename)
	if err != nil {
		log.Fatal(err)
	}
	fileScanner := bufio.NewScanner(originFile)
	defer originFile.Close()

	langs := Langs{}

	// read the file to create the map
	for fileScanner.Scan() {

		// read a line and make it lowercase
		line := strings.ToLower(fileScanner.Text())

		// split using tabs
		countryInfo := strings.Split(line, "\t")

		// remove special characters from the translation
		translation := textProcessor.CleanText(countryInfo[2])

		lang := countryInfo[3]
		// if the lang is empty or in a wrong format we move the translations to "other" section
		// TODO: fill or fix the lang codes in the source file
		if lang == "" || len(lang) > 2 {
			lang = "other"
		}
		code := countryInfo[0]
		// check if we already have a map for the language and add the translation and its country code
		if codes, exists := langs[lang]; exists {
			codes[translation] = code
		} else {
			langs[lang] = Translations{translation: code}
		}
	}
	return langs
}

// writeTranslations writes in a .go file a map where the keys are  ISO 639-1 lang codes and the values are maps with
// every country name translated to to the corresponding language
func writeTranslations(filename string) {

	// creates a langs map
	langs := groupTranslationsByLang(filename)

	// create and open destination file
	destinationFileName := strings.Replace(filename, "txt", "go", 1)
	destinationFile, err := os.Create(destinationFolder + destinationFileName)
	if err != nil {
		log.Fatal(err)
	}
	w := bufio.NewWriter(destinationFile)
	defer destinationFile.Close()

	// start to write
	w.WriteString("package locations\n\n")
	w.WriteString("var Translations = map[string]map[string]string{\n")
	for lang, transalations := range langs {
		w.WriteString(fmt.Sprintf("\t\"%s\": {\n", lang))
		for translation, country := range transalations {
			w.WriteString(fmt.Sprintf("\t\t\"%s\": \"%s\",\n", translation, country))
		}
		w.WriteString("\t},\n")
	}
	w.WriteString("}\n")
	w.Flush()
}

// writetranslations writes in a .go file a map where the keys are the ISO Alfa-2 codes of each
// country available in the countries file and the values are the information of the countries
func writeCountries(filename string) {

	// open origin file
	originFile, err := os.Open(originFolder + filename)
	if err != nil {
		log.Fatal(err)
	}
	fileScanner := bufio.NewScanner(originFile)
	defer originFile.Close()

	// create and open destination file
	destinationFileName := strings.Replace(filename, "txt", "go", 1)
	destinationFile, err := os.Create(destinationFolder + destinationFileName)
	if err != nil {
		log.Fatal(err)
	}
	w := bufio.NewWriter(destinationFile)
	defer destinationFile.Close()

	// start to write
	w.WriteString("package locations\n\n")
	w.WriteString("var Countries map[string]Place = map[string]Place{\n")
	for fileScanner.Scan() {
		line := strings.ToLower(fileScanner.Text())
		countryInfo := strings.Split(line, "\t")
		countryName := textProcessor.CleanText(countryInfo[3])
		w.WriteString(fmt.Sprintf("\t\"%s\": {", countryInfo[0]))         // ISO code
		w.WriteString(fmt.Sprintf("Latitude: \"%s\", ", countryInfo[1]))  // Latitude
		w.WriteString(fmt.Sprintf("Longitude: \"%s\", ", countryInfo[2])) // Longitude
		w.WriteString(fmt.Sprintf("Name: \"%s\"},\n", countryName))       // Name
	}
	w.WriteString("}\n")
	w.Flush()
}

// writeMunicipalitiesFromCountry writes in a .go file a map where the keys are the variants of the
// names of the municipalities availables for a certain country, and the values are the information
// of the municipalities.
func writeMunicipalitiesFromCountry(countryCode string) {

	// open origin file
	originFile, err := os.Open(fmt.Sprintf("%s%s.txt", municipalitiesFolder, countryCode))
	if err != nil {
		log.Fatal(err)
	}
	fileScanner := bufio.NewScanner(originFile)
	defer originFile.Close()

	// create and open destination file
	destinationFile, err := os.Create(fmt.Sprintf("%smunicipalities_%s.go", destinationFolder, countryCode))
	if err != nil {
		log.Fatal(err)
	}
	w := bufio.NewWriter(destinationFile)
	defer destinationFile.Close()

	// start to write
	w.WriteString("package locations\n\n")
	w.WriteString(fmt.Sprintf("var %s map[string]Place = map[string]Place{\n", countryCode))
	for fileScanner.Scan() {
		line := strings.ToLower(fileScanner.Text())
		municipalityInfo := strings.Split(line, "\t")
		municipalityShortName := textProcessor.CleanText(municipalityInfo[0])
		municipalityName := textProcessor.CleanText(municipalityInfo[3])
		w.WriteString(fmt.Sprintf("\t\"%s\": {", municipalityShortName))       // Short or alternative name
		w.WriteString(fmt.Sprintf("Latitude: \"%s\", ", municipalityInfo[1]))  // Latitude
		w.WriteString(fmt.Sprintf("Longitude: \"%s\", ", municipalityInfo[2])) // Longitude
		w.WriteString(fmt.Sprintf("Name: \"%s\"},\n", municipalityName))       // Name
	}
	w.WriteString("}\n")
	w.Flush()
}

// writeMunicipalities writes in .go files the maps for each municipality available in the
// municipalities folder and add the maps to a init.go file to be initialized properly
func writeMunicipalities() {
	// open origin folder
	files, err := ioutil.ReadDir(municipalitiesFolder)
	if err != nil {
		log.Fatal(err)
	}

	// create init file
	initFile, err := os.Create(destinationFolder + "municipalities_init.go")
	if err != nil {
		log.Fatal(err)
	}
	initWriter := bufio.NewWriter(initFile)
	defer initFile.Close()

	// write municipalities_init.go to initialize a map with available municipalities
	initWriter.WriteString("package locations\n\n")
	initWriter.WriteString("func init() {\n")

	for _, file := range files {
		// get the country code from the filename
		code := strings.Replace(file.Name(), ".txt", "", 1)

		// write municipalities map in its own file
		writeMunicipalitiesFromCountry(code)

		// add the municipalities map in municipalities_init.go
		initWriter.WriteString(fmt.Sprintf("\tMunicipalities[\"%s\"] = %s\n", code, code))
	}
	initWriter.WriteString("}\n")
	initWriter.Flush()
}

func main() {
	writeTranslations("translations.txt")
	writeCountries("countries.txt")
	writeMunicipalities()
}
