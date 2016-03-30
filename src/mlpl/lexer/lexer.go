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
	"strconv"
)

type lexBuffer struct {
	token  types.Token
	index  int
	tokens []types.Token
}

func syntaxError(token types.Token) {
	fmt.Printf("Syntax error at line %d, unexpected token -> ", token.Lineno)

	switch token.TokenType {
	case types.IF, types.THEN, types.ELSE, types.END, types.REPEAT, types.UNTIL, types.READ, types.WRITE:
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

func newStmtNode(kind types.StmtKind, lineno int) *types.TreeNode {
	node := new(types.TreeNode)

	node.Children = make([]*types.TreeNode, 0, 0)
	node.Sibling = nil
	node.Node = types.StmtK
	node.Stmt = kind
	node.Lineno = lineno

	return node
}

func newExpNode(kind types.ExpKind, lineno int) *types.TreeNode {
	node := new(types.TreeNode)

	node.Children = make([]*types.TreeNode, 0, 0)
	node.Sibling = nil
	node.Node = types.ExpK
	node.Exp = kind
	node.Lineno = lineno

	return node
}

func factor(buffer *lexBuffer) *types.TreeNode {
	var node *types.TreeNode
	var err error

	switch buffer.token.TokenType {
	case types.NUM:
		node = newExpNode(types.ConstK, buffer.token.Lineno)
		if buffer.token.TokenType == types.NUM {
			node.Val, err = strconv.Atoi(buffer.token.TokenString)
			if err != nil {
				panic(err)
			}
		}
		match(types.NUM, buffer)
	case types.ID:
		node = newExpNode(types.IdK, buffer.token.Lineno)
		if buffer.token.TokenType == types.ID {
			node.Name = buffer.token.TokenString
		}
		match(types.ID, buffer)
	case types.LPAREN:
		match(types.LPAREN, buffer)
		node = exp(buffer)
		match(types.RPAREN, buffer)
	default:
		syntaxError(buffer.token)
	}

	return node
}

func term(buffer *lexBuffer) *types.TreeNode {
	node := factor(buffer)

	for buffer.token.TokenType == types.TIMES || buffer.token.TokenType == types.OVER {
		p := newExpNode(types.OpK, buffer.token.Lineno)
		p.Children = append(p.Children, node)
		p.Op = buffer.token.TokenType
		node = p
		match(buffer.token.TokenType, buffer)
		node.Children = append(node.Children, factor(buffer))
	}

	return node
}

func simple_exp(buffer *lexBuffer) *types.TreeNode {
	node := term(buffer)

	for buffer.token.TokenType == types.PLUS || buffer.token.TokenType == types.MINUS {
		p := newExpNode(types.OpK, buffer.token.Lineno)
		p.Children = append(p.Children, node)
		p.Op = buffer.token.TokenType
		node = p
		match(buffer.token.TokenType, buffer)
		node.Children = append(node.Children, term(buffer))
	}

	return node
}

func exp(buffer *lexBuffer) *types.TreeNode {
	node := simple_exp(buffer)

	if buffer.token.TokenType == types.LT || buffer.token.TokenType == types.EQ {
		p := newExpNode(types.OpK, buffer.token.Lineno)
		p.Children = append(p.Children, node)
		p.Op = buffer.token.TokenType
		node = p
		match(buffer.token.TokenType, buffer)
		node.Children = append(node.Children, simple_exp(buffer))
	}

	return node
}

func if_stmt(buffer *lexBuffer) *types.TreeNode {
	node := newStmtNode(types.IfK, buffer.token.Lineno)

	match(types.IF, buffer)
	node.Children = append(node.Children, exp(buffer))
	match(types.THEN, buffer)
	node.Children = append(node.Children, stmt_sequence(buffer))
	if buffer.token.TokenType == types.ELSE {
		match(types.ELSE, buffer)
		node.Children = append(node.Children, stmt_sequence(buffer))
	}

	return node
}

func repeat_stmt(buffer *lexBuffer) *types.TreeNode {
	node := newStmtNode(types.RepeatK, buffer.token.Lineno)

	match(types.REPEAT, buffer)
	node.Children = append(node.Children, stmt_sequence(buffer))
	match(types.UNTIL, buffer)
	node.Children = append(node.Children, exp(buffer))

	return node
}

func assign_stmt(buffer *lexBuffer) *types.TreeNode {
	node := newStmtNode(types.AssignK, buffer.token.Lineno)

	if buffer.token.TokenType == types.ID {
		node.Name = buffer.token.TokenString
	}
	match(types.ID, buffer)
	match(types.ASSIGN, buffer)
	node.Children = append(node.Children, exp(buffer))

	return node
}

func read_stmt(buffer *lexBuffer) *types.TreeNode {
	node := newStmtNode(types.ReadK, buffer.token.Lineno)

	match(types.READ, buffer)
	if buffer.token.TokenType == types.ID {
		node.Name = buffer.token.TokenString
	}
	match(types.ID, buffer)

	return node
}

func write_stmt(buffer *lexBuffer) *types.TreeNode {
	node := newStmtNode(types.WriteK, buffer.token.Lineno)

	match(types.WRITE, buffer)
	node.Children = append(node.Children, exp(buffer))

	return node
}

func statement(buffer *lexBuffer) *types.TreeNode {
	var node *types.TreeNode

	switch buffer.token.TokenType {
	case types.IF:
		node = if_stmt(buffer)
	case types.REPEAT:
		node = repeat_stmt(buffer)
	case types.ID:
		node = assign_stmt(buffer)
	case types.READ:
		node = read_stmt(buffer)
	case types.WRITE:
		node = write_stmt(buffer)
	default:
		syntaxError(buffer.token)
	}

	return node
}

func stmt_sequence(buffer *lexBuffer) *types.TreeNode {
	node := statement(buffer)
	p := node

	for buffer.token.TokenType != types.ENDFILE &&
		buffer.token.TokenType != types.END &&
		buffer.token.TokenType != types.ELSE &&
		buffer.token.TokenType != types.UNTIL {
		match(types.SEMI, buffer)
		q := statement(buffer)
		if q != nil {
			if node == nil {
				p = q
				node = p
			} else {
				// Now p cannot be nil either.
				p.Sibling = q
				p = q
			}
		}
	}

	return node
}

func Lex(tokens []types.Token) *types.TreeNode {
	buffer := lexBuffer{tokens[0], 0, tokens}

	return stmt_sequence(&buffer)
}
