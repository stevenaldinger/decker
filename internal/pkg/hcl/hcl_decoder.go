package hcl

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

// DecodeHCLResourceBlock calls BuildEvalContext() with the plugin results aggregated from each
// iterative run and attempts to decode a Resource block with the latest
// context (this is the magic allowing outputs from one plugin to work
// as inputs to another)
func DecodeHCLResourceBlock(block *hcl.Block, runningVals *map[string]*map[string]cty.Value) *ResourceConfig {
	var c ResourceConfig

	// will return evalcontext with environment variables
	ctx := BuildEvalContext(runningVals)

	// print diagnostics if debugging
	// diags := gohcl.DecodeBody(block.Body, ctx, &c)
	gohcl.DecodeBody(block.Body, ctx, &c)

	return &c
}

// DecodeHCLAttribute calls BuildEvalContext() with the plugin results aggregated from each
// iterative run and attempts to decode a Block's Attribute's Expression
// using the context
func DecodeHCLAttribute(attribute *hcl.Attribute, runningVals *map[string]*map[string]cty.Value) string {
	// will return evalcontext with environment variables
	ctx := BuildEvalContext(runningVals)

	ctyVal, _ := attribute.Expr.Value(ctx)

	var decodedVal string
	var decodedBool bool
	// err := gocty.FromCtyValue(actual_val, &val_test)
	err := gocty.FromCtyValue(ctyVal, &decodedVal)

	if err != nil {
		boolErr := gocty.FromCtyValue(ctyVal, &decodedBool)
		if boolErr != nil {
			fmt.Println("Decoding error for both string and bool:", boolErr)
		} else {
			decodedVal = strconv.FormatBool(decodedBool)
		}
	}

	return decodedVal
}
