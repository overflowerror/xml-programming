package parser

import "encoding/xml"

type ProgramElement struct {
	XMLName    xml.Name          `xml:"program"`
	Statements []StatementElement `xml:",any"`
}

type StatementElement struct {
	XMLName xml.Name

	OutputStatementElement
	VariableDeclarationElement
	VariableAssignmentElement
	FunctionElement
	ConditionStatementElement
	LoopStatementElement
	ForStatementElement

	Name string `xml:"name,attr"`
	Type string `xml:"type,attr"`

	Exprs []ExpressionElement `xml:",any"`

	Body StatementBody
}

const OutputStatementElementName = "output"
type OutputStatementElement struct {
}

const VariableDeclarationElementName = "declare"
type VariableDeclarationElement struct {
}

const VariableAssignmentElementName = "assign"
type VariableAssignmentElement struct {
}

const LiteralExpressionStringElementName = "string"
const LiteralExpressionBoolElementName = "bool"
const LiteralExpressionIntElementName = "int"
const LiteralExpressionFloatElementName = "float"
const VariableExpressionElementName = "var"
const OperatorExpressionAddElementName = "add"
const OperatorExpressionSubElementName = "sub"
const OperatorExpressionMulElementName = "mul"
const OperatorExpressionDivElementName = "div"
const OperatorExpressionModElementName = "mod"
const OperatorExpressionConcatElementName = "concat"
const OperatorExpressionEqualElementName = "equal"
const OperatorExpressionGreaterThanElementName = "gt"
const OperatorExpressionLessThanElementName = "lt"
const OperatorExpressionNotElementName = "not"
const OperatorExpressionAndElementName = "and"
const OperatorExpressionOrElementName = "or"
const FunctionCallExpressionElementName = "call"
type ExpressionElement struct {
	XMLName xml.Name

	Name string `xml:"name,attr"`

	Content string `xml:",chardata"`
	Exprs []ExpressionElement `xml:",any"`
}

const FunctionElementName = "func"
type FunctionElement struct {
	Args FunctionArgs `xml:"args"`
}

type FunctionArgs struct {
	Returns FunctionReturnType `xml:"returns"`
	Args    []FunctionArg      `xml:"arg"`
}

type FunctionArg struct {
	Name string `xml:"name,attr"`
	Type string `xml:"type,attr"`
}

type FunctionReturnType struct {
	Type string `xml:"type,attr"`
}

type StatementBody struct {
	XMLName xml.Name `xml:"body"`
	Statements []StatementElement `xml:",any"`
}

const FunctionReturnElementName = "return"
type FunctionReturnElement struct {
}

const FunctionCallStatementElementName = "call"
type FunctionCallStatementElement struct {
}

const ConditionStatementElementName = "switch"
type ConditionStatementElement struct {
	Ifs []ConditionStatementIfClauseElement `xml:"if"`
	Else *ConditionStatementIfClauseElement `xml:"else"`
}

type ConditionStatementIfClauseElement struct {
	Condition ConditionStatementConditionElement `xml:"cond"`
	Then      ConditionStatementThenElement
}

type ConditionStatementElseClauseElement struct {
	Then ConditionStatementThenElement
}

type ConditionStatementConditionElement struct {
	Expr ExpressionElement `xml:",any"`
}

type ConditionStatementThenElement struct {
	XMLName xml.Name `xml:"then"`
	Statements []StatementElement `xml:",any"`
}

const LoopStatementElementName = "loop"
type LoopStatementElement struct {
	Condition LoopConditionElement `xml:"cond"`
}

type LoopConditionElement struct {
	Expr ExpressionElement `xml:",any"`
}

type LoopBodyElement struct {
	Statements []StatementElement `xml:",any"`
}

const ForStatementElementName = "for"
type ForStatementElement struct {
	To int `xml:"to,attr"`
	From int `xml:"from,attr"`
}