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
func parsePluginHCLFromFile(hclFilePath string) (*PluginConfig, *hcl.BodySchema) {
	hclConfig := parsePluginHCL(hclFilePath)

	pluginInputAttributes := []hcl.AttributeSchema{}

	for _, element := range hclConfig.Inputs {
		// should eventually use element.Type and element.Default to support types
		// other than string and to support unrequired inputs with default values
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

// GetPluginContent takes a *hcl.Block and a path to an HCL config file and
// returns the BodyContent
func GetPluginContent(block *hcl.Block, hclFilePath string) *hcl.BodyContent {
	// hcl_plugin_config, pluginInputSchema
	_, pluginInputSchema := parsePluginHCLFromFile(hclFilePath)

	pluginContent, pluginDiags := block.Body.Content(pluginInputSchema)

	if pluginDiags.HasErrors() {
		// should possibly exit the program here
		fmt.Println("Error getting Plugin block's content.", pluginDiags)
	}

	return pluginContent
}

// CreateInputsMap decodes the HCL attributes with an evaluation context
// consisting of the outputs of all previously run plugins
func CreateInputsMap(attributes hcl.Attributes, evalVals *map[string]*map[string]cty.Value) map[string]string {
	// declare inputsMap with default "plugin_enabled" = true
	// the rest is pulled from the specific plugin's HCL file "input" blocks
	// ex: "internal/app/decker/plugins/nslookup/nslookup.hcl"
	var inputsMap = map[string]string{
		"plugin_enabled": "true",
	}

	for key, attribute := range attributes {
		inputsMap[key] = DecodeHCLAttribute(attribute, evalVals)
	}

	return inputsMap
}
