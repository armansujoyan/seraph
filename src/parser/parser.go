package parser

import (
	"bufio"
	"fmt"
	"os"
	"seraph/src/generator"
	"seraph/src/scanner"
	"seraph/src/utils"
)

type SymbolTable = map[string]string

var (
	mathOperands = []string{"+", "-"}
)

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

	parseProgramHeader(iterator)

	// Q: Should this part be in definitions?
	fileWriter.Write([]byte(".section .bss\n"))
	if iterator.ViewNext().Value == "var" {
		parseVariableDefinitions(iterator)
	}

	if token, ok := iterator.Next(); ok && token.Value == "begin" {
		fileWriter.Write([]byte(".section .text\n"))
    fileWriter.Write([]byte("  .globl _start\n_start:\n"))
		for iterator.ViewNext().Value != "end" {
			parseStatementSequence(iterator)
		}
	} else {
		panic("Invalid program begin")
	}

  parseProgramEnd(iterator);
}

func parseProgramHeader(iterator *scanner.TokenIterator) {
	if token, ok := iterator.Next(); !ok || token.Category != "term" || token.Value != "program" {
		panic("Invalid program header")
	}

	if token, ok := iterator.Next(); ok && token.Category == "ident" {
		symbolTable[token.Value] = "programHeader"
		programName = token.Value
	} else {
		panic("Invalid program header")
	}

	if token, ok := iterator.Next(); !ok || token.Category != "term" || token.Value != ";" {
		panic("Invalid program header")
	}
}

func parseVariableDefinitions(iterator *scanner.TokenIterator) {
	if token, ok := iterator.Next(); ok && token.Value == "var" {
		for iterator.ViewNext().Value != "begin" {
			parseVariableSequence(iterator)
		}
	} else {
		panic("Invalide start of variable defintion. Expected 'var'")
	}
}

// Refactor the whole method
func parseVariableSequence(iterator *scanner.TokenIterator) {
	// Parse variable identifier
	variables := make(map[string]string)
	if token, ok := iterator.Next(); ok && token.Category == "ident" {
		if token.Value == programName {
			panic(fmt.Sprint("Variable name cannot be the same as program name"))
		}
		if ok := validateIdentifier(token.Value); !ok {
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

func parseStatementSequence(iterator *scanner.TokenIterator) {
	targetIdent, _ := iterator.Next()
	parseIdentifier(targetIdent)
	if assignmentExpr, ok := iterator.Next(); assignmentExpr.Value != ":=" || !ok {
		panic("Invalid statement")
	}
	operand, _ := iterator.Next()
	parseOperand(operand)

	if iterator.ViewNext().Value != ";" {
		mathOperation, _ := iterator.Next()
		if !utils.Contains(mathOperands, mathOperation.Value) {
			panic("Invalid operator: " + mathOperation.Value)
		}
		operandTwo, _ := iterator.Next()
		parseOperand(operandTwo)
		generator.GenerateComplexStatement(targetIdent, operand, mathOperation, operandTwo, fileWriter)
    if closingToken, ok := iterator.Next(); closingToken.Value != ";" || !ok {
      panic("Invalid token: expectd ; got " + closingToken.Value)
    }
	} else {
		generator.GenerateStatement(targetIdent, operand, fileWriter)
    if closingToken, ok := iterator.Next(); closingToken.Value != ";" || !ok {
      panic("Invalid token: expectd ; got " + closingToken.Value)
    }
	}
}

func parseProgramEnd(iterator *scanner.TokenIterator) {
	if token, ok := iterator.Next(); !ok || token.Category != "term" || token.Value != "end" {
		panic("Invalid program header")
	}

	if token, ok := iterator.Next(); !ok || token.Category != "term" || token.Value != "." {
		panic("Invalid program header")
	}

  generator.GenerateProgramEnd(fileWriter)
}

func parseOperand(operand *scanner.Token) {
	if operand.Category == "ident" {
		if ok := validateIdentifier(operand.Value); !ok {
			panic("Invalid identifier: " + operand.Value)
		}
    modularizeToken(operand);
    if symbolTable[operand.Value] != "integer" {
      panic("Invalid identifier type: " + operand.Value)
    }
	} else if operand.Category != "number" {
		panic("Invalid operand, expected number")
	}
}

func parseIdentifier(ident *scanner.Token) {
	if ok := validateIdentifier(ident.Value); !ok || ident.Category != "ident" {
		panic("Invalid identifier: " + ident.Value)
	}
  modularizeToken(ident)
}

func validateIdentifier(identifier string) bool {
	_, ok := symbolTable[programName+"_"+identifier]
	return ok
}

func getVariableFromToken(token *scanner.Token) string {
	return programName + "_" + token.Value
}

// TODO: Remove this, it is a bad idea
func modularizeToken(token *scanner.Token) {
  token.Value = programName + "_" + token.Value
}
