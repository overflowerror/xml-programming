package vm

import (
	"xml-programming/internal/ast"
	"xml-programming/internal/scope"
)

func Run(program *ast.Program) {
	localScope := scope.New()
	_ = executeStatements(program.Statements, localScope)
}
