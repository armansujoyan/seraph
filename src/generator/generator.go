package generator

import (
	"bufio"
	"seraph/src/scanner"
	"strings"
)

func GenerateVariables(variables map[string]string, writer *bufio.Writer) {
	for key := range variables {
		declaration := "  .lcomm " + key + ", 4\n"
		writer.Write([]byte(declaration))
	}
	writer.Flush()
}

// TODO: Refactor
func GenerateComplexStatement(targetIdent, operand, operator, operandTwo *scanner.Token, writer *bufio.Writer) {
	loadOperandOne := move(operand, "eax")
	loadOperandTwo := move(operandTwo, "ebx")
	executeOperation := performArithmeticOperation("eax", "ebx", operator)
	storeOperaion := storeRegister("ebx", targetIdent.Value)
	statement := []string{loadOperandOne, loadOperandTwo, executeOperation, storeOperaion, "\n"}
	writer.Write([]byte(strings.Join(statement, "\n")))
	writer.Flush()
}

func GenerateStatement(targetIdent, operand *scanner.Token, writer *bufio.Writer) {
	loadOperand := move(operand, "eax")
	storeOperation := storeRegister("eax", targetIdent.Value)
	statement := []string{loadOperand, storeOperation, "\n"}
	writer.Write([]byte(strings.Join(statement, "\n")))
	writer.Flush()
}

func GenerateProgramEnd(writer *bufio.Writer) {
	returnStatement := `
  mov $60, %rax
  xor %rdi, %rdi
  syscall
  `
	writer.Write([]byte(returnStatement))
	writer.Flush()
}

// TODO: Generalize to handle reg to reg?
func move(loc1 *scanner.Token, loc2 string) string {
	t1 := "$" + loc1.Value
	if loc1.Category != "number" {
		t1 = loc1.Value
	}
	return "  movl " + t1 + ", %" + loc2
}

func performArithmeticOperation(reg1, reg2 string, operator *scanner.Token) string {
	if operator.Value == "+" {
		return "  addl %" + reg1 + ", %" + reg2
	} else {
		return "  imul %" + reg1 + ", %" + reg2
	}
}

func storeRegister(reg, variable string) string {
	return "  movl %" + reg + ", " + variable
}
