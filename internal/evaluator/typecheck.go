package evaluator

import "fmt"

// coerceType attempts to coerce a value to match a type hint.
// It returns an error if incompatible.
func coerceType(hint string, val *Value) (*Value, error) {
	if hint == "" || hint == "any" {
		return val, nil
	}
	switch hint {
	case "int":
		switch val.Type {
		case IntVal:
			return val, nil
		case FloatVal:
			return MakeInt(int64(val.FVal)), nil
		}
		return nil, fmt.Errorf("type error: expected int, got %s", val.Type)
	case "float":
		switch val.Type {
		case FloatVal:
			return val, nil
		case IntVal:
			return MakeFloat(float64(val.IVal)), nil
		}
		return nil, fmt.Errorf("type error: expected float, got %s", val.Type)
	case "string":
		if val.Type == StringVal {
			return val, nil
		}
		return nil, fmt.Errorf("type error: expected string, got %s", val.Type)
	case "bool":
		if val.Type == BoolVal {
			return val, nil
		}
		return nil, fmt.Errorf("type error: expected bool, got %s", val.Type)
	case "void":
		return MakeNull(), nil
	}
	// for custom types, just pass through
	return val, nil
}

func validateFields(typeDef *TypeDef, fields map[string]*Value) error {
	for _, f := range typeDef.Fields {
		v, ok := fields[f.Name]
		if !ok {
			if !f.Optional {
				return fmt.Errorf("missing required field '%s' for type '%s'", f.Name, typeDef.Name)
			}
			continue
		}
		if _, err := coerceType(f.TypeHint, v); err != nil {
			return fmt.Errorf("field '%s': %w", f.Name, err)
		}
	}
	return nil
}
