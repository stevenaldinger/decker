package main

import (
	"fmt"
	"github.com/stevenaldinger/decker/internal/pkg/plugins"
	"github.com/zclconf/go-cty/cty"
	ctyjson "github.com/zclconf/go-cty/cty/json"

	"encoding/json"
	// "encoding/xml"
	// "os"
	// "os/exec"
	"github.com/clbanning/mxj"
	// "bytes"
)

type plugin string

// --- individual json / xml for outputs ---

// returns json string, xml string
func convert(input *map[string]cty.Value) (string, string) {
	simpleJSONVal := ctyjson.SimpleJSONValue{cty.ObjectVal(*input)}

	buf, err := json.Marshal(simpleJSONVal)
	if err != nil {
		fmt.Println("unexpected error from json.Marshal:", err)
		return "", ""
	}

	var f interface{}
	unmarshalErr := json.Unmarshal(buf, &f)
	if unmarshalErr != nil {
		fmt.Println("unexpected error from json.Unmarshal:", unmarshalErr)
		return "", ""
	}

	m, nmjErr := mxj.NewMapJson(buf)
	if nmjErr != nil {
		fmt.Println("unexpected error from mxj.NewMapJson:", nmjErr)
		return "", ""
	}

	xmlValue, xmlErr := m.Xml()
	if xmlErr != nil {
		fmt.Println("unexpected error from m.Xml():", xmlErr)
		return "", ""
	}

	return string(buf), string(xmlValue)
}

func (p plugin) Run(inputsMap *map[string]*map[string]cty.Value, outputsResults *plugins.OutputsResults) {
	var outputMap = map[string]cty.Value{}

	results := map[string]plugins.OutputResults{}

	// key = unique resource name
	for key, val := range *inputsMap {
		outputMap[key] = cty.ObjectVal(*val)
		jsonStr, xmlStr := convert(val)
		results[key] = plugins.OutputResults{
			XML:  xmlStr,
			JSON: jsonStr,
		}
	}

	allJSON, allXML := convert(&outputMap)

	*outputsResults = plugins.OutputsResults{
		AllXML:  allXML,
		AllJSON: allJSON,
		Results: results,
	}
}

// Plugin is an implementation of github.com/stevenaldinger/decker/pkg/plugins.Plugin
// All this includes is a single function, "Run(*map[string]string, *map[string]string)"
// which takes a map of inputs and an empty map of outputs that the Plugin
// is expected to populate
var Plugin plugin
