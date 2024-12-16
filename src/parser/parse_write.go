package parser

import (
	"errors"
	"seraph/src/generator"
	"seraph/src/scanner"
)

func parseWrite(iterator *scanner.TokenIterator) error {
	token, err := iterator.Next()
	if errors.Is(err, scanner.ErrExhaustedInput) {
		return NewParserError("Expected write", token)
	}

	if !token.IsEqual(scanner.WriteToken) {
		return NewParserError("Expected 'write' found "+token.Value, token)
	}

	token, err = iterator.Next()
	if errors.Is(err, scanner.ErrExhaustedInput) {
		return NewParserError("Expected write", token)
	}

	if !token.IsEqual(scanner.OpenParenthesisToken) {
		return NewParserError("Expected '(' found "+token.Value, token)
	}

  token, err = iterator.Next()
	if errors.Is(err, scanner.ErrExhaustedInput) {
		return NewParserError("Expected identifier", token)
	}

	if !symbolTable.Exists(modularizeIdentifer(token.Value)) {
		return NewParserError("Unknown identifier: "+token.Value, token)
	}
  identifier := modularizeIdentifer(token.Value)
  idType := symbolTable[identifier].TypeDef

	token, err = iterator.Next()
	if errors.Is(err, scanner.ErrExhaustedInput) {
		return NewParserError("Expected write", token)
	}

	if !token.IsEqual(scanner.CloseParenthesisToken) {
		return NewParserError("Expected ')' found "+token.Value, token)
	}

	token, err = iterator.Next()
	if errors.Is(err, scanner.ErrExhaustedInput) {
		return NewParserError("Expected ';'", token)
	}

	if !token.IsEqual(scanner.SemicolonToken) {
		return NewParserError("Expected ';' found "+token.Value, token)
	}

  generator.GenerateWriteCall(identifier, idType)

  return nil
}
