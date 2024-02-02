package scope

import (
	"xml-programming/internal/ast"
	"xml-programming/internal/values"
)

type Function struct {
	Name string
	Args []ast.FunctionArg
	Return ast.Type
	Body []ast.Statement
}

type Variable struct {
	Name string
	values.Value
}

type Scope struct {
	functions []Function
	variables []Variable

	parentScope *Scope
}

func (s *Scope) CurrentScopeHas(name string) bool {
	for _, f := range s.functions {
		if f.Name == name {
			return true
		}
	}

	for _, v := range s.variables {
		if v.Name == name {
			return true
		}
	}

	return false
}

func (s *Scope) GetFunction(name string) *Function {
	for _, f := range s.functions {
		if f.Name == name {
			return &f
		}
	}

	if s.parentScope != nil {
		return s.parentScope.GetFunction(name)
	} else {
		return nil
	}
}

func (s *Scope) GetVariable(name string) *Variable {
	for i, v := range s.variables {
		if v.Name == name {
			return &s.variables[i]
		}
	}

	if s.parentScope != nil {
		return s.parentScope.GetVariable(name)
	} else {
		return nil
	}
}

func (s *Scope) AddVariable(variable Variable) {
	s.variables = append(s.variables, variable)
}

func (s *Scope) AddFunction(function Function) {
	s.functions = append(s.functions, function)
}

func New() *Scope {
	return &Scope{}
}

func FromParent(scope *Scope) *Scope {
	return &Scope{parentScope: scope}
}

