package source

import (
	"context"
	"github.com/peanut-io/peanut/config/encoder"
	_ "github.com/peanut-io/peanut/config/encoder/hcl"
	_ "github.com/peanut-io/peanut/config/encoder/json"
	_ "github.com/peanut-io/peanut/config/encoder/toml"
	_ "github.com/peanut-io/peanut/config/encoder/xml"
	_ "github.com/peanut-io/peanut/config/encoder/yaml"
)

type Options struct {
	// Encoder
	Encoder encoder.Encoder

	// for alternative data
	Context context.Context
}

var defaultEncoder = encoder.JSON

func NewOptions(format string) *Options {

	if len(format) == 0 {
		format = defaultEncoder
	}
	options := &Options{
		Encoder: encoder.Encoders[format](),
		Context: context.Background(),
	}
	return options
}
