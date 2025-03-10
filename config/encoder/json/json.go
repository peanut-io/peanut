package json

import (
	jsoniter "github.com/json-iterator/go"

	"github.com/peanut-io/peanut/config/encoder"
)

func init() {
	encoder.RegisterEncoder(encoder.JSON, NewEncoder)
}

type jsonEncoder struct{}

func (j jsonEncoder) Encode(v interface{}) ([]byte, error) {
	return jsoniter.Marshal(v)
}

func (j jsonEncoder) Decode(d []byte, v interface{}) error {
	return jsoniter.Unmarshal(d, v)
}

func (j jsonEncoder) String() string {
	return encoder.JSON
}

func NewEncoder() encoder.Encoder {
	return jsonEncoder{}
}
