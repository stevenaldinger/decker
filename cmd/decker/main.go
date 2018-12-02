package main

import (
	"fmt"
	"github.com/zclconf/go-cty/cty"

	"github.com/stevenaldinger/decker/internal/pkg/dependencies"
	"github.com/stevenaldinger/decker/internal/pkg/hcl"
	"github.com/stevenaldinger/decker/internal/pkg/paths"
	"github.com/stevenaldinger/decker/internal/pkg/plugins"
	"github.com/stevenaldinger/decker/internal/pkg/reports"
)

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

		hclConfig, pluginContent := hcl.GetPluginContent(block, pluginHCLPath)

		inputsMap := hcl.CreateInputsMap(hclConfig.Inputs, pluginContent.Attributes, &evalCtxVals)

		// declare a new empty map to be passed into the plugin
		var resultsMap = map[string]string{}

		pluginEnabled := plugins.RunIfEnabled(resourcePlugin, &inputsMap, &resultsMap)

		if pluginEnabled {
			fmt.Println(fmt.Sprintf("DECKER: Ran plugin %d of %d: %s (%s)", index+1, len(resBlocksSorted), resourcePlugin, resourceName))
		} else {
			fmt.Println(fmt.Sprintf("DECKER: [Disabled] Did not run plugin %d of %d: %s (%s)", index+1, len(resBlocksSorted), resourcePlugin, resourceName))
		}

		// build eval context from plugin results and add it to the ongoing map
		evalCtxVals[resourceName] = hcl.BuildEvalContextFromMap(&resultsMap)

		reports.WriteStringToFile(paths.GetReportFilePath(resourceName), resultsMap["raw_output"])
	}
}
