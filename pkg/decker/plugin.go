package decker

import (
	"github.com/zclconf/go-cty/cty"
)

func NewPlugin(key string) Plugin {
  p := Plugin{}
  // declare a new empty map to be passed into the plugin
  p.ResultsMap = map[string]cty.Value{}
  p.ResultsListMap = map[string][]cty.Value{}
  p.InputsMap = map[string]cty.Value{}
	p.Key = key

  return p
}

type Plugin struct {
  ResultsMap       map[string]cty.Value
  ResultsListMap   map[string][]cty.Value
  InputsMap        map[string]cty.Value
  Enabled          bool
	Key 						 string
}
