package parser

import (
	"errors"
	"fmt"
	"seraph/src/generator"
	"seraph/src/scanner"
)

type Symbol struct {
	TypeDef   string
	IsDefined bool
	Value     string
}

type VariableDefintion struct {
	Name   string
	Symbol *Symbol
}

type ParserError struct {
	message string
	row     int
	column  int
}

func (err *ParserError) Error() string {
	return fmt.Sprintf("Parser error: %s at %d, %d", err.message, err.row, err.column)
}

// TODO: Add position
func NewParserError(msg string, token scanner.Token) error {
	return &ParserError{msg, token.Row, token.Column}
}

type SymbolTable = map[string]*Symbol

var (
	moduleName   string
	mathOperands = []scanner.Token{scanner.PlusToken, scanner.MinusToken}
	symbolTable  = make(SymbolTable)
)

// TODO: Maybe the iterator interface has to be different
// TODO: Maybe I need linked lists with handlers in it?
func Parse(iterator *scanner.TokenIterator) (string, error) {
	moduleSymbol, err := parseProgramHeader(iterator)
	if err != nil {
		return "", fmt.Errorf("Error parsing program header: %w", err)
	}
	symbolTable[moduleSymbol.Value] = &moduleSymbol
	moduleName = moduleSymbol.Value
	generator.Init(moduleSymbol.Value)

	nextToken, err := iterator.ViewNext()
	if err == nil && nextToken.IsEqual(scanner.VarToken) {
		generator.GenerateGlobalVarSection()
		err := parseVariableDefinitions(iterator)
		if err != nil {
			return "", fmt.Errorf("Error parsing variable definitions: %w", err)
		}
	}

	token, err := iterator.Next()
	if errors.Is(err, scanner.ErrExhaustedInput) {
		return "", NewParserError("Expected 'begin'", token)
	}

	if !token.IsEqual(scanner.BeginToken) {
		return "", NewParserError("Expected 'begin' found "+token.Value, token)
	}

	generator.GenerateTextSection()
	for {
		nextToken, err := iterator.ViewNext()
		if err != nil || nextToken.IsEqual(scanner.EndToken) {
			break
		}
		parseStatementSequence(iterator)
	}

	parseProgramEnd(iterator)

	return moduleName, nil
}

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

func parseVariableSequence(iterator *scanner.TokenIterator) error {
	variables := make([]*VariableDefintion, 0)
	variable, err := parseVariable(iterator)
	variables = append(variables, variable)
	if err != nil {
		return fmt.Errorf("Unable to define variable : %w", err)
	}

	for {
		token, err := iterator.ViewNext()
		if errors.Is(err, scanner.ErrExhaustedInput) || token.IsEqual(scanner.ColonToken) {
			break
		}

		if !token.IsEqual(scanner.CommaToken) {
			return NewParserError("Expected ',' found"+token.Value, token)
		}

		iterator.Next()
		variable, err = parseVariable(iterator)
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

	if !token.IsEqual(scanner.IntegerToken) || token.IsEqual(scanner.StringToken) {
		return NewParserError("Expected 'integer' or 'string' found "+token.Value, token)
	}

	variablesTypeMap := make(map[string]string)
	for _, varDef := range variables {
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

	token, err = iterator.Next()
	if errors.Is(err, scanner.ErrExhaustedInput) {
		return NewParserError("Expected ':='", token)
	}

	if !token.IsEqual(scanner.AssignmentToken) {
		return NewParserError("Expected identifier found "+token.Value, token)
	}

	operand, err := parseOperand(iterator)

	token, err = iterator.Next()
	if errors.Is(err, scanner.ErrExhaustedInput) {
		return NewParserError("Expected operator or ';'", token)
	}

	if token.IsEqual(scanner.SemicolonToken) {
		generator.GenerateStatement(&target, &operand)
		return nil
	}

	// Parse + or -
	if !token.IsEqual(scanner.PlusToken) && !token.IsEqual(scanner.MinusToken) {
		return NewParserError("Expected operator found "+token.Value, token)
	}

	operator := token

	// Parse complex
	operand_two, err := parseOperand(iterator)
	generator.GenerateComplexStatement(&target, &operand, &operator, &operand_two)

	// Scan last semicolon
	token, err = iterator.Next()
	if errors.Is(err, scanner.ErrExhaustedInput) {
		return NewParserError("Expected ';'", token)
	}

	if token.IsEqual(scanner.SemicolonToken) {
		return NewParserError("Expected ';' found "+token.Value, token)
	}

	return nil
}

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
		return token, nil
	case "number":
		return token, nil
	default:
		return scanner.Token{}, NewParserError("Expected identifier or number found "+token.Category, token)
	}
}

func parseVariable(iterator *scanner.TokenIterator) (*VariableDefintion, error) {
	token, err := iterator.Next()
	if errors.Is(err, scanner.ErrExhaustedInput) {
		return nil, NewParserError("Expected identifier", token)
	}
	if token.Category != "ident" {
		return nil, NewParserError("Expected identifier found "+token.Value, token)
	}

	_, isDuplicate := symbolTable[token.Value]
	if isDuplicate {
		return nil, NewParserError("Duplicate identifier "+token.Value, token)
	}
	variableDefinition := &VariableDefintion{
		Name:   modularizeIdentifer(token.Value),
		Symbol: &Symbol{"unknown", false, ""},
	}
	return variableDefinition, nil
}

func modularizeIdentifer(ident string) string {
	return moduleName + "_" + ident
}
