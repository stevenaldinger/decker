package main

import (
	// "encoding/json"
	"fmt"
	"github.com/stevenaldinger/decker/internal/pkg/gocty"
	"github.com/zclconf/go-cty/cty"
	"os"
	"os/exec"
)

type plugin string

// --- wants input: ---
// inputsMap{
//   "exploit": "example.com",
//   "options": "{ ... }".
//   "plugin_enabled": "true",
//   "db_enabled": "false",
// }
func (p plugin) Run(inputsMap, resultsMap *map[string]cty.Value) {
	var (
		cmdOut []byte
		err    error
	)

	decoder := gocty.Decoder{}
	encoder := gocty.Encoder{}

	options := decoder.GetMap((*inputsMap)["options"])

	exploit := decoder.GetString((*inputsMap)["exploit"])
	// dbEnabled := (*inputsMap)["db_enabled"]

	cmdName := "msfconsole"

	// if dbEnabled == "true" {
	// 	dbStatusCmd := []string{"-x", "db_status"}
	//
	// 	if cmdOut, err = exec.Command(cmdName, dbStatusCmd...).Output(); err != nil {
	// 		fmt.Fprintln(os.Stderr, "There was an error fetching database status:", err)
	// 		return
	// 	}
	//
	// 	// search for this string in results
	// 	// "postgresql selected, no connection"
	// 	if notInitialized {
	// 		initCmd := "msfdb"
	// 		initArgs := []string{"init"}
	// 		if dbCmdOut, dbErr = exec.Command(initCmd, initArgs...).Output; dbErr != nil {
	// 			fmt.Fprintln(os.Stderr, "There was an error initializing metasploit database:", dbErr)
	// 			return
	// 		}
	//
	// 		if cmdOut, err = exec.Command(cmdName, dbStatusCmd...).Output(); err != nil {
	// 			fmt.Fprintln(os.Stderr, "There was an error fetching database status:", err)
	// 			return
	// 		}
	// 		// search for this string in results
	// 		// "postgresql selected, no connection"
	// 	}
	// }

	cmdStr := "use " + exploit + ";"

	for key, val := range *options {
		cmdStr = cmdStr + "set " + key + " " + val + ";"
	}

	cmdStr = cmdStr + "run"

	cmdArgs := []string{"-x", cmdStr}

	if cmdOut, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error running metasploit command: ", err)
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
