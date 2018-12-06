package hcl

import (
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/hcl2/hcl"
	"github.com/zclconf/go-cty/cty"
)

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// GetHCLEvalContextVarsFromEnv gets all environment variables with prefix "DECKER_" and creates a map
// with keys equal to the environment variable name but with "DECKER_" prefix
// stripped away, and the rest of the name is set to lower case. The values
// of the environment variables remain untouched. This allows HCL configuration
// blocks to use environment variables like "${var.my_environmment_variable}".
func GetHCLEvalContextVarsFromEnv(varNames []string) *map[string]cty.Value {
	var Var = map[string]cty.Value{}

	// gets DECKER_* env vars for var.* values in HCL
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		if strings.HasPrefix(pair[0], "DECKER_") {
			strippedName := strings.ToLower(strings.TrimPrefix(pair[0], "DECKER_"))
			containsVar := contains(varNames, strippedName)

			if containsVar {
				Var[strippedName] = cty.StringVal(pair[1])
			}
		}
	}

	// check that every variable name that should be defined is defined.
	// this should probably just be combined with the for loop above instead
	// of using os.Environ() and looping through all of them.
	for _, e := range varNames {
		if _, ok := Var[e]; !ok {
			fmt.Println("Environment variable", "DECKER_"+strings.ToUpper(e), "not defined.")
			os.Exit(1)
		}
	}

	return &Var
}

// BuildEvalContextFromMap takes a *map[string]string and returns a map of
// cty.Value to be used in an hcl EvalContext
func BuildEvalContextFromMap(m *map[string]string, lm *map[string][]string) *map[string]cty.Value {
	var variables = map[string]cty.Value{}

	for key, value := range *m {
		variables[key] = cty.StringVal(value)
	}

	for key, value := range *lm {
		var listVars = []cty.Value{}
		for _, listVal := range value {
			listVars = append(listVars, cty.StringVal(listVal))
		}

		variables[key] = cty.ListVal(listVars)
	}

	return &variables
}

// BuildEvalContext builds an HCL evaluation context with all "DECKER_" environment variables
// available using "var" prefix in config files, and also loops over all the
// aggregated results maps from plugins that have run and makes them available
// for the next round of HCL decoding.
func BuildEvalContext(envVarsCtx *map[string]cty.Value, runningVals *map[string]*map[string]cty.Value, runningValsNested *map[string]*map[string]*map[string]cty.Value) *hcl.EvalContext {
	var Variables = map[string]cty.Value{}

	Variables["var"] = cty.ObjectVal(*envVarsCtx)

	for key, element := range *runningVals {
		Variables[key] = cty.ObjectVal(*element)
	}

	for key, nestedVal := range *runningValsNested {
		for nestedKey, element := range *nestedVal {
			if _, ok := Variables[key]; ok {
				// key already exists in map, need to merge them
				Variables[key] = cty.ObjectVal(map[string]cty.Value{
					key:       Variables[key],
					nestedKey: cty.ObjectVal(*element),
				})
			} else {
				Variables[key] = cty.ObjectVal(map[string]cty.Value{
					nestedKey: cty.ObjectVal(*element),
				})
			}
		}
	}

	ctx := &hcl.EvalContext{
		Variables: Variables,
	}

	return ctx
}
