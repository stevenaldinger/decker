package main

import (
	"fmt"
	"github.com/stevenaldinger/decker/pkg/gocty"
	"github.com/zclconf/go-cty/cty"
	"os"
	"os/exec"
)

type plugin string

// --- wants input: ---
// inputsMap{
//   "say_hello_to": "world",
// }
//
// --- gives output: ---
// resultsMap{
//  "raw_output": "...",
// }
func (p plugin) Run(inputsMap, resultsMap *map[string]cty.Value) {
	var (
		cmdOut []byte
		err    error
	)

	decoder := gocty.Decoder{}
	encoder := gocty.Encoder{}

	sayHelloTo := decoder.GetString((*inputsMap)["say_hello_to"])

	cmdName := "echo"
	cmdArgs := []string{"Hello,", sayHelloTo}

	if cmdOut, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error running hello_world plugin: ", err)
		return
	}

	output := string(cmdOut)

	(*resultsMap)["said_hello_to"] = encoder.StringVal(sayHelloTo)
	(*resultsMap)["raw_output"] = encoder.StringVal(output)
}

var Plugin plugin
