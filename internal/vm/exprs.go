package vm

import (
	"fmt"
	"strconv"
	"strings"
	"xml-programming/internal/ast"
	"xml-programming/internal/scope"
	"xml-programming/internal/values"
)

func evaluateAddExpression(expression ast.OperatorExpression, localScope *scope.Scope) values.Value {
	isFloat := false

	var args []values.Value

	for _, _arg := range expression.Exprs {
		arg := evaluateExpression(_arg, localScope)
		if arg.Type == ast.Float {
			isFloat = true
		}
		args = append(args, arg)
	}

	if isFloat {
		var sum float32 = 0
		for _, arg := range args {
			if arg.Type == ast.Int {
				sum += float32(arg.Int)
			} else {
				sum += arg.Float
			}
		}
		return values.Value{
			Type:  ast.Float,
			Float: sum,
		}
	} else {
		var sum int = 0
		for _, arg := range args {
			// all ints
			sum += arg.Int
		}
		return values.Value{
			Type: ast.Int,
			Int:  sum,
		}
	}
}

func evaluateSubExpression(expression ast.OperatorExpression, localScope *scope.Scope) values.Value {
	isFloat := false

	arg1 := evaluateExpression(expression.Exprs[0], localScope)
	if arg1.Type == ast.Float {
		isFloat = true
	}
	arg2 := evaluateExpression(expression.Exprs[1], localScope)
	if arg2.Type == ast.Float {
		isFloat = true
	}

	if isFloat {
		var result float32

		if arg1.Type == ast.Int {
			result = float32(arg1.Int)
		} else {
			result = arg1.Float
		}

		if arg2.Type == ast.Int {
			result -= float32(arg1.Int)
		} else {
			result -= arg2.Float
		}

		return values.Value{
			Type:  ast.Float,
			Float: result,
		}
	} else {
		var result int = arg1.Int
		result -= arg2.Int

		return values.Value{
			Type: ast.Int,
			Int:  result,
		}
	}
}

func evaluateMulExpression(expression ast.OperatorExpression, localScope *scope.Scope) values.Value {
	isFloat := false

	var args []values.Value

	for _, _arg := range expression.Exprs {
		arg := evaluateExpression(_arg, localScope)
		if arg.Type == ast.Float {
			isFloat = true
		}
		args = append(args, arg)
	}

	if isFloat {
		var product float32 = 1
		for _, arg := range args {
			if arg.Type == ast.Int {
				product *= float32(arg.Int)
			} else {
				product *= arg.Float
			}
		}
		return values.Value{
			Type:  ast.Float,
			Float: product,
		}
	} else {
		var product int = 1
		for _, arg := range args {
			// all ints
			product *= arg.Int
		}
		return values.Value{
			Type: ast.Int,
			Int:  product,
		}
	}
}

func evaluateDivExpression(expression ast.OperatorExpression, localScope *scope.Scope) values.Value {
	isFloat := false

	arg1 := evaluateExpression(expression.Exprs[0], localScope)
	if arg1.Type == ast.Float {
		isFloat = true
	}
	arg2 := evaluateExpression(expression.Exprs[1], localScope)
	if arg2.Type == ast.Float {
		isFloat = true
	}

	if isFloat {
		var result float32

		if arg1.Type == ast.Int {
			result = float32(arg1.Int)
		} else {
			result = arg1.Float
		}

		if arg2.Type == ast.Int {
			result /= float32(arg1.Int)
		} else {
			result /= arg2.Float
		}

		return values.Value{
			Type:  ast.Float,
			Float: result,
		}
	} else {
		var result int = arg1.Int
		result /= arg2.Int

		return values.Value{
			Type: ast.Int,
			Int:  result,
		}
	}
}

func evaluateModExpression(expression ast.OperatorExpression, localScope *scope.Scope) values.Value {
	arg1 := evaluateExpression(expression.Exprs[0], localScope)
	arg2 := evaluateExpression(expression.Exprs[1], localScope)

	var result int = arg1.Int
	result %= arg2.Int

	return values.Value{
		Type: ast.Int,
		Int:  result,
	}
}

func evaluateConcatExpression(expression ast.OperatorExpression, localScope *scope.Scope) values.Value {
	builder := strings.Builder{}

	for _, _arg := range expression.Exprs {
		arg := evaluateExpression(_arg, localScope)
		switch arg.Type {
		case ast.String:
			builder.WriteString(arg.String)
		case ast.Bool:
			builder.WriteString(strconv.FormatBool(arg.Bool))
		case ast.Int:
			builder.WriteString(strconv.FormatInt(int64(arg.Int), 10))
		case ast.Float:
			builder.WriteString(strconv.FormatFloat(float64(arg.Float), 'e', 4, 32))
		}
	}

	return values.Value{
		Type: ast.String,
		String: builder.String(),
	}
}

func evaluateEqualExpression(expression ast.OperatorExpression, localScope *scope.Scope) values.Value {
	arg1 := evaluateExpression(expression.Exprs[0], localScope)
	arg2 := evaluateExpression(expression.Exprs[1], localScope)

	if arg1.Type == ast.Int {
		if arg2.Type == ast.Float {
			return values.Value{
				Type: ast.Bool,
				Bool: float32(arg1.Int) == arg2.Float,
			}
		} else {
			return values.Value{
				Type: ast.Bool,
				Bool: arg1.Int == arg2.Int,
			}
		}
	} else {
		if arg2.Type == ast.Float {
			return values.Value{
				Type: ast.Bool,
				Bool: arg1.Float == arg2.Float,
			}
		} else {
			return values.Value{
				Type: ast.Bool,
				Bool: arg1.Float == float32(arg2.Int),
			}
		}
	}
}

func evaluateGreaterThanExpression(expression ast.OperatorExpression, localScope *scope.Scope) values.Value {
	arg1 := evaluateExpression(expression.Exprs[0], localScope)
	arg2 := evaluateExpression(expression.Exprs[1], localScope)

	if arg1.Type == ast.Int {
		if arg2.Type == ast.Float {
			return values.Value{
				Type: ast.Bool,
				Bool: float32(arg1.Int) > arg2.Float,
			}
		} else {
			return values.Value{
				Type: ast.Bool,
				Bool: arg1.Int > arg2.Int,
			}
		}
	} else {
		if arg2.Type == ast.Float {
			return values.Value{
				Type: ast.Bool,
				Bool: arg1.Float > arg2.Float,
			}
		} else {
			return values.Value{
				Type: ast.Bool,
				Bool: arg1.Float > float32(arg2.Int),
			}
		}
	}
}

func evaluateLessThanExpression(expression ast.OperatorExpression, localScope *scope.Scope) values.Value {
	arg1 := evaluateExpression(expression.Exprs[0], localScope)
	arg2 := evaluateExpression(expression.Exprs[1], localScope)

	if arg1.Type == ast.Int {
		if arg2.Type == ast.Float {
			return values.Value{
				Type: ast.Bool,
				Bool: float32(arg1.Int) < arg2.Float,
			}
		} else {
			return values.Value{
				Type: ast.Bool,
				Bool: arg1.Int < arg2.Int,
			}
		}
	} else {
		if arg2.Type == ast.Float {
			return values.Value{
				Type: ast.Bool,
				Bool: arg1.Float < arg2.Float,
			}
		} else {
			return values.Value{
				Type: ast.Bool,
				Bool: arg1.Float < float32(arg2.Int),
			}
		}
	}
}

func evaluateNotExpression(expression ast.OperatorExpression, localScope *scope.Scope) values.Value {
	arg := evaluateExpression(expression.Exprs[0], localScope)

	return values.Value{
		Type: ast.Bool,
		Bool: !arg.Bool,
	}
}

func evaluateAndExpression(expression ast.OperatorExpression, localScope *scope.Scope) values.Value {
	for _, expr := range expression.Exprs {
		arg := evaluateExpression(expr, localScope)
		if !arg.Bool {
			return values.Value{
				Type: ast.Bool,
				Bool: false,
			}
		}
	}

	return values.Value{
		Type: ast.Bool,
		Bool: true,
	}
}

func evaluateOrExpression(expression ast.OperatorExpression, localScope *scope.Scope) values.Value {
	for _, expr := range expression.Exprs {
		arg := evaluateExpression(expr, localScope)
		if arg.Bool {
			return values.Value{
				Type: ast.Bool,
				Bool: true,
			}
		}
	}

	return values.Value{
		Type: ast.Bool,
		Bool: false,
	}
}

func evaluateExpression(expression ast.Expression, localScope *scope.Scope) values.Value {
	switch v := expression.(type) {
	case ast.LiteralExpression:
		return values.FromLiteralExpression(v)
	case ast.VariableExpression:
		return localScope.GetVariable(v.Name).Value
	case ast.OperatorExpression:
		switch v.Operator {
		case ast.Add:
			return evaluateAddExpression(v, localScope)
		case ast.Sub:
			return evaluateSubExpression(v, localScope)
		case ast.Mul:
			return evaluateMulExpression(v, localScope)
		case ast.Div:
			return evaluateDivExpression(v, localScope)
		case ast.Mod:
			return evaluateModExpression(v, localScope)
		case ast.Concat:
			return evaluateConcatExpression(v, localScope)
		case ast.Equal:
			return evaluateEqualExpression(v, localScope)
		case ast.GreaterThan:
			return evaluateGreaterThanExpression(v, localScope)
		case ast.LessThan:
			return evaluateLessThanExpression(v, localScope)
		case ast.Not:
			return evaluateNotExpression(v, localScope)
		case ast.And:
			return evaluateAndExpression(v, localScope)
		case ast.Or:
			return evaluateOrExpression(v, localScope)
		default:
			fmt.Println(v)
			panic("not yet implemented")
		}
	case ast.FunctionCall:
		var args []values.Value
		for _, arg := range v.Args {
			args = append(args, evaluateExpression(arg, localScope))
		}

		return callFunction(v.Name, localScope, args)
	default:
		fmt.Println(expression)
		panic("not yet implemented")
	}
}
