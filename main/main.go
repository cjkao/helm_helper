package main

import (
	// "go/ast"
	// "go/token"
	"fmt"
	"log"

	_ "github.com/cjkao/helmHelper"
	hhelper "github.com/cjkao/helmHelper"
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
- Pluto     # You call this a planet?
- a
- c
`
	tokens2 := lexer.Tokenize(src2)
	f2, _ := parser.Parse(tokens2, 0)
	return f2.Docs[0].Body
}

var srcNode ast.Node

func main() {

	src := `
# Ordered sequence of nodes in YAML STRUCTURE
Block style: 
- Pluto     # You call this a planet?
- Mercury   # Rotates - no light/dark sides.
- Venus     # Deadliest. Aptly named.
- Earth     # Mostly dirt.
- Mars      # Seems empty.
- Jupiter   # The king.
- Saturn    # Pretty.
- Uranus    # Where the sun hardly shines.
- Neptune   # Boring. No rings.
Flow style: [ Mercury, Venus, Earth, Mars,      # Rocks
              Jupiter, Saturn, Uranus, Neptune, # Gas
              Pluto ]                           # Overrated
`

	srcNode = getNode()
	tokens := lexer.Tokenize(src)
	f, err := parser.Parse(tokens, 0)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	// print(node)
	var v Visitor
	for _, doc := range f.Docs {
		Walk(&v, doc.Body)
		printYaml(doc.Body)
	}
}
func printYaml(v ast.Node) {
	bytes, err := yaml.Marshal(v)
	if err != nil {
		panic(err)
	}
	fmt.Println("\n--------\n" + string(bytes))
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
	case *ast.MappingNode:
		walkComment(v, n.BaseNode)
		for _, value := range n.Values {
			Walk(v, value)
		}
	case *ast.MappingKeyNode:
		walkComment(v, n.BaseNode)
		Walk(v, n.Value)
	case *ast.MappingValueNode:
		walkComment(v, n.BaseNode)
		Walk(v, n.Key)
		Walk(v, n.Value)
	case *ast.SequenceNode:
		walkComment(v, n.BaseNode)
		nv := n.Values
		for i := len(n.Values) - 1; i >= 0; i-- {

			if hhelper.IsLeaf(n.Values[i]) && hhelper.IsExistInGolden(n.Values[i], srcNode) {
				fmt.Printf("current value is leaf")
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
