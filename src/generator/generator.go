package generator

import (
	"bufio"
	"seraph/src/scanner"
)

func GenerateVariables(variables map[string]string, writer *bufio.Writer) {
	for key := range variables {
		declaration := "  .lcomm " + key + ", 4\n"
		writer.Write([]byte(declaration))
	}
	writer.Flush()
}

// TODO: Refactor
func GenerateComplexStatement(targetIdent *scanner.Token, operand string, mathOperation *scanner.Token, operandTwo *scanner.Token) {
}

func GenerateStatement(targetIdent *scanner.Token, operand *scanner.Token, writer *bufio.Writer) {
	if operand.Category == "number" {
		writer.Write([]byte("  movl $" + operand.Value + ", %eax\n"))
		writer.Write([]byte("  movl %eax, " + targetIdent.Value + "\n"))
	} else {
		writer.Write([]byte("  movl " + operand.Value + ", %eax\n"))
		writer.Write([]byte("  movl %eax, " + targetIdent.Value + "\n"))
	}
	writer.Flush()
}
