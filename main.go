package main

import (
	// "go/ast"
	// "go/token"
	"fmt"
	"log"

	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/lexer"
	"github.com/goccy/go-yaml/parser"
	"github.com/goccy/go-yaml/token"
)

//flags,
//read yaml, remove duplicate block

func main() {

	src := `
defaults: &defaults
   adapter:  postgres
   host:     localhost
`
	tokens := lexer.Tokenize(src)
	f, err := parser.Parse(tokens, 0)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	var v Visitor
	for _, doc := range f.Docs {
		ast.Walk(&v, doc.Body)
	}

}

type pathCapturer struct {
	capturedNum   int
	orderedPaths  []string
	orderedTypes  []ast.NodeType
	orderedTokens []*token.Token
}

func (c *pathCapturer) Visit(node ast.Node) ast.Visitor {
	c.capturedNum++
	c.orderedPaths = append(c.orderedPaths, node.GetPath())
	c.orderedTypes = append(c.orderedTypes, node.Type())
	c.orderedTokens = append(c.orderedTokens, node.GetToken())
	return c
}

type Visitor struct {
}

func (v *Visitor) Visit(node ast.Node) ast.Visitor {
	fmt.Printf("level %d %s %v\n", node.GetToken().Position.IndentLevel, node.GetToken().Type, node.GetToken().Value)
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
		for _, value := range n.Values {
			Walk(v, value)
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
