package dependencies

import (
	"testing"

	"github.com/hashicorp/hcl2/hcl"
)

func TestResourceEdge(t *testing.T) {
	nodeID1 := int64(1)
	nodeID2 := int64(2)
	nodeName1 := "node_name"
	nodeName2 := "node_name_2"
	nodeBlock1 := &hcl.Block{}
	nodeBlock2 := &hcl.Block{}

	node1 := ResourceNode{
		name:  nodeName1,
		id:    nodeID1,
		block: nodeBlock1,
	}

	node2 := ResourceNode{
		name:  nodeName2,
		id:    nodeID2,
		block: nodeBlock2,
	}

	edge := ResourceEdge{
		from: node1,
		to:   node2,
	}

	if edge.From().ID() != node1.ID() {
		t.Errorf("FROM node ID was incorrect, got: %d, want: %d.", edge.From().ID(), node1.ID())
	}

	if edge.To().ID() != node2.ID() {
		t.Errorf("TO node ID was incorrect, got: %d, want: %d.", edge.To().ID(), node2.ID())
	}

	if edge.Weight() != 0.0 {
		t.Errorf("Node weight was incorrect, got: %f, want: %f.", edge.Weight(), 0.0)
	}
}
