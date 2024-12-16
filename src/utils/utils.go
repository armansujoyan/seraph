package utils

import (
	"fmt"
	"io"
	"os"
)

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

func PrependToFile(filename, content string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	originalContent, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	file, err = os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file for writing: %v", err)
	}
	defer file.Close()

	newContent := content + string(originalContent)
	_, err = file.Write([]byte(newContent))
	if err != nil {
		return fmt.Errorf("failed to write to file: %v", err)
	}

	return nil
}
