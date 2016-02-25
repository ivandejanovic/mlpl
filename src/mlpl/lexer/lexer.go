package lexer

import (
	"mlpl/types"
)

func stmt_sequence(tokens []types.Token) types.TreeNode {
	return nil
}

func Lex(tokens []types.Token) types.TreeNode {
	return stmt_sequence(tokens)
}
