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

package analyze

import (
	"errors"
	"fmt"
	"mlpl/types"
)

type procNode func(buf *buffer, node *types.TreeNode)

type buffer struct {
	location  int
	bucketMap map[string]types.Bucket
}

func (buf *buffer) st_insert(name string, lineno int) {
	bucket, ok := buf.bucketMap[name]

	if ok {
		line := bucket.Lines
		for line != nil {
			line = line.Next
		}
		line = &types.LineList{lineno, nil}
	} else {
		line := types.LineList{lineno, nil}
		bucket = types.Bucket{name, &line, buf.location}
		buf.location = buf.location + 1
		buf.bucketMap[name] = bucket
	}
}

func (buf *buffer) st_lookup(name string) int {
	bucket, ok := buf.bucketMap[name]

	if ok {
		return bucket.MemLoc
	}

	return -1
}

func typeError(lineno int, message string) {
	errorMessage := fmt.Sprintf("Type error at line %d: %s\n", lineno, message)
	err := errors.New(errorMessage)
	panic(err)
}

func insertNode(buf *buffer, node *types.TreeNode) {
	switch node.Node {
	case types.StmtK:
		if node.Stmt == types.AssignK || node.Stmt == types.ReadK {
			if buf.st_lookup(node.Name) == -1 {
				buf.st_insert(node.Name, node.Lineno)
			} else {
				buf.st_insert(node.Name, 0)
			}
		}
	case types.ExpK:
		if node.Exp == types.IdK {
			if buf.st_lookup(node.Name) == -1 {
				buf.st_insert(node.Name, node.Lineno)
			} else {
				buf.st_insert(node.Name, 0)
			}
		}
	}
}

func nullProc(buf *buffer, node *types.TreeNode) {
	return
}

func checkNode(buf *buffer, node *types.TreeNode) {
	switch node.Node {
	case types.ExpK:
		if node.Exp == types.OpK {
			if node.Children[0].Type != types.Integer || node.Children[1].Type != types.Integer {
				typeError(node.Lineno, "Op applied to non-integer")
			}
			if node.Op == types.EQ || node.Op == types.LT {
				node.Type = types.Boolean
			} else {
				node.Type = types.Integer
			}
		} else if node.Exp == types.ConstK || node.Exp == types.IdK {
			node.Type = types.Integer
		}
	case types.StmtK:
		switch node.Stmt {
		case types.IfK:
			if node.Children[0].Type == types.Integer {
				typeError(node.Lineno, "if test is not Boolean")
			}
		case types.AssignK:
			if node.Children[0].Type != types.Integer {
				typeError(node.Lineno, "assignment of non-integer value")
			}
		case types.WriteK:
			if node.Children[0].Type != types.Integer {
				typeError(node.Lineno, "write of non-integer value")
			}
		case types.RepeatK:
			if node.Children[0].Type == types.Integer {
				typeError(node.Lineno, "repeat test is not Boolean")
			}
		}
	}
}

func transverse(buf *buffer, node *types.TreeNode, preProc procNode, postProc procNode) {
	preProc(buf, node)
	for index := 0; index < len(node.Children); index++ {
		transverse(buf, node.Children[index], preProc, postProc)
	}
	postProc(buf, node)
	if node.Sibling != nil {
		transverse(buf, node.Sibling, preProc, postProc)
	}

}

func BuildSymtab(node *types.TreeNode) map[string]types.Bucket {
	buf := buffer{0, make(map[string]types.Bucket)}
	buf.location = 0
	transverse(&buf, node, insertNode, nullProc)
	return buf.bucketMap
}

func TypeCheck(node *types.TreeNode) {
	transverse(nil, node, nullProc, checkNode)
}
