package mod

const Version = "v0.0.1"

const LiteralComment = "#"

type Type string

const (
	TypeNumber     = Type("NUMBER")
	TypeIdentifier = Type("IDENTIFIER")
	TypeComment    = Type("COMMENT")
)

type Token struct {
	Type    Type
	Literal string
	File    string
	Row     int
	Col     int
	Meta    interface{}
}

type Identifier interface {
	Parse(literal string, file string, row int, col int) (Token, error)
	Transpile(token Token) ([]string, string, string, error)
}

type TranspileData struct {
	Version     string
	Imports     []string
	Definitions []string
	Statements  []string
}

const TranspileBase = `// Generated automatically by YAY {{.Version}}
package main

import (
	{{- range .Imports}}
	"{{. -}}"
	{{- end}}
)

func unused(variable interface{}) {}

func push(stack *[]interface{}, value interface{}) {
	*stack = append(*stack, value)
}

func pop(stack *[]interface{}) interface{} {
	value := (*stack)[len(*stack)-1]
	*stack = (*stack)[:len(*stack)-1]
	return value
}

{{range .Definitions}}
{{. -}}
{{- end}}

func main() {
	var stack []interface{}
	unused(stack)
	var a interface{}
	unused(a)
	var b interface{}
	unused(b)
	var aInt int
	unused(aInt)
	var bInt int
	unused(bInt)

	{{range .Statements}}
	{{. -}}
	{{- end -}}
}
`