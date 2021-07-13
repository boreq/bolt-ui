package tracker

import (
	"io"
	"math"

	"github.com/boreq/errors"
	"github.com/boreq/velo/domain"
	"github.com/tormoder/fit"
)

type RouteFileParserFit struct {
}

func NewRouteFileParserFit() *RouteFileParserFit {
	return &RouteFileParserFit{}
}

func (p *RouteFileParserFit) Parse(r io.Reader) ([]domain.Point, error) {
	f, err := fit.Decode(r)
	if err != nil {
		return nil, errors.Wrap(err, "parsing failed")
	}

	activity, err := f.Activity()
	if err != nil {
		return nil, errors.Wrap(err, "could not extract the activity")
	}

	var points []domain.Point

	for _, record := range activity.Records {
		if p.skipPoint(record) {
			continue
		}

		point, err := p.toPoint(record)
		if err != nil {
			return nil, errors.Wrap(err, "could not create a point")
		}

		points = append(points, point)
	}

	return points, nil
}

func (p *RouteFileParserFit) skipPoint(record *fit.RecordMsg) bool {
	if math.IsNaN(record.PositionLong.Degrees()) {
		return true
	}

	if math.IsNaN(record.PositionLat.Degrees()) {
		return true
	}

	return false
}

func (p *RouteFileParserFit) toPoint(record *fit.RecordMsg) (domain.Point, error) {
	longitude, err := domain.NewLongitude(record.PositionLong.Degrees())
	if err != nil {
		return domain.Point{}, errors.Wrap(err, "invalid longitude")
	}

	latitude, err := domain.NewLatitude(record.PositionLat.Degrees())
	if err != nil {
		return domain.Point{}, errors.Wrap(err, "invalid latitude")
	}

	position := domain.NewPosition(latitude, longitude)

	altitude, err := domain.NewAltitude(float64(record.Altitude))
	if err != nil {
		return domain.Point{}, errors.Wrap(err, "invalid altitude")
	}

	return domain.NewPoint(record.Timestamp, position, altitude)
}
