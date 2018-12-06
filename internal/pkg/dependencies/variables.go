package dependencies

import (
	hashicorpHCL "github.com/hashicorp/hcl2/hcl"
)

// GetVariableNames takes a list of all blocks and returns the names with
// type == "variable"
func GetVariableNames(blocks []*hashicorpHCL.Block) []string {
	var varBlockNames = []string{}

	for _, block := range blocks {
		if block.Type == "variable" {
			varBlockNames = append(varBlockNames, block.Labels[0])
		}
	}

	return varBlockNames
}
