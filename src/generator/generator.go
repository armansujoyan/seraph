package generator

import (
	"bufio"
	"fmt"
	"os"
	"seraph/src/allocator"
	"seraph/src/common"
)

var (
	writer *bufio.Writer
  typeSizeMap = map[string]int{
    "string": 128,
    "integer": 8,
  }
)

func Init(name string) {
	file, err := os.Create("./" + name + ".s")
	if err != nil {
		panic("Unable to create file")
	}
	writer = bufio.NewWriter(file)
}

func GenerateGlobalVarSection() {
	writer.Write([]byte(".section .rodata\n"))
  writer.Write([]byte("digitfmt:\n"))
  writer.Write([]byte("  .asciz \"%d\\n\"\n"))
	writer.Write([]byte(".section .bss\n"))
	writer.Flush()
}

func GenerateTextSection() {
	writer.Write([]byte(".section .text\n"))
	writer.Write([]byte("  .globl _start\n_start:\n"))
	writer.Flush()
}

func GenerateVariables(variables map[string]string) {
	for variable, varType := range variables {
    size := typeSizeMap[varType]
		declaration := fmt.Sprintf("  .lcomm %s, %d\n", variable, size)
		writer.Write([]byte(declaration))
	}
	writer.Flush()
}

func GenerateArithmeticStatement(reg1, reg2 *allocator.Register, op *common.Operator, al *allocator.Allocator) *allocator.Register {
	if !reg1.GetIsLoaded() {
		LoadRegister(reg1)
	}
	if !reg2.GetIsLoaded() {
		LoadRegister(reg2)
	}
	generateArithmeticOperation(reg1, reg2, op)
	writer.Flush()
	al.Release(reg1)
	return reg2
}

func generateArithmeticOperation(reg1, reg2 *allocator.Register, op *common.Operator) {
	switch op.Value {
	case "+":
		writer.Write([]byte("  addq %" + reg1.GetName() + ", %" + reg2.GetName() + "\n"))
	case "-":
		writer.Write([]byte("  subq %" + reg1.GetName() + ", %" + reg2.GetName() + "\n"))
	default:
		writer.Write([]byte("  imulq %" + reg1.GetName() + ", %" + reg2.GetName() + "\n"))
	}
}

func LoadRegister(reg *allocator.Register) {
	statement := "  movq "
	if reg.GetType() == allocator.VariableRegister {
		statement += reg.GetContent()
	} else {
		statement += "$" + reg.GetContent()
	}
	reg.SetIsLoaded(true)
	statement += ", %" + reg.GetName() + "\n"
	writer.Write([]byte(statement))
	writer.Flush()
}

func StoreRegister(target string, reg *allocator.Register) {
	statement := "  movq %" + reg.GetName() + ", " + target + "\n\n"
	writer.Write([]byte(statement))
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

func GenerateWriteCall(identifier string, bytes int) {
  writeCall := fmt.Sprintf("  mov %s, %%rsi\n", identifier)
  writeCall += "  mov $digitfmt %rdi\n"
	writeCall += "  xor %rax, %rax\n"
	writeCall += "  call printf\n"
	writer.Write([]byte(writeCall))
	writer.Flush()
}
