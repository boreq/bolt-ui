package display

import (
	"github.com/boreq/errors"
)

type ContentType struct {
	s string
}

var (
	ContentTypeJSON   ContentType = ContentType{"json"}
	ContentTypeCBOR   ContentType = ContentType{"cbor"}
	ContentTypeString ContentType = ContentType{"string"}
)

type Prettifier interface {
	Prettify(b []byte) (string, error)
}

type Prettified struct {
	Type  ContentType
	Value string
}

type prettifier struct {
	Prettifier  Prettifier
	ContentType ContentType
}

type Pretty struct {
	prettifiers []prettifier
}

func NewPretty() *Pretty {
	return &Pretty{prettifiers: []prettifier{
		{
			Prettifier:  NewPrettifierCBOR(),
			ContentType: ContentTypeCBOR,
		},
		{
			Prettifier:  NewPrettifierJSON(),
			ContentType: ContentTypeJSON,
		},
		{
			Prettifier:  NewPrettifierString(),
			ContentType: ContentTypeString,
		},
	}}
}

func (p *Pretty) Print(b []byte) (Prettified, error) {
	for _, prettifier := range p.prettifiers {
		v, err := prettifier.Prettifier.Prettify(b)
		if err == nil {
			return Prettified{
				Type:  prettifier.ContentType,
				Value: v,
			}, nil
		}
	}
	return Prettified{}, errors.New("no prettifiers completed successfully")
}
