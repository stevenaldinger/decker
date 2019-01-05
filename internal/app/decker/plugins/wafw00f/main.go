package main

import (
	"fmt"
	"github.com/stevenaldinger/decker/internal/pkg/gocty"
	"github.com/zclconf/go-cty/cty"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type plugin string

// --- wants input: ---
// inputsMap{
//   "host": "example.com",
// }
//
// --- gives output: ---
// resultsMap{
//  "raw_output": "...",
// 	"waf_detected": "true",
// }
func (p plugin) Run(inputsMap, resultsMap *map[string]cty.Value) {
	var (
		cmdOut []byte
		err    error
	)

	decoder := gocty.Decoder{}
	encoder := gocty.Encoder{}

	targetHost := decoder.GetString((*inputsMap)["host"])

	cmdName := "wafw00f"
	cmdArgs := []string{targetHost}

	if cmdOut, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error running wafw00f command: ", err)
		return
	}

	output := string(cmdOut)

	(*resultsMap)["raw_output"] = encoder.StringVal(output)
	(*resultsMap)["waf_detected"] = encoder.StringVal(strconv.FormatBool(!strings.Contains(output, "No WAF detected by the generic detection")))
}

// Plugin is an implementation of github.com/stevenaldinger/decker/pkg/plugins.Plugin
// All this includes is a single function, "Run(*map[string]string, *map[string]string)"
// which takes a map of inputs and an empty map of outputs that the Plugin
// is expected to populate
var Plugin plugin
