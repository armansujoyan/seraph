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

  err = os.RemoveAll("out.*")
  if err != nil {
    return err
  }

  return nil
}
