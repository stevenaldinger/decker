package plugins

import (
	"fmt"
	"github.com/stevenaldinger/decker/pkg/paths"
	"github.com/stevenaldinger/decker/pkg/gocty"
	"github.com/zclconf/go-cty/cty"
	"os"
	"plugin"
)

// Plugin is an interface for each Plugin to implement. All this includes is a single
// function, "Run(*map[string]string, *map[string]string)" which takes a map
// of inputs and an empty map of outputs that the Plugin is expected to populate
type Plugin interface {
	// Run(inputsMap, resultsMap)
	Run(*map[string]cty.Value, *map[string]cty.Value)
}

// OutputPlugin is an interface for each Plugin to implement. All this includes is a single
// function, "Run(*map[string]string, *map[string]string)" which takes a map
// of inputs and an empty map of outputs that the Plugin is expected to populate
type OutputPlugin interface {
	// Run(inputsMap)
	Run(*map[string]*map[string]cty.Value, *OutputsResults)
}

// runPlugin takes the name of a plugin, a map of inputs, and an empty map for results,
// this function will load the plugin based on Decker convention (plugin is
// expected to be in ./internal/app/decker/plugins/PLUGIN_NAME and the .so file
// is expected to be PLUGIN_NAME.so) and call it with the maps supplied to
// it as arguments.
func runPlugin(name string, inputsMap, resultsMap *map[string]cty.Value, resultsListMap *map[string][]cty.Value) {
	mod := paths.GetPluginPath(name) + "/" + name + ".so"

	// load module - open the .so file to load the symbols
	plug, err := plugin.Open(mod)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// look up a symbol (exported variable - Plugin)
	symPlugin, err := plug.Lookup("Plugin")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// 3. Assert that loaded symbol is of interface type - Plugin
	var plugin Plugin
	plugin, ok := symPlugin.(Plugin)
	if !ok {
		fmt.Println("unexpected type from module symbol Plugin")
		os.Exit(1)
	}

	plugin.Run(inputsMap, resultsMap)
}

// RunIfEnabled takes the name of a plugin, a map of inputs, and an empty map for results,
// and will run the plugin if its inputs map does not contain plugin_enabled == "false".
func RunIfEnabled(pluginName string, inputsMap, resultsMap *map[string]cty.Value, resultsListMap *map[string][]cty.Value) bool {
	decoder := gocty.Decoder{}

	// only run the plugin if plugin_enabled != "false"
	if decoder.GetStringOrBool((*inputsMap)["plugin_enabled"]) != "false" {
		runPlugin(pluginName, inputsMap, resultsMap, resultsListMap)

		return true
	}

	return false
}

// OutputResults holds XML and JSON string representations of a plugin's outputs
type OutputResults struct {
	XML  string
	JSON string
}

// OutputsResults holds XML and JSON string representations of the combined
// plugin's outputs as well as a map of the XML/JSON string representations of
// individual plugins' outputs
type OutputsResults struct {
	Results map[string]OutputResults
	AllXML  string
	AllJSON string
}

// RunOutputsPlugin takes a map of inputs, runs the "outputs" plugin, and
// returns JSON/XML string representations of the inputs.
func RunOutputsPlugin(inputsMap *map[string]*map[string]cty.Value) OutputsResults {
	mod := paths.GetPluginPath("outputs") + "/outputs.so"

	// load module - open the .so file to load the symbols
	plug, err := plugin.Open(mod)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// look up a symbol (exported variable - Plugin)
	symPlugin, err := plug.Lookup("Plugin")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// 3. Assert that loaded symbol is of interface type - Plugin
	var plugin OutputPlugin
	plugin, ok := symPlugin.(OutputPlugin)
	if !ok {
		fmt.Println("unexpected type from module symbol OutputPlugin")
		os.Exit(1)
	}

	outputsResults := OutputsResults{
		Results: map[string]OutputResults{},
		AllXML:  "",
		AllJSON: "",
	}

	plugin.Run(inputsMap, &outputsResults)

	return outputsResults
}
