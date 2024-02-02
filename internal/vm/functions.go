package vm

import (
	"xml-programming/internal/ast"
	"xml-programming/internal/scope"
	"xml-programming/internal/values"
)

func callFunction(name string, localScope *scope.Scope, args []values.Value) values.Value {
	function := localScope.GetFunction(name)
	functionScope := scope.FromParent(localScope)

	for i, arg := range args {
		functionScope.AddVariable(scope.Variable{
			Name:  function.Args[i].Name,
			Value: arg,
		})
	}

	result := executeStatements(function.Body, functionScope)
	if result != nil {
		return *result
	} else {
		return values.Value{
			Type: ast.Void,
		}
	}
}