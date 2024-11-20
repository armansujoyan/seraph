package assembler

import (
	"os"
	"os/exec"
	"strings"
)

func AssembleExecutable(src string) error {
	name := strings.Split(src, ".")[0]

	assembleCmd := exec.Command("as", "-o", name+".o", name+".s")
  _, err := assembleCmd.Output()
  if err != nil {
    return err
  }

	linkCmd := exec.Command("ld", "-o", name, name+".o")
  _, err = linkCmd.Output()
  if err != nil {
    return err
  }

  // TODO: This will change to include custom names
  // based on input file name
  err = os.Remove("./out.s")
  if err != nil {
    return err
  }

  err = os.Remove("./out.o")
  if err != nil {
    return err
  }

  return nil
}
