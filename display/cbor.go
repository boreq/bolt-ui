package display

import (
	"bufio"
	"bytes"

	"github.com/acarl005/stripansi"
	"github.com/boreq/errors"
	"github.com/fxamacker/cbor/v2"
	refmtcbor "github.com/polydawn/refmt/cbor"
	refmtpretty "github.com/polydawn/refmt/pretty"
	refmtshared "github.com/polydawn/refmt/shared"
)

type PrettifierCBOR struct {
}

func NewPrettifierCBOR() *PrettifierCBOR {
	return &PrettifierCBOR{}
}

func (p PrettifierCBOR) Prettify(b []byte) (string, error) {
	if err := cbor.Wellformed(b); err != nil {
		return "", errors.Wrap(err, "invalid cbor")
	}
	return cborToText(b)
}

// from https://github.com/boreq/bolt-ui/pull/2
func cborToText(dataCBOR []byte) (string, error) {
	var buf bytes.Buffer
	bufWriter := bufio.NewWriter(&buf)
	err := refmtshared.TokenPump{
		TokenSource: refmtcbor.NewDecoder(refmtcbor.DecodeOptions{}, bytes.NewReader(dataCBOR)),
		TokenSink:   refmtpretty.NewEncoder(bufWriter),
	}.Run()
	if err != nil {
		return "", errors.Wrap(err, "tokenpump run failed")
	}
	err = bufWriter.Flush()
	if err != nil {
		return "", errors.Wrap(err, "error flushing the buffer")
	}
	return stripansi.Strip(buf.String()), nil
}
