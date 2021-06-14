package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/JhonSalgado/text-processor/processor"
)

type Codes map[string]int

var (
	originFolder      = "./static/"
	destinationFolder = "./detector/locations/"
)

func buildTranslations(filename string) map[string]Codes {
	// get processor
	filter := processor.Filter{
		OnlyCustom: false,
	}
	texProcessor, _ := processor.GetTextProcessorWithStopWordsFilter(filter)

	// create map
	translationsMap := make(map[string]Codes)

	// open origin file
	originFile, err := os.Open(originFolder + filename)
	if err != nil {
		log.Fatal(err)
	}
	fileScanner := bufio.NewScanner(originFile)
	defer originFile.Close()

	for fileScanner.Scan() {
		line := fileScanner.Text()
		transaltionInfo := strings.Split(line, "\t")
		keys := texProcessor.GetWordsSet(transaltionInfo[2])
		nKeys := len(keys)
		code := transaltionInfo[0]
		for _, key := range keys {
			codes, keyExists := translationsMap[key]
			if keyExists {
				n, codeExists := codes[code]
				if !codeExists || nKeys < n {
					codes[code] = nKeys
				}
			} else {
				translationsMap[key] = Codes{code: nKeys}
			}
		}
	}
	return translationsMap
}

func writeTranslations(filename string) {
	translationsMap := buildTranslations(filename)

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
	w.WriteString("var Translations map[string]Codes = map[string]Codes{\n")
	for key, codes := range translationsMap {
		w.WriteString(fmt.Sprintf("\t\"%s\": {", key))
		for code, n := range codes {
			w.WriteString(fmt.Sprintf("\"%s\": %d, ", code, n))
		}
		w.WriteString("},\n")
	}
	w.WriteString("}\n")
	w.Flush()
}

func writeHeavyTranslations(filename string) {

	// open origin file
	originFile, err := os.Open(originFolder + filename)
	if err != nil {
		log.Fatal(err)
	}
	fileScanner := bufio.NewScanner(originFile)
	defer originFile.Close()

	// create and open destination file
	destinationFileName := strings.Replace(filename, "txt", "go", 1)
	destinationFile, err := os.Create(destinationFolder + "heavy_" + destinationFileName)
	if err != nil {
		log.Fatal(err)
	}
	w := bufio.NewWriter(destinationFile)
	defer destinationFile.Close()

	// map to check if a translation is repeated
	isRepeated := make(map[string]bool)

	// start to write
	w.WriteString("package locations\n\n")
	w.WriteString("var HeavyTranslations map[string]string = map[string]string{\n")
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

func main() {
	writeTranslations("translations.txt")
	// writeHeavyTranslations("translations.txt")
	// writeCountries("countries.txt")
}
