package main

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/neoxelox/yay/mod"
)

func transpile(program []mod.Token, filepath string, comments bool) (string, error) {
	var transpilation bytes.Buffer
	template := template.Must(template.New("").Parse(mod.TranspileBase))

	var imports []string
	var definitions []string
	var statements []string
	importSet := make(map[string]struct{})

	err := beginTranspile()
	if err != nil {
		return "", err
	}

	for _, token := range program {
		switch token.Type {
		case mod.TypeNumber:
			statements = append(statements, fmt.Sprintf(`push(int64(%s))`, token.Literal))
		case mod.TypeString:
			statements = append(statements, fmt.Sprintf(`push("%s")`, token.Literal))
		case mod.TypeComment:
			if comments {
				statements = append(statements, fmt.Sprintf(`//%s%c`, token.Literal, '\n'))
			}
		case mod.TypeIdentifier:
			iImports, iDefinitions, iStatements, err := Identifiers[token.Literal].Transpile(token)
			if err != nil {
				return "", fmt.Errorf("cannot transpile identifier '%s' at %s:%d:%d", token.Literal, token.File, token.Row+1, token.Col+1)
			}
			for _, iImport := range iImports {
				importSet[iImport] = struct{}{}
			}
			if len(iDefinitions) > 0 {
				definitions = append(definitions, iDefinitions)
			}
			if len(iStatements) > 0 {
				statements = append(statements, iStatements)
			}
		default:
			return "", fmt.Errorf("unknown token type '%s'", token.Type)
		}
	}

	for imprt := range importSet {
		imports = append(imports, imprt)
	}

	err = endTranspile()
	if err != nil {
		return "", err
	}

	err = template.Execute(&transpilation, mod.TranspileData{
		Version:     mod.Version,
		Imports:     imports,
		Definitions: definitions,
		Statements:  statements,
	})

	if err != nil {
		return "", fmt.Errorf("cannot transpile file '%s': %w", filepath, err)
	}

	return transpilation.String(), nil
}
