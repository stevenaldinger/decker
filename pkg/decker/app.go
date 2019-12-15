package decker

import (
	"github.com/zclconf/go-cty/cty"

	"github.com/stevenaldinger/decker/pkg/gocty"
	"github.com/stevenaldinger/decker/pkg/plugins"

	hashicorpHCL "github.com/hashicorp/hcl2/hcl"
)

// Status is the resulting latency on a specific queue.
type Status struct {
	Queue   string
	Latency float64
}

type Env struct {
	Outputs Outputs
}

type Outputs struct {
	XML bool
	JSON bool
}

// App is the main object for decker
type App struct {
	// map to keep track of all the values returned from plugins
	ResultsMap 			 map[string]*map[string]cty.Value
	ResultsMapNested map[string]*map[string]*map[string]cty.Value
	Channel 				 chan Status
	Decoder					 gocty.Decoder
	Encoder					 gocty.Encoder
	Environment			 *map[string]cty.Value
	Env							 Env
	Blocks					 []*hashicorpHCL.Block
	OutputsResults	 plugins.OutputsResults
	ActiveBlock			 Block
}
