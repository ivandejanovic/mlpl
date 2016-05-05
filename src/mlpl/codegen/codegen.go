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

package codegen

import (
	"errors"
	"fmt"
	"mlpl/types"
)

const (
	pc  int = 7 // pc = program counter
	mp  int = 6 // mp = "memory pointer" point to top of memory (for temp storage)
	gp  int = 5 // gp = "global pointer" points to bottom of memory for (global) variable storage
	ac  int = 0 // accumulator
	ac1 int = 1 // 2nd accumulator
)

type codeBuffer struct {
	code        []string
	tmpOffset   int // tmpOffset is the memory offset for temps. It is decremented each time a temp is stored, and incremeted when loaded again.
	emitLoc     int // TM location number for current instruction emission
	highEmitLoc int // Highest TM location emitted so far. For use in conjunction with emitSkip, emitBackup, and emitRestore
}

/* Procedure emitSO emits a string-only TM instruction
   op = the opcode
   s = string
*/
func (codeBuf *codeBuffer) emitSO(op string, s string) {
	code := fmt.Sprintf("%3d: %5s %s", codeBuf.emitLoc, op, s)
	codeBuf.emitLoc += 1
	if codeBuf.highEmitLoc < codeBuf.emitLoc {
		codeBuf.highEmitLoc = codeBuf.emitLoc
	}
	codeBuf.code = append(codeBuf.code, code)
}

/* Procedure emitRO emits a register-only TM instruction
   op = the opcode
   r = target register
   s = 1st source register
   t = 2nd source register
*/
func (codeBuf *codeBuffer) emitRO(op string, r int, s int, t int) {
	code := fmt.Sprintf("%3d: %5s %d, %d, %d", codeBuf.emitLoc, op, r, s, t)
	codeBuf.emitLoc += 1
	if codeBuf.highEmitLoc < codeBuf.emitLoc {
		codeBuf.highEmitLoc = codeBuf.emitLoc
	}
	codeBuf.code = append(codeBuf.code, code)
}

/* Procedure emitRM emits a register-to-memory TM instruction
   op = the opcode
   r = target register
   d = the offset
   s = the base register
*/
func (codeBuf *codeBuffer) emitRM(op string, r int, d int, s int) {
	code := fmt.Sprintf("%3d: %5s %d, %d(%d)", codeBuf.emitLoc, op, r, d, s)
	codeBuf.emitLoc += 1
	if codeBuf.highEmitLoc < codeBuf.emitLoc {
		codeBuf.highEmitLoc = codeBuf.emitLoc
	}
	codeBuf.code = append(codeBuf.code, code)
}

// Function emitSkip skips "howMany" code locations for later backpatch. It also returns the current code position
func (codeBuf *codeBuffer) emitSkip(howMany int) int {
	i := codeBuf.emitLoc
	codeBuf.emitLoc += howMany
	if codeBuf.highEmitLoc < codeBuf.emitLoc {
		codeBuf.highEmitLoc = codeBuf.emitLoc
	}

	return i
}

// Procedure emitBackup backs up to loc = a previously skipped location
func (codeBuf *codeBuffer) emitBackup(loc int) {
	codeBuf.emitLoc = loc
}

// Procedure emitRestore restores the current code position to the highest previously unemitted position
func (codeBuf *codeBuffer) emitRestore() {
	codeBuf.emitLoc = codeBuf.highEmitLoc
}

/* Procedure emitRM_Abs converts an absolute reference to a pc-relative reference when emitting a register-to-memory TM instruction
   op = the opcode
   r = target register
   a = the absolute location in memory
*/
func (codeBuf *codeBuffer) emitRM_Abs(op string, r int, a int) {
	abs := a - (codeBuf.emitLoc + 1)
	code := fmt.Sprintf("%3d: %5s %d, %d(%d)", codeBuf.emitLoc, op, r, abs, pc)
	codeBuf.emitLoc += 1
	if codeBuf.highEmitLoc < codeBuf.emitLoc {
		codeBuf.highEmitLoc = codeBuf.emitLoc
	}
	codeBuf.code = append(codeBuf.code, code)
}

func findLoc(bucketMap map[string]types.Bucket, name string) int {
	bucket, ok := bucketMap[name]

	if ok {
		return bucket.MemLoc
	}

	return -1
}

// Procedure genStmt generates code at a statement node
func genStmt(treeNode *types.TreeNode, bucketMap map[string]types.Bucket, codeBuf *codeBuffer) {
	var p1, p2, p3 *types.TreeNode = nil, nil, nil
	var savedLoc1, savedLoc2, loc int

	switch treeNode.Stmt {
	case types.IfK:
		p1 = treeNode.Children[0]
		p2 = treeNode.Children[1]
		if len(treeNode.Children) == 3 {
			p3 = treeNode.Children[2]
		}

		// Generate code for test expression
		cGen(p1, bucketMap, codeBuf)
		savedLoc1 = codeBuf.emitSkip(1)

		// Recurse on then part
		cGen(p2, bucketMap, codeBuf)
		savedLoc2 = codeBuf.emitSkip(1)
		loc = codeBuf.emitSkip(0)
		codeBuf.emitBackup(savedLoc1)
		codeBuf.emitRM_Abs("JEQ", ac, loc)
		codeBuf.emitRestore()

		// Recurse on else part
		cGen(p3, bucketMap, codeBuf)
		loc = codeBuf.emitSkip(0)
		codeBuf.emitBackup(savedLoc2)
		codeBuf.emitRM_Abs("LDA", pc, loc)
		codeBuf.emitRestore()
	case types.RepeatK:
		p1 = treeNode.Children[0]
		p2 = treeNode.Children[1]
		loc = codeBuf.emitSkip(0)

		// Generate code for body
		cGen(p1, bucketMap, codeBuf)
		// Generate code for test
		cGen(p2, bucketMap, codeBuf)

		codeBuf.emitRM_Abs("JEQ", ac, loc)
	case types.AssignK:
		// Generate code for rhs
		p1 = treeNode.Children[0]
		cGen(p1, bucketMap, codeBuf)
		// Now store value
		loc = findLoc(bucketMap, treeNode.Name)
		codeBuf.emitRM("ST", ac, loc, gp)
	case types.ReadK:
		codeBuf.emitRO("IN", ac, 0, 0)
		loc = findLoc(bucketMap, treeNode.Name)
		codeBuf.emitRM("ST", ac, loc, gp)
	case types.WriteK:
		//Get child
		p1 = treeNode.Children[0]
		//Check if we output string or id
		if p1.Type == types.String {
			//Generate print code
			codeBuf.emitSO("PRINT", p1.ValString)
		} else {
			// Generate code for expression to write
			p1 = treeNode.Children[0]
			cGen(p1, bucketMap, codeBuf)
			// Now output it
			codeBuf.emitRO("OUT", ac, 0, 0)
		}
	}
}

// Procedure genExp generates code at an expression node
func genExp(treeNode *types.TreeNode, bucketMap map[string]types.Bucket, codeBuf *codeBuffer) {
	var p1, p2 *types.TreeNode
	var loc int

	switch treeNode.Exp {
	case types.ConstK:
		// Gen code to load integer constant using LDC
		codeBuf.emitRM("LDC", ac, treeNode.Val, 0)
	case types.IdK:
		loc = findLoc(bucketMap, treeNode.Name)
		codeBuf.emitRM("LD", ac, loc, gp)
	case types.OpK:
		p1 = treeNode.Children[0]
		p2 = treeNode.Children[1]
		// Gen code for ac = left arg
		cGen(p1, bucketMap, codeBuf)
		// Gen code to push left operand
		codeBuf.emitRM("ST", ac, codeBuf.tmpOffset, mp)
		codeBuf.tmpOffset -= 1
		// Gen code for ac = right operand
		cGen(p2, bucketMap, codeBuf)
		// Now load left operand
		codeBuf.tmpOffset += 1
		codeBuf.emitRM("LD", ac1, codeBuf.tmpOffset, mp)
		switch treeNode.Op {
		case types.PLUS:
			codeBuf.emitRO("ADD", ac, ac1, ac)
		case types.MINUS:
			codeBuf.emitRO("SUB", ac, ac1, ac)
		case types.TIMES:
			codeBuf.emitRO("MUL", ac, ac1, ac)
		case types.OVER:
			codeBuf.emitRO("DIV", ac, ac1, ac)
		case types.LT:
			codeBuf.emitRO("SUB", ac, ac1, ac)
			codeBuf.emitRM("JLT", ac, 2, pc)
			codeBuf.emitRM("LDC", ac, 0, ac)
			codeBuf.emitRM("LDA", pc, 1, pc)
			codeBuf.emitRM("LDC", ac, 1, ac)
		case types.EQ:
			codeBuf.emitRO("SUB", ac, ac1, ac)
			codeBuf.emitRM("JEQ", ac, 2, pc)
			codeBuf.emitRM("LDC", ac, 0, ac)
			codeBuf.emitRM("LDA", pc, 1, pc)
			codeBuf.emitRM("LDC", ac, 1, ac)
		default:
			panic(errors.New("Unknown operator"))
		}
	}
}

//Procedure cGen recursively generates code by tree traversal
func cGen(treeNode *types.TreeNode, bucketMap map[string]types.Bucket, codeBuf *codeBuffer) {
	if treeNode != nil {
		switch treeNode.Node {
		case types.StmtK:
			genStmt(treeNode, bucketMap, codeBuf)
		case types.ExpK:
			genExp(treeNode, bucketMap, codeBuf)
		default:
			err := errors.New("Unknow type for code generation")
			panic(err)
		}
		cGen(treeNode.Sibling, bucketMap, codeBuf)
	}
}

func CodeGen(treeNode *types.TreeNode, bucketMap map[string]types.Bucket) []string {
	codeBuf := &codeBuffer{make([]string, 0, 0), 0, 0, 0}

	codeBuf.emitRM("LD", mp, 0, ac)
	codeBuf.emitRM("ST", ac, 0, ac)
	cGen(treeNode, bucketMap, codeBuf)
	codeBuf.emitRO("HALT", 0, 0, 0)

	return codeBuf.code
}
