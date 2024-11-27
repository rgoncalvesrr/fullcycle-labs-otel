package weather

import "math"

type Celsius float64

func (c Celsius) ToKelvin() float64 {
	return math.Trunc(float64((c+273.15)*100)) / 100
}

func (c Celsius) ToFahrenheit() float64 {
	return math.Trunc(float64((c*1.8+32)*100)) / 100
}
