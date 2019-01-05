package dependencies

import (
	"fmt"

	hashicorpHCL "github.com/hashicorp/hcl2/hcl"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/topo"

	"github.com/stevenaldinger/decker/internal/pkg/hcl"
	"github.com/stevenaldinger/decker/internal/pkg/paths"
)

func addNodesToGraph(dag *simple.DirectedGraph, blocks []*hashicorpHCL.Block, resourceIDs *map[string]int64, resourceNodes *map[int64]ResourceNode) {
	for _, block := range blocks {
		switch block.Type {
		case "resource":
			nodeUniqueID := int64(len((*resourceIDs)))
			nodeUniqueName := block.Labels[1]
			(*resourceIDs)[nodeUniqueName] = nodeUniqueID
			resourceNode := ResourceNode{
				name:  nodeUniqueName,
				id:    int64(nodeUniqueID),
				block: block,
			}
			(*resourceNodes)[nodeUniqueID] = resourceNode
			(*dag).AddNode(resourceNode)

			// case "variable":
			//   fmt.Println("Havent written variable decoder yet")
		}
	}
}

// loop over list of strings and return true if list contains a given string
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func addEdgesToGraph(dag *simple.DirectedGraph, blocks []*hashicorpHCL.Block, resourceIDs *map[string]int64, resourceNodes *map[int64]ResourceNode) {
	for _, block := range blocks {
		switch block.Type {
		case "resource":
			// block.Labels[0] == plugin
			// block.Labels[1] == unique name
			nodePluginName := block.Labels[0]
			nodeUniqueName := block.Labels[1]
			nodeUniqueID := (*resourceIDs)[nodeUniqueName]
			resourceNode := (*resourceNodes)[nodeUniqueID]

			exprVars := hcl.GetExprVars(block)

			for attr := range exprVars {
				for _, exprVar := range exprVars[attr] {
					dependentOn := exprVar.RootName()

					// add node edge if the dependency is on something other than var
					if dependentOn != "var" {
						// fmt.Println(fmt.Sprintf("Dependency detected: %s dependent on %s", block.Labels[1], dependentOn))

						dag.SetEdge(ResourceEdge{
							from: (*resourceNodes)[(*resourceIDs)[dependentOn]],
							to:   resourceNode,
						})
					}
				}
			}

			// handle for_each attributes
			pluginHCLPath := paths.GetPluginHCLFilePath(nodePluginName)
			pluginAttrs := hcl.GetPluginAttributes(block)
			containsForEach := contains(pluginAttrs, "for_each")

			if containsForEach {
				_, pluginContent := hcl.GetPluginContent(containsForEach, block, pluginHCLPath)
				dependentOn := pluginContent.Attributes["for_each"].Expr.Variables()[0].RootName()
				if dependentOn != "var" {
					dag.SetEdge(ResourceEdge{
						from: (*resourceNodes)[(*resourceIDs)[dependentOn]],
						to:   resourceNode,
					})
				}
			}

			// case "variable":
			//   fmt.Println("Havent written variable decoder yet")
		}
	}
}

// Sort takes an array of *hcl.Block and returns an array of the same blocks,
// topologically sorted by their dependency graph. This ensures all blocks
// are evaluated after the blocks they depend on.
func Sort(blocks []*hashicorpHCL.Block) []*hashicorpHCL.Block {
	// key: resource unique ID
	// value: unique int id for Node
	var resourceIDsForDepGraph = map[string]int64{}
	var resourceNodesForDepGraph = map[int64]ResourceNode{}
	//   // initialize graph
	dag := simple.NewDirectedGraph()
	addNodesToGraph(dag, blocks, &resourceIDsForDepGraph, &resourceNodesForDepGraph)
	addEdgesToGraph(dag, blocks, &resourceIDsForDepGraph, &resourceNodesForDepGraph)

	sortedDag, err := topo.Sort(dag)

	if err != nil {
		fmt.Println("Error sorting topology:", err)
	}

	var orderedBlocks = []*hashicorpHCL.Block{}

	for i := 0; i < len(sortedDag); i++ {
		orderedBlocks = append(orderedBlocks, resourceNodesForDepGraph[sortedDag[i].ID()].Block())
	}

	return orderedBlocks
}
