package hhelper

import (
	"fmt"
	"regexp"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
	"golang.org/x/exp/slices"
)

// what is leaf
// it is value node
func IsLeaf(node ast.Node) bool {
	switch node.(type) {
	case *ast.CommentNode,
		*ast.IntegerNode,
		*ast.FloatNode,
		*ast.StringNode,

		*ast.BoolNode,
		*ast.InfinityNode,
		*ast.NanNode,
		*ast.MergeKeyNode,
		*ast.NullNode:
		return true
	case *ast.LiteralNode:
	case *ast.DirectiveNode:
	case *ast.TagNode:
	case *ast.DocumentNode:
	case *ast.MappingNode:
	case *ast.MappingKeyNode:
	case *ast.MappingValueNode:
	case *ast.SequenceNode:
	case *ast.AnchorNode:
	case *ast.AliasNode:
		return false
	}
	return false
}
func IsExistInGolden(snode ast.Node, golden ast.Node) bool {
	pathStr := snode.GetPath()
	fmt.Print(snode.String())
	sampleRegex := regexp.MustCompile(`\[\d+\]$`)
	widePath := sampleRegex.ReplaceAllString(pathStr, "[*]")
	//path, err := yaml.PathString(pathStr)
	path, err := yaml.PathString(widePath)

	if err != nil {
		fmt.Printf("path error,%v\n", err)
		return false
	}
	node, err := path.FilterNode(golden)
	if err != nil { //we are ok for not found
		// fmt.Printf("filter fail: %v\n", err)
		return false
	}
	if node != nil {
		switch n := node.(type) {
		case *ast.SequenceNode:
			for _, val := range n.Values {
				if val.String() == snode.String() {
					return true
				}
			}
		default:
			if node.String() == snode.String() {
				return true
			}
		}

	}
	return false
}
func SliceRemove(elem ast.Node, arr []ast.Node) []ast.Node {
	idx := slices.IndexFunc(arr, func(c ast.Node) bool { return c == elem })
	arr = append(arr[:idx], arr[idx+1:]...)
	return arr
}
