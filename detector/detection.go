package detector

type placeInfo struct {
	country string
	code    string
	location
	commune place
}

func (detector countryDetector) DetectFromText(text string) (placeInfo, bool) {
	place := placeInfo{}
	return place, false
}

func (detector countryDetector) DetectFromList(tokens []string) (placeInfo, bool) {
	place := placeInfo{}
	return place, false
}

func (detector countryDetector) DetectFromMap(tokens map[string]interface{}) (placeInfo, bool) {
	place := placeInfo{}
	return place, false
}
