package analysis

import (
	"fmt"
	"xml-programming/internal/ast"
	"xml-programming/internal/scope"
	"xml-programming/internal/values"
)

func analyseArithmetic(types []ast.Type) (ast.Type, error) {
	isFloat := false
	for _, _type := range types {
		if !_type.IsNumber() {
			return ast.Void, fmt.Errorf("can't do arithmetic on non-number types")
		}
		if _type == ast.Float {
			isFloat = true
		}
	}
	if isFloat {
		return ast.Float, nil
	} else {
		return ast.Int, nil
	}
}

func analyseComparison(types []ast.Type, operator ast.Operator) (ast.Type, error) {
	if operator != ast.Equal {
		for _, _type := range types {
			if !_type.IsNumber() {
				return ast.Void, fmt.Errorf("can not compare non-number types")
			}
		}
	}
	return ast.Bool, nil
}

func analyseLogic(types []ast.Type) (ast.Type, error) {
	for _, _type := range types {
		if _type != ast.Bool {
			return ast.Void, fmt.Errorf("can only do logic on bools")
		}
	}
	return ast.Bool, nil
}

func analyseExpression(expression ast.Expression, localScope *scope.Scope) (ast.Type, error) {
	switch v := expression.(type) {
	case ast.LiteralExpression:
		return v.Type, nil
	case ast.VariableExpression:
		variable := localScope.GetVariable(v.Name)
		if variable == nil {
			return ast.Void, fmt.Errorf("name %s not found in local scope", v.Name)
		}
		return variable.Type, nil
	case ast.FunctionCall:
		function := localScope.GetFunction(v.Name)
		if function == nil {
			return ast.Void, fmt.Errorf("name %s not found in local scope", v.Name)
		}
		if len(function.Args) != len(v.Args) {
			return ast.Void, fmt.Errorf("mismatched number of arguments for name %s", v.Name)
		}
		for i, _ := range function.Args {
			expectedType := function.Args[i]
			actualType, err := analyseExpression(v.Args[i], localScope)
			if err != nil {
				return ast.Void, err
			}
			if expectedType.Type != actualType {
				return ast.Void, fmt.Errorf("mismatched types for argument %d of name %s", i + 1, v.Name)
			}
		}

		return function.Return, nil
	case ast.OperatorExpression:
		types, err := analyseExpressionist(v.Exprs, localScope)
		if err != nil {
			return ast.Void, err
		}

		if len(types) < 1 {
			return ast.Void, fmt.Errorf("operator with no arguments is not valid")
		}

		switch v.Operator {
		case ast.Add:
			return analyseArithmetic(types)
		case ast.Sub:
			if len(types) != 2 {
				return ast.Void, fmt.Errorf("sub is binary")
			}
			return analyseArithmetic(types)
		case ast.Mul:
			return analyseArithmetic(types)
		case ast.Div:
			if len(types) != 2 {
				return ast.Void, fmt.Errorf("div is binary")
			}
			return analyseArithmetic(types)
		case ast.Mod:
			if len(types) != 2 {
				return ast.Void, fmt.Errorf("mod is binary")
			}
			if types[0] == ast.Float || types[1] == ast.Float {
				return ast.Void, fmt.Errorf("mod only supports int")
			}
			return ast.Int, nil
		case ast.Concat:
			return ast.String, nil
		case ast.Equal:
			if len(types) != 2 {
				return ast.Void, fmt.Errorf("equal is binary")
			}
			return analyseComparison(types, v.Operator)
		case ast.LessThan:
			if len(types) != 2 {
				return ast.Void, fmt.Errorf("less than is binary")
			}
			return analyseComparison(types, v.Operator)
		case ast.GreaterThan:
			if len(types) != 2 {
				return ast.Void, fmt.Errorf("greater than is binary")
			}
			return analyseComparison(types, v.Operator)
		case ast.Not:
			if len(types) != 1 {
				return ast.Void, fmt.Errorf("not is unary")
			}
			return analyseLogic(types)
		case ast.And:
			return analyseLogic(types)
		case ast.Or:
			return analyseLogic(types)
		default:
			return ast.Void, fmt.Errorf("not yet implemented")
		}
	default:
		return ast.Void, fmt.Errorf("not yet implemented")
	}
}

func analyseExpressionist(expressions []ast.Expression, localScope *scope.Scope) ([]ast.Type, error) {
	var types []ast.Type
	for _, expr := range expressions {
		_type, err := analyseExpression(expr, localScope)
		if err != nil {
			return nil, err
		}
		types = append(types, _type)
	}
	return types, nil
}

func analyseStatement(statement ast.Statement, localScope *scope.Scope, currentFunction *scope.Function) error {
	switch v := statement.(type) {
	case ast.OutputStatement:
		_, err := analyseExpressionist(v.Exprs, localScope)
		if err != nil {
			return err
		}
	case ast.VariableDeclarationStatement:
		if localScope.CurrentScopeHas(v.Name) {
			return fmt.Errorf("name %s already exists in local scope", v.Name)
		}
		localScope.AddVariable(scope.Variable{
			Name: v.Name,
			Value: values.Value{
				Type: v.Type,
			},
		})
	case ast.VariableAssignmentStatement:
		variable := localScope.GetVariable(v.Name)
		if variable == nil {
			return fmt.Errorf("name %s not found in local scope", v.Name)
		}
		_type, err := analyseExpression(v.Expr, localScope)
		if err != nil {
			return err
		}
		if _type != variable.Type {
			return fmt.Errorf("type mismatch on assignment")
		}
	case ast.FunctionStatement:
		if localScope.CurrentScopeHas(v.Name) {
			return fmt.Errorf("name %s already exists in local scope", v.Name)
		}
		function := scope.Function{
			Name:   v.Name,
			Args:   v.Args,
			Return: v.Returns,
		}
		localScope.AddFunction(function)

		functionScope := scope.FromParent(localScope)
		for _, arg := range v.Args {
			functionScope.AddVariable(scope.Variable{
				Name: arg.Name,
				Value: values.Value{
					Type: arg.Type,
				},
			})
		}

		err := analyseStatements(v.Body, functionScope, &function)
		if err != nil {
			return err
		}
	case ast.FunctionReturnStatement:
		if currentFunction == nil {
			return fmt.Errorf("can not return outside of function")
		}

		_type, err := analyseExpression(v.Expr, localScope)
		if err != nil {
			return err
		}
		if _type != currentFunction.Return {
			return fmt.Errorf("return type mismatch in name %v", currentFunction.Name)
		}
	case ast.ConditionalStatement:
		for _, _if := range v.Ifs {
			_type, err := analyseExpression(_if.Expr, localScope)
			if err != nil {
				return err
			}
			if _type != ast.Bool {
				return fmt.Errorf("condition type has to be bool")
			}

			err = analyseStatements(_if.Then, localScope, currentFunction)
			if err != nil {
				return err
			}
		}
		if v.Else != nil {
			err := analyseStatements(v.Else, localScope, currentFunction)
			if err != nil {
				return err
			}
		}
	case ast.LoopStatement:
		_type, err := analyseExpression(v.LoopCondition, localScope)
		if err != nil {
			return err
		}
		if _type != ast.Bool {
			return fmt.Errorf("condition type has to be bool")
		}

		err = analyseStatements(v.Body, localScope, currentFunction)
		if err != nil {
			return err
		}
	case ast.ForStatement:
		if localScope.CurrentScopeHas(v.Name) {
			return fmt.Errorf("name %s already in local scope", v.Name)
		}
		localScope.AddVariable(scope.Variable{
			Name: v.Name,
			Value: values.Value{
				Type: ast.Int,
			},
		})

		err := analyseStatements(v.Body, localScope, currentFunction)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("not yet implemented")
	}

	return nil
}

func analyseStatements(statements []ast.Statement, scope *scope.Scope, currentFunction *scope.Function) error {
	for _, statement := range statements {
		err := analyseStatement(statement, scope, currentFunction)
		if err != nil {
			return err
		}
	}

	return nil
}

func StaticAnalysis(program *ast.Program) error {
	scope := scope.New()

	return analyseStatements(program.Statements, scope, nil)
}