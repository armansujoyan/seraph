package scanner

import (
	"errors"
)

var ErrExhaustedInput = errors.New("No more tokens")

type Token struct {
	Category string
	Value    string
	Row      int
	Column   int
}

func (t *Token) IsEqual(token Token) bool {
	return t.Category == token.Category && t.Value == token.Value
}

type TokenIterator struct {
	slice []Token
	index int
}

func NewTokenIterator(tokens []Token) *TokenIterator {
	return &TokenIterator{tokens, 0}
}

func (iter *TokenIterator) Next() (Token, error) {
	if iter.HasMore() {
		next := iter.slice[iter.index]
		iter.index++
		return next, nil
	}
	return Token{}, ErrExhaustedInput
}

func (iter *TokenIterator) ViewNext() (Token, error) {
	if iter.HasMore() {
		return iter.slice[iter.index], nil
	} else {
		return Token{}, ErrExhaustedInput
	}
}

func (iter *TokenIterator) HasMore() bool {
	return iter.index < len(iter.slice)
}
