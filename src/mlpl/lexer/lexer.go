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

package lexer

import (
	"errors"
	"fmt"
	"mlpl/types"
)

type lexBuffer struct {
	token  types.Token
	index  int
	tokens []types.Token
}

func syntaxError(token types.Token) {
	fmt.Printf("Syntax error at line %d, unexpected token -> ", token.Lineno)

	switch token.TokenType {
	case types.IF:
		fmt.Printf("reserved word: %s\n", token.TokenString)
	case types.THEN:
		fmt.Printf("reserved word: %s\n", token.TokenString)
	case types.ELSE:
		fmt.Printf("reserved word: %s\n", token.TokenString)
	case types.END:
		fmt.Printf("reserved word: %s\n", token.TokenString)
	case types.REPEAT:
		fmt.Printf("reserved word: %s\n", token.TokenString)
	case types.UNTIL:
		fmt.Printf("reserved word: %s\n", token.TokenString)
	case types.READ:
		fmt.Printf("reserved word: %s\n", token.TokenString)
	case types.WRITE:
		fmt.Printf("reserved word: %s\n", token.TokenString)
	case types.ASSIGN:
		fmt.Printf(":=\n")
	case types.LT:
		fmt.Printf("<\n")
	case types.EQ:
		fmt.Printf("=\n")
	case types.LPAREN:
		fmt.Printf("(\n")
	case types.RPAREN:
		fmt.Printf(")\n")
	case types.SEMI:
		fmt.Printf(";\n")
	case types.PLUS:
		fmt.Printf("+\n")
	case types.MINUS:
		fmt.Printf("-\n")
	case types.TIMES:
		fmt.Printf("*\n")
	case types.OVER:
		fmt.Printf("/\n")
	case types.ENDFILE:
		fmt.Printf("EOF\n")
	case types.NUM:
		fmt.Printf("NUM, name= %s\n", token.TokenString)
	case types.ID:
		fmt.Printf("ID, name= %s\n", token.TokenString)
	case types.ERROR:
		fmt.Printf("ERROR: %s\n", token.TokenString)
	default:
		// Should never happen.
		fmt.Printf("Unknown token: %d\n", token.TokenType)
	}

	err := errors.New("Aborting.")
	if err != nil {
		panic(err)
	}
}

func getToken(buffer *lexBuffer) {
	if buffer.index < len(buffer.tokens) {
		buffer.index += 1
		buffer.token = buffer.tokens[buffer.index]
	}
}

func match(expected types.TokenType, buffer *lexBuffer) {
	if buffer.token.TokenType == expected {
		getToken(buffer)
	} else {
		syntaxError(buffer.token)
	}
}

func statement(buffer *lexBuffer) *types.TreeNode {
	return nil
}

func stmt_sequence(buffer *lexBuffer) *types.TreeNode {
	t := statement(buffer)
	p := t

	for buffer.token.TokenType != types.ENDFILE &&
		buffer.token.TokenType != types.END &&
		buffer.token.TokenType != types.ELSE &&
		buffer.token.TokenType != types.UNTIL {
		q := statement(buffer)
		if q != nil {
			if t == nil {
				p = q
				t = p
			} else {
				// Now p cannot be nil either.
				p.Sibling = q
				p = q
			}
		}
	}

	return t
}

func Lex(tokens []types.Token) *types.TreeNode {
	buffer := lexBuffer{tokens[0], 0, tokens}

	return stmt_sequence(&buffer)
}
