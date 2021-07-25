package tests

import (
	"math/rand"
	"sort"
	"testing"

	"github.com/boreq/velo/application"
	"github.com/boreq/velo/internal/fixture"
	"github.com/boreq/velo/internal/wire"
	"github.com/stretchr/testify/require"
	"go.etcd.io/bbolt"
)

var nilValue = application.MustNewValue(nil)

func TestBrowseRoot(t *testing.T) {
	testApp := NewTracker(t)

	var keys []string
	for i := 0; i < 30; i++ {
		keys = append(keys, randString())
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	err := testApp.DB.Update(func(tx *bbolt.Tx) error {
		for _, key := range keys {
			_, err := tx.CreateBucketIfNotExists([]byte(key))
			if err != nil {
				return err
			}
		}
		return nil
	})
	require.NoError(t, err)

	firstPage := keys[0:10]
	secondPage := keys[10:20]
	thirdPage := keys[20:30]

	// initial
	entries, err := testApp.Application.Browse.Execute(
		application.Browse{
			Path: nil,
		},
	)
	require.NoError(t, err)
	require.Equal(t, bucketEntries(firstPage), entries)

	// first page
	entries, err = testApp.Application.Browse.Execute(
		application.Browse{
			Path:   nil,
			Before: keyPointer(firstPage[0]),
		},
	)
	require.NoError(t, err)
	require.Empty(t, entries)

	entries, err = testApp.Application.Browse.Execute(
		application.Browse{
			Path:  nil,
			After: keyPointer(firstPage[len(firstPage)-1]),
		},
	)
	require.NoError(t, err)
	require.Equal(t, bucketEntries(secondPage), entries)

	// second page
	entries, err = testApp.Application.Browse.Execute(
		application.Browse{
			Path:   nil,
			Before: keyPointer(secondPage[0]),
		},
	)
	require.NoError(t, err)
	require.Equal(t, bucketEntries(firstPage), entries)

	entries, err = testApp.Application.Browse.Execute(
		application.Browse{
			Path:  nil,
			After: keyPointer(secondPage[len(secondPage)-1]),
		},
	)
	require.NoError(t, err)
	require.Equal(t, bucketEntries(thirdPage), entries)

	// third page
	entries, err = testApp.Application.Browse.Execute(
		application.Browse{
			Path:   nil,
			Before: keyPointer(thirdPage[0]),
		},
	)
	require.NoError(t, err)
	require.Equal(t, bucketEntries(secondPage), entries)

	entries, err = testApp.Application.Browse.Execute(
		application.Browse{
			Path:  nil,
			After: keyPointer(thirdPage[len(thirdPage)-1]),
		},
	)
	require.NoError(t, err)
	require.Empty(t, entries)
}

func keyPointer(s string) *application.Key {
	v := application.MustNewKey([]byte(s))
	return &v
}

func NewTracker(t *testing.T) wire.TestApplication {
	db, cleanup := fixture.Bolt(t)
	t.Cleanup(cleanup)

	application, err := wire.BuildApplicationForTest(db)
	if err != nil {
		t.Fatal(err)
	}

	return application
}

func bucketEntries(keys []string) (result []application.Entry) {
	for _, key := range keys {
		entry := application.Entry{
			Key:   application.MustNewKey([]byte(key)),
			Value: nilValue,
		}
		result = append(result, entry)
	}
	return result
}

//var (
//	a = []byte("bucket")
//	b = mustDecodeString("c328")
//)
//
//func populate(t *testing.T, db *bbolt.DB) {
//	err := db.Update(func(tx *bbolt.Tx) error {
//		_, err := tx.CreateBucketIfNotExists(a)
//		if err != nil {
//			return err
//		}
//
//		_, err = tx.CreateBucketIfNotExists(b)
//		if err != nil {
//			return err
//		}
//
//		return nil
//	})
//	require.NoError(t, err)
//}
//
//func mustDecodeString(h string) []byte {
//	b, err := hex.DecodeString(h)
//	if err != nil {
//		panic(err)
//	}
//	return b
//}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randString() string {
	b := make([]rune, 32)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
