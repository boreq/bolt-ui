package domain

type Altitude struct {
	altitude float64
}

func NewAltitude(altitude float64) Altitude {
	return Altitude{
		altitude: altitude,
	}
}
