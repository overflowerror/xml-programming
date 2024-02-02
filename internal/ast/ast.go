package ast

type Program struct {
	Statements []Statement
}

type Statement interface {
}

type OutputStatement struct {
	Exprs []Expression
}
var _ Statement = OutputStatement{}

type VariableDeclarationStatement struct {
	Name string
	Type Type
}
var _ Statement = VariableDeclarationStatement{}

type VariableAssignmentStatement struct {
	Name string
	Expr Expression
}
var _ Statement = VariableAssignmentStatement{}

type FunctionStatement struct {
	Name string
	Returns Type
	Args []FunctionArg
	Body []Statement
}
var _ Statement = FunctionStatement{}

type FunctionArg struct {
	Name string
	Type Type
}

type FunctionReturnStatement struct {
	Expr Expression
}
var _ Statement = FunctionReturnStatement{}

type ConditionalStatement struct {
	Ifs []ConditionIf
	Else []Statement
}
var _ Statement = ConditionalStatement{}

type ConditionIf struct {
	Expr Expression
	Then []Statement
}

type LoopStatement struct {
	LoopCondition Expression
	Body []Statement
}
var _ Statement = LoopStatement{}

type ForStatement struct {
	Name string
	From int
	To int
	Body []Statement
}
var _ Statement = ForStatement{}

type Expression interface {
}

type LiteralExpression struct {
	Type Type
	String string
	Bool bool
	Int int
	Float float32
}
var _ Expression = LiteralExpression{}

type VariableExpression struct {
	Name string
}
var _ Expression = VariableExpression{}

type OperatorExpression struct {
	Operator Operator
	Exprs []Expression
}
var _ Expression = OperatorExpression{}

type FunctionCall struct {
	Name string
	Args []Expression
}
var _ Expression = FunctionCall{}
var _ Statement = FunctionCall{}
