package values

import (
	"fmt"
	"xml-programming/internal/ast"
)

type Value struct {
	Type ast.Type

	String string
	Bool bool
	Int int
	Float float32
}

func FromLiteralExpression(expression ast.LiteralExpression) Value {
	value := Value{
		Type: expression.Type,
	}

	switch expression.Type {
	case ast.String:
		value.String = expression.String
	case ast.Int:
		value.Int = expression.Int
	case ast.Bool:
		value.Bool = expression.Bool
	case ast.Float:
		value.Float = expression.Float
	default:
		fmt.Println("unknown type: ", expression.Type)
	}

	return value
}