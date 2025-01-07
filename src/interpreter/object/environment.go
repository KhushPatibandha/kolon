package object

type VariableType int

const (
	VAR VariableType = iota
	CONST
	FUNCTION
)

type Variable struct {
	Value Object
	Type  VariableType
}

type Environment struct {
	store map[string]*Variable
	outer *Environment
}

func NewEnvironment() *Environment {
	s := make(map[string]*Variable)
	return &Environment{store: s, outer: nil}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func (e *Environment) Get(name string) (*Variable, bool) {
	variable, ok := e.store[name]
	if !ok && e.outer != nil {
		variable, ok = e.outer.Get(name)
	}
	return variable, ok
}

func (e *Environment) Set(name string, val Object, valType VariableType) {
	e.store[name] = &Variable{
		Value: val,
		Type:  valType,
	}
}

func (e *Environment) Update(name string, newVal Object, valType VariableType) {
	variable, ok := e.store[name]
	if ok {
		if variable.Type == CONST || variable.Type == FUNCTION {
			return
		}
		variable.Value = newVal
		variable.Type = valType
	} else if !ok && e.outer != nil {
		e.outer.Update(name, newVal, valType)
	} else {
		return
	}
}
