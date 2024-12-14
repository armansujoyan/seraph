package parser

import (
	"errors"
	"seraph/src/generator"
	"seraph/src/scanner"
)

func parseProgramEnd(iterator *scanner.TokenIterator) error {
	token, err := iterator.Next()
	if errors.Is(err, scanner.ErrExhaustedInput) {
		return NewParserError("Expected end", token)
	}

	if !token.IsEqual(scanner.EndToken) {
		return NewParserError("Expected end found "+token.Value, token)
	}
	//

	token, err = iterator.Next()
	if errors.Is(err, scanner.ErrExhaustedInput) {
		return NewParserError("Expected .", token)
	}

	if !token.IsEqual(scanner.DotToken) {
		return NewParserError("Expected .", token)
	}

	generator.GenerateProgramEnd()
	return nil
}


