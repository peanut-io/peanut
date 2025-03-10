package encoder

const (
	YAML = "yaml"
	JSON = "json"
	TOML = "toml"
	XML  = "xml"
	HCL  = "hcl"
)

type Encoder interface {
	Encode(interface{}) ([]byte, error)
	Decode([]byte, interface{}) error
	String() string
}

var Encoders map[string]func() Encoder

func RegisterEncoder(key string, encoder func() Encoder) {
	if Encoders == nil {
		Encoders = make(map[string]func() Encoder)
	}
	Encoders[key] = encoder
}
