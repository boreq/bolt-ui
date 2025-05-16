package display

import (
	"bytes"
	"encoding/json"

	"github.com/boreq/errors"
)

type PrettifierJSON struct {
}

func NewPrettifierJSON() *PrettifierJSON {
	return &PrettifierJSON{}
}

func (p PrettifierJSON) Prettify(b []byte) (string, error) {
	if json.Valid(b) {
		buf := &bytes.Buffer{}
		if err := json.Indent(buf, b, "", "  "); err != nil {
			return "", errors.Wrap(err, "error indenting")
		}
		return buf.String(), nil
	}
	return "", errors.New("invalid json")
}
