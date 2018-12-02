package hcl

import (
	"github.com/hashicorp/hcl2/hcl"
	"github.com/stevenaldinger/decker/internal/pkg/paths"
)

// PluginInputConfig is used in conjunction with gohcl to decode HCL attributes
// into native Go structs.
type PluginInputConfig struct {
	Name    string   `hcl:"name,label"`
	Type    string   `hcl:"type"`
	Default string   `hcl:"default"`
	Remain  hcl.Body `hcl:",remain"`
}

// PluginOutputConfig is used in conjunction with gohcl to decode HCL attributes
// into native Go structs.
type PluginOutputConfig struct {
	Name   string   `hcl:"name,label"`
	Type   string   `hcl:"type"`
	Remain hcl.Body `hcl:",remain"`
}

// PluginConfig is used in conjunction with gohcl to decode HCL attributes
// into native Go structs.
type PluginConfig struct {
	Inputs  []PluginInputConfig  `hcl:"input,block"`
	Outputs []PluginOutputConfig `hcl:"output,block"`
}

// VariableConfig is used in conjunction with gohcl to decode HCL attributes
// into native Go structs.
type VariableConfig struct {
	Name    string   `hcl:"name,label"`
	Type    string   `hcl:"type"`
	Default string   `hcl:"default"`
	Remain  hcl.Body `hcl:",remain"`
}

// ResourceConfig is used in conjunction with gohcl to decode HCL attributes
// into native Go structs.
type ResourceConfig struct {
	PluginName    string   `hcl:"plugin_name,label"`
	UniqueID      string   `hcl:"unique_id,label"`
	Host          string   `hcl:"host"`
	PluginEnabled string   `hcl:"plugin_enabled"`
	Remain        hcl.Body `hcl:",remain"`
}

// Config is used in conjunction with gohcl to decode HCL attributes
// into native Go structs.
type Config struct {
	TargetHost string           `hcl:"target_host"`
	Variables  []VariableConfig `hcl:"variable,block"`
	Resources  []ResourceConfig `hcl:"resource,block"`
}

// Variable represents a "variable" block in a module or file.
// type Resource struct {
// 	Name        string
// 	Host        string
// 	DeclRange hcl.Range
// }

// GetConfigFileSchema returns an hcl BodySchema that can be used to decode
// the top level HCL config file
func GetConfigFileSchema() *hcl.BodySchema {
	// configFileSchema is the schema for the top-level of a config file. We use
	// the low-level HCL API for this level so we can easily deal with each
	// block type separately with its own decoding logic.
	configFileSchema := &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type: "terraform",
			},
			{
				Type:       "resource",
				LabelNames: []string{"plugin_name", "unique_id"},
			},
			{
				Type:       "variable",
				LabelNames: []string{"name"},
			},
		},
	}

	return configFileSchema
}

// GetResourceBlockSchema takes a plugin name, determines its schema based
// on its HCL file's inputs, and returns the hcl.BodySchema which can be used
// to decode the "resource" blocks in an HCL config file
func GetResourceBlockSchema(pluginName string) *hcl.BodySchema {
	// can probably handle this better / in a difference spot so the HCL file
	// doesn't need to be read in multiple times
	hclFilePath := paths.GetPluginHCLFilePath(pluginName)
	pluginConfig := parsePluginHCL(hclFilePath)

	inputNames := GetPluginInputNames(pluginConfig)

	pluginInputAttributes := []hcl.AttributeSchema{}

	for _, element := range inputNames {
		pluginInputAttributes = append(pluginInputAttributes, hcl.AttributeSchema{Name: element})
	}

	resourceBlockSchema := &hcl.BodySchema{
		Attributes: pluginInputAttributes,
	}

	return resourceBlockSchema
}

// GetPluginConfigFileSchema returns an hcl BodySchema that can be used to decode
// the top level Plugin HCL config files (specifying inputs and outputs for that
// plugin).
func GetPluginConfigFileSchema() *hcl.BodySchema {
	// configFileSchema is the schema for the top-level of a config file. We use
	// the low-level HCL API for this level so we can easily deal with each
	// block type separately with its own decoding logic.
	configFileSchema := &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type:       "input",
				LabelNames: []string{"input_name"},
			},
			{
				Type:       "output",
				LabelNames: []string{"name"},
			},
		},
	}

	return configFileSchema
}

// GetPluginInputSchema returns an hcl BodySchema that can be used to decode
// a plugin's inputs in a plugin HCL config file
func GetPluginInputSchema() *hcl.BodySchema {
	inputBlockSchema := &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{
				Name: "type",
			},
			{
				Name: "default",
			},
		},
	}

	return inputBlockSchema
}

// GetPluginOutputSchema returns an hcl BodySchema that can be used to decode
// a plugin's outputs in a plugin HCL config file
func GetPluginOutputSchema() *hcl.BodySchema {
	outputBlockSchema := &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{
				Name: "type",
			},
		},
	}

	return outputBlockSchema
}
