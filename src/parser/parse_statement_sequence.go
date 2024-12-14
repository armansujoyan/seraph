package parser

import (
	"errors"
	"seraph/src/scanner"
)

func parseStatementSequence(iterator *scanner.TokenIterator) error {
	token, err := iterator.ViewNext()
	if errors.Is(err, scanner.ErrExhaustedInput) {
		return NewParserError("Expected identifier or write", token)
	}

	if token.Value == "write" {
		return parseWrite(iterator)
	}

  return parseAssignment(iterator)
}
