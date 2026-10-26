package evaluator

import (
	"fmt"
	"morot-lang/internal/ast"
	"strings"
)

type ValueType string

const (
	IntVal    ValueType = "int"
	FloatVal  ValueType = "float"
	StringVal ValueType = "string"
	BoolVal   ValueType = "bool"
	NullVal   ValueType = "null"
	ArrayVal  ValueType = "array"
	ObjectVal ValueType = "object"
	FnVal     ValueType = "fn"
	BuiltinVal ValueType = "builtin"
	ReturnVal ValueType = "return"
	BreakVal  ValueType = "break"
	ContinueVal ValueType = "continue"
	ErrorVal  ValueType = "error"
)

type Value struct {
	Type ValueType
	// primitives
	IVal int64
	FVal float64
	SVal string
	BVal bool
	// composite
	Elements []*Value            // array
	Fields   map[string]*Value   // object/struct
	TypeName string              // struct type name
	// function
	Params  []ast.Param
	Body    *ast.BlockStmt
	Env     *Environment
	// builtin
	Builtin func(args []*Value) (*Value, error)
	// wrapped
	Inner *Value // return/break/continue wrapper
	// error
	Err string
}

func MakeInt(v int64) *Value    { return &Value{Type: IntVal, IVal: v} }
func MakeFloat(v float64) *Value { return &Value{Type: FloatVal, FVal: v} }
func MakeString(v string) *Value { return &Value{Type: StringVal, SVal: v} }
func MakeBool(v bool) *Value    { return &Value{Type: BoolVal, BVal: v} }
func MakeNull() *Value          { return &Value{Type: NullVal} }
func MakeArray(els []*Value) *Value { return &Value{Type: ArrayVal, Elements: els} }
func MakeObject(fields map[string]*Value) *Value {
	return &Value{Type: ObjectVal, Fields: fields}
}
func MakeError(msg string) *Value { return &Value{Type: ErrorVal, Err: msg} }
func MakeReturn(v *Value) *Value  { return &Value{Type: ReturnVal, Inner: v} }
func MakeBreak() *Value           { return &Value{Type: BreakVal} }
func MakeContinue() *Value        { return &Value{Type: ContinueVal} }

func (v *Value) Truthy() bool {
	switch v.Type {
	case BoolVal:
		return v.BVal
	case NullVal:
		return false
	case IntVal:
		return v.IVal != 0
	case FloatVal:
		return v.FVal != 0
	case StringVal:
		return v.SVal != ""
	case ArrayVal:
		return len(v.Elements) > 0
	case ObjectVal:
		return true
	case FnVal, BuiltinVal:
		return true
	}
	return false
}

func (v *Value) String() string {
	switch v.Type {
	case IntVal:
		return fmt.Sprintf("%d", v.IVal)
	case FloatVal:
		f := fmt.Sprintf("%g", v.FVal)
		if !strings.Contains(f, ".") && !strings.Contains(f, "e") {
			f += ".0"
		}
		return f
	case StringVal:
		return v.SVal
	case BoolVal:
		if v.BVal {
			return "true"
		}
		return "false"
	case NullVal:
		return "null"
	case ArrayVal:
		parts := make([]string, len(v.Elements))
		for i, el := range v.Elements {
			parts[i] = el.Inspect()
		}
		return "[" + strings.Join(parts, ", ") + "]"
	case ObjectVal:
		parts := make([]string, 0, len(v.Fields))
		for k, fv := range v.Fields {
			parts = append(parts, k+": "+fv.Inspect())
		}
		name := v.TypeName
		if name == "" {
			name = "object"
		}
		return name + "{" + strings.Join(parts, ", ") + "}"
	case FnVal:
		return "<fn>"
	case BuiltinVal:
		return "<builtin>"
	case ErrorVal:
		return "Error: " + v.Err
	case ReturnVal:
		if v.Inner == nil {
			return "null"
		}
		return v.Inner.String()
	}
	return "?"
}

func (v *Value) Inspect() string {
	if v.Type == StringVal {
		return fmt.Sprintf("%q", v.SVal)
	}
	return v.String()
}

func (v *Value) ToFloat() float64 {
	switch v.Type {
	case IntVal:
		return float64(v.IVal)
	case FloatVal:
		return v.FVal
	}
	return 0
}
