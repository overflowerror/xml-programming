package vm

import (
	"fmt"
	"xml-programming/internal/ast"
	"xml-programming/internal/scope"
	"xml-programming/internal/values"
)

func printValue(value values.Value) {
	switch value.Type {
	case ast.String:
		fmt.Print(value.String)
	case ast.Int:
		fmt.Print(value.Int)
	case ast.Float:
		fmt.Print(value.Float)
	case ast.Bool:
		fmt.Print(value.Bool)
	}
}

func executeStatement(statement ast.Statement, localScope *scope.Scope) *values.Value {
	switch v := statement.(type) {
	case ast.OutputStatement:
		for _, _arg := range v.Exprs {
			arg := evaluateExpression(_arg, localScope)
			printValue(arg)
		}
		fmt.Println()
	case ast.VariableDeclarationStatement:
		localScope.AddVariable(scope.Variable{
			Name:  v.Name,
			Value: values.Value{
				Type: v.Type,
			},
		})
	case ast.VariableAssignmentStatement:
		arg := evaluateExpression(v.Expr, localScope)
		variable := localScope.GetVariable(v.Name)
		variable.Value = arg
	case ast.FunctionStatement:
		localScope.AddFunction(scope.Function{
			Name:   v.Name,
			Args:   v.Args,
			Return: v.Returns,
			Body:   v.Body,
		})
	case ast.FunctionReturnStatement:
		arg := evaluateExpression(v.Expr, localScope)
		return &arg
	case ast.FunctionCall:
		var args []values.Value
		for _, arg := range v.Args {
			args = append(args, evaluateExpression(arg, localScope))
		}

		_ = callFunction(v.Name, localScope, args)
	case ast.ConditionalStatement:
		for _, _if := range v.Ifs {
			arg := evaluateExpression(_if.Expr, localScope)
			if arg.Bool {
				return executeStatements(_if.Then, localScope)
			}
		}
		return executeStatements(v.Else, localScope)
	case ast.LoopStatement:
		for {
			arg := evaluateExpression(v.LoopCondition, localScope)
			if !arg.Bool {
				break
			}

			result := executeStatements(v.Body, localScope)
			if result != nil {
				return result
			}
		}
	case ast.ForStatement:
		for i := v.From; i < v.To; i++ {
			forScope := scope.FromParent(localScope)
			forScope.AddVariable(scope.Variable{
				Name:  v.Name,
				Value: values.Value{
					Type: ast.Int,
					Int: i,
				},
			})
			result := executeStatements(v.Body, forScope)
			if result != nil {
				return result
			}
		}
	default:
		fmt.Println(v)
		panic("not yet implemented")
	}

	return nil
}

func executeStatements(statements []ast.Statement, localScope *scope.Scope) *values.Value {
	for _, statement := range statements {
		returnValue := executeStatement(statement, localScope)
		if returnValue != nil {
			return returnValue
		}
	}

	return nil
}