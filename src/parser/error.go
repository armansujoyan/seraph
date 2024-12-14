package parser

import (
	"fmt"
	"seraph/src/scanner"
)

type ParserError struct {
	message string
	row     int
	column  int
}

func (err *ParserError) Error() string {
	return fmt.Sprintf("Parser error: %s at line %d, character %d", err.message, err.row, err.column)
}

func NewParserError(msg string, token scanner.Token) error {
	return &ParserError{msg, token.Row, token.Column}
}
