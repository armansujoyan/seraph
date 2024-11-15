package scanner

type Token struct {
	Category string
	Value    string
}

type TokenIterator struct {
	slice []Token
	index int
}

func NewTokenIterator(tokens []Token) *TokenIterator {
	return &TokenIterator{tokens, 0}
}

func (iter *TokenIterator) Next() (bool, *Token) {
	if iter.HasMore() {
		next := &iter.slice[iter.index]
		iter.index++
		return true, next
	}
	return false, nil
}

func (iter *TokenIterator) ViewNext() *Token {
	if iter.HasMore() {
		return &iter.slice[iter.index]
	} else {
		return nil
	}
}

func (iter *TokenIterator) HasMore() bool {
	return iter.index < len(iter.slice)
}
