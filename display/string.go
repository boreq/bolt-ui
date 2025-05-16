package display

import (
	"unicode"

	"github.com/boreq/errors"
)

type PrettifierString struct {
}

func NewPrettifierString() *PrettifierString {
	return &PrettifierString{}
}

func (p *PrettifierString) Prettify(b []byte) (string, error) {
	if p.canDisplayAsString(b) {
		return string(b), nil
	}
	return "", errors.New("can't display as string")
}

func (p *PrettifierString) canDisplayAsString(b []byte) bool {
	for _, rne := range string(b) {
		if !unicode.IsGraphic(rne) && !unicode.IsSpace(rne) {
			return false
		}
	}
	return true
}
