package gocty

import (
	"github.com/zclconf/go-cty/cty"
)

// Encoder contains helpers for encoding native go types to go-cty values.
type Encoder struct{}

// IntVal takes an int64 and returns a cty value.
func (*Encoder) IntVal(val int64) cty.Value {
	return cty.NumberIntVal(val)
}

// StringVal takes a string and returns a cty value.
func (*Encoder) StringVal(val string) cty.Value {
	return cty.StringVal(val)
}

// ListVal takes a list of cty values and returns a cty value of type list.
func (*Encoder) ListVal(val []cty.Value) cty.Value {
	return cty.ListVal(val)
}

// ListVal takes a map of cty values and returns a cty value of type list.
func (*Encoder) MapVal(val map[string]cty.Value) cty.Value {
	return cty.ObjectVal(val)
}

// BoolVal takes a boolean and returns a cty value.
func (*Encoder) BoolVal(val bool) cty.Value {
	return cty.BoolVal(val)
}
