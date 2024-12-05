package generator

import (
	"bufio"
	"os"
	"seraph/src/scanner"
	"strings"
)

var (
	writer *bufio.Writer
)

func Init(name string) {
	file, err := os.Create("./" + name + ".s")
	if err != nil {
		panic("Unable to create file")
	}
	writer = bufio.NewWriter(file)
}

func GenerateGlobalVarSection() {
	writer.Write([]byte(".section .bss\n"))
  writer.Flush()
}

func GenerateTextSection() {
	writer.Write([]byte(".section .text\n"))
	writer.Write([]byte("  .globl _start\n_start:\n"))
  writer.Flush()
}

func GenerateVariables(variables map[string]string) {
	for variable := range variables {
		declaration := "  .lcomm " + variable + ", 4\n"
		writer.Write([]byte(declaration))
	}
	writer.Flush()
}

// TODO: Refactor
func GenerateComplexStatement(targetIdent, operand, operator, operandTwo *scanner.Token) {
	loadOperandOne := move(operand, "eax")
	loadOperandTwo := move(operandTwo, "ebx")
	executeOperation := performArithmeticOperation("eax", "ebx", operator)
	storeOperaion := storeRegister("ebx", targetIdent.Value)
	statement := []string{loadOperandOne, loadOperandTwo, executeOperation, storeOperaion, "\n"}
	writer.Write([]byte(strings.Join(statement, "\n")))
	writer.Flush()
}

func GenerateStatement(targetIdent, operand *scanner.Token) {
	loadOperand := move(operand, "eax")
	storeOperation := storeRegister("eax", targetIdent.Value)
	statement := []string{loadOperand, storeOperation, "\n"}
	writer.Write([]byte(strings.Join(statement, "\n")))
	writer.Flush()
}

func GenerateProgramEnd() {
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
		return "  subl %" + reg1 + ", %" + reg2
	}
}

func storeRegister(reg, variable string) string {
	return "  movl %" + reg + ", " + variable
}
