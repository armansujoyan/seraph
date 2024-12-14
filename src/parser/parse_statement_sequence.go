package parser

import (
	"errors"
	"seraph/src/scanner"
)

func parseStatementSequence(iterator *scanner.TokenIterator) error {
	token, err := iterator.Next()
	if errors.Is(err, scanner.ErrExhaustedInput) {
		return NewParserError("Expected identifier", token)
	}

	if token.Category != "ident" {
		return NewParserError("Expected identifier found "+token.Value, token)
	}

	if !symbolTable.Exists(modularizeIdentifer(token.Value)) {
		return NewParserError("Unknown identifier: "+token.Value, token)
	}

	target := token
	target.Value = modularizeIdentifer(target.Value)

	token, err = iterator.Next()
	if errors.Is(err, scanner.ErrExhaustedInput) {
		return NewParserError("Expected ':='", token)
	}

	if !token.IsEqual(scanner.AssignmentToken) {
		return NewParserError("Expected identifier found "+token.Value, token)
	}

	err = parseExpression(&target, iterator)
	if err != nil {
		return NewParserError(err.Error(), token)
	}

	// Scan last semicolon
	token, err = iterator.Next()
	if errors.Is(err, scanner.ErrExhaustedInput) {
		return NewParserError("Expected ';'", token)
	}

	if !token.IsEqual(scanner.SemicolonToken) {
		return NewParserError("Expected ';' found "+token.Value, token)
	}

	symbolTable[target.Value].IsDefined = true
	return nil
}
