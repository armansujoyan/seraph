package parser

type Symbol struct {
  TypeDef   string
  IsDefined bool
  Value     string
}

func (table SymbolTable) Exists(name string) bool {
  _, ok := table[name]
  return ok
}

type SymbolTable map[string]*Symbol
