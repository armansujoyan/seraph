package assembler

import (
	"fmt"
	"os/exec"
	"strings"
)

func Assemble(src string) error {
	name := strings.Split(src, ".")[0]
	assembleCmd := exec.Command("as", "-o", name+".o", name+".s")
  stdout, err := assembleCmd.Output()
  if err != nil {
    return err
  }

  fmt.Println(stdout)
  return nil
}
