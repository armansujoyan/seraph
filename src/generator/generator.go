package generator

import "bufio"

func GenerateVariables(variables map[string]string, writer *bufio.Writer) {
	for key := range variables {
		declaration := "  .lcomm " + key + ", 4\n"
		writer.Write([]byte(declaration))
	}
	writer.Flush()
}
