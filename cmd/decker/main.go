package main

import (
	"encoding/json"
	"fmt"
	"github.com/zclconf/go-cty/cty"
	"os"

	"github.com/stevenaldinger/decker/internal/pkg/dependencies"
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

	resBlocks := hcl.GetResourceBlocksFromConfig(hclConfigFile)
	resBlocksSorted := dependencies.Sort(resBlocks)

	// map to keep track of all the values returned from plugins
	evalCtxVals := map[string]*map[string]cty.Value{}

	for index, block := range resBlocksSorted {
		resourcePlugin, resourceName := block.Labels[0], block.Labels[1]

		pluginHCLPath := paths.GetPluginHCLFilePath(resourcePlugin)

		pluginAttrs := hcl.GetPluginAttributes(block)

		containsForEach := contains(pluginAttrs, "for_each")

		hclConfig, pluginContent := hcl.GetPluginContent(containsForEach, block, pluginHCLPath)

		if containsForEach {
			// returns JSON, not sure why
			forEachDecoded := hcl.DecodeHCLListAttribute(pluginContent.Attributes["for_each"], &evalCtxVals)

			var forEachList []string

			jsonUnmarshalErr := json.Unmarshal([]byte(forEachDecoded), &forEachList)

			if jsonUnmarshalErr != nil {
				fmt.Println("Error unmarshaling json", jsonUnmarshalErr)
				os.Exit(1)
			}

			for _, eachKey := range forEachList {
				var forEachMap = map[string]string{
					"key": string(eachKey),
				}
				// not used for anything right now, just need it for the function call
				var forEachListMap = map[string][]string{}

				evalCtxVals["each"] = hcl.BuildEvalContextFromMap(&forEachMap, &forEachListMap)

				inputsMap := hcl.CreateInputsMap(hclConfig.Inputs, pluginContent.Attributes, &evalCtxVals)

				// declare a new empty map to be passed into the plugin
				var resultsMap = map[string]string{}
				var resultsListMap = map[string][]string{}

				pluginEnabled := plugins.RunIfEnabled(resourcePlugin, &inputsMap, &resultsMap, &resultsListMap)

				if pluginEnabled {
					fmt.Println(fmt.Sprintf("DECKER: Ran plugin %d[%s] of %d: %s (%s)", index+1, string(eachKey), len(resBlocksSorted), resourcePlugin, resourceName))
				} else {
					fmt.Println(fmt.Sprintf("DECKER: [Disabled] Did not run plugin %d of %d: %s (%s)", index+1, len(resBlocksSorted), resourcePlugin, resourceName))
				}

				// build eval context from plugin results and add it to the ongoing map
				evalCtxVals[resourceName] = hcl.BuildEvalContextFromMap(&resultsMap, &resultsListMap)

				reports.WriteStringToFile(paths.GetReportFilePath(resourceName+"["+string(eachKey)+"]"), resultsMap["raw_output"])
			}
		} else {
			inputsMap := hcl.CreateInputsMap(hclConfig.Inputs, pluginContent.Attributes, &evalCtxVals)

			// declare a new empty map to be passed into the plugin
			var resultsMap = map[string]string{}
			var resultsListMap = map[string][]string{}

			pluginEnabled := plugins.RunIfEnabled(resourcePlugin, &inputsMap, &resultsMap, &resultsListMap)

			if pluginEnabled {
				fmt.Println(fmt.Sprintf("DECKER: Ran plugin %d of %d: %s (%s)", index+1, len(resBlocksSorted), resourcePlugin, resourceName))
			} else {
				fmt.Println(fmt.Sprintf("DECKER: [Disabled] Did not run plugin %d of %d: %s (%s)", index+1, len(resBlocksSorted), resourcePlugin, resourceName))
			}

			// build eval context from plugin results and add it to the ongoing map
			evalCtxVals[resourceName] = hcl.BuildEvalContextFromMap(&resultsMap, &resultsListMap)

			reports.WriteStringToFile(paths.GetReportFilePath(resourceName), resultsMap["raw_output"])
		}
	}
}
