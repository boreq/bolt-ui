package domain

type Altitude struct {
	altitude float64
}

func NewAltitude(altitude float64) Altitude {
	return Altitude{
		altitude: altitude,
	}
}

func (a Altitude) Float64() float64 {
	return a.altitude
}
