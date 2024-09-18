package display_test

import (
	"github.com/boreq/bolt-ui/display"
	"github.com/fxamacker/cbor/v2"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPretty(t *testing.T) {
	cbor, err := cbor.Marshal(struct {
		Field1 string `cbor:"1,keyasint"`
		Field2 int    `cbor:"2,keyasint"`
	}{
		Field1: "string",
		Field2: 123,
	})
	require.NoError(t, err)

	testCases := []struct {
		Name   string
		Bytes  []byte
		Result display.Prettified
	}{
		{
			Name:  "json",
			Bytes: []byte(`{"some":"json"}`),
			Result: display.Prettified{
				Type: display.ContentTypeJSON,
				Value: `{
  "some": "json"
}`,
			},
		},
		{
			Name:  "cbor",
			Bytes: cbor,
			Result: display.Prettified{
				Type:  display.ContentTypeCBOR,
				Value: "Map<len:2> {\n\r\t1: \"string\"\n\r\t2: 123\n\r}\n\r",
			},
		},
		{
			Name:  "string",
			Bytes: []byte("some_string"),
			Result: display.Prettified{
				Type:  display.ContentTypeString,
				Value: "some_string",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			p := display.NewPretty()
			result, err := p.Print(testCase.Bytes)
			require.NoError(t, err)
			require.Equal(t, testCase.Result, result)
		})
	}
}
