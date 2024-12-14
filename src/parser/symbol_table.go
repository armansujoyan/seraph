package parser

type Symbol struct {
  TypeDef   string
  IsDefined bool
  Value     string
}


type SymbolTable = map[string]*Symbol
