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

	expectedEntries := bucketEntries(30)

	err := testApp.DB.Update(func(tx *bbolt.Tx) error {
		for _, entry := range expectedEntries {
			_, err := tx.CreateBucketIfNotExists(entry.Key.Bytes())
			if err != nil {
				return err
			}
		}
		return nil
	})
	require.NoError(t, err)

	firstPage := expectedEntries[0:10]
	secondPage := expectedEntries[10:20]
	thirdPage := expectedEntries[20:30]

	// initial
	tree, err := testApp.Application.Browse.Execute(
		application.Browse{
			Path: nil,
		},
	)
	require.NoError(t, err)
	require.Equal(t, (firstPage), tree.Entries)

	// first page
	tree, err = testApp.Application.Browse.Execute(
		application.Browse{
			Path:   nil,
			Before: keyPointer(firstPage[0].Key),
		},
	)
	require.NoError(t, err)
	require.Empty(t, tree.Entries)

	tree, err = testApp.Application.Browse.Execute(
		application.Browse{
			Path:  nil,
			After: keyPointer(firstPage[len(firstPage)-1].Key),
		},
	)
	require.NoError(t, err)
	require.Equal(t, (secondPage), tree.Entries)

	// second page
	tree, err = testApp.Application.Browse.Execute(
		application.Browse{
			Path:   nil,
			Before: keyPointer(secondPage[0].Key),
		},
	)
	require.NoError(t, err)
	require.Equal(t, (firstPage), tree.Entries)

	tree, err = testApp.Application.Browse.Execute(
		application.Browse{
			Path:  nil,
			After: keyPointer(secondPage[len(secondPage)-1].Key),
		},
	)
	require.NoError(t, err)
	require.Equal(t, (thirdPage), tree.Entries)

	// third page
	tree, err = testApp.Application.Browse.Execute(
		application.Browse{
			Path:   nil,
			Before: keyPointer(thirdPage[0].Key),
		},
	)
	require.NoError(t, err)
	require.Equal(t, (secondPage), tree.Entries)

	tree, err = testApp.Application.Browse.Execute(
		application.Browse{
			Path:  nil,
			After: keyPointer(thirdPage[len(thirdPage)-1].Key),
		},
	)
	require.NoError(t, err)
	require.Empty(t, tree.Entries)
}

func TestBrowse(t *testing.T) {
	testApp := NewTracker(t)

	bucketNameA := "bucket1"
	bucketNameB := "bucket2"

	expectedEntries := mixedBucketEntries(30)

	err := testApp.DB.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucket([]byte(bucketNameA))
		if err != nil {
			return err
		}

		bucket, err = bucket.CreateBucket([]byte(bucketNameB))
		if err != nil {
			return err
		}

		for _, entry := range expectedEntries {
			if entry.Value.IsEmpty() {
				_, err := bucket.CreateBucket(entry.Key.Bytes())
				if err != nil {
					return err
				}
			} else {
				if err = bucket.Put(entry.Key.Bytes(), entry.Value.Bytes()); err != nil {
					return err
				}
			}
		}
		return nil
	})
	require.NoError(t, err)

	firstPage := expectedEntries[0:10]
	secondPage := expectedEntries[10:20]
	thirdPage := expectedEntries[20:30]

	path := []application.Key{
		application.MustNewKey([]byte(bucketNameA)),
		application.MustNewKey([]byte(bucketNameB)),
	}

	// initial
	tree, err := testApp.Application.Browse.Execute(
		application.Browse{
			Path: path,
		},
	)
	require.NoError(t, err)
	require.Equal(t, (firstPage), tree.Entries)

	// first page
	tree, err = testApp.Application.Browse.Execute(
		application.Browse{
			Path:   path,
			Before: keyPointer(firstPage[0].Key),
		},
	)
	require.NoError(t, err)
	require.Empty(t, tree.Entries)

	tree, err = testApp.Application.Browse.Execute(
		application.Browse{
			Path:  path,
			After: keyPointer(firstPage[len(firstPage)-1].Key),
		},
	)
	require.NoError(t, err)
	require.Equal(t, (secondPage), tree.Entries)

	// second page
	tree, err = testApp.Application.Browse.Execute(
		application.Browse{
			Path:   path,
			Before: keyPointer(secondPage[0].Key),
		},
	)
	require.NoError(t, err)
	require.Equal(t, (firstPage), tree.Entries)

	tree, err = testApp.Application.Browse.Execute(
		application.Browse{
			Path:  path,
			After: keyPointer(secondPage[len(secondPage)-1].Key),
		},
	)
	require.NoError(t, err)
	require.Equal(t, (thirdPage), tree.Entries)

	// third page
	tree, err = testApp.Application.Browse.Execute(
		application.Browse{
			Path:   path,
			Before: keyPointer(thirdPage[0].Key),
		},
	)
	require.NoError(t, err)
	require.Equal(t, (secondPage), tree.Entries)

	tree, err = testApp.Application.Browse.Execute(
		application.Browse{
			Path:  nil,
			After: keyPointer(thirdPage[len(thirdPage)-1].Key),
		},
	)
	require.NoError(t, err)
	require.Empty(t, tree.Entries)
}

func keyPointer(v application.Key) *application.Key {
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

func bucketEntries(n int) (result []application.Entry) {
	var keys []string
	for i := 0; i < n; i++ {
		keys = append(keys, randString())
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	for _, key := range keys {
		entry := application.Entry{
			Key:   application.MustNewKey([]byte(key)),
			Value: nilValue,
		}
		result = append(result, entry)
	}

	return result
}

func mixedBucketEntries(n int) (result []application.Entry) {
	var keys []string
	for i := 0; i < n; i++ {
		keys = append(keys, randString())
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	for i, key := range keys {
		if i%2 == 0 {
			entry := application.Entry{
				Key:   application.MustNewKey([]byte(key)),
				Value: nilValue,
			}
			result = append(result, entry)
		} else {
			entry := application.Entry{
				Key:   application.MustNewKey([]byte(key)),
				Value: application.MustNewValue([]byte(key + "_value")),
			}
			result = append(result, entry)
		}
	}

	return result
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randString() string {
	b := make([]rune, 32)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
