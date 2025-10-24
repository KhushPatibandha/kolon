package ktype

var typePool = make(map[string]*Type)

func InternType(t *Type) *Type {
	if t == nil {
		return nil
	}
	key := t.String()
	if existing, ok := typePool[key]; ok {
		return existing
	}
	typePool[key] = t
	return t
}

func ResetTypePool() {
	typePool = make(map[string]*Type)
}
