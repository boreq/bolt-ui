package tracker

import (
	"archive/zip"
	"compress/gzip"
	"encoding/csv"
	"fmt"
	"io"
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

		route, err := s.loadRoute(r, record)
		if err != nil {
			return errors.Wrap(err, "could not load the route file")
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

func (s *StravaExportFileParser) loadRoute(r *zip.Reader, record []string) ([]domain.Point, error) {
	filename := record[stravaActivitiesCSVFieldActivityFilename]

	var activityFile io.ReadCloser
	var err error

	activityFile, err = r.Open(filename)
	if err != nil {
		return nil, errors.Wrap(err, "could not open the activity file")
	}

	defer activityFile.Close()

	if strings.HasSuffix(filename, ".gz") {
		gzipReader, err := gzip.NewReader(activityFile)
		if err != nil {
			return nil, errors.Wrap(err, "could not open a gzip reader")
		}

		defer gzipReader.Close()

		activityFile = gzipReader
	}

	return s.routeFileParser.Parse(activityFile)
}

func (s *StravaExportFileParser) getTitle(record []string) (domain.ActivityTitle, error) {
	title := record[stravaActivitiesCSVFieldActivityName]

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
