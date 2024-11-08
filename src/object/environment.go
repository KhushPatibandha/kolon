package object

type VariableType int

const (
	VAR VariableType = iota
	CONST
)

type Variable struct {
	Value Object
	Type  VariableType
}

type Environment struct {
	store map[string]*Variable
}

func NewEnvironment() *Environment {
	s := make(map[string]*Variable)
	return &Environment{store: s}
}

func (e *Environment) Get(name string) (*Variable, bool) {
	variable, ok := e.store[name]
	return variable, ok
}

func (e *Environment) Set(name string, val Object, valType VariableType) {
	e.store[name] = &Variable{
		Value: val,
		Type:  valType,
	}
}

func (e *Environment) Update(name string, newVal Object, valType VariableType) {
	delete(e.store, name)
	e.Set(name, newVal, valType)
}
