package hcl

import (
	"fmt"

	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hclparse"
	"github.com/zclconf/go-cty/cty"
)

func parsePluginHCL(hclFilePath string) *PluginConfig {
	parser := hclparse.NewParser()

	f, diags := parser.ParseHCLFile(hclFilePath)

	var c PluginConfig

	// var build_context_empty_eval_map = map[string]*map[string]cty.Value{}
	// will return evalcontext with environment variables
	// ctx := BuildEvalContext(&build_context_empty_eval_map)

	moreDiags := gohcl.DecodeBody(f.Body, nil, &c)
	// print diags if debugging
	_ = append(diags, moreDiags...)

	// configFileSchema is the schema for the top-level of a config file. We use
	// the low-level HCL API for this level so we can easily deal with each
	// block type separately with its own decoding logic.
	configFileSchema := GetPluginConfigFileSchema()
	// content, content_diags := f.Body.Content(configFileSchema)
	f.Body.Content(configFileSchema)

	return &c
}

// parsePluginHCLFromFile builds Input and Output body schemas from parsed plugin HCL
func parsePluginHCLFromFile(forEach bool, hclFilePath string) (*PluginConfig, *hcl.BodySchema) {
	hclConfig := parsePluginHCL(hclFilePath)

	pluginInputAttributes := []hcl.AttributeSchema{}

	// this keeps plugins from having to specify a for_each input in their hcl
	// and allows it to be optional
	if forEach {
		pluginInputAttributes = append(pluginInputAttributes, hcl.AttributeSchema{Name: "for_each"})
	}

	for _, element := range hclConfig.Inputs {
		pluginInputAttributes = append(pluginInputAttributes, hcl.AttributeSchema{Name: element.Name})
	}

	pluginInputSchema := &hcl.BodySchema{
		Attributes: pluginInputAttributes,
	}

	pluginOutputAttributes := []hcl.AttributeSchema{}

	for _, element := range hclConfig.Outputs {
		pluginOutputAttributes = append(pluginOutputAttributes, hcl.AttributeSchema{Name: element.Name})
	}

	// === NOT YET USED, SHOULD BE USED LATER FOR CHECKING VALIDITY OF CONFIG FILES
	// plugin_output_schema := &hcl.BodySchema{
	//   Attributes: pluginOutputAttributes,
	// }

	return hclConfig, pluginInputSchema
}

// GetPluginInputNames takes an hcl.PluginConfig and returns an array of the
// names of its inputs (for building a resource block's schema)
func GetPluginInputNames(pluginConfig *PluginConfig) []string {
	inputNames := []string{}

	for _, element := range pluginConfig.Inputs {
		inputNames = append(inputNames, element.Name)
	}

	return inputNames
}

// GetPluginAttributes returns a list of the plugin's attributes. This is used
// to determine whether or not for_each is set so main.go loops over the plugin
// instead of running it once.
func GetPluginAttributes(block *hcl.Block) []string {
	var attributes = []string{}

	attrs, diags := block.Body.JustAttributes()

	for attr := range attrs {
		attributes = append(attributes, attr)
	}

	if diags.HasErrors() {
		fmt.Println("Error getting plugin attributes:", diags)
	}

	return attributes
}

// GetPluginContent takes a *hcl.Block and a path to an HCL config file and
// returns the BodyContent
func GetPluginContent(forEach bool, block *hcl.Block, hclFilePath string) (*PluginConfig, *hcl.BodyContent) {
	// hclConfig, pluginInputSchema
	hclConfig, pluginInputSchema := parsePluginHCLFromFile(forEach, hclFilePath)

	pluginContent, pluginDiags := block.Body.Content(pluginInputSchema)

	if pluginDiags.HasErrors() {
		// should possibly exit the program here
		fmt.Println("Error getting Plugin block's content.", pluginDiags)
	}

	return hclConfig, pluginContent
}

// CreateInputsMap decodes the HCL attributes with an evaluation context
// consisting of the outputs of all previously run plugins
func CreateInputsMap(inputs []PluginInputConfig, attributes hcl.Attributes, evalVals *map[string]*map[string]cty.Value) map[string]string {
	// declare inputsMap with default "plugin_enabled" = true
	// the rest is pulled from the specific plugin's HCL file "input" blocks
	// ex: "internal/app/decker/plugins/nslookup/nslookup.hcl"
	var inputsMap = map[string]string{
		"plugin_enabled": "true",
	}

	var hclInputs = map[string]PluginInputConfig{}

	for _, element := range inputs {
		hclInputs[element.Name] = element
	}

	// create a map of attribute names to inputs and get Input.Type to determine
	// which DecodeHCL...Attribute function to call
	// pass in default in case parsing fails
	for key, attribute := range attributes {
		inputType := hclInputs[key].Type
		inputDefault := hclInputs[key].Default

		if key != "for_each" {
			if inputType == "list" {
				inputsMap[key] = DecodeHCLListAttribute(attribute, evalVals)
			} else if inputType == "map" {
				inputsMap[key] = DecodeHCLMapAttribute(attribute, evalVals)
			} else {
				// strings and booleans
				inputsMap[key] = DecodeHCLAttribute(attribute, evalVals, inputDefault)
			}
		}
	}

	return inputsMap
}
