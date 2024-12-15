package assembler

import (
	"os"
	"os/exec"
	"strings"
)

func AssembleExecutable(src string) error {
	name := strings.Split(src, ".")[0]
  output := name + ".o"
  input := name + ".s"

	assembleCmd := exec.Command("as", "-o", output, input)
	_, err := assembleCmd.Output()
	if err != nil {
		return err
	}

	linkCmd := exec.Command("ld", "-o", name, output, "-lc", "-dynamic-linker", "/lib64/ld-linux-x86-64.so.2")
	_, err = linkCmd.Output()
	if err != nil {
		return err
	}

	err = os.Remove("./" + input)
	if err != nil {
		return err
	}

	err = os.Remove("./" + output)
	if err != nil {
		return err
	}

	return nil
}
