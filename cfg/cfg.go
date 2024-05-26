/*
The MIT License (MIT)

Copyright (c) 2016-2024 Ivan Dejanovic

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
	"encoding/json"
	"fmt"
	"github.com/ivandejanovic/mlpl/locale"
	"io/ioutil"
	"os"
	"strings"
)

const (
	minus       = "-"
	doubleMinus = "--"
	empty       = ""
	usage       = "Usage: mlpl <codefilename> [configurationfilename]"
)

func getLocaleFromConfig(configFile string) {
	config, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(err)
	}

	json.Unmarshal(config, locale.Locale)
	locale.AssembleReserved()
}

func HandleArgs() (bool, string) {
	var abort bool = true
	var codeFile string

	args := os.Args[1:]
	argc := len(args)

	for index := 0; index < argc; index++ {
		var flag string = empty
		var flagArg string = args[index]

		if strings.HasPrefix(flagArg, doubleMinus) {
			flag = strings.TrimPrefix(flagArg, doubleMinus)
		} else if strings.HasPrefix(flagArg, minus) {
			flag = strings.TrimPrefix(flagArg, minus)
		}

		if flag != empty {
			switch flag {
			case "h", "help":
				fmt.Println()
				fmt.Println(usage)
				fmt.Println()
				fmt.Println("Options:")
				fmt.Println("  -h, --help       Prints help")
				fmt.Println("  -v, --version    Prints version")
			case "v", "version":
				fmt.Println("MLPL interpreter version 1.1.1")
			default:
				fmt.Println("Invalid usage. For correct usage examples please try: mlpl -h")
			}
			return abort, codeFile
		}
	}

	if argc < 1 || argc > 2 {
		fmt.Println(usage)
		return abort, codeFile
	}

	if argc == 2 {
		getLocaleFromConfig(args[1])
	} else {
		locale.AssembleReserved()
	}

	//If we get this far we have good data to process
	abort = false
	codeFile = args[0]

	return abort, codeFile
}
