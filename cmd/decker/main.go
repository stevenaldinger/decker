package main

import (
	"encoding/json"
	"fmt"
	"github.com/zclconf/go-cty/cty"
	"os"

	"github.com/stevenaldinger/decker/internal/pkg/dependencies"
	"github.com/stevenaldinger/decker/internal/pkg/gocty"
	"github.com/stevenaldinger/decker/internal/pkg/hcl"
	"github.com/stevenaldinger/decker/internal/pkg/paths"
	"github.com/stevenaldinger/decker/internal/pkg/plugins"
	"github.com/stevenaldinger/decker/internal/pkg/reports"
)

// loop over list of strings and return true if list contains a given string
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func main() {
	// get the path to a given config file passed in via CLI
	hclConfigFile := paths.GetConfigFilePath()

	blocks := hcl.GetBlocksFromConfig(hclConfigFile)
	varBlockNames := dependencies.GetVariableNames(blocks)
	resBlocksSorted := dependencies.Sort(blocks)

	// get environment variable values for each declared variable block
	var envVarCtx = hcl.GetHCLEvalContextVarsFromEnv(varBlockNames)

	// map to keep track of all the values returned from plugins
	resultsMapCty := map[string]*map[string]cty.Value{}
	resultsMapCtyNested := map[string]*map[string]*map[string]cty.Value{}

	decoder := gocty.Decoder{}
	encoder := gocty.Encoder{}

	for index, block := range resBlocksSorted {
		resourcePlugin, resourceName := block.Labels[0], block.Labels[1]

		pluginHCLPath := paths.GetPluginHCLFilePath(resourcePlugin)

		pluginAttrs := hcl.GetPluginAttributes(block)

		containsForEach := contains(pluginAttrs, "for_each")

		hclConfig, pluginContent := hcl.GetPluginContent(containsForEach, block, pluginHCLPath)

		if containsForEach {
			// returns JSON, not sure why
			forEachDecoded := hcl.DecodeHCLListAttribute(pluginContent.Attributes["for_each"], envVarCtx, &resultsMapCty, &resultsMapCtyNested)

			var forEachList []string

			jsonUnmarshalErr := json.Unmarshal([]byte(forEachDecoded), &forEachList)

			if jsonUnmarshalErr != nil {
				fmt.Println("Error unmarshaling json", jsonUnmarshalErr)
				os.Exit(1)
			}

			for _, eachKey := range forEachList {
				var forEachMap = &map[string]cty.Value{
					"key": encoder.StringVal(string(eachKey)),
				}

				resultsMapCty["each"] = forEachMap

				inputsMap := hcl.CreateInputsMapCty(hclConfig.Inputs, pluginContent.Attributes, envVarCtx, &resultsMapCty, &resultsMapCtyNested)

				// declare a new empty map to be passed into the plugin
				var resultsMap = map[string]cty.Value{}
				var resultsListMap = map[string][]cty.Value{}

				pluginEnabled := plugins.RunIfEnabled(resourcePlugin, &inputsMap, &resultsMap, &resultsListMap)

				if pluginEnabled {
					fmt.Println(fmt.Sprintf("DECKER: Ran plugin %d[%s] of %d: %s (%s)", index+1, string(eachKey), len(resBlocksSorted), resourcePlugin, resourceName))
				} else {
					fmt.Println(fmt.Sprintf("DECKER: [Disabled] Did not run plugin %d of %d: %s (%s)", index+1, len(resBlocksSorted), resourcePlugin, resourceName))

					resultsMap["raw_output"] = encoder.StringVal("")
				}

				// initialize if map doesn't exist yet
				if _, ok := resultsMapCtyNested[resourceName]; !ok {
					var initMap = &map[string]*map[string]cty.Value{}
					// var initMap = &map[string]cty.Value {}
					resultsMapCtyNested[resourceName] = initMap
				}

				// build eval context from plugin results and add it to the ongoing map
				(*resultsMapCtyNested[resourceName])[string(eachKey)] = &resultsMap

				if pluginEnabled {
					reports.WriteStringToFile(paths.GetReportFilePath(resourceName+"["+string(eachKey)+"]"), decoder.GetString(resultsMap["raw_output"]))
				}
			}
		} else {
			inputsMap := hcl.CreateInputsMapCty(hclConfig.Inputs, pluginContent.Attributes, envVarCtx, &resultsMapCty, &resultsMapCtyNested)

			// declare a new empty map to be passed into the plugin
			var resultsMap = map[string]cty.Value{}
			var resultsListMap = map[string][]cty.Value{}

			pluginEnabled := plugins.RunIfEnabled(resourcePlugin, &inputsMap, &resultsMap, &resultsListMap)

			if pluginEnabled {
				fmt.Println(fmt.Sprintf("DECKER: Ran plugin %d of %d: %s (%s)", index+1, len(resBlocksSorted), resourcePlugin, resourceName))

				// build eval context from plugin results and add it to the ongoing map
				resultsMapCty[resourceName] = &resultsMap

				reports.WriteStringToFile(paths.GetReportFilePath(resourceName), decoder.GetString(resultsMap["raw_output"]))
			} else {
				fmt.Println(fmt.Sprintf("DECKER: [Disabled] Did not run plugin %d of %d: %s (%s)", index+1, len(resBlocksSorted), resourcePlugin, resourceName))

				resultsMap["raw_output"] = encoder.StringVal("")

				// build eval context from plugin results and add it to the ongoing map
				resultsMapCty[resourceName] = &resultsMap
			}
		}
	}
}
