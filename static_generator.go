package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

type Codes map[string]int

var (
	originFolder         = "./static/"
	municipalitiesFolder = "./static/municipalities/"
	destinationFolder    = "./detector/locations/"
)

// writetranslations writes in a .go file a map where the keys are every transalation of every
// country available in the translations file, and the values are the Alfa-2 codes of the countries
func writeTranslations(filename string) {

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

	// map to check if a translation is repeated
	isRepeated := make(map[string]bool)

	// start to write
	w.WriteString("package locations\n\n")
	w.WriteString("var Translations map[string]string = map[string]string{\n")
	for fileScanner.Scan() {
		line := strings.ToLower(fileScanner.Text())
		countryInfo := strings.Split(line, "\t")
		symbols := regexp.MustCompile(`[^\p{L}\s]`)
		countryName := symbols.ReplaceAllString(countryInfo[2], "")
		if repeated := isRepeated[countryName]; repeated {
			continue
		}
		w.WriteString(fmt.Sprintf("\t\"%s\": \"%s\",\n", countryName, countryInfo[0]))
		isRepeated[countryName] = true
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
		symbols := regexp.MustCompile(`[^\p{L}\s]`)
		countryName := symbols.ReplaceAllString(countryInfo[3], "")
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
		symbols := regexp.MustCompile(`[^\p{L}\s]`)
		municipalityShortName := symbols.ReplaceAllString(municipalityInfo[0], "")
		municipalityName := symbols.ReplaceAllString(municipalityInfo[3], "")
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
