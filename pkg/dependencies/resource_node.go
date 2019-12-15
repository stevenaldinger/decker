package dependencies

import (
	"github.com/hashicorp/hcl2/hcl"
)

// ResourceNode implements graph.Node and represents a resource node in the
// dependency graph DAG
type ResourceNode struct {
	id    int64
	name  string
	block *hcl.Block
}

// ID returns the unique int64 id of a ResourceNode
func (N ResourceNode) ID() int64 {
	return N.id
}

// Block returns the hcl.Block of a ResourceNode
func (N ResourceNode) Block() *hcl.Block {
	return N.block
}

// Name returns the unique string name of a ResourceNode
func (N ResourceNode) Name() string {
	return N.name
}
