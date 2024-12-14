package parser

import (
	"errors"
	"seraph/src/scanner"
)

var expectedTokens = []*scanner.Token{&scanner.ProgramToken, &scanner.SemicolonToken}

func parseProgramHeader(iterator *scanner.TokenIterator) (Symbol, error) {
	ProgramNameSymbol := Symbol{}
	token, err := iterator.Next()
	if errors.Is(err, scanner.ErrExhaustedInput) {
		return ProgramNameSymbol, NewParserError("Expected program header", token)
	}

	if !token.IsEqual(scanner.ProgramToken) {
		return ProgramNameSymbol, NewParserError("Expected program header, found "+token.Value, token)
	}

	token, err = iterator.Next()
	if errors.Is(err, scanner.ErrExhaustedInput) {
		return ProgramNameSymbol, NewParserError("Expected identfier", token)
	}

	if token.Category != "ident" {
		return ProgramNameSymbol, NewParserError("Expected identifier found "+token.Category, token)
	} else {
		ProgramNameSymbol = Symbol{TypeDef: "moduleName", IsDefined: true, Value: token.Value}
	}

	token, err = iterator.Next()
	if errors.Is(err, scanner.ErrExhaustedInput) {
		return ProgramNameSymbol, NewParserError("Expected ';'", token)
	}

	if !token.IsEqual(scanner.SemicolonToken) {
		return ProgramNameSymbol, NewParserError("Expected ';' found "+token.Value, token)
	}

	return ProgramNameSymbol, nil
}
