package main

import (
	"bytes"
	"fmt"
	"github.com/stevenaldinger/decker/pkg/gocty"
	"github.com/zclconf/go-cty/cty"
	"os"
	"os/exec"
	"text/template"
)

type plugin string

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// writeStringToFile takes a file path and a content string and writes the
// content to a file with the given path
func writeStringToFile(filePath, str string) {
	f, err := os.Create(filePath)
	check(err)

	defer f.Close()

	// written, err := f.WriteString("writes\n")
	f.WriteString(str)
	// fmt.Printf("wrote string:\n", written)

	// flush
	f.Sync()
}

// ScriptVariables contains variables that will be injected into the w3af script.
type ScriptVariables struct {
	TargetHost string
	Verbose    string
}

const w3afScript = `
# all usage demo!

plugins
output console,text_file
output
output config text_file
set output_file /tmp/reports/output-w3af.txt
set verbose True
back
output config console
set verbose False
back

crawl all, !bing_spider, !google_spider, !spider_man
crawl

grep all
grep

audit all
audit

back

target
set target http://{{.TargetHost}}/
back

start

exit
`

// --- wants input: ---
// inputsMap{
//   "host": "example.com",
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

	targetHost := decoder.GetString((*inputsMap)["host"])
	filePath := "/tmp/w3af-script"

	cmdName := "w3af_console"
	cmdArgs := []string{"-s", filePath}

	scriptVariables := ScriptVariables{
		TargetHost: targetHost,
	}

	t := template.New("w3af script template")

	t, err = t.Parse(w3afScript)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error parsing w3af script:", err)
		return
	}

	var tpl bytes.Buffer
	if tErr := t.Execute(&tpl, scriptVariables); tErr != nil {
		fmt.Fprintln(os.Stderr, "Error executing w3af script template:", tErr)
		return
	}

	result := tpl.String()

	writeStringToFile(filePath, result)

	if cmdOut, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error running the w3af tool: ", err)
		return
	}

	output := string(cmdOut)

	(*resultsMap)["raw_output"] = encoder.StringVal(output)
}

// Plugin is an implementation of github.com/stevenaldinger/decker/pkg/plugins.Plugin
// All this includes is a single function, "Run(*map[string]string, *map[string]string)"
// which takes a map of inputs and an empty map of outputs that the Plugin
// is expected to populate
var Plugin plugin
