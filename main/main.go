package main

import (
	// "go/ast"
	// "go/token"
	"fmt"
	"log"

	hhelper "github.com/cjkao/helm-helper"
	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/lexer"
	"github.com/goccy/go-yaml/parser"
)

// flags,
// read yaml, remove duplicate block
// delete a line

func getNode() ast.Node {
	src2 := `
Block style: 
  a : a1
Flow style: [ Mercury, Venus, Earth, Mars]
`
	tokens2 := lexer.Tokenize(src2)
	f2, _ := parser.Parse(tokens2, 0)
	return f2.Docs[0].Body
}

var srcNode ast.Node

// fix me : flow style + comment in flow will break document if print with comment
func main() {

	src := `
Block style: 
  a : a1
Flow style: [ Mercury, Venus, Earth, Mars,      
              Jupiter, Saturn, Uranus, Neptune, 
              Pluto ]                           
`
	srcNode = getNode()
	tokens := lexer.Tokenize(src)
	f, err := parser.Parse(tokens, parser.ParseComments)
	// f, err := parser.Parse(tokens, 0)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	print(f.Docs[0].String())
	v := Visitor{}
	for _, doc := range f.Docs {
		Walk(&v, doc.Body)
		printYaml(doc)
	}
}
func printYaml(v ast.Node) {
	bytes, err := yaml.Marshal(v)
	if err != nil {
		panic(err)
	}
	fmt.Println("\n--------\n" + string(bytes))
}
func IsNullNode(node ast.Node) bool {
	switch n := node.(type) {
	case *ast.MappingValueNode:
		if n.Value == nil {
			return true
		}
	}
	return false
}
func WalkToNoNil(node ast.Node) {
	switch n := node.(type) {

	case *ast.MappingNode: // ":"
		nv := n.Values
		for i := len(n.Values) - 1; i >= 0; i-- {
			if IsNullNode(n.Values[i]) {
				nv = hhelper.SliceMappingRemove(i, nv)
			}
		}
		for _, value := range n.Values {
			WalkToNoNil(value)
		}
		// case *ast.MappingValueNode:
		// 	if hhelper.IsLeaf(n.Value) && hhelper.IsExistInGolden(n.Value, srcNode) {
		// 		fmt.Printf("_+_")
		// 		n.Key = nil
		// 		n.Value = nil
		// 	} else {
		// 		Walk(v, n.Key)
		// 		Walk(v, n.Value)
		// 	}
		// case *ast.SequenceNode:
		// 	nv := n.Values
		// 	for i := len(n.Values) - 1; i >= 0; i-- {
		// 		println(n.Values[i].String())
		// 		if hhelper.IsLeaf(n.Values[i]) && hhelper.IsExistInGolden(n.Values[i], srcNode) {
		// 			fmt.Printf("+")
		// 			nv = hhelper.SliceRemove(n.Values[i], nv)
		// 		}
		// 		// Walk(v, value)
		// 	}
		// 	n.Values = nv
		// 	for _, value := range nv {
		// 		if hhelper.IsLeaf(value) { //skip leaf
		// 			continue
		// 		}
		// 		Walk(v, value) //drill down other types
		// 	}
		// case *ast.AnchorNode:
		// 	walkComment(v, n.BaseNode)
		// 	Walk(v, n.Name)
		// 	Walk(v, n.Value)
		// case *ast.AliasNode:
		// 	walkComment(v, n.BaseNode)
		// 	Walk(v, n.Value)
	}
}

type Visitor struct {
}

func (v *Visitor) Visit(node ast.Node) ast.Visitor {
	fmt.Printf("level %s > %s > %v\n", node.GetPath(), node.GetToken().Type, node.GetToken().Value)
	return v
}
func Walk(v ast.Visitor, node ast.Node) {
	if v = v.Visit(node); v == nil {
		return
	}

	switch n := node.(type) {
	case *ast.CommentNode:
	case *ast.NullNode:
		walkComment(v, n.BaseNode)
	case *ast.IntegerNode:
		walkComment(v, n.BaseNode)
	case *ast.FloatNode:
		walkComment(v, n.BaseNode)
	case *ast.StringNode:
		walkComment(v, n.BaseNode)
	case *ast.MergeKeyNode:
		walkComment(v, n.BaseNode)
	case *ast.BoolNode:
		walkComment(v, n.BaseNode)
	case *ast.InfinityNode:
		walkComment(v, n.BaseNode)
	case *ast.NanNode:
		walkComment(v, n.BaseNode)
	case *ast.LiteralNode:
		walkComment(v, n.BaseNode)
		Walk(v, n.Value)
	case *ast.DirectiveNode:
		walkComment(v, n.BaseNode)
		Walk(v, n.Value)
	case *ast.TagNode:
		walkComment(v, n.BaseNode)
		Walk(v, n.Value)
	case *ast.DocumentNode:
		walkComment(v, n.BaseNode)
		Walk(v, n.Body)
	case *ast.MappingNode: // ":"
		walkComment(v, n.BaseNode)
		for _, value := range n.Values {
			Walk(v, value)
		}
		nv := n.Values
		for i := len(n.Values) - 1; i >= 0; i-- {
			if IsNullNode(n.Values[i]) {
				nv = hhelper.SliceMappingRemove(i, nv)
			} else {
				Walk(v, n.Values[i])
			}
		}
		n.Values = nv

	case *ast.MappingKeyNode:
		walkComment(v, n.BaseNode)
		Walk(v, n.Value)
	case *ast.MappingValueNode:
		walkComment(v, n.BaseNode)
		if hhelper.IsLeaf(n.Value) && hhelper.IsExistInGolden(n.Value, srcNode) {
			fmt.Printf("_+_")
			n.Key = nil
			n.Value = nil
			// n = nil
		} else {

			Walk(v, n.Key)
			Walk(v, n.Value)
		}
	case *ast.SequenceNode:
		walkComment(v, n.BaseNode)
		nv := n.Values
		for i := len(n.Values) - 1; i >= 0; i-- {
			println(n.Values[i].String())
			if hhelper.IsLeaf(n.Values[i]) && hhelper.IsExistInGolden(n.Values[i], srcNode) {
				fmt.Printf("+")
				nv = hhelper.SliceRemove(n.Values[i], nv)
			}
			// Walk(v, value)
		}
		n.Values = nv
		for _, value := range nv {
			if hhelper.IsLeaf(value) { //skip leaf
				continue
			}
			Walk(v, value) //drill down other types
		}
	case *ast.AnchorNode:
		walkComment(v, n.BaseNode)
		Walk(v, n.Name)
		Walk(v, n.Value)
	case *ast.AliasNode:
		walkComment(v, n.BaseNode)
		Walk(v, n.Value)
	}
}

func walkComment(v ast.Visitor, base *ast.BaseNode) {

	if base == nil {
		return
	}
	if base.Comment == nil {
		return
	}
	Walk(v, base.Comment)
}
