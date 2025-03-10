package hcl

import (
	"github.com/hashicorp/hcl"
	jsoniter "github.com/json-iterator/go"

	"github.com/peanut-io/peanut/config/encoder"
)

func init() {
	encoder.RegisterEncoder(encoder.HCL, NewEncoder)
}

type hclEncoder struct{}

func (h hclEncoder) Encode(v interface{}) ([]byte, error) {
	return jsoniter.Marshal(v)
}

func (h hclEncoder) Decode(d []byte, v interface{}) error {
	return hcl.Unmarshal(d, v)
}

func (h hclEncoder) String() string {
	return encoder.HCL
}

func NewEncoder() encoder.Encoder {
	return hclEncoder{}
}
