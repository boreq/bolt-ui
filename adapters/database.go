package adapters

import (
	"github.com/boreq/errors"
	"github.com/boreq/velo/application"
	"go.etcd.io/bbolt"
)

const perPage = 10

type Database struct {
	tx *bbolt.Tx
}

func NewDatabase(tx *bbolt.Tx) *Database {
	return &Database{
		tx: tx,
	}
}

func (d *Database) Browse(path []application.Key, before *application.Key, after *application.Key) ([]application.Entry, error) {
	if len(path) == 0 {
		c := d.tx.Cursor()
		return d.iterate(c, before, after)
	}

	bucket, err := d.getBucket(path)
	if err != nil {
		return nil, errors.Wrap(err, "could not get the bucket")
	}

	c := bucket.Cursor()
	return d.iterate(c, before, after)
}

func (d *Database) iterate(c *bbolt.Cursor, before *application.Key, after *application.Key) ([]application.Entry, error) {
	if before != nil {
		return iterBefore(c, *before)
	}

	if after != nil {
		return iterAfter(c, *after)
	}

	return iter(c)
}

func (d *Database) getBucket(path []application.Key) (*bbolt.Bucket, error) {
	bucket := d.tx.Bucket(path[0].Bytes())
	if bucket == nil {
		return nil, application.ErrBucketNotFound
	}

	for i := 1; i < len(path); i++ {
		bucket = bucket.Bucket(path[i].Bytes())
		if bucket == nil {
			return nil, application.ErrBucketNotFound
		}
	}

	return bucket, nil
}

func iterBefore(c *bbolt.Cursor, before application.Key) ([]application.Entry, error) {
	var entries []application.Entry

	c.Seek(before.Bytes())

	for key, value := c.Prev(); key != nil; key, value = c.Prev() {
		entry, err := newEntry(key, value)
		if err != nil {
			return nil, errors.Wrap(err, "could not create an entry")
		}

		entries = append(entries, entry)

		if len(entries) >= perPage {
			break
		}
	}

	for i, j := 0, len(entries)-1; i < j; i, j = i+1, j-1 {
		entries[i], entries[j] = entries[j], entries[i]
	}

	return entries, nil
}

func iterAfter(c *bbolt.Cursor, after application.Key) ([]application.Entry, error) {
	var entries []application.Entry

	c.Seek(after.Bytes())

	for key, value := c.Next(); key != nil; key, value = c.Next() {
		entry, err := newEntry(key, value)
		if err != nil {
			return nil, errors.Wrap(err, "could not create an entry")
		}

		entries = append(entries, entry)

		if len(entries) >= perPage {
			break
		}
	}

	return entries, nil
}

func iter(c *bbolt.Cursor) ([]application.Entry, error) {
	var entries []application.Entry

	for key, value := c.First(); key != nil; key, value = c.Next() {
		entry, err := newEntry(key, value)
		if err != nil {
			return nil, errors.Wrap(err, "could not create an entry")
		}

		entries = append(entries, entry)

		if len(entries) >= perPage {
			break
		}
	}

	return entries, nil
}

func newEntry(k, v []byte) (application.Entry, error) {
	key, err := application.NewKey(k)
	if err != nil {
		return application.Entry{}, errors.Wrap(err, "could not create a key")
	}

	value, err := application.NewValue(v)
	if err != nil {
		return application.Entry{}, errors.Wrap(err, "could not create a value")
	}

	return application.Entry{
		Key:   key,
		Value: value,
	}, nil
}
