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
		":=":      {},
		":":       {},
		";":       {},
		".":       {},
	}
)

func Scan(reader *bufio.Reader) ([]Token, error) {
	output := make([]Token, 0)

	for {
		rune, _, err := reader.ReadRune()

		if err == io.EOF {
			break
		} else if err != nil {
      return nil, fmt.Errorf("Unable to read next character: %w", err)
		}

		if utils.IsCharacter(rune) {
			lexem := ""
			for ; utils.IsCharacter(rune) || utils.IsDigit(rune); rune, _, err = reader.ReadRune() {
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
				output = append(output, Token{"term", lexem})
			} else {
				output = append(output, Token{"ident", lexem})
			}
		}

		if utils.IsDigit(rune) {
			lexem := ""
			for ; utils.IsDigit(rune); rune, _, err = reader.ReadRune() {
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
				output = append(output, Token{"unknown", lexem})
			} else {
				output = append(output, Token{"number", lexem})
        if _, ok := Term[string(rune)]; ok {
          output = append(output, Token{"term", string(rune)})
        }
			}
			continue
		}

		if rune == ':' {
			rune, _, _ := reader.ReadRune()
			if rune == '=' {
				output = append(output, Token{"term", ":="})
			} else {
				output = append(output, Token{"term", ":"})
			}
			continue
		}

		if _, ok := Term[string(rune)]; ok {
			output = append(output, Token{"term", string(rune)})
		} else if rune != ' ' && rune != '\n' && rune != '\t' {
			output = append(output, Token{"unknown", string(rune)})
		}
	}

	return output, nil
}
