package parser

import (
	"encoding/xml"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"xml-programming/internal/ast"
)

func Parse(content []byte) (*ast.Program, error) {
	var programElement ProgramElement
	err := xml.Unmarshal(content, &programElement)
	if err != nil {
		return nil, err
	}

	//json, _ := json.MarshalIndent(programElement, "", "  ")
	//fmt.Println(string(json))

	var program ast.Program
	program.Statements, err = ParseStatements(programElement.Statements)
	if err != nil {
		return nil, err
	}

	return &program, nil
}

func ParseStatements(statementsElements []StatementElement) ([]ast.Statement, error) {
	var statements []ast.Statement
	for _, element := range statementsElements {
		statement, err := ParseStatement(element)
		if err != nil {
			return nil, err
		}
		statements = append(statements, statement)
	}

	return statements, nil
}

func ParseType(str string) (ast.Type, error) {
	switch str {
	case "string":
		return ast.String, nil
	case "bool":
		return ast.Bool, nil
	case "int":
		return ast.Int, nil
	case "float":
		return ast.Float, nil
	default:
		return ast.Void, fmt.Errorf("unknown type: %v", str)
	}
}

func ParseStatement(statement StatementElement) (ast.Statement, error) {
	switch statement.XMLName.Local {
	case OutputStatementElementName:
		var exprs []ast.Expression
		for _, exprElement := range statement.Exprs {
			expr, err := ParseExpression(exprElement)
			if err != nil {
				return nil, err
			}
			exprs = append(exprs, expr)
		}
		return ast.OutputStatement{
			Exprs: exprs,
		}, nil
	case VariableDeclarationElementName:
		t, err := ParseType(statement.Type)
		if err != nil {
			return nil, err
		}
		return ast.VariableDeclarationStatement{
			Name: statement.Name,
			Type: t,
		}, nil
	case VariableAssignmentElementName:
		if len(statement.Exprs) != 1 {
			return nil, errors.New("variable assignments must have exactly one expression")
		}
		expr, err := ParseExpression(statement.Exprs[0])
		if err != nil {
			return nil, err
		}
		return ast.VariableAssignmentStatement{
			Name: statement.Name,
			Expr: expr,
		}, nil
	case FunctionElementName:
		var err error
		var args []ast.FunctionArg
		for _, argElement := range statement.Args.Args {
			arg := ast.FunctionArg{
				Name: argElement.Name,
			}
			arg.Type, err = ParseType(argElement.Type)
			if err != nil {
				return nil, fmt.Errorf("couldn't parse function argument list: %w", err)
			}
			args = append(args, arg)
		}
		returns, err := ParseType(statement.Args.Returns.Type)
		if  err != nil {
			return nil, fmt.Errorf("couldn't parse function return type: %w", err)
		}
		var body []ast.Statement
		for _, statementElement := range statement.Body.Statements {
			bodyStatement, err := ParseStatement(statementElement)
			if err != nil {
				return nil, err
			}
			body = append(body, bodyStatement)
		}
		return ast.FunctionStatement{
			Name:    statement.Name,
			Returns: returns,
			Args:    args,
			Body:    body,
		}, nil
	case FunctionReturnElementName:
		if len(statement.Exprs) != 1 {
			return nil, errors.New("a return statement must have exactly one expression")
		}
		expr, err := ParseExpression(statement.Exprs[0])
		if err != nil {
			return nil, err
		}
		return ast.FunctionReturnStatement{
			Expr: expr,
		}, nil
	case FunctionCallStatementElementName:
		exprs, err := ParseExpressionList(statement.Exprs)
		if err != nil {
			return nil, err
		}
		return ast.FunctionCall{
			Name: statement.Name,
			Args: exprs,
		}, nil
	case ConditionStatementElementName:
		conditional := ast.ConditionalStatement{}

		for _, _if := range statement.Ifs {
			then, err := ParseStatements(_if.Then.Statements)
			if err != nil {
				return nil, err
			}
			expr, err := ParseExpression(_if.Condition.Expr)
			if err != nil {
				return nil, err
			}
			conditional.Ifs = append(conditional.Ifs, ast.ConditionIf{
				Expr: expr,
				Then: then,
			})
		}

		if statement.Else != nil {
			then, err := ParseStatements(statement.Else.Then.Statements)
			if err != nil {
				return nil, err
			}
			conditional.Else = then
		}

		return conditional, nil
	case LoopStatementElementName:
		expr, err := ParseExpression(statement.LoopStatementElement.Condition.Expr)
		if err != nil {
			return nil, err
		}
		body, err := ParseStatements(statement.Body.Statements)
		if err != nil {
			return nil, err
		}
		return ast.LoopStatement{
			LoopCondition: expr,
			Body:          body,
		}, nil
	case ForStatementElementName:
		body, err := ParseStatements(statement.Body.Statements)
		if err != nil {
			return nil, err
		}
		return ast.ForStatement{
			From: statement.From,
			To: statement.To,
			Name: statement.Name,
			Body:          body,
		}, nil
	default:
		return nil, fmt.Errorf("unknown statement: <%v>", statement.XMLName.Local)
	}
}

func ParseExpressionList(exprList []ExpressionElement) ([]ast.Expression, error) {
	exprs := make([]ast.Expression, len(exprList))

	var err error
	for i, expr := range exprList {
		exprs[i], err = ParseExpression(expr)
		if err != nil {
			return nil, fmt.Errorf("unable to parse sub expression: %w", err)
		}
	}

	return exprs, nil
}

func ParseOperatorExpression(exprElement ExpressionElement, operator ast.Operator) (ast.Expression, error) {
	exprs, err := ParseExpressionList(exprElement.Exprs)
	if err != nil {
		return nil, err
	}

	return ast.OperatorExpression{
		Operator: operator,
		Exprs:    exprs,
	}, nil
}

func ParseExpression(exprElement ExpressionElement) (ast.Expression, error) {
	switch exprElement.XMLName.Local {
	case LiteralExpressionStringElementName:
		return ast.LiteralExpression{
			Type:   ast.String,
			String: exprElement.Content,
		}, nil
	case LiteralExpressionBoolElementName:
		val, err := strconv.ParseBool(strings.TrimSpace(exprElement.Content))
		if err != nil {
			return nil, fmt.Errorf("unable to parse int value: %w", err)
		}
		return ast.LiteralExpression{
			Type:   ast.Bool,
			Bool: val,
		}, nil
	case LiteralExpressionIntElementName:
		val, err := strconv.Atoi(strings.TrimSpace(exprElement.Content))
		if err != nil {
			return nil, fmt.Errorf("unable to parse int value: %w", err)
		}
		return ast.LiteralExpression{
			Type:   ast.Int,
			Int: val,
		}, nil
	case LiteralExpressionFloatElementName:
		val, err := strconv.ParseFloat(strings.TrimSpace(exprElement.Content), 32)
		if err != nil {
			return nil, fmt.Errorf("unable to parse float value: %w", err)
		}
		return ast.LiteralExpression{
			Type:   ast.Float,
			Float: float32(val),
		}, nil
	case VariableExpressionElementName:
		return ast.VariableExpression{
			Name: exprElement.Name,
		}, nil
	case OperatorExpressionAddElementName:
		return ParseOperatorExpression(exprElement, ast.Add)
	case OperatorExpressionSubElementName:
		return ParseOperatorExpression(exprElement, ast.Sub)
	case OperatorExpressionMulElementName:
		return ParseOperatorExpression(exprElement, ast.Mul)
	case OperatorExpressionDivElementName:
		return ParseOperatorExpression(exprElement, ast.Div)
	case OperatorExpressionModElementName:
		return ParseOperatorExpression(exprElement, ast.Mod)
	case OperatorExpressionConcatElementName:
		return ParseOperatorExpression(exprElement, ast.Concat)
	case OperatorExpressionEqualElementName:
		return ParseOperatorExpression(exprElement, ast.Equal)
	case OperatorExpressionGreaterThanElementName:
		return ParseOperatorExpression(exprElement, ast.GreaterThan)
	case OperatorExpressionLessThanElementName:
		return ParseOperatorExpression(exprElement, ast.LessThan)
	case OperatorExpressionNotElementName:
		return ParseOperatorExpression(exprElement, ast.Not)
	case OperatorExpressionAndElementName:
		return ParseOperatorExpression(exprElement, ast.And)
	case OperatorExpressionOrElementName:
		return ParseOperatorExpression(exprElement, ast.Or)
	case FunctionCallExpressionElementName:
		exprs, err := ParseExpressionList(exprElement.Exprs)
		if err != nil {
			return nil, err
		}
		return ast.FunctionCall{
			Name: exprElement.Name,
			Args: exprs,
		}, nil
	default:
		return nil, errors.New("unknown expression type")
	}
}