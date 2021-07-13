package tracker

import (
	"archive/zip"
	"compress/gzip"
	"encoding/csv"
	"fmt"
	"io"
	"path"
	"strings"

	"github.com/boreq/errors"
	"github.com/boreq/velo/application/tracker"
	"github.com/boreq/velo/domain"
)

const (
	stravaActivitiesCSVFieldActivityName     = 2
	stravaActivitiesCSVFieldActivityFilename = 10
)

type StravaExportFileParser struct {
	routeFileParser *RouteFileParser
}

func NewStravaExportFileParser(routeFileParser *RouteFileParser) *StravaExportFileParser {
	return &StravaExportFileParser{
		routeFileParser: routeFileParser,
	}
}

func (s *StravaExportFileParser) Parse(ra io.ReaderAt, size int64) (<-chan tracker.StravaActivity, error) {
	r, err := zip.NewReader(ra, size)
	if err != nil {
		return nil, errors.Wrap(err, "opening a zip reader failed")
	}

	ch := make(chan tracker.StravaActivity)

	go func() {
		defer close(ch)

		err = s.loadActivities(r, ch)
		if err != nil {
			ch <- tracker.StravaActivity{
				Err: err,
			}
		}
	}()

	return ch, nil
}

func (s *StravaExportFileParser) loadActivities(r *zip.Reader, ch chan tracker.StravaActivity) error {
	f, err := r.Open("activities.csv")
	if err != nil {
		return errors.Wrap(err, "could not open the activities file")
	}

	defer f.Close()

	fr := csv.NewReader(f)

	// validate that we are loading appropriate fields
	header, err := fr.Read()
	if err != nil {
		return errors.Wrap(err, "could not read the activities file header")
	}

	err = s.validateHeader(header)
	if err != nil {
		return errors.Wrap(err, "invalid activities file")
	}

	for {
		record, err := fr.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return errors.Wrap(err, "failed to read the csv file")
		}

		filename := record[stravaActivitiesCSVFieldActivityFilename]
		route, err := s.loadRoute(r, filename)
		if err != nil {
			return errors.Wrapf(err, "could not load the route file '%s'", filename)
		}

		title, err := s.getTitle(record)
		if err != nil {
			return errors.Wrap(err, "could not get the activity title")
		}

		ch <- tracker.StravaActivity{
			Title: title,
			Route: route,
		}
	}

	return nil
}

const gzSuffix = ".gz"

func (s *StravaExportFileParser) loadRoute(r *zip.Reader, filename string) ([]domain.Point, error) {

	var activityFile io.ReadCloser
	var err error

	activityFile, err = r.Open(filename)
	if err != nil {
		return nil, errors.Wrap(err, "could not open the activity file")
	}

	defer activityFile.Close()

	if ext := path.Ext(filename); strings.EqualFold(ext, gzSuffix) {
		gzipReader, err := gzip.NewReader(activityFile)
		if err != nil {
			return nil, errors.Wrap(err, "could not open a gzip reader")
		}

		defer gzipReader.Close()

		activityFile = gzipReader

		filename = filename[:len(filename)-len(gzSuffix)]
	}

	format, err := tracker.NewRouteFileFormatFromExtension(path.Ext(filename))
	if err != nil {
		return nil, errors.New("could not determine route file format")
	}

	return s.routeFileParser.Parse(activityFile, format)
}

func (s *StravaExportFileParser) getTitle(record []string) (domain.ActivityTitle, error) {
	title := record[stravaActivitiesCSVFieldActivityName]

	if len(title) > 50 {
		title = title[0:47] + "..."
	}

	return domain.NewActivityTitle(title)
}

func (s *StravaExportFileParser) validateHeader(header []string) error {
	if l := len(header); l < 78 {
		return fmt.Errorf("invalid row length '%d'", l)
	}

	if f := header[stravaActivitiesCSVFieldActivityName]; f != "Activity Name" {
		return fmt.Errorf("expected activity name but got '%s'", f)
	}

	if f := header[stravaActivitiesCSVFieldActivityFilename]; f != "Filename" {
		return fmt.Errorf("expected filename but got '%s'", f)
	}

	return nil
}
