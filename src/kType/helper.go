package ktype

func NewBaseType(t string) *Type {
	ty := &Type{
		Kind: TypeBase,
		Name: t,
	}
	return InternType(ty)
}

func NewArrayType(ele *Type) *Type {
	ty := &Type{
		Kind:        TypeArray,
		ElementType: ele,
	}
	return InternType(ty)
}

func NewHashMapType(kt, vt *Type) *Type {
	ty := &Type{
		Kind:      TypeHashMap,
		KeyType:   kt,
		ValueType: vt,
	}
	return InternType(ty)
}

func (t *Type) Equals(other *Type) bool {
	if t == other {
		return true
	}
	// fmt.Println("curr type and other type check miss")
	if other == nil || t.Kind != other.Kind {
		return false
	}
	switch t.Kind {
	case TypeBase:
		if other.Kind != TypeBase {
			return false
		}
		return t.Name == other.Name
	case TypeArray:
		if other.Kind != TypeArray {
			return false
		}
		if other.ElementType == nil {
			other.ElementType = t.ElementType
			return true
		}
		return t.ElementType.Equals(other.ElementType)
	default:
		if other.Kind != TypeHashMap {
			return false
		}
		if other.KeyType == nil && other.ValueType == nil {
			other.KeyType = t.KeyType
			other.ValueType = t.ValueType
			return true
		}
		return t.KeyType.Equals(other.KeyType) && t.ValueType.Equals(other.ValueType)
	}
}

func (t *Type) TypeKindToString() string {
	switch t.Kind {
	case TypeBase:
		return "TypeBase"
	case TypeArray:
		return "TypeArray"
	case TypeHashMap:
		return "TypeHashMap"
	default:
		return "UnknownTypeKind"
	}
}
