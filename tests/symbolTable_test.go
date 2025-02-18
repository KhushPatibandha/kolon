package tests

import (
	"testing"

	"github.com/KhushPatibandha/Kolon/src/compiler/compiler"
)

func Test53(t *testing.T) {
	expected := map[string]compiler.Symbol{
		"a": {Name: "a", Scope: compiler.GlobalScope, Index: 0},
		"b": {Name: "b", Scope: compiler.GlobalScope, Index: 1},
	}
	global := compiler.NewSymTable()
	a := global.Define("a")
	if a != expected["a"] {
		t.Errorf("expected a=%+v, got=%+v", expected["a"], a)
	}
	b := global.Define("b")
	if b != expected["b"] {
		t.Errorf("expected b=%+v, got=%+v", expected["b"], b)
	}
}

func Test54(t *testing.T) {
	global := compiler.NewSymTable()
	global.Define("a")
	global.Define("b")

	expected := []compiler.Symbol{
		{Name: "a", Scope: compiler.GlobalScope, Index: 0},
		{Name: "b", Scope: compiler.GlobalScope, Index: 1},
	}

	for _, sym := range expected {
		res, ok := global.Resolve(sym.Name)
		if !ok {
			t.Errorf("name %s not resolvable", sym.Name)
			continue
		}
		if res != sym {
			t.Errorf("expected %s to resolve to %+v, got=%+v", sym.Name, sym, res)
		}
	}
}
