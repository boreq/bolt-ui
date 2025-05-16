package tests

import (
	"math/rand"
	"sort"
	"testing"

	"github.com/boreq/bolt-ui/application"
	"github.com/boreq/bolt-ui/internal/fixture"
	"github.com/boreq/bolt-ui/internal/wire"
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
		application.MustNewBrowse(nil, nil, nil, nil),
	)
	require.NoError(t, err)
	require.Equal(t, (firstPage), tree.Entries)

	// first page
	tree, err = testApp.Application.Browse.Execute(
		application.MustNewBrowse(nil, keyPointer(firstPage[0].Key), nil, nil),
	)
	require.NoError(t, err)
	require.Empty(t, tree.Entries)

	tree, err = testApp.Application.Browse.Execute(
		application.MustNewBrowse(nil, nil, keyPointer(firstPage[len(firstPage)-1].Key), nil),
	)
	require.NoError(t, err)
	require.Equal(t, (secondPage), tree.Entries)

	// second page
	tree, err = testApp.Application.Browse.Execute(
		application.MustNewBrowse(nil, keyPointer(secondPage[0].Key), nil, nil),
	)
	require.NoError(t, err)
	require.Equal(t, (firstPage), tree.Entries)

	tree, err = testApp.Application.Browse.Execute(
		application.MustNewBrowse(nil, nil, keyPointer(secondPage[len(secondPage)-1].Key), nil),
	)
	require.NoError(t, err)
	require.Equal(t, (thirdPage), tree.Entries)

	// third page
	tree, err = testApp.Application.Browse.Execute(
		application.MustNewBrowse(nil, keyPointer(thirdPage[0].Key), nil, nil),
	)
	require.NoError(t, err)
	require.Equal(t, (secondPage), tree.Entries)

	tree, err = testApp.Application.Browse.Execute(
		application.MustNewBrowse(nil, nil, keyPointer(thirdPage[len(thirdPage)-1].Key), nil),
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
		application.MustNewBrowse(path, nil, nil, nil),
	)
	require.NoError(t, err)
	require.Equal(t, (firstPage), tree.Entries)

	// first page
	tree, err = testApp.Application.Browse.Execute(
		application.MustNewBrowse(path, keyPointer(firstPage[0].Key), nil, nil),
	)
	require.NoError(t, err)
	require.Empty(t, tree.Entries)

	tree, err = testApp.Application.Browse.Execute(
		application.MustNewBrowse(path, nil, keyPointer(firstPage[len(firstPage)-1].Key), nil),
	)
	require.NoError(t, err)
	require.Equal(t, (secondPage), tree.Entries)

	tree, err = testApp.Application.Browse.Execute(
		application.MustNewBrowse(path, nil, nil, keyPointer(firstPage[0].Key)),
	)
	require.NoError(t, err)
	require.Equal(t, firstPage, tree.Entries)

	// second page
	tree, err = testApp.Application.Browse.Execute(
		application.MustNewBrowse(path, keyPointer(secondPage[0].Key), nil, nil),
	)
	require.NoError(t, err)
	require.Equal(t, (firstPage), tree.Entries)

	tree, err = testApp.Application.Browse.Execute(
		application.MustNewBrowse(path, nil, keyPointer(secondPage[len(secondPage)-1].Key), nil),
	)
	require.NoError(t, err)
	require.Equal(t, (thirdPage), tree.Entries)

	// third page
	tree, err = testApp.Application.Browse.Execute(
		application.MustNewBrowse(path, keyPointer(thirdPage[0].Key), nil, nil),
	)
	require.NoError(t, err)
	require.Equal(t, (secondPage), tree.Entries)

	tree, err = testApp.Application.Browse.Execute(
		application.MustNewBrowse(nil, nil, keyPointer(thirdPage[len(thirdPage)-1].Key), nil),
	)
	require.NoError(t, err)
	require.Empty(t, tree.Entries)
}

func TestBrowseNilValues(t *testing.T) {
	testApp := NewTracker(t)

	bucketName := "bucket"

	key1 := []byte("a")
	key2 := []byte("b")
	key3 := []byte("c")

	err := testApp.DB.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucket([]byte(bucketName))
		if err != nil {
			return err
		}

		_, err = bucket.CreateBucket(key1)
		if err != nil {
			return err
		}

		if err = bucket.Put(key2, key2); err != nil {
			return err
		}

		if err = bucket.Put(key3, nil); err != nil {
			return err
		}

		return nil
	})
	require.NoError(t, err)

	path := []application.Key{
		application.MustNewKey([]byte(bucketName)),
	}

	tree, err := testApp.Application.Browse.Execute(
		application.MustNewBrowse(path, nil, nil, nil),
	)
	require.NoError(t, err)
	require.Equal(t,
		[]application.Entry{
			{
				Bucket: true,
				Key:    application.MustNewKey(key1),
				Value:  nilValue,
			},
			{
				Bucket: false,
				Key:    application.MustNewKey(key2),
				Value:  application.MustNewValue(key2),
			},
			{
				Bucket: false,
				Key:    application.MustNewKey(key3),
				Value:  nilValue,
			},
		},
		tree.Entries)
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
			Bucket: true,
			Key:    application.MustNewKey([]byte(key)),
			Value:  nilValue,
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
				Bucket: true,
				Key:    application.MustNewKey([]byte(key)),
				Value:  nilValue,
			}
			result = append(result, entry)
		} else {
			entry := application.Entry{
				Bucket: false,
				Key:    application.MustNewKey([]byte(key)),
				Value:  application.MustNewValue([]byte(key + "_value")),
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

func keyPointer(v application.Key) *application.Key {
	return &v
}
