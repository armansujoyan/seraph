package main

import (
	"bufio"
	"fmt"
	"os"
	"seraph/src/assembler"
	"seraph/src/io"
	"seraph/src/parser"
	"seraph/src/scanner"
)

func main() {
  sourceFileName, err := io.ParseArgs()
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }

  file, err := io.OpenSourceFile(sourceFileName)
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
	defer file.Close()
	bufferedReader := bufio.NewReader(file)

	tokens, err := scanner.Scan(bufferedReader)
	if err != nil {
		fmt.Println(err)
    os.Exit(1)
	}

  // TODO: Need to check for the error, refine later
	tokenIterator := scanner.NewTokenIterator(tokens)
  name, err := parser.Parse(tokenIterator)
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }

  // TODO: Default to filename.s
  err = assembler.AssembleExecutable(name + ".s")
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}
