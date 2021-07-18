# country-detector
This golang package is used to detect countries mentioned in a text, in many languages. If it detects a country, it returns the name in English, its ISO Alpha-2 code and its geographical coordinates.

It also allows detecting the municipalities of some countries, returning their full name, coordinates and the information of their country, but at the moment this feature is only available for Chile.

## Install
With Go installed:
`go get github.com/JhonSalgado/country-detector`

## Usage
In order to use the detection method you need to create a country detector. For this you have 3 builders:

- GetDetector: Returns a detector that only has country information.
- GetDetectorWithMunicipalities: It receives as a parameter an Alpha-2 country code and returns a detector that has information of countries and also information about the municipalities of the country to which the code belongs. This builder is designed for when this project has information about the municipalities of more countries. 
- GetDetectorChile: Just a short way to call GetDetectorWithMunicipalities with the code "cl".

The detector has two main methods called Detect and DetectInAnyLang, they both receive a text message and return a bool to indicate if a country was detected along with the information for that country. If the detector has municipalities loaded and the detected country corresponds to the country they belong to, or if no country was found, it will also try to detect the municipality. If found, it will be delivered along with country information, as shown in the example below.

The main difference between both methods, as indicated by their names, is that DetectInAnyLang will check if the text contains the name of a country translated into up to 140 different languages, while Detect receives a list of languages ​​(their ISO 639-1 codes), and it will only use those translations to detect. In case any provided language code is invalid or not supported, the method will return an error as the third result. I highly recommend using Detect if you have knowledge of the language of the text to be parsed, as it can be dozens of times faster.

Also available in the package is a helper method called ContainsSentence that indicates whether a text contains a sentence (this method is used by the main methods mentioned above).

### Example
```
package main

import (
	"encoding/json"
	"log"

	// importing the package
	"github.com/JhonSalgado/country-detector/detector"
)

func main() {
	
	// creating our detectors
	normalDetector := detector.GetDetector()
	chileDetector := detector.GetDetectorChile()

	// This phrase is in Vietnamese and it means "Portugal is beautiful"
	text1 := "Bồ Đào Nha xinh đẹp"

	// This sentence mentions Santiago, which is a municipality of Chile.
	text2 := "I'm going to spend Christmas in Santiago, Chile"

	// analyzing the first text with the normal detector
	detectedPlace1, found := normalDetector.DetectInAnyLang(text1)
	if found {
		// print the result as json just for readability
		bytes, _ := json.MarshalIndent(detectedPlace1, "", "    ")
		log.Println("Detection from text1:\n", string(bytes))
	}

	// analyzing the secong text with the chilean detector
	detectedPlace2, found2, err := chileDetector.Detect(text2, []string{"en"})
	if found2 && err == nil {
		bytes, _ := json.MarshalIndent(detectedPlace2, "", "    ")
		log.Println("Detection from text2:\n", string(bytes))
	}
}
```
### Output
```
Detection from text1:
{
    "Name": "portugal",
    "Code": "pt",
    "Longitude": "-8.224454",
    "Latitude": "39.399872",
    "Municipality": {
        "Latitude": "",
        "Longitude": "",
        "Name": ""
    }
}
Detection from text2:
{
    "Name": "chile",
    "Code": "cl",
    "Longitude": "-71.542969",
    "Latitude": "-35.675147",
    "Municipality": {
        "Latitude": "-33.4489",
        "Longitude": "-70.6693",
        "Name": "santiago"
    }
}
```
## To contribute:

If you want to contribute by adding countries or their translations you just have to edit countries.txt or translations.txt respectively in the static folder, respecting their format (see files format subsection), and execute the last step of this section.

If you want to add municipalities for a country, you must name the file as code.txt, where code is the ISO Alpha-2 code of the country, and add the file to static/municipalities/ folder.

To make the changes effective, you must execute the static\_generator.go file, with the following command: `go run static_generator.go`. This will make the static files available to the package, regardless of where it is being used, through Go files, which are stored in detector/locations/ (this folder should not be edited manually), and then create a pull request to develop with the generated changes.

### Files format
There are 3 types of files:
- translations.txt: It has 5 columns separated by tabs, which correspond respectively to: alpha-2 code of the country, name of the country in english, translated name, ISO 639-1 code of the language to which it is translated, the language name in english. Currently only the first and third are being used but the other columns could be used in the future.
- countries.txt: It has 4 columns separated by tab, which correspond respectively to: alpha-2 country code, latitude, longitude, country name in english.
- municipality/: Each file on that folder must have 4 columns separated by tab, which are respectively: detection name, latitude, longitude, full name. The first column is the one used to detect so I recommend adding a row for each alternative name that the municipality has, including its full name.

## License and Copyright
Copyright (c) 2021 Jhon Salgado, contributors. Released under the MIT license.
