package xml

import (
	"encoding/xml"

	"github.com/peanut-io/peanut/config/encoder"
)

func init() {
	encoder.RegisterEncoder(encoder.XML, NewEncoder)
}

type xmlEncoder struct{}

func (x xmlEncoder) Encode(v interface{}) ([]byte, error) {
	return xml.Marshal(v)
}

func (x xmlEncoder) Decode(d []byte, v interface{}) error {
	return xml.Unmarshal(d, v)
}

func (x xmlEncoder) String() string {
	return encoder.XML
}

func NewEncoder() encoder.Encoder {
	return xmlEncoder{}
}
