package parser

import (
	"errors"
	"fmt"
	"seraph/src/generator"
	"seraph/src/scanner"
)

func parseVariableSequence(iterator *scanner.TokenIterator) error {
	variables := make([]*VariableDefintion, 0)
	variable, err := parseVariableDefinition(iterator)
	if err != nil {
		return fmt.Errorf("Unable to define variable : %w", err)
	}
	variables = append(variables, variable)

	for {
		token, err := iterator.ViewNext()
		if errors.Is(err, scanner.ErrExhaustedInput) || token.IsEqual(scanner.ColonToken) {
			break
		}

		if !token.IsEqual(scanner.CommaToken) {
			return NewParserError("Expected ',' found"+token.Value, token)
		}

		iterator.Next()
		variable, err = parseVariableDefinition(iterator)
		variables = append(variables, variable)
		if err != nil {
			return fmt.Errorf("Unable to define variable : %w", err)
		}
	}

	token, err := iterator.Next()
	if errors.Is(err, scanner.ErrExhaustedInput) {
		return NewParserError("Expected ':'", token)
	}

	if !token.IsEqual(scanner.ColonToken) {
		return NewParserError("Expected ':' found "+token.Value, token)
	}

	token, err = iterator.Next()
	if errors.Is(err, scanner.ErrExhaustedInput) {
		return NewParserError("Expected variable type declaration", token)
	}

	if !token.IsEqual(scanner.IntegerToken) && !token.IsEqual(scanner.StringToken) {
		return NewParserError("Expected 'integer' or 'string' found "+token.Value, token)
	}

	variablesTypeMap := make(map[string]string)
	for _, varDef := range variables {
    varDef.Symbol.TypeDef = token.Value
		symbolTable[varDef.Name] = varDef.Symbol
		variablesTypeMap[varDef.Name] = varDef.Symbol.TypeDef
	}

	token, err = iterator.Next()
	if errors.Is(err, scanner.ErrExhaustedInput) {
		return NewParserError("Expected ';'", token)
	}

	if !token.IsEqual(scanner.SemicolonToken) {
		return NewParserError("Expected ';' found "+token.Value, token)
	}

	generator.GenerateVariables(variablesTypeMap)
	return nil
}
