package compiler

type SymbolScope string 

const (
	GlobalScope SymbolScope = "GLOBAL"
)

type Symbol struct {
	Name string
	Scope SymbolScope
	Index int
}
type SymbolTable struct {
	store map[string]*Symbol
	numDefinitions int
}

func NewSymbolTable() *SymbolTable {
	s := &SymbolTable{}
	s.store = make(map[string]*Symbol)
	return s
}