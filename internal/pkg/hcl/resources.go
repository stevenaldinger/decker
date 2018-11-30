package hcl

import (
	"fmt"

	"github.com/hashicorp/hcl2/hcl"
)

func getResourceContent(block *hcl.Block) *hcl.BodyContent {
	resourceBlockSchema := GetResourceBlockSchema()

	// content, remain, diags
	content, _, diags := block.Body.PartialContent(resourceBlockSchema)

	if diags.HasErrors() {
		fmt.Println("Error getting resource body partial content:", diags)
	}

	return content
}

// GetExprVars takes a block and a list of attribute names and will return a
// map of all the expression variables for those attributes.
func GetExprVars(block *hcl.Block) map[string][]hcl.Traversal {
	content := getResourceContent(block)

	var exprVars = map[string][]hcl.Traversal{}

	for attr := range content.Attributes {
		exprVars[attr] = content.Attributes[attr].Expr.Variables()
	}

	return exprVars
}
