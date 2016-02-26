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

package cfg

import (
	"bufio"
	"fmt"
	"mlpl/types"
	"os"
)

func GetDefaultReserved() []types.ReservedWord {
	reserved := make([]types.ReservedWord, 0, 8)

	reserved = append(reserved, types.ReservedWord{types.IF, "if"})
	reserved = append(reserved, types.ReservedWord{types.THEN, "then"})
	reserved = append(reserved, types.ReservedWord{types.ELSE, "else"})
	reserved = append(reserved, types.ReservedWord{types.END, "end"})
	reserved = append(reserved, types.ReservedWord{types.REPEAT, "repeat"})
	reserved = append(reserved, types.ReservedWord{types.UNTIL, "until"})
	reserved = append(reserved, types.ReservedWord{types.READ, "read"})
	reserved = append(reserved, types.ReservedWord{types.WRITE, "write"})

	return reserved
}

func GetConfigReservedWords(configFile string) []types.ReservedWord {
	reserved := make([]types.ReservedWord, 0, 8)
	var localization []string
	const length = 8

	config, err := os.Open(configFile)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(config)
	for scanner.Scan() {
		localization = append(localization, scanner.Text())
	}

	defer config.Close()

	if len(localization) != length {
		fmt.Println("Configuration file must contain localizations for eight key word.")
	}

	reserved = append(reserved, types.ReservedWord{types.IF, localization[0]})
	reserved = append(reserved, types.ReservedWord{types.THEN, localization[1]})
	reserved = append(reserved, types.ReservedWord{types.ELSE, localization[2]})
	reserved = append(reserved, types.ReservedWord{types.END, localization[3]})
	reserved = append(reserved, types.ReservedWord{types.REPEAT, localization[4]})
	reserved = append(reserved, types.ReservedWord{types.UNTIL, localization[5]})
	reserved = append(reserved, types.ReservedWord{types.READ, localization[6]})
	reserved = append(reserved, types.ReservedWord{types.WRITE, localization[7]})

	return reserved
}
