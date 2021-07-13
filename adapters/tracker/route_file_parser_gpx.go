package tracker

import (
	"bytes"
	"io"

	"github.com/boreq/errors"
	"github.com/boreq/velo/domain"
	"github.com/tkrajina/gpxgo/gpx"
)

type RouteFileParserGpx struct {
}

func NewRouteFileParserGpx() *RouteFileParserGpx {
	return &RouteFileParserGpx{}
}

func (r *RouteFileParserGpx) Parse(f io.Reader) ([]domain.Point, error) {
	buf := &bytes.Buffer{}

	if _, err := io.Copy(buf, f); err != nil {
		return nil, errors.Wrap(err, "could not copy the bytes")
	}

	gpxFile, err := gpx.ParseBytes(buf.Bytes())
	if err != nil {
		return nil, errors.Wrap(err, "parsing failed")
	}

	var points []domain.Point

	for _, track := range gpxFile.Tracks {
		for _, segment := range track.Segments {
			for _, gpxPoint := range segment.Points {
				point, err := r.toPoint(gpxPoint)
				if err != nil {
					return nil, errors.Wrap(err, "could not create a point")
				}

				points = append(points, point)
			}
		}
	}

	return points, nil
}

func (r *RouteFileParserGpx) toPoint(gpxPoint gpx.GPXPoint) (domain.Point, error) {
	longitude, err := domain.NewLongitude(gpxPoint.Longitude)
	if err != nil {
		return domain.Point{}, errors.Wrap(err, "invalid longitude")
	}

	latitude, err := domain.NewLatitude(gpxPoint.Latitude)
	if err != nil {
		return domain.Point{}, errors.Wrap(err, "invalid latitude")
	}

	position := domain.NewPosition(latitude, longitude)

	altitude, err := domain.NewAltitude(gpxPoint.Elevation.Value())
	if err != nil {
		return domain.Point{}, errors.Wrap(err, "invalid altitude")
	}

	return domain.NewPoint(gpxPoint.Timestamp, position, altitude)
}
