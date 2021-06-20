package locations

type Place struct {
	Latitude  string
	Longitude string
	Name      string
}

var Communes map[string]map[string]Place = make(map[string]map[string]Place)
