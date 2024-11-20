package utils

func IsCharacter(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c > 'A' && c < 'Z')
}

func IsDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func Contains[S ~[]E, E comparable](list S, element E) bool {
	return Index(list, element) >= 0
}

func Index[S ~[]E, E comparable](list S, element E) int {
	for i := range list {
		if list[i] == element {
			return i
		}
	}
	return -1
}
