/*
The MIT License (MIT)

Copyright (c) 2016 Ivan Dejanovic

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

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
