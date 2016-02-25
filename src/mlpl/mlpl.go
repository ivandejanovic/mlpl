package main

import (
	"fmt"
	"mlpl/cfg"
	"mlpl/lexer"
	"mlpl/parse"
	"mlpl/types"
	"os"
)

func main() {
	args := os.Args[1:]
	argc := len(args)

	if argc < 1 || argc > 2 {
		fmt.Println("Usage: <codefilename> [configurationfilename].")
		return
	}

	var reserved []types.ReservedWord

	if argc == 2 {
		reserved = cfg.GetConfigReservedWords(args[1])
	} else {
		reserved = cfg.GetDefaultReserved()
	}

	fmt.Println(reserved)

	tokens := parse.Parse(args[0], reserved)

	fmt.Println(tokens)

	treeNode := lexer.Lex(tokens)

	fmt.Println(treeNode)
}
