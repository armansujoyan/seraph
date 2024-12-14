package parser

import (
	"errors"
	"seraph/src/scanner"
)

func parseVariableDefinition(iterator *scanner.TokenIterator) (*VariableDefintion, error) {
	token, err := iterator.Next()
	if errors.Is(err, scanner.ErrExhaustedInput) {
		return nil, NewParserError("Expected identifier", token)
	}

	if token.Category != "ident" {
		return nil, NewParserError("Expected identifier found "+token.Value, token)
	}

	if symbolTable.Exists(token.Value) {
		return nil, NewParserError("Duplicate identifier "+token.Value, token)
	}

	variableDefinition := &VariableDefintion{
		Name:   modularizeIdentifer(token.Value),
		Symbol: &Symbol{"unknown", false, ""},
	}
	return variableDefinition, nil
}
