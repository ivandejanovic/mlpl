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
