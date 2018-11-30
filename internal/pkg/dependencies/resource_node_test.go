package dependencies

import (
	"testing"

	"github.com/hashicorp/hcl2/hcl"
)

func TestResourceNode(t *testing.T) {
	nodeID := int64(1)
	nodeName := "nodeName"
	nodeBlock := &hcl.Block{}

	node := ResourceNode{
		name:  nodeName,
		id:    nodeID,
		block: nodeBlock,
	}

	if node.ID() != nodeID {
		t.Errorf("Node ID was incorrect, got: %d, want: %d.", node.ID(), nodeID)
	}

	if node.Name() != nodeName {
		t.Errorf("Node name was incorrect, got: %s, want: %s.", node.Name(), nodeName)
	}

	if node.Block() != nodeBlock {
		t.Errorf("Node block was incorrect, got: %s, want: %s.", node.Block(), nodeBlock)
	}
}
