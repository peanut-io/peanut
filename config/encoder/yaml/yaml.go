package yaml

import (
	"github.com/ghodss/yaml"

	"github.com/peanut-io/peanut/config/encoder"
)

func init() {
	encoder.RegisterEncoder(encoder.YAML, NewEncoder)
}

type yamlEncoder struct{}

func (y yamlEncoder) Encode(v interface{}) ([]byte, error) {
	return yaml.Marshal(v)
}

func (y yamlEncoder) Decode(d []byte, v interface{}) error {
	return yaml.Unmarshal(d, v)
}

func (y yamlEncoder) String() string {
	return encoder.YAML
}

func NewEncoder() encoder.Encoder {
	return yamlEncoder{}
}
