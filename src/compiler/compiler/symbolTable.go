package compiler

type SymbolScope string

const GlobalScope SymbolScope = "GLOBAL"

type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

type SymbolTable struct {
	store  map[string]Symbol
	numDef int
}

func NewSymTable() *SymbolTable {
	s := make(map[string]Symbol)
	return &SymbolTable{store: s}
}

func (s *SymbolTable) Define(name string) Symbol {
	symbol := Symbol{Name: name, Index: s.numDef, Scope: GlobalScope}
	s.store[name] = symbol
	s.numDef++
	return symbol
}

func (s *SymbolTable) Resolve(name string) (Symbol, bool) {
	sym, ok := s.store[name]
	return sym, ok
}
