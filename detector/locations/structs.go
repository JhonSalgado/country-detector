package locations

// Place is a structure to store information of a place
type Place struct {
	Latitude  string
	Longitude string
	Name      string
}

var Municipalities map[string]map[string]Place = make(map[string]map[string]Place)
