package tracker

import (
	"fmt"
	"io"

	"github.com/boreq/velo/application/tracker"
	"github.com/boreq/velo/domain"
)

type RouteFileParser struct {
	parserGpx *RouteFileParserGpx
	parserFit *RouteFileParserFit
}

func NewRouteFileParser(
	parserGpx *RouteFileParserGpx,
	parserFit *RouteFileParserFit,
) *RouteFileParser {
	return &RouteFileParser{
		parserGpx: parserGpx,
		parserFit: parserFit,
	}
}

func (r *RouteFileParser) Parse(
	f io.Reader,
	format tracker.RouteFileFormat,
) ([]domain.Point, error) {
	switch format {
	case tracker.RouteFileFormatGpx:
		return r.parserGpx.Parse(f)
	case tracker.RouteFileFormatFit:
		return r.parserFit.Parse(f)
	default:
		return nil, fmt.Errorf("unsupported route file format '%s'", format)
	}
}
