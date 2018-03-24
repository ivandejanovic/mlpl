/*
The MIT License (MIT)

Copyright (c) 2016-2018 Ivan Dejanovic

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

package vm

import (
	"fmt"
	"github.com/ivandejanovic/mlpl/locale"
	"strconv"
	"strings"
)

const (
	iaddr_size int = 4096
	daddr_size int = 4096
	no_regs    int = 8
	pc_reg     int = 7
)

type opclass int

const (
	opcLRR opclass = 1 + iota // reg operands r,s,t
	opclRM                    // reg r, mem d+s
	opclRA                    // reg r, int d+s
)

type opcode int

const (
	// RR instructions
	opHALT opcode = 1 + iota // RR     halt, operands are ignored
	opPRNT                   // RR     print, print operant to console
	opIN                     // RR     read into reg(r); s and t are ignored
	opOUT                    // RR     write from reg(r), s and t are ignored
	opADD                    // RR     reg(r) = reg(s)+reg(t)
	opSUB                    // RR     reg(r) = reg(s)-reg(t)
	opMUL                    // RR     reg(r) = reg(s)*reg(t)
	opDIV                    // RR     reg(r) = reg(s)/reg(t)

	// RM instructions
	opLD // RM     reg(r) = mem(d+reg(s))
	opST // RM     mem(d+reg(s)) = reg(r)

	// RA instructions
	opLDA // RA     reg(r) = d+reg(s)
	opLDC // RA     reg(r) = d ; reg(s) is ignored
	opJLT // RA     if reg(r)<0 then reg(7) = d+reg(s)
	opJLE // RA     if reg(r)<=0 then reg(7) = d+reg(s)
	opJGT // RA     if reg(r)>0 then reg(7) = d+reg(s)
	opJGE // RA     if reg(r)>=0 then reg(7) = d+reg(s)
	opJEQ // RA     if reg(r)==0 then reg(7) = d+reg(s)
	opJNE // RA     if reg(r)!=0 then reg(7) = d+reg(s)
)

type stepRESULT int

const (
	srOKAY stepRESULT = 1 + iota
	srHALT
	srIMEM_ERR
	srDMEM_ERR
	srZERODIVIDE
)

type instruction struct {
	iop    opcode
	iarg1  int
	iarg2  int
	iarg3  int
	iargs1 string
}

type vmMem struct {
	iMem [iaddr_size]instruction
	dMem [daddr_size]int
	reg  [no_regs]int
}

func (vm *vmMem) loadCode(code []string) bool {
	var (
		lineNo           int = 0
		loc              int
		op               opcode
		arg1, arg2, arg3 int
		args1            string
		err              error
		ok               bool
		opcodeMap        map[string]opcode = map[string]opcode{
			// RR opcodes
			"HALT":  opHALT,
			"PRINT": opPRNT,
			"IN":    opIN,
			"OUT":   opOUT,
			"ADD":   opADD,
			"SUB":   opSUB,
			"MUL":   opMUL,
			"DIV":   opDIV,

			// RM instructions
			"LD": opLD,
			"ST": opST,

			// RA instructions
			"LDA": opLDA,
			"LDC": opLDC,
			"JLT": opJLT,
			"JLE": opJLE,
			"JGT": opJGT,
			"JGE": opJGE,
			"JEQ": opJEQ,
			"JNE": opJNE,
		}
	)

	for _, inst := range code {
		lineNo++

		instSlice := strings.Split(strings.Trim(inst, " "), ":")
		if len(instSlice) < 2 {
			fmt.Printf(locale.Locale.VmMissingColonError, lineNo)
			return false
		}

		loc, err = strconv.Atoi(strings.Trim(instSlice[0], " "))
		if err != nil {
			fmt.Printf(locale.Locale.VmMemoryLocationError, instSlice[0], lineNo)
			return false
		}
		if loc > iaddr_size {
			fmt.Printf(locale.Locale.VmMemoryToLargeError, loc, lineNo)
			return false
		}

		opValue := strings.Trim(instSlice[1], " ")
		opIndex := strings.Index(opValue, " ")
		if opIndex == -1 {
			fmt.Printf(locale.Locale.VmMissingOpcodeError, loc, lineNo)
			return false
		}

		opCodeKey := opValue[0:opIndex]
		args := strings.Trim(opValue[opIndex+1:len(opValue)], " ")
		op, ok = opcodeMap[opCodeKey]
		if !ok {
			fmt.Println(inst)
			fmt.Println(opCodeKey)
			fmt.Printf(locale.Locale.VmInvalidOpcodeError, loc, lineNo)
			return false
		}

		switch op {
		case opHALT, opIN, opOUT, opADD, opSUB, opMUL, opDIV:
			argsSlice := strings.Split(args, ",")
			if len(argsSlice) != 3 {
				fmt.Printf(locale.Locale.VmInvalidNumberOfArgumentsError, loc, lineNo)
				return false
			}

			arg1, err = strconv.Atoi(strings.Trim(argsSlice[0], " "))
			if err != nil {
				fmt.Printf(locale.Locale.VmInvalidFirstArgumentError, loc, lineNo)
				return false
			}

			arg2, err = strconv.Atoi(strings.Trim(argsSlice[1], " "))
			if err != nil {
				fmt.Printf(locale.Locale.VmInvalidSecondArgumentError, loc, lineNo)
				return false
			}

			arg3, err = strconv.Atoi(strings.Trim(argsSlice[2], " "))
			if err != nil {
				fmt.Printf(locale.Locale.VmInvalidThirdArgumentError, loc, lineNo)
				return false
			}

			vm.iMem[loc].iop = op
			vm.iMem[loc].iarg1 = arg1
			vm.iMem[loc].iarg2 = arg2
			vm.iMem[loc].iarg3 = arg3
		case opLD, opST, opLDA, opLDC, opJLT, opJLE, opJGT, opJGE, opJEQ, opJNE:
			argsSlice1 := strings.Split(args, ",")
			if len(argsSlice1) != 2 {
				fmt.Printf(locale.Locale.VmInvalidNumberOfArgumentsError, loc, lineNo)
				return false
			}

			argsSlice2 := strings.Split(argsSlice1[1], "(")
			if len(argsSlice2) != 2 {
				fmt.Printf(locale.Locale.VmInvalidNumberOfArgumentsError, loc, lineNo)
				return false
			}

			arg1, err = strconv.Atoi(strings.Trim(argsSlice1[0], " "))
			if err != nil {
				fmt.Printf(locale.Locale.VmInvalidFirstArgumentError, loc, lineNo)
				return false
			}

			arg2, err = strconv.Atoi(strings.Trim(argsSlice2[0], " "))
			if err != nil {
				fmt.Printf(locale.Locale.VmInvalidSecondArgumentError, loc, lineNo)
				return false
			}

			arg3, err = strconv.Atoi(strings.Trim(argsSlice2[1], ")"))
			if err != nil {
				fmt.Printf(locale.Locale.VmInvalidThirdArgumentError, loc, lineNo)
				return false
			}
		case opPRNT:
			args1 = args
		}

		vm.iMem[loc].iop = op
		vm.iMem[loc].iarg1 = arg1
		vm.iMem[loc].iarg2 = arg2
		vm.iMem[loc].iarg3 = arg3
		vm.iMem[loc].iargs1 = args1
	}
	return true
}

func (vm *vmMem) executeCode() {
	var execute bool = true

	for execute {
		var r, s, t, m int = 0, 0, 0, 0
		var str string = ""
		pc := vm.reg[pc_reg]
		if pc < 0 || pc > iaddr_size {
			fmt.Printf(locale.Locale.VmInvalidProgramCounterError, pc)
			return
		}

		vm.reg[pc_reg] = pc + 1
		inst := vm.iMem[pc]

		//Setup instruction arguments
		switch inst.iop {
		case opHALT, opIN, opOUT, opADD, opSUB, opMUL, opDIV:
			r = inst.iarg1
			s = inst.iarg2
			t = inst.iarg3
		case opLD, opST:
			r = inst.iarg1
			s = inst.iarg3
			m = inst.iarg2 + vm.reg[s]

			if m < 0 || m > daddr_size {
				fmt.Printf(locale.Locale.VmInvalidMemoryAddressError, m)
				return
			}
		case opLDA, opLDC, opJLT, opJLE, opJGT, opJGE, opJEQ, opJNE:
			r = inst.iarg1
			s = inst.iarg3
			m = inst.iarg2 + vm.reg[s]
		case opPRNT:
			str = inst.iargs1
		}

		//Execute instruction
		switch inst.iop {
		case opHALT:
			return
		case opPRNT:
			fmt.Println(str)
		case opIN:
			var num int = 0
			_, err := fmt.Scanf("%d", &num)
			if err != nil {
				fmt.Println(locale.Locale.VmNonIntegerEnteredError)
				return
			}
			vm.reg[r] = num
		case opOUT:
			fmt.Println(vm.reg[r])
		case opADD:
			vm.reg[r] = vm.reg[s] + vm.reg[t]
		case opSUB:
			vm.reg[r] = vm.reg[s] - vm.reg[t]
		case opMUL:
			vm.reg[r] = vm.reg[s] * vm.reg[t]
		case opDIV:
			if vm.reg[t] == 0 {
				fmt.Println(locale.Locale.VmDivisionWIthZeroError)
				return
			}
			vm.reg[r] = vm.reg[s] / vm.reg[t]
		case opLD:
			vm.reg[r] = vm.dMem[m]
		case opST:
			vm.dMem[m] = vm.reg[r]
		case opLDA:
			vm.reg[r] = m
		case opLDC:
			vm.reg[r] = inst.iarg2
		case opJLT:
			if vm.reg[r] < 0 {
				vm.reg[pc_reg] = m
			}
		case opJLE:
			if vm.reg[r] <= 0 {
				vm.reg[pc_reg] = m
			}
		case opJGT:
			if vm.reg[r] > 0 {
				vm.reg[pc_reg] = m
			}
		case opJGE:
			if vm.reg[r] >= 0 {
				vm.reg[pc_reg] = m
			}
		case opJEQ:
			if vm.reg[r] == 0 {
				vm.reg[pc_reg] = m
			}
		case opJNE:
			if vm.reg[r] != 0 {
				vm.reg[pc_reg] = m
			}
		}
	}
}

func Execute(code []string) {
	vm := new(vmMem)
	vm.dMem[0] = daddr_size - 1

	if !vm.loadCode(code) {
		return
	}

	vm.executeCode()
}
