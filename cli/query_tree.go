package cli

import (
	"fmt"
)

type QueryTreeNode struct {
	CurrentNode string                    // name of the current node
	SubNodes    map[string]*QueryTreeNode // list of the next categories
	Leaves      []string                  // list of the queries at this level in the tree
}

var QueryTree = QueryTreeNode{
	SubNodes: map[string]*QueryTreeNode{},
	Leaves:   []string{},
}

func (node *QueryTreeNode) AddQueryName(queryNameParts []string, method string) {
	if len(queryNameParts) == 1 {
		leaveName := fmt.Sprintf("%s (%s)", queryNameParts[0], method)
		node.Leaves = append(node.Leaves, leaveName)
		return
	}

	existingNode, ok := node.SubNodes[queryNameParts[0]]
	if ok {
		existingNode.AddQueryName(queryNameParts[1:], method)
	} else {
		newNode := QueryTreeNode{
			CurrentNode: queryNameParts[0],
			SubNodes:    map[string]*QueryTreeNode{},
			Leaves:      []string{},
		}
		node.SubNodes[queryNameParts[0]] = &newNode
		newNode.AddQueryName(queryNameParts[1:], method)
	}

}

// print all the node in the tree until the requested depth is reached
// pass depth = -1 for printing all the tree
func (node *QueryTreeNode) Print(depth int, prefix string) {
	fmt.Printf("%s%s\n", prefix, node.CurrentNode)
	for _, leaf := range node.Leaves {
		fmt.Printf("%s%s\n", prefix+"    ", leaf)
	}
	if depth == 0 {
		return
	} else if len(node.SubNodes) == 0 {
		return
	} else {
		if depth > 0 {
			depth -= 1
		}
		for _, subNode := range node.SubNodes {
			subNode.Print(depth, prefix+"    ")
		}
	}
}
