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
	"mlpl/locale"
	"mlpl/types"
	"strconv"
)

type lexBuffer struct {
	token  types.Token
	index  int
	tokens []types.Token
}

func syntaxError(token types.Token) {
	fmt.Printf(locale.Locale.LexerSyntaxError, token.Lineno)

	switch token.TokenType {
	case types.IF, types.THEN, types.ELSE, types.END, types.REPEAT, types.UNTIL, types.READ, types.WRITE:
		fmt.Printf(locale.Locale.LexerReservedWordError, token.TokenString)
	case types.ASSIGN:
		fmt.Printf(locale.Locale.LexerAssignError)
	case types.LT:
		fmt.Printf(locale.Locale.LexerLTError)
	case types.EQ:
		fmt.Printf(locale.Locale.LexerEQError)
	case types.LPAREN:
		fmt.Printf(locale.Locale.LexerLPARENError)
	case types.RPAREN:
		fmt.Printf(locale.Locale.LexerRPARENError)
	case types.SEMI:
		fmt.Printf(locale.Locale.LexerSEMIError)
	case types.PLUS:
		fmt.Printf(locale.Locale.LexerPLUSError)
	case types.MINUS:
		fmt.Printf(locale.Locale.LexerMINUSError)
	case types.TIMES:
		fmt.Printf(locale.Locale.LexerTIMESError)
	case types.OVER:
		fmt.Printf(locale.Locale.LexerOVERError)
	case types.ENDFILE:
		fmt.Printf(locale.Locale.LexerENDFILEError)
	case types.NUM:
		fmt.Printf(locale.Locale.LexerNUMError, token.TokenString)
	case types.ID:
		fmt.Printf(locale.Locale.LexerIDError, token.TokenString)
	case types.ERROR:
		fmt.Printf(locale.Locale.LexerERRORError, token.TokenString)
	default:
		// Should never happen.
		fmt.Printf(locale.Locale.LexerDEFAULTError, token.TokenType)
	}

	err := errors.New(locale.Locale.LexerABORTINGError)
	if err != nil {
		panic(err)
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

func (buffer *lexBuffer) nextToken() {
	if buffer.index < len(buffer.tokens) {
		buffer.index += 1
		buffer.token = buffer.tokens[buffer.index]
	}
}

func (buffer *lexBuffer) match(expected types.TokenType) {
	if buffer.token.TokenType == expected {
		buffer.nextToken()
	} else {
		syntaxError(buffer.token)
	}
}

func (buffer *lexBuffer) factor() *types.TreeNode {
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
		buffer.match(types.NUM)
	case types.ID:
		node = newExpNode(types.IdK, buffer.token.Lineno)
		if buffer.token.TokenType == types.ID {
			node.Name = buffer.token.TokenString
		}
		buffer.match(types.ID)
	case types.STRING:
		node = newExpNode(types.StringK, buffer.token.Lineno)
		node.ValString = buffer.token.TokenString
		buffer.match(types.STRING)
	case types.LPAREN:
		buffer.match(types.LPAREN)
		node = buffer.exp()
		buffer.match(types.RPAREN)
	default:
		syntaxError(buffer.token)
	}

	return node
}

func (buffer *lexBuffer) term() *types.TreeNode {
	node := buffer.factor()

	for buffer.token.TokenType == types.TIMES || buffer.token.TokenType == types.OVER {
		p := newExpNode(types.OpK, buffer.token.Lineno)
		p.Children = append(p.Children, node)
		p.Op = buffer.token.TokenType
		node = p
		buffer.match(buffer.token.TokenType)
		node.Children = append(node.Children, buffer.factor())
	}

	return node
}

func (buffer *lexBuffer) simpleExp() *types.TreeNode {
	node := buffer.term()

	for buffer.token.TokenType == types.PLUS || buffer.token.TokenType == types.MINUS {
		p := newExpNode(types.OpK, buffer.token.Lineno)
		p.Children = append(p.Children, node)
		p.Op = buffer.token.TokenType
		node = p
		buffer.match(buffer.token.TokenType)
		node.Children = append(node.Children, buffer.term())
	}

	return node
}

func (buffer *lexBuffer) exp() *types.TreeNode {
	node := buffer.simpleExp()

	if buffer.token.TokenType == types.LT || buffer.token.TokenType == types.EQ {
		p := newExpNode(types.OpK, buffer.token.Lineno)
		p.Children = append(p.Children, node)
		p.Op = buffer.token.TokenType
		node = p
		buffer.match(buffer.token.TokenType)
		node.Children = append(node.Children, buffer.simpleExp())
	}

	return node
}

func (buffer *lexBuffer) ifStmt() *types.TreeNode {
	node := newStmtNode(types.IfK, buffer.token.Lineno)

	buffer.match(types.IF)
	node.Children = append(node.Children, buffer.exp())
	buffer.match(types.THEN)
	node.Children = append(node.Children, buffer.stmtSequence())
	if buffer.token.TokenType == types.ELSE {
		buffer.match(types.ELSE)
		node.Children = append(node.Children, buffer.stmtSequence())
	}

	return node
}

func (buffer *lexBuffer) repeatStmt() *types.TreeNode {
	node := newStmtNode(types.RepeatK, buffer.token.Lineno)

	buffer.match(types.REPEAT)
	node.Children = append(node.Children, buffer.stmtSequence())
	buffer.match(types.UNTIL)
	node.Children = append(node.Children, buffer.exp())

	return node
}

func (buffer *lexBuffer) assignStmt() *types.TreeNode {
	node := newStmtNode(types.AssignK, buffer.token.Lineno)

	if buffer.token.TokenType == types.ID {
		node.Name = buffer.token.TokenString
	}
	buffer.match(types.ID)
	buffer.match(types.ASSIGN)
	node.Children = append(node.Children, buffer.exp())

	return node
}

func (buffer *lexBuffer) readStmt() *types.TreeNode {
	node := newStmtNode(types.ReadK, buffer.token.Lineno)

	buffer.match(types.READ)
	if buffer.token.TokenType == types.ID {
		node.Name = buffer.token.TokenString
	}
	buffer.match(types.ID)

	return node
}

func (buffer *lexBuffer) writeStmt() *types.TreeNode {
	node := newStmtNode(types.WriteK, buffer.token.Lineno)

	buffer.match(types.WRITE)
	node.Children = append(node.Children, buffer.exp())

	return node
}

func (buffer *lexBuffer) statement() *types.TreeNode {
	var node *types.TreeNode

	switch buffer.token.TokenType {
	case types.IF:
		node = buffer.ifStmt()
	case types.REPEAT:
		node = buffer.repeatStmt()
	case types.ID:
		node = buffer.assignStmt()
	case types.READ:
		node = buffer.readStmt()
	case types.WRITE:
		node = buffer.writeStmt()
	default:
		syntaxError(buffer.token)
	}

	return node
}

func (buffer *lexBuffer) stmtSequence() *types.TreeNode {
	node := buffer.statement()
	p := node

	for buffer.token.TokenType != types.ENDFILE &&
		buffer.token.TokenType != types.END &&
		buffer.token.TokenType != types.ELSE &&
		buffer.token.TokenType != types.UNTIL {
		buffer.match(types.SEMI)
		q := buffer.statement()
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

func (buffer *lexBuffer) lexSequence() *types.TreeNode {
	var node, p, q *types.TreeNode = nil, nil, nil

	for buffer.token.TokenType != types.ENDFILE {
		if buffer.token.TokenType == types.END ||
			buffer.token.TokenType == types.ELSE ||
			buffer.token.TokenType == types.UNTIL {
			buffer.nextToken()
		}
		if buffer.token.TokenType == types.ENDFILE {
			break
		}
		p = buffer.stmtSequence()

		if node == nil {
			node = p
		} else {
			q = node

			for q.Sibling != nil {
				q = q.Sibling
			}
			q.Sibling = p
		}
	}

	return node
}

func Lex(tokens []types.Token) *types.TreeNode {
	buffer := &lexBuffer{tokens[0], 0, tokens}

	return buffer.lexSequence()
}
