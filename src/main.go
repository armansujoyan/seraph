package main

import (
	"fmt"
	"log"
	"os"
	"seraph/src/parser"
	"seraph/src/scanner"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		log.Fatal("No file name is provided")
	}
	fileName := args[0]
	tokens, err := scanner.Scan(fileName)

	if err != nil {
		fmt.Errorf("Syntax error")
	}

	tokenIterator := scanner.NewTokenIterator(tokens)
	parser.Parse(tokenIterator)
}
