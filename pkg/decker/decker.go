package decker

import (
	"encoding/json"
	"fmt"
	"github.com/zclconf/go-cty/cty"
	"os"

	"github.com/stevenaldinger/decker/pkg/dependencies"
	"github.com/stevenaldinger/decker/pkg/hcl"
	"github.com/stevenaldinger/decker/pkg/paths"
	"github.com/stevenaldinger/decker/pkg/reports"
	"github.com/stevenaldinger/decker/pkg/gocty"
	"github.com/stevenaldinger/decker/pkg/plugins"

	log "github.com/sirupsen/logrus"
)

// Init returns a new App object
func Init() *App {
	resultsMap := map[string]*map[string]cty.Value{}
	resultsMapNested := map[string]*map[string]*map[string]cty.Value{}
	c := make(chan Status)
	decoder := gocty.Decoder{}
	encoder := gocty.Encoder{}

	a := &App{
		ResultsMap: resultsMap,
		ResultsMapNested: resultsMapNested,
		Channel: c,
		Decoder: decoder,
		Encoder: encoder,
	}

  makeDirectory(paths.GetReportsDir())

	// get the path to a given config file passed in via CLI
	hclConfigFile := paths.GetConfigFilePath()

	blocks := hcl.GetBlocksFromConfig(hclConfigFile)

	dependencies.ValidateConfig(blocks)

	varBlockNames := dependencies.GetVariableNames(blocks)
	a.Blocks = dependencies.Sort(blocks)

	// get environment variable values for each declared variable block
	a.Environment = hcl.GetHCLEvalContextVarsFromEnv(varBlockNames)
	log.Trace("Got environment context")

	a.Env = Env{
		Outputs: Outputs{
			XML: os.Getenv("DECKER_OUTPUTS_XML") == "true",
			JSON: os.Getenv("DECKER_OUTPUTS_JSON") == "true",
		},
	}

	return a
}

// GetResults gets the plugin results in XML and JSON formats
func (a *App) GetResults() {
	a.OutputsResults = plugins.RunOutputsPlugin(&a.ResultsMap)

	if a.Env.Outputs.XML {
		reports.WriteStringToFile(paths.GetReportFilePath("All", "xml"), a.OutputsResults.AllXML)

		for uniqueName, outputVal := range a.OutputsResults.Results {
			reports.WriteStringToFile(paths.GetReportFilePath(uniqueName, "xml"), outputVal.XML)
		}
	}

	if a.Env.Outputs.JSON {
		reports.WriteStringToFile(paths.GetReportFilePath("All", "json"), a.OutputsResults.AllJSON)

		for uniqueName, outputVal := range a.OutputsResults.Results {
			reports.WriteStringToFile(paths.GetReportFilePath(uniqueName, "json"), outputVal.JSON)
		}
	}
}

func (a *App) RunPlugin(b Block, p Plugin) {
	p.InputsMap = hcl.CreateInputsMapCty(b.HCLConfig.Inputs, p.Key, b.PluginContent.Attributes, a.Environment, &a.ResultsMap, &a.ResultsMapNested)
  p.Enabled = plugins.RunIfEnabled(b.PluginName, &p.InputsMap, &p.ResultsMap, &p.ResultsListMap)

	if b.ForEach {
		// initialize if map doesn't exist yet
		if _, ok := a.ResultsMapNested[b.ResourceName]; !ok {
			var initMap = &map[string]*map[string]cty.Value{}
			// var initMap = &map[string]cty.Value {}
			a.ResultsMapNested[b.ResourceName] = initMap
		}

		// build eval context from plugin results and add it to the ongoing map
		(*a.ResultsMapNested[b.ResourceName])[p.Key] = &p.ResultsMap

		if p.Enabled {
			// unaware of its index if I do it this way, should invert control later
			// fmt.Println(fmt.Sprintf("DECKER: Ran plugin %d of %d: %s (%s)", index+1, len(a.Blocks), b.PluginName, b.ResourceName))

			reports.WriteStringToFile(paths.GetReportFilePath(b.ResourceName+"["+p.Key+"]", "txt"), a.Decoder.GetString(p.ResultsMap["raw_output"]))
		}
	} else {
		if p.Enabled {
		 	a.ResultsMap[b.ResourceName] = &p.ResultsMap
		} else {
			p.ResultsMap["raw_output"] = a.Encoder.StringVal("")
		}
		reports.WriteStringToFile(paths.GetReportFilePath(b.ResourceName, "txt"), a.Decoder.GetString(p.ResultsMap["raw_output"]))
	}
}

// RunPlugins runs all the plugins
func (a *App) RunPlugins() {
	// for index, block := range a.Blocks {
	for _, block := range a.Blocks {
    b := NewBlock(block)

		if b.ForEach {
			// returns JSON, not sure why
			forEachDecoded := hcl.DecodeHCLListAttribute(b.PluginContent.Attributes["for_each"], a.Environment, &a.ResultsMap, &a.ResultsMapNested)

			var forEachList []string

			jsonUnmarshalErr := json.Unmarshal([]byte(forEachDecoded), &forEachList)

			if jsonUnmarshalErr != nil {
				fmt.Println("Error unmarshaling json", jsonUnmarshalErr)
				os.Exit(1)
			}

			for _, eachKey := range forEachList {
        p := NewPlugin(eachKey)
        a.RunPlugin(b, p)
			}
		} else {
      p := NewPlugin("0")
      a.RunPlugin(b, p)
		}
	}
}
