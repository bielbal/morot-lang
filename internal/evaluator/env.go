package evaluator

import "fmt"

type Environment struct {
	store  map[string]*binding
	parent *Environment
}

type binding struct {
	value   *Value
	isConst bool
}

func NewEnv() *Environment {
	return &Environment{store: make(map[string]*binding)}
}

func NewEnvWithParent(parent *Environment) *Environment {
	return &Environment{store: make(map[string]*binding), parent: parent}
}

func (e *Environment) Get(name string) (*Value, bool) {
	if b, ok := e.store[name]; ok {
		return b.value, true
	}
	if e.parent != nil {
		return e.parent.Get(name)
	}
	return nil, false
}

func (e *Environment) Set(name string, val *Value, isConst bool) {
	e.store[name] = &binding{value: val, isConst: isConst}
}

func (e *Environment) Assign(name string, val *Value) error {
	if b, ok := e.store[name]; ok {
		if b.isConst {
			return fmt.Errorf("cannot reassign const '%s'", name)
		}
		b.value = val
		return nil
	}
	if e.parent != nil {
		return e.parent.Assign(name, val)
	}
	return fmt.Errorf("undefined variable '%s'", name)
}

func (e *Environment) SetTypes(types map[string]*TypeDef) {
	for k, v := range types {
		e.store["__type__"+k] = &binding{value: &Value{Type: ObjectVal, Fields: nil, TypeName: k, SVal: k}, isConst: true}
		_ = v
	}
}

type TypeDef struct {
	Name   string
	Fields []TypeFieldDef
}

type TypeFieldDef struct {
	Name     string
	TypeHint string
	Optional bool
}
