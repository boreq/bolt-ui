package http

import (
	"encoding/hex"
	"encoding/json"
	"unicode"

	"github.com/boreq/bolt-ui/application"
)

type Tree struct {
	Path    []Key   `json:"path"`
	Entries []Entry `json:"entries"`
}

type Key struct {
	Hex string `json:"hex"`
	Str string `json:"str,omitempty"`
}

type Value struct {
	Hex string `json:"hex"`
	Str string `json:"str,omitempty"`
}

type Entry struct {
	Bucket bool   `json:"bucket"`
	Key    Key    `json:"key"`
	Value  *Value `json:"value,omitempty"`
}

func toTree(tree application.Tree) Tree {
	return Tree{
		toKeys(tree.Path),
		toEntries(tree.Entries),
	}
}

func toKeys(keys []application.Key) []Key {
	result := make([]Key, 0)
	for _, key := range keys {
		result = append(result, toKey(key))
	}
	return result
}

func toEntries(entries []application.Entry) []Entry {
	result := make([]Entry, 0)
	for _, entry := range entries {
		result = append(result, toEntry(entry))
	}
	return result
}

func toEntry(entry application.Entry) Entry {
	return Entry{
		Bucket: entry.Bucket,
		Key:    toKey(entry.Key),
		Value:  toValue(entry.Value),
	}
}

func toKey(key application.Key) Key {
	b := key.Bytes()

	result := Key{
		Hex: hex.EncodeToString(b),
	}

	if canDisplayAsString(b) {
		result.Str = string(b)
	}

	return result
}

func toValue(value application.Value) *Value {
	if value.IsEmpty() {
		return nil
	}

	b := value.Bytes()

	result := &Value{
		Hex: hex.EncodeToString(b),
	}

	if canDisplayAsString(b) {
		result.Str = string(b)
	}

	return result
}

func canDisplayAsString(b []byte) bool {
	if json.Valid(b) {
		return true
	}

	for _, rne := range string(b) {
		if !unicode.IsGraphic(rne) {
			return false
		}
	}

	return true
}
