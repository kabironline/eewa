package compiler

type SymbolScope string

const (
	GlobalScope SymbolScope = "GLOBAL"
)

type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}
type SymbolTable struct {
	store          map[string]*Symbol
	numDefinitions int
}

func NewSymbolTable() *SymbolTable {
	s := &SymbolTable{}
	s.store = make(map[string]*Symbol)
	return s
}

func (s *SymbolTable) Define(name string) Symbol {
	if _, ok := s.store[name]; ok {
		return *s.store[name]
	}
	sym := &Symbol{
		Name:  name,
		Scope: GlobalScope,
		Index: s.numDefinitions,
	}
	s.store[name] = sym
	s.numDefinitions++
	return *sym
}

func (s *SymbolTable) Resolve(name string) (Symbol, bool) {
	obj, ok := s.store[name]
	return *obj, ok
}
