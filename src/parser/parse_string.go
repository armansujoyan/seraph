package parser

import (
	"errors"
	"seraph/src/scanner"
)

func parseString(target *scanner.Token, iterator *scanner.TokenIterator) error {
	token, err := iterator.Next()
	if errors.Is(err, scanner.ErrExhaustedInput) {
		return NewParserError("Expected identifier or write", token)
	}

	if !token.IsEqual(scanner.QuotationMarkToken) {
		return NewParserError("Expected '\"' found "+token.Value, token)
	}

	token, err = iterator.Next()
	if token.Category != "string" {
		return NewParserError("Expected string found "+token.Value, token)
	}

	symbolTable[target.Value].IsDefined = true
	symbolTable[target.Value].Value = token.Value

	token, err = iterator.Next()
	if errors.Is(err, scanner.ErrExhaustedInput) {
		return NewParserError("Expected identifier or write", token)
	}

	if !token.IsEqual(scanner.QuotationMarkToken) {
		return NewParserError("Expected '\"' found "+token.Value, token)
	}

	return nil
}
