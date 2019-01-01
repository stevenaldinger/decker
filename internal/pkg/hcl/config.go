package hcl

import (
	"fmt"
	"os"

	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hclparse"
)

// GetBlocksFromConfig takes an HCL file path and return an go native Config,
// ordered array of the plugin names that should be run, and an ordered array
// of hcl.Block that need to be parsed
func GetBlocksFromConfig(hclFilePath string) []*hcl.Block {
	jsonHCL := []byte(os.Getenv("DECKER_RUN_CONFIGURATION"))
	parser := hclparse.NewParser()

	var f *hcl.File
	var diags hcl.Diagnostics

	if len(jsonHCL) == 0 {
		f, diags = parser.ParseHCLFile(hclFilePath)
	} else {
		f, diags = parser.ParseJSON(jsonHCL, "cache-file.cache")
	}

	if diags.HasErrors() {
		// should possibly exit the program here
		fmt.Println("Error parsing config HCL:", diags)
	}

	// configFileSchema is the schema for the top-level of a config file. We use
	// the low-level HCL API for this level so we can easily deal with each
	// block type separately with its own decoding logic.
	configFileSchema := GetConfigFileSchema()
	content, contentDiags := f.Body.Content(configFileSchema)

	if contentDiags.HasErrors() {
		// should possibly exit the program here
		fmt.Println("Error getting config file content:", contentDiags)
	}

	return content.Blocks
}
