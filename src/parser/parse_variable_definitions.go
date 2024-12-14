package parser

import (
	"errors"
	"seraph/src/scanner"
)

func parseVariableDefinitions(iterator *scanner.TokenIterator) error {
	token, err := iterator.Next()
	if errors.Is(err, scanner.ErrExhaustedInput) {
		return NewParserError("Expected 'var'", token)
	}

	if !token.IsEqual(scanner.VarToken) {
		return NewParserError("Expected 'var' found "+token.Value, token)
	}

	for {
		token, err := iterator.ViewNext()
		if err != nil || token.IsEqual(scanner.BeginToken) {
			break
		}
		err = parseVariableSequence(iterator)
		if err != nil {
			return err
		}
	}

	return nil
}


