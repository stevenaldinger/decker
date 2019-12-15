package dependencies

import (
	"fmt"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/stevenaldinger/decker/pkg/plugins"
	"os"
)

// ValidateConfig checks the given config for some common errors and exits with
// error messages and how to fix the issues.
func ValidateConfig(blocks []*hcl.Block) {
	availablePlugins := plugins.AvailablePlugins()

	pluginNames := []string{}
	invalidPluginNames := []string{}
	resourceNames := []string{}
	duplicateResourceNames := []string{}

	for _, block := range blocks {
		switch block.Type {
		case "resource":
			resourcePlugin, resourceName := block.Labels[0], block.Labels[1]
			validPlugin := contains(availablePlugins, resourcePlugin)
			if validPlugin {
				pluginNames = append(pluginNames, resourcePlugin)
			} else {
				invalidPluginNames = append(invalidPluginNames, resourcePlugin)
			}

			duplicateResourceName := contains(resourceNames, resourceName)
			if duplicateResourceName {
				duplicateResourceNames = append(duplicateResourceNames, resourceName)
			} else {
				resourceNames = append(resourceNames, resourceName)
			}
		}
	}

	if len(invalidPluginNames) != 0 {
		fmt.Println("[ERROR]: Invalid plugin names were referenced in your config. The following plugins could not be found by decker:", invalidPluginNames)
	}

	if len(duplicateResourceNames) != 0 {
		fmt.Println("[ERROR]: Duplicate resource names were found in your config. The following resource names were used more than once:", duplicateResourceNames)
		fmt.Println("[ERROR]: Please give each resource a unique name and try again.")
	}

	if len(invalidPluginNames) != 0 || len(duplicateResourceNames) != 0 {
		os.Exit(1)
	}
}
