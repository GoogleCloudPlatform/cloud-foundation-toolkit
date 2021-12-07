package bptest

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

// getTestFuncsFromFile parses a go source file and returns slice of test function names
func getTestFuncsFromFile(filePath string) ([]string, error) {
	fileSet := token.NewFileSet()
	f, err := parser.ParseFile(fileSet, filePath, nil, parser.AllErrors)
	if err != nil {
		return nil, err
	}
	testFuncs := make([]string, 0)
	for _, decl := range f.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		// not a function declaration
		if !ok {
			continue
		}
		if strings.HasPrefix(funcDecl.Name.Name, "Test") {
			testFuncs = append(testFuncs, funcDecl.Name.Name)
		}
	}
	return testFuncs, nil
}
