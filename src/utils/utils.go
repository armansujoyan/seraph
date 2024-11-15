package utils

func IsCharacter(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c > 'A' && c < 'Z')
}

func IsDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

