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
	// multicharacter tokens
	ID
	NUM
	// special symbols
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
)

type ExpType int

const (
	Void ExpType = 1 + iota
	Integer
	Boolean
)

type TreeNode struct {
	Children []TreeNode
	Sibling  TreeNode
	Lineno   int
	Node     NodeKind
	Stmt     StmtKind
	Exp      ExpKind
	Op       TokenType
	Val      int
	Name     string
	Type     ExpType
}
