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
	originFolder      = "./static/"
	communesFolder    = "./static/communes/"
	destinationFolder = "./detector/locations/"
)

// writetranslations writes in a .go file a map where the keys are every transalation of every
// country available in the translations file, and the values are the ISO codes of the countries
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

// writetranslations writes in a .go file a map where the keys are the iso codes of each country
// available in the countries file and the values are the information of the countries
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

// writeCommunesFromCountry writes in a .go file a map where the keys are the names of the communes
// availables for a certain country and the values are the information of the communes
func writeCommunesFromCountry(countryCode string) {

	// open origin file
	originFile, err := os.Open(fmt.Sprintf("%s%s.txt", communesFolder, countryCode))
	if err != nil {
		log.Fatal(err)
	}
	fileScanner := bufio.NewScanner(originFile)
	defer originFile.Close()

	// create and open destination file
	destinationFile, err := os.Create(fmt.Sprintf("%scommunes_%s.go", destinationFolder, countryCode))
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
		communeInfo := strings.Split(line, "\t")
		symbols := regexp.MustCompile(`[^\p{L}\s]`)
		communeShortName := symbols.ReplaceAllString(communeInfo[0], "")
		communeName := symbols.ReplaceAllString(communeInfo[3], "")
		w.WriteString(fmt.Sprintf("\t\"%s\": {", communeShortName))       // Short or alternative name
		w.WriteString(fmt.Sprintf("Latitude: \"%s\", ", communeInfo[1]))  // Latitude
		w.WriteString(fmt.Sprintf("Longitude: \"%s\", ", communeInfo[2])) // Longitude
		w.WriteString(fmt.Sprintf("Name: \"%s\"},\n", communeName))       // Name
	}
	w.WriteString("}\n")
	w.Flush()
}

// writeCommunes writes in .go files the maps for each commune available in the communes folder
// and add the maps to a init.go file to be initialized properly
func writeCommunes() {
	// open origin folder
	files, err := ioutil.ReadDir(communesFolder)
	if err != nil {
		log.Fatal(err)
	}

	// create init file
	initFile, err := os.Create(destinationFolder + "communes_init.go")
	if err != nil {
		log.Fatal(err)
	}
	initWriter := bufio.NewWriter(initFile)
	defer initFile.Close()

	// write communes_init.go to initialize a map with available communes
	initWriter.WriteString("package locations\n\n")
	initWriter.WriteString("func init() {\n")

	for _, file := range files {
		// get the country code from the filename
		code := strings.Replace(file.Name(), ".txt", "", 1)

		// write communes map in its own file
		writeCommunesFromCountry(code)

		// add the communes map in communes_init.go
		initWriter.WriteString(fmt.Sprintf("\tCommunes[\"%s\"] = %s\n", code, code))
	}
	initWriter.WriteString("}\n")
	initWriter.Flush()
}

func main() {
	writeTranslations("translations.txt")
	writeCountries("countries.txt")
	writeCommunes()
}
