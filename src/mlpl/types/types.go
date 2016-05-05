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

package types

import ()

type TokenType int

const (
	ENDFILE TokenType = 1 + iota
	ERROR
	// reserved words
	IF
	THEN
	ELSE
	END
	REPEAT
	UNTIL
	READ
	WRITE
	// Multicharacter tokens.
	ID
	NUM
	STRING
	// Special symbols.
	ASSIGN
	EQ
	LT
	PLUS
	MINUS
	TIMES
	OVER
	LPAREN
	RPAREN
	SEMI
)

type ReservedWord struct {
	TokenType TokenType
	Str       string
}

type Token struct {
	TokenType   TokenType
	TokenString string
	Lineno      int
}

type NodeKind int

const (
	StmtK NodeKind = 1 + iota
	ExpK
)

type StmtKind int

const (
	IfK StmtKind = 1 + iota
	RepeatK
	AssignK
	ReadK
	WriteK
)

type ExpKind int

const (
	OpK ExpKind = 1 + iota
	ConstK
	IdK
	StringK
)

type ExpType int

const (
	Void ExpType = 1 + iota
	Integer
	Boolean
	String
)

type TreeNode struct {
	Children  []*TreeNode
	Sibling   *TreeNode
	Lineno    int
	Node      NodeKind
	Stmt      StmtKind
	Exp       ExpKind
	Op        TokenType
	Val       int
	Name      string
	ValString string
	Type      ExpType
}

type LineList struct {
	Lineno int
	Next   *LineList
}

type Bucket struct {
	Name   string
	Lines  *LineList
	MemLoc int
}
