package models

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
func (node *QueryTreeNode) Print(depth int, prefix string) string {
	res := fmt.Sprintf("%s/%s\n", prefix, node.CurrentNode)
	// fmt.Println(res)
	for _, leaf := range node.Leaves {
		res += fmt.Sprintf("%s/%s\n", prefix+"    ", leaf)
	}
	if depth == 0 || len(node.SubNodes) == 0 {
		return res
	} else {
		if depth > 0 {
			depth -= 1
		}
		for _, subNode := range node.SubNodes {
			res += subNode.Print(depth, prefix+"    ")
		}
	}
	return res
}
