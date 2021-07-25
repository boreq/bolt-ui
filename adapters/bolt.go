package adapters

import (
	"os"
	"time"

	"github.com/boreq/errors"
	bolt "go.etcd.io/bbolt"
)

func NewBolt(path string) (*bolt.DB, error) {
	_, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, errors.Wrap(err, "database file does not exist")
		}

		return nil, errors.Wrap(err, "could not stat the database file")
	}

	options := &bolt.Options{
		Timeout: 5 * time.Second,
	}

	db, err := bolt.Open(path, 0600, options)
	if err != nil {
		if errors.Is(err, bolt.ErrTimeout) {
			return nil, errors.Wrap(err, "error opening the database (is another instance of the program running?)")
		}
		return nil, errors.Wrap(err, "error opening the database")
	}

	return db, nil
}
