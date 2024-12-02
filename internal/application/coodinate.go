package application

type Coordinate struct {
	City      string
	Latitude  string
	Longitude string
}

func NewCoordinate(city, lat, lng string) *Coordinate {
	return &Coordinate{
		City:      city,
		Latitude:  lat,
		Longitude: lng,
	}
}
