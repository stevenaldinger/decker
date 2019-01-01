package main

import (
	"encoding/json"
	"fmt"
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
func (p plugin) Run(inputsMap, resultsMap *map[string]string, resultsListMap *map[string][]string) {
	var (
		cmdOut []byte
		err    error
	)
	var result map[string]interface{}

	// https://www.sohamkamani.com/blog/2017/10/18/parsing-json-in-golang/
	jsonUnmarshalErr := json.Unmarshal([]byte((*inputsMap)["options"]), &result)

	if jsonUnmarshalErr != nil {
		fmt.Println("Error unmarshaling json", jsonUnmarshalErr)
	}

	exploit := (*inputsMap)["exploit"]
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

	for key, val := range result {
		if str, ok := val.(string); ok {
			cmdStr = cmdStr + "set " + key + " " + str + ";"
		} else {
			fmt.Println("Option value is not a string for "+key+":", val)
		}
	}

	cmdStr = cmdStr + "run"

	cmdArgs := []string{"-x", cmdStr}

	if cmdOut, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error running metasploit command: ", err)
		return
	}

	output := string(cmdOut)

	(*resultsMap)["raw_output"] = output
}

// Plugin is an implementation of github.com/stevenaldinger/decker/pkg/plugins.Plugin
// All this includes is a single function, "Run(*map[string]string, *map[string]string)"
// which takes a map of inputs and an empty map of outputs that the Plugin
// is expected to populate
var Plugin plugin
