package parser

import (
	"errors"
	"seraph/src/generator"
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

	_, isListed := symbolTable[modularizeIdentifer(token.Value)]
	if !isListed {
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

	operand, err := parseOperand(iterator)
	if err != nil {
		return err
	}

	token, err = iterator.Next()
	if errors.Is(err, scanner.ErrExhaustedInput) {
		return NewParserError("Expected operator or ';'", token)
	}

	if token.IsEqual(scanner.SemicolonToken) {
		generator.GenerateStatement(&target, &operand)
		symbolTable[target.Value].IsDefined = true
		symbolTable[target.Value].Value = operand.Value
		return nil
	}

	// Parse + or -
	if !token.IsEqual(scanner.PlusToken) && !token.IsEqual(scanner.MinusToken) {
		return NewParserError("Expected operator found "+token.Value, token)
	}

	operator := token

	// Parse complex
	operand_two, err := parseOperand(iterator)
	if err != nil {
		return err
	}
	generator.GenerateComplexStatement(&target, &operand, &operator, &operand_two)

	// Scan last semicolon
	token, err = iterator.Next()
	if errors.Is(err, scanner.ErrExhaustedInput) {
		return NewParserError("Expected ';'", token)
	}

	if !token.IsEqual(scanner.SemicolonToken) {
		return NewParserError("Expected ';' found "+token.Value, token)
	}

	symbolTable[target.Value].IsDefined = true
	symbolTable[target.Value].Value = operand.Value
	return nil
}
