package io

import (
	"errors"
	"os"
	"strings"
)

func ParseArgs() (string, error) {
	args := os.Args[1:]
	if len(args) == 0 {
		return "", errors.New("No-command line arguments supplied")
	}

	sourceFileName := args[0]
	err := parseFileName(sourceFileName)
	if err != nil {
		return "", err
	}

	return sourceFileName, nil
}

func OpenSourceFile(sourceFileName string) (*os.File, error) {
  _, err := os.Stat(sourceFileName)
	if err != nil {
		if os.IsNotExist(err) {
		  return nil, errors.New("File does not exist")
		}
	}

	file, err := os.Open(sourceFileName)
	if err != nil {
		return nil, errors.New("Unable to read the file")
	}

  return file, nil
}

func parseFileName(fileName string) error {
	fileNameParts := strings.Split(fileName, ".")
	if len(fileNameParts) != 2 {
		return errors.New("Invalid file name. Should be of format *.pas")
	}
	return nil
}
