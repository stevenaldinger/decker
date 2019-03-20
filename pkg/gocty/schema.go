package gocty

import (
	"github.com/zclconf/go-cty/cty"
)

// OutputValue is not yet used but organizes plugins outputs. The end goal is
// going to be to make it even easier for plugins to deal with inputs/outputs
// with a layer of abstraction in between them and go-cty.
type OutputValue struct {
	Name  string
	Type  string
	Value cty.Value
}
