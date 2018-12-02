package hcl

import (
	"os"
	"strings"

	"github.com/hashicorp/hcl2/hcl"
	"github.com/zclconf/go-cty/cty"
)

// Gets all environment variables with prefix "DECKER_" and creates a map
// with keys equal to the environment variable name but with "DECKER_" prefix
// stripped away, and the rest of the name is set to lower case. The values
// of the environment variables remain untouched. This allows HCL configuration
// blocks to use environment variables like "${var.my_environmment_variable}".
func getHCLEvalContextVarsFromEnv() *map[string]cty.Value {
	var Var = map[string]cty.Value{}

	// gets DECKER_* env vars for var.* values in HCL
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		if strings.HasPrefix(pair[0], "DECKER_") {
			Var[strings.ToLower(strings.TrimPrefix(pair[0], "DECKER_"))] = cty.StringVal(pair[1])
		}
	}

	return &Var
}

// BuildEvalContextFromMap takes a *map[string]string and returns a map of
// cty.Value to be used in an hcl EvalContext
func BuildEvalContextFromMap(m *map[string]string) *map[string]cty.Value {
	var variables = map[string]cty.Value{}

	for key, value := range *m {
		variables[key] = cty.StringVal(value)
	}

	return &variables
}

// BuildEvalContext builds an HCL evaluation context with all "DECKER_" environment variables
// available using "var" prefix in config files, and also loops over all the
// aggregated results maps from plugins that have run and makes them available
// for the next round of HCL decoding.
func BuildEvalContext(runningVals *map[string]*map[string]cty.Value) *hcl.EvalContext {
	// func BuildEvalContext() (*hcl.EvalContext) {
	var Variables = map[string]cty.Value{}

	evalContextVariables := getHCLEvalContextVarsFromEnv()

	Variables["var"] = cty.ObjectVal(*evalContextVariables)

	for key, element := range *runningVals {
		Variables[key] = cty.ObjectVal(*element)
	}

	ctx := &hcl.EvalContext{
		Variables: Variables,
	}

	return ctx
}
