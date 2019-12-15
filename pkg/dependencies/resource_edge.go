package dependencies

import (
	"gonum.org/v1/gonum/graph"
)

// ResourceEdge implements graph.Edge and represents an edge between two
// resources in the dependency graph DAG
type ResourceEdge struct {
	from ResourceNode
	to   ResourceNode
}

// From returns the graph.Node (ResourceNode) that represents the "from" node
// in an edge of the dependency graph DAG
func (E ResourceEdge) From() graph.Node {
	return E.from
}

// To returns the graph.Node (ResourceNode) that represents the "to" node
// in an edge of the dependency graph DAG
func (E ResourceEdge) To() graph.Node {
	return E.to
}

// Weight returns a hardcoded weight for an edge of the dependency graph DAG
// This isn't actually used for anything right now but needs to be implemented,
// so it returns 0.0 for every edge.
func (E ResourceEdge) Weight() float64 {
	return 0.0
}
