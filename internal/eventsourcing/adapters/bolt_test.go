package adapters_test

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/boreq/eggplant/internal/eventsourcing"
	"github.com/boreq/eggplant/internal/eventsourcing/adapters"
	"github.com/stretchr/testify/require"
	bolt "go.etcd.io/bbolt"
)

type CleanupFunc func()

func RunTestBolt(t *testing.T, test Test) {
	db, cleanup := FixtureBolt(t)
	defer cleanup()

	err := db.Update(func(tx *bolt.Tx) error {
		adapter := adapters.NewBoltPersistenceAdapter(tx, func(uuid eventsourcing.AggregateUUID) []adapters.BucketName {
			return []adapters.BucketName{
				[]byte("events"),
				[]byte(uuid),
			}
		})

		test(t, adapter)

		return nil
	})
	require.NoError(t, err)
}

func FixtureFile(t *testing.T) (string, CleanupFunc) {
	file, err := ioutil.TempFile("", "eventsourcing_test")
	if err != nil {
		t.Fatal(err)
	}

	cleanup := func() {
		err := os.Remove(file.Name())
		if err != nil {
			t.Fatal(err)
		}
	}

	return file.Name(), cleanup
}

func FixtureBolt(t *testing.T) (*bolt.DB, CleanupFunc) {
	file, fileCleanup := FixtureFile(t)

	db, err := bolt.Open(file, 0600, &bolt.Options{Timeout: 5 * time.Second})
	if err != nil {
		t.Fatal(err)
	}

	cleanup := func() {
		defer fileCleanup()

		err := db.Close()
		if err != nil {
			t.Fatal(err)
		}
	}

	return db, cleanup
}
