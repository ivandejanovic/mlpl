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

package locale

import (
	"errors"
	"fmt"
	"mlpl/types"
)

type LocaleType struct {
	ReservedArray []string
	Reserved      []types.ReservedWord

	ParseError string

	LexerSyntaxError       string
	LexerReservedWordError string
	LexerAssignError       string
	LexerLTError           string
	LexerEQError           string
	LexerLPARENError       string
	LexerRPARENError       string
	LexerSEMIError         string
	LexerPLUSError         string
	LexerMINUSError        string
	LexerTIMESError        string
	LexerOVERError         string
	LexerENDFILEError      string
	LexerNUMError          string
	LexerIDError           string
	LexerERRORError        string
	LexerDEFAULTError      string
	LexerABORTINGError     string

	AnalyzeTypePrefixError string
	AnalyzeTypeOpError     string
	AnalyzeTypeIfError     string
	AnalyzeTypeAssignError string
	AnalyzeTypeWriteError  string
	AnalyzeTypeRepeatError string

	CodegenUnknownOperatorError string
	CodegenUnknownTypeError     string

	VmMissingColonError             string
	VmMemoryLocationError           string
	VmMemoryToLargeError            string
	VmMissingOpcodeError            string
	VmInvalidOpcodeError            string
	VmInvalidNumberOfArgumentsError string
	VmInvalidFirstArgumentError     string
	VmInvalidSecondArgumentError    string
	VmInvalidThirdArgumentError     string
	VmInvalidProgramCounterError    string
	VmInvalidMemoryAddressError     string
	VmNonIntegerEnteredError        string
	VmDivisionWIthZeroError         string
}

var Locale *LocaleType = new(LocaleType)

const ReservedLength int = 8

func init() {
	var reserved []string

	reserved = append(reserved, "if")
	reserved = append(reserved, "then")
	reserved = append(reserved, "else")
	reserved = append(reserved, "end")
	reserved = append(reserved, "repeat")
	reserved = append(reserved, "until")
	reserved = append(reserved, "read")
	reserved = append(reserved, "write")

	Locale.ReservedArray = reserved

	Locale.ParseError = "Scanner bug: state= %d\n"

	Locale.LexerSyntaxError = "Syntax error at line %d, unexpected token -> "
	Locale.LexerReservedWordError = "reserved word: %s\n"
	Locale.LexerAssignError = ":=\n"
	Locale.LexerLTError = "<\n"
	Locale.LexerEQError = "=\n"
	Locale.LexerLPARENError = "(\n"
	Locale.LexerRPARENError = ")\n"
	Locale.LexerSEMIError = ";\n"
	Locale.LexerPLUSError = "+\n"
	Locale.LexerMINUSError = "-\n"
	Locale.LexerTIMESError = "*\n"
	Locale.LexerOVERError = "/\n"
	Locale.LexerENDFILEError = "EOF\n"
	Locale.LexerNUMError = "NUM, name= %s\n"
	Locale.LexerIDError = "ID, name= %s\n"
	Locale.LexerERRORError = "ERROR: %s\n"
	Locale.LexerDEFAULTError = "Unknown token: %d\n"
	Locale.LexerABORTINGError = "Aborting\n"

	Locale.AnalyzeTypePrefixError = "Type error at line %d: %s\n"
	Locale.AnalyzeTypeOpError = "Op applied to non-integer"
	Locale.AnalyzeTypeIfError = "if test is not Boolean"
	Locale.AnalyzeTypeAssignError = "assignment of non-integer value"
	Locale.AnalyzeTypeWriteError = "write of non-integer or non-string value"
	Locale.AnalyzeTypeRepeatError = "repeat test is not Boolean"

	Locale.CodegenUnknownOperatorError = "Unknown operator for code generation"
	Locale.CodegenUnknownTypeError = "Unknown type for code generation"

	Locale.VmMissingColonError = "Missing colon on line: %d\n"
	Locale.VmMemoryLocationError = "Invalid memory location %s on line: %d\n"
	Locale.VmMemoryToLargeError = "To large memory location %d on line: %d\n"
	Locale.VmMissingOpcodeError = "Missing opcode on location %d and line: %d\n"
	Locale.VmInvalidOpcodeError = "Invalid opcode on location %d and line: %d\n"
	Locale.VmInvalidNumberOfArgumentsError = "Invalid number of arguments on location %d and line: %d\n"
	Locale.VmInvalidFirstArgumentError = "Invalid first argument on location %d and line: %d\n"
	Locale.VmInvalidSecondArgumentError = "Invalid second argument on location %d and line: %d\n"
	Locale.VmInvalidThirdArgumentError = "Invalid third argument on location %d and line: %d\n"
	Locale.VmInvalidProgramCounterError = "Invalid program counter value: %d\n"
	Locale.VmInvalidMemoryAddressError = "Invalid memory address value: %d\n"
	Locale.VmNonIntegerEnteredError = "Non integer entered."
	Locale.VmDivisionWIthZeroError = "Division with zero."
}

func AssembleReserved() {
	if len(Locale.ReservedArray) != ReservedLength {
		errorMessage := fmt.Sprintf("Configuration file must contain localizations for eight key word.\n")
		err := errors.New(errorMessage)
		panic(err)
	}

	reserved := make([]types.ReservedWord, 0, ReservedLength)

	reserved = append(reserved, types.ReservedWord{types.IF, Locale.ReservedArray[0]})
	reserved = append(reserved, types.ReservedWord{types.THEN, Locale.ReservedArray[1]})
	reserved = append(reserved, types.ReservedWord{types.ELSE, Locale.ReservedArray[2]})
	reserved = append(reserved, types.ReservedWord{types.END, Locale.ReservedArray[3]})
	reserved = append(reserved, types.ReservedWord{types.REPEAT, Locale.ReservedArray[4]})
	reserved = append(reserved, types.ReservedWord{types.UNTIL, Locale.ReservedArray[5]})
	reserved = append(reserved, types.ReservedWord{types.READ, Locale.ReservedArray[6]})
	reserved = append(reserved, types.ReservedWord{types.WRITE, Locale.ReservedArray[7]})

	Locale.Reserved = reserved
}
