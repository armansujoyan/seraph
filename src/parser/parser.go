package parser

import (
	"bufio"
	"fmt"
	"os"
	"seraph/src/generator"
	"seraph/src/scanner"
)

type SymbolTable = map[string]string

var (
	programName string
	symbolTable = make(SymbolTable)
	fileWriter  *bufio.Writer
)

// TODO: Remove the concrete iterator
func Parse(iterator *scanner.TokenIterator) string {
	prog := parseProgram(iterator)
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Parser error")
		}
	}()
	return prog
}

// TODO: Create buffer and strat writing to file
func parseProgram(iterator *scanner.TokenIterator) string {
	targetTransaltion := ""
	file, err := os.Create("./out.s")
	defer file.Close()
	if err != nil {
		panic("Unable to create file")
	}
	fileWriter = bufio.NewWriter(file)

	programHeader(iterator)
	// Q: Should this part be in definitions?

	fileWriter.Write([]byte(".data\n"))
	if ok, token := iterator.Next(); ok && token.Value == "var" {
		variableDefinitions(iterator)
	}
	return targetTransaltion
}

func programHeader(iterator *scanner.TokenIterator) {
	if ok, token := iterator.Next(); !ok || token.Category != "term" || token.Value != "program" {
		panic(fmt.Sprint("Invalid program header"))
	}

	if ok, token := iterator.Next(); ok && token.Category == "ident" {
		symbolTable[token.Value] = "programHeader"
		programName = token.Value
	} else {
		panic(fmt.Sprint("Invalid program header"))
	}

	if ok, token := iterator.Next(); !ok || token.Category != "term" || token.Value != ";" {
		panic(fmt.Sprint("Invalid program header"))
	}
}

func variableDefinitions(iterator *scanner.TokenIterator) {
	for iterator.ViewNext().Value != "begin" {
		variableSequence(iterator)
	}
}

// Refactor the whole method
func variableSequence(iterator *scanner.TokenIterator) {
	// Parse variable identifier
	variables := make(map[string]string)
	if ok, token := iterator.Next(); ok && token.Category == "ident" {
		if _, ok := symbolTable[token.Value]; !ok {
			variables[programName + "_" + token.Value] = "undefined"
		} else {
			panic(fmt.Sprint("Duplicate identifier"))
		}
	} else {
		panic(fmt.Sprint("Invalid variable declaration"))
	}

	// One or more definition
	ok, token := iterator.Next()
	for ; ok && token.Category == "term" && token.Value == ","; ok, token = iterator.Next() {
		if ok, token := iterator.Next(); ok && token.Category == "ident" {
			if _, ok := symbolTable[token.Value]; !ok {
				variables[programName + "_" + token.Value] = "undefined"
			} else {
				panic(fmt.Sprint("Duplicate identifier"))
			}
		} else {
			panic(fmt.Sprint("Invalid variable declaration"))
		}
	}

	// Get the type and assign to the variables
	if token.Category == "term" && token.Value == ":" {
		// TODO: Remove this mess
		if ok, t := iterator.Next(); ok && t.Category == "term" && (t.Value == "integer" || t.Value == "string") {
			for key := range variables {
				variables[key] = t.Value
				symbolTable[key] = t.Value
			}
		} else {
			panic(fmt.Sprint("Invalid variable declaration"))
		}
		if ok, token := iterator.Next(); !ok || token.Category != "term" || token.Value != ";" {
			panic(fmt.Sprint("Invalid variable declaration"))
		}
	}

  generator.GenerateVariables(variables, fileWriter)
}
