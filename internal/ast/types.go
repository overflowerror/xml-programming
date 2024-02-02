package ast

type Type int

const (
	Void Type = iota
	String
	Bool
	Int
	Float
)

func (t Type) IsNumber() bool {
	return t == Int || t == Float
}