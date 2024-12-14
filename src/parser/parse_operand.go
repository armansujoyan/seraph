package parser

import (
	"errors"
	"seraph/src/scanner"
)

func parseOperand(iterator *scanner.TokenIterator) (scanner.Token, error) {
	token, err := iterator.Next()
	if errors.Is(err, scanner.ErrExhaustedInput) {
		return scanner.Token{}, NewParserError("Expected identifier or number", token)
	}

	switch token.Category {
	case "ident":
		sym, isPresent := symbolTable[modularizeIdentifer(token.Value)]
		if !isPresent {
			return scanner.Token{}, NewParserError("Unknown identifier: "+token.Value, token)
		}
		if !sym.IsDefined {
			return scanner.Token{}, NewParserError("Undefined symbol: "+token.Value, token)
		}
    token.Value = modularizeIdentifer(token.Value)
		return token, nil
	case "number":
		return token, nil
	default:
		return scanner.Token{}, NewParserError("Expected identifier or number found "+token.Category, token)
	}
}


