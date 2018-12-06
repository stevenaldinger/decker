package hcl

import (
	"encoding/json"
	"fmt"
	"os"
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

// DecodeHCLAttribute calls BuildEvalContext() with the plugin results aggregated from each
// iterative run and attempts to decode a Block's Attribute's Expression
// using the context
func DecodeHCLAttribute(attribute *hcl.Attribute, envVals *map[string]cty.Value, runningVals *map[string]*map[string]cty.Value, runningValsNested *map[string]*map[string]*map[string]cty.Value, defVal string) string {
	// will return evalcontext with environment variables
	ctx := BuildEvalContext(envVals, runningVals, runningValsNested)

	ctyVal, _ := attribute.Expr.Value(ctx)

	var decodedVal string
	var decodedBool bool
	err := gocty.FromCtyValue(ctyVal, &decodedVal)

	if err != nil {
		boolErr := gocty.FromCtyValue(ctyVal, &decodedBool)
		if boolErr != nil {
			if defVal != "" {
				decodedVal = defVal
			} else {
				fmt.Println("Decoding error for string and bool and default value was an empty string:", boolErr)
				os.Exit(1)
			}
		} else {
			decodedVal = strconv.FormatBool(decodedBool)
		}
	}

	return decodedVal
}

// DecodeHCLListAttribute calls BuildEvalContext() with the plugin results aggregated from each
// iterative run and attempts to decode a Block's Attribute's Expression
// using the context
func DecodeHCLListAttribute(attribute *hcl.Attribute, envVals *map[string]cty.Value, runningVals *map[string]*map[string]cty.Value, runningValsNested *map[string]*map[string]*map[string]cty.Value) string {
	// will return evalcontext with environment variables
	ctx := BuildEvalContext(envVals, runningVals, runningValsNested)

	ctyVal, _ := attribute.Expr.Value(ctx)

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

// DecodeHCLMapAttribute calls BuildEvalContext() with the plugin results aggregated from each
// iterative run and attempts to decode a Block's Attribute's Expression
// using the context
func DecodeHCLMapAttribute(attribute *hcl.Attribute, envVals *map[string]cty.Value, runningVals *map[string]*map[string]cty.Value, runningValsNested *map[string]*map[string]*map[string]cty.Value) string {
	// will return evalcontext with environment variables
	ctx := BuildEvalContext(envVals, runningVals, runningValsNested)

	ctyVal, _ := attribute.Expr.Value(ctx)

	var decodedMapVal string
	var decodedBool bool
	var decodedMap = map[string]string{}

	// if this errors out use default?
	ctyValMap := ctyVal.AsValueMap()

	for key, val := range ctyValMap {
		err := gocty.FromCtyValue(val, &decodedMapVal)
		if err != nil {
			boolErr := gocty.FromCtyValue(val, &decodedBool)
			if boolErr != nil {
				fmt.Println("Error trying to decode cty val for string and bool in map:", boolErr)
				// exit the program here?
			} else {
				decodedMap[key] = strconv.FormatBool(decodedBool)
			}
		} else {
			decodedMap[key] = decodedMapVal
		}
	}

	jsonBytes, jsonErr := json.Marshal(decodedMap)

	if jsonErr != nil {
		fmt.Println("json.Marshal(decodedMap) err:", jsonErr)
		// exit the program here?
		return ""
	}

	jsonString := string(jsonBytes)
	return jsonString
}
