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
american:
  - Boston Red Sox
  - Detroit Tigers
  - New York Yankees
national:
  - New York Mets
  - Chicago Cubs
  - Atlanta Braves
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
	// expect := fmt.Sprintf("\n%+v\n", f)
	// if test.expect != expect {
	// tokens.Dump()
	// t.Fatalf("unexpected output: [%s] != [%s]", test.expect, expect)
	// }
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
	tk := node.GetToken()
	fmt.Printf("level %d %s %v\n", node.GetToken().Position.IndentLevel, node.GetToken().Type, node.GetToken().Value)
	tk.Prev = nil
	tk.Next = nil
	return v
}
