package application

type Coordinate struct {
	Latitude  string
	Longitude string
}

func NewCoordinate(lat, lng string) *Coordinate {
	return &Coordinate{
		Latitude:  lat,
		Longitude: lng,
	}
}
