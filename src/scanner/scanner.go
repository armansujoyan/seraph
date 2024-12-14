package scanner

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"seraph/src/utils"
)

var (
	err                 error
	max_identifier_size = 512
	row                 = 0
	column              = 0
	Term                = map[string]struct{}{
		"program": {},
		"var":     {},
		"begin":   {},
		"end":     {},
		"integer": {},
		"string":  {},
		",":       {},
		"+":       {},
		"-":       {},
		"*":       {},
		"(":       {},
		")":       {},
		":=":      {},
		":":       {},
		";":       {},
		".":       {},
	}
)

func Scan(reader *bufio.Reader) ([]Token, error) {
	output := make([]Token, 0)

	for {
		rune, err := readRune(reader)

		if err == io.EOF {
			break
		} else if err != nil {
			return nil, fmt.Errorf("Unable to read next character: %w", err)
		}

		if utils.IsCharacter(rune) {
			lexem := ""
			for ; utils.IsCharacter(rune) || utils.IsDigit(rune); rune, err = readRune(reader) {
				if len(lexem) > max_identifier_size {
					log.Fatal("Max token size exceeded")
				}
				if err == io.EOF {
					break
				} else if err != nil {
					log.Fatal("Unable to read next character", err)
				}
				lexem += string(rune)
			}

			_, ok := Term[lexem]
			if ok {
				output = append(output, Token{"term", lexem, row, column})
			} else {
				output = append(output, Token{"ident", lexem, row, column})
			}
		}

		if utils.IsDigit(rune) {
			lexem := ""
			for ; utils.IsDigit(rune); rune, err = readRune(reader) {
				if len(lexem) > max_identifier_size {
					log.Fatal("Max token size exceeded")
				}
				if err == io.EOF {
					break
				} else if err != nil {
					log.Fatal("Unable to read next character", err)
				}
				lexem += string(rune)
			}

			if utils.IsCharacter(rune) {
				lexem += string(rune)
				output = append(output, Token{"unknown", lexem, row, column})
			} else {
				output = append(output, Token{"number", lexem, row, column})
				if _, ok := Term[string(rune)]; ok {
					output = append(output, Token{"term", string(rune), row, column})
				}
			}
			continue
		}

		if rune == ':' {
			rune, _ := readRune(reader)
			if rune == '=' {
				output = append(output, Token{"term", ":=", row, column})
			} else {
				output = append(output, Token{"term", ":", row, column})
			}
			continue
		}

		if _, ok := Term[string(rune)]; ok {
			output = append(output, Token{"term", string(rune), row, column})
		} else if rune != ' ' && rune != '\n' && rune != '\t' {
			output = append(output, Token{"unknown", string(rune), row, column})
		}
	}

	return output, nil
}

func readRune(reader *bufio.Reader) (rune, error) {
	rune, _, err := reader.ReadRune()
	if err != nil {
		return 0, err
	}
  column++
	if rune == '\n' {
    row++
	}
	return rune, nil
}
