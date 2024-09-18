package http

import (
	"encoding/hex"
	"github.com/boreq/bolt-ui/display"
	"github.com/boreq/errors"
	"unicode"

	"github.com/boreq/bolt-ui/application"
)

type Tree struct {
	Path    []Key   `json:"path"`
	Entries []Entry `json:"entries"`
}

type Entry struct {
	Bucket bool   `json:"bucket"`
	Key    Key    `json:"key"`
	Value  *Value `json:"value,omitempty"`
}

type Key struct {
	Hex string `json:"hex"`
	Str string `json:"str,omitempty"`
}

type Value struct {
	Hex    string  `json:"hex"`
	Pretty *Pretty `json:"pretty"`
}

type Pretty struct {
	ContentType string `json:"content_type"`
	Value       string `json:"value"`
}

func toTree(tree application.Tree) (Tree, error) {
	entries, err := toEntries(tree.Entries)
	if err != nil {
		return Tree{}, errors.Wrap(err, "error converting to entries")
	}
	return Tree{
		Path:    toKeys(tree.Path),
		Entries: entries,
	}, nil
}

func toKeys(keys []application.Key) []Key {
	result := make([]Key, 0)
	for _, key := range keys {
		result = append(result, toKey(key))
	}
	return result
}

func toEntries(entries []application.Entry) ([]Entry, error) {
	result := make([]Entry, 0)
	for _, entry := range entries {
		v, err := toEntry(entry)
		if err != nil {
			return nil, errors.Wrap(err, "error converting to an entry")
		}
		result = append(result, v)
	}
	return result, nil
}

func toEntry(entry application.Entry) (Entry, error) {
	value, err := toValue(entry.Value)
	if err != nil {
		return Entry{}, errors.Wrap(err, "error converting to a value")
	}

	return Entry{
		Bucket: entry.Bucket,
		Key:    toKey(entry.Key),
		Value:  value,
	}, nil
}

func toKey(key application.Key) Key {
	b := key.Bytes()

	result := Key{
		Hex: hex.EncodeToString(b),
	}

	if canDisplayKeyAsString(b) {
		result.Str = string(b)
	}

	return result
}

func toValue(value application.Value) (*Value, error) {
	if value.IsEmpty() {
		return nil, nil
	}

	b := value.Bytes()
	hexB := hex.EncodeToString(b)
	pretty, err := toPretty(value)
	if err != nil {
		return nil, errors.Wrap(err, "error converting to a pretty value")
	}

	return &Value{
		Hex:    hexB,
		Pretty: pretty,
	}, nil
}

func toPretty(value application.Value) (*Pretty, error) {
	b := value.Bytes()
	pretty := display.NewPretty()
	prettyPrinted, err := pretty.Print(b)
	if err == nil {
		encodedContentType, err := encodeContentType(prettyPrinted.Type)
		if err != nil {
			return nil, errors.New("error encoding content type")
		}

		return &Pretty{
			ContentType: encodedContentType,
			Value:       prettyPrinted.Value,
		}, nil
	}
	return nil, nil
}

func encodeContentType(t display.ContentType) (string, error) {
	switch t {
	case display.ContentTypeJSON:
		return "json", nil
	case display.ContentTypeString:
		return "string", nil
	case display.ContentTypeCBOR:
		return "cbor", nil
	default:
		return "", errors.New("unknown content type")
	}
}

func canDisplayKeyAsString(b []byte) bool {
	for _, rne := range string(b) {
		if !unicode.IsGraphic(rne) || !unicode.IsSpace(rne) {
			return false
		}
	}
	return true
}
