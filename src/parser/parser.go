package parser

import (
	"errors"
	"fmt"
	"seraph/src/allocator"
	"seraph/src/generator"
	"seraph/src/scanner"
)

type VariableDefintion struct {
	Name   string
	Symbol *Symbol
}

var (
	moduleName   string
	mathOperands = []scanner.Token{scanner.PlusToken, scanner.MinusToken}
	symbolTable  = make(SymbolTable)
  registerAllocator = allocator.NewAllocator(allocator.RegistersX86[:])
)

// TODO: Maybe the iterator interface has to be different
// TODO: Maybe I need linked lists with handlers in it?
func Parse(iterator *scanner.TokenIterator) (string, error) {
	moduleSymbol, err := parseProgramHeader(iterator)
	if err != nil {
		return "", fmt.Errorf("Error parsing program header: %w", err)
	}
	symbolTable[moduleSymbol.Value] = &moduleSymbol
	moduleName = moduleSymbol.Value
	generator.Init(moduleSymbol.Value)

	nextToken, err := iterator.ViewNext()
	if err == nil && nextToken.IsEqual(scanner.VarToken) {
		generator.GenerateGlobalVarSection()
		err := parseVariableDefinitions(iterator)
		if err != nil {
			return "", fmt.Errorf("Error parsing variable definitions: %w", err)
		}
	}

	token, err := iterator.Next()
	if errors.Is(err, scanner.ErrExhaustedInput) {
		return "", NewParserError("Expected 'begin'", token)
	}

	if !token.IsEqual(scanner.BeginToken) {
		return "", NewParserError("Expected 'begin' found "+token.Value, token)
	}

	generator.GenerateTextSection()
	for {
		nextToken, err := iterator.ViewNext()
		if err != nil || nextToken.IsEqual(scanner.EndToken) {
			break
		}
    err = parseStatementSequence(iterator)
    if err != nil {
      return "", err
    }
	}

  err = parseProgramEnd(iterator)
  if err != nil {
    return "", err
  }

  for varName, symbol := range symbolTable {
    if symbol.TypeDef == "string" && symbol.IsDefined {
      generator.GenerateStaticString(varName, symbol.Value);
    }
  }
  generator.GenerateStaticSection()

	return moduleName, nil
}

