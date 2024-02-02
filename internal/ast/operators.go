package ast

type Operator int

const (
	Add Operator = iota
	Sub
	Mul
	Div
	Mod
	Concat

	Equal
	LessThan
	GreaterThan
	Not
	Or
	And
)