package hcl

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

// DecodeHCLResourceBlock calls BuildEvalContext() with the plugin results aggregated from each
// iterative run and attempts to decode a Resource block with the latest
// context (this is the magic allowing outputs from one plugin to work
// as inputs to another)
func DecodeHCLResourceBlock(block *hcl.Block, envVals *map[string]cty.Value, runningVals *map[string]*map[string]cty.Value, runningValsNested *map[string]*map[string]*map[string]cty.Value) *ResourceConfig {
	var c ResourceConfig

	// will return evalcontext with environment variables
	ctx := BuildEvalContext(envVals, runningVals, runningValsNested)

	diags := gohcl.DecodeBody(block.Body, ctx, &c)

	if diags.HasErrors() {
		fmt.Println("Error decoding resource block body:", diags)
	}

	return &c
}

// DecodeHCLAttributeCty calls BuildEvalContext() with the plugin results aggregated from each
// iterative run and attempts to decode a Block's Attribute's Expression
// using the context
func DecodeHCLAttributeCty(attribute *hcl.Attribute, envVals *map[string]cty.Value, runningVals *map[string]*map[string]cty.Value, runningValsNested *map[string]*map[string]*map[string]cty.Value, defVal string) cty.Value {
	// will return evalcontext with environment variables
	ctx := BuildEvalContext(envVals, runningVals, runningValsNested)

	ctyVal, _ := attribute.Expr.Value(ctx)

	return ctyVal
}

// DecodeHCLListAttribute calls BuildEvalContext() with the plugin results aggregated from each
// iterative run and attempts to decode a Block's Attribute's Expression
// using the context. This is used in the "for_each" logic, not sure if its still needed.
func DecodeHCLListAttribute(attribute *hcl.Attribute, envVals *map[string]cty.Value, runningVals *map[string]*map[string]cty.Value, runningValsNested *map[string]*map[string]*map[string]cty.Value) string {
	// will return evalcontext with environment variables
	ctx := BuildEvalContext(envVals, runningVals, runningValsNested)

	ctyVal, _ := attribute.Expr.Value(ctx)

	fmt.Println("Cty val:", ctyVal)

	var decodedArrVal string
	var decodedArray = []string{}

	ctyValArr := ctyVal.AsValueSlice()

	for _, val := range ctyValArr {
		err := gocty.FromCtyValue(val, &decodedArrVal)
		if err != nil {
			fmt.Println("Error trying to decode cty val in arr.", err)
		} else {
			decodedArray = append(decodedArray, decodedArrVal)
		}
	}

	jsonBytes, jsonErr := json.Marshal(decodedArray)
	if jsonErr != nil {
		fmt.Println("json.Marshal(decodedArray) err:", jsonErr)
		// exit the program here?
		return ""
	}

	jsonString := string(jsonBytes)
	return jsonString
}
