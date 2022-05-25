package http

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"unicode"

	"github.com/acarl005/stripansi"
	"github.com/boreq/bolt-ui/application"
	"github.com/fxamacker/cbor/v2"
	refmtcbor "github.com/polydawn/refmt/cbor"
	refmtpretty "github.com/polydawn/refmt/pretty"
	refmtshared "github.com/polydawn/refmt/shared"
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

func toTree(tree application.Tree, encodingCBOR bool) Tree {
	return Tree{
		toKeys(tree.Path),
		toEntries(tree.Entries, encodingCBOR),
	}
}

func toKeys(keys []application.Key) []Key {
	result := make([]Key, 0)
	for _, key := range keys {
		result = append(result, toKey(key))
	}
	return result
}

func toEntries(entries []application.Entry, encodingCBOR bool) []Entry {
	result := make([]Entry, 0)
	for _, entry := range entries {
		result = append(result, toEntry(entry, encodingCBOR))
	}
	return result
}

func toEntry(entry application.Entry, encodingCBOR bool) Entry {
	return Entry{
		Bucket: entry.Bucket,
		Key:    toKey(entry.Key),
		Value:  toValue(entry.Value, encodingCBOR),
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

func toValue(value application.Value, encodingCBOR bool) *Value {
	if value.IsEmpty() {
		return nil
	}

	b := value.Bytes()

	result := &Value{
		Hex: hex.EncodeToString(b),
	}

	if encodingCBOR {
		if err := cbor.Valid(b); err == nil {
			t, err := cborToText(b)
			if err == nil && isDisplayableString(t) {
				result.Str = t
				return result
			}
		}
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
	return isDisplayableString(string(b))
}

func isDisplayableString(str string) bool {
	for _, rne := range str {
		if !unicode.IsGraphic(rne) && !unicode.IsSpace(rne) {
			return false
		}
	}
	return true
}

func cborToText(dataCBOR []byte) (string, error) {
	var buf bytes.Buffer
	bufWriter := bufio.NewWriter(&buf)
	err := refmtshared.TokenPump{
		refmtcbor.NewDecoder(refmtcbor.DecodeOptions{}, bytes.NewReader(dataCBOR)),
		refmtpretty.NewEncoder(bufWriter),
	}.Run()
	if err != nil {
		return "", fmt.Errorf("shared.TokenPump.Run() failed: %s", err)
	}
	err = bufWriter.Flush()
	if err != nil {
		return "", fmt.Errorf("bufWriter.Flush() failed: %s", err)
	}
	return stripansi.Strip(buf.String()), nil
}
