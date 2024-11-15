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
func Parse(iterator *scanner.TokenIterator) {
	parseProgram(iterator)
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Parser error")
		}
	}()
}

// TODO: Create buffer and strat writing to file
func parseProgram(iterator *scanner.TokenIterator) {
	file, err := os.Create("./out.s")
	defer file.Close()
	if err != nil {
		panic("Unable to create file")
	}
	fileWriter = bufio.NewWriter(file)

	programHeader(iterator)

	// Q: Should this part be in definitions?
	fileWriter.Write([]byte(".data\n"))
	if iterator.ViewNext().Value == "var" {
		variableDefinitions(iterator)
	}
}

func programHeader(iterator *scanner.TokenIterator) {
	if token, ok := iterator.Next(); !ok || token.Category != "term" || token.Value != "program" {
		panic(fmt.Sprint("Invalid program header"))
	}

	if token, ok := iterator.Next(); ok && token.Category == "ident" {
		symbolTable[token.Value] = "programHeader"
		programName = token.Value
	} else {
		panic(fmt.Sprint("Invalid program header"))
	}

	if token, ok := iterator.Next(); !ok || token.Category != "term" || token.Value != ";" {
		panic(fmt.Sprint("Invalid program header"))
	}
}

func variableDefinitions(iterator *scanner.TokenIterator) {
	if token, ok := iterator.Next(); ok && token.Value == "var" {
    for iterator.ViewNext().Value != "begin" {
			variableSequence(iterator)
		}
	} else {
    panic("Invalide start of variable defintion. Expected 'var'")
  }
}

// Refactor the whole method
func variableSequence(iterator *scanner.TokenIterator) {
	// Parse variable identifier
	variables := make(map[string]string)
	if token, ok := iterator.Next(); ok && token.Category == "ident" {
		if _, ok := symbolTable[token.Value]; !ok {
			variables[programName+"_"+token.Value] = "undefined"
		} else {
      panic(fmt.Sprint("Duplicate identifier :" + token.Value))
		}
	} else {
		panic(fmt.Sprint("Invalid variable declaration"))
	}

	// One or more definition
	token, ok := iterator.Next()
	for ; ok && token.Category == "term" && token.Value == ","; token, ok = iterator.Next() {
		if token, ok := iterator.Next(); ok && token.Category == "ident" {
			if _, ok := symbolTable[token.Value]; !ok {
				variables[programName+"_"+token.Value] = "undefined"
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
		if t, ok := iterator.Next(); ok && t.Category == "term" && (t.Value == "integer" || t.Value == "string") {
			for key := range variables {
				variables[key] = t.Value
				symbolTable[key] = t.Value
			}
		} else {
			panic(fmt.Sprint("Invalid variable declaration"))
		}
		if token, ok := iterator.Next(); !ok || token.Category != "term" || token.Value != ";" {
			panic(fmt.Sprint("Invalid variable declaration"))
		}
	}

	generator.GenerateVariables(variables, fileWriter)
}
