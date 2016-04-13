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

package parse

import (
	"bufio"
	"fmt"
	"io"
	"mlpl/types"
	"os"
	"unicode"
)

const (
	newLine    rune = '\n'
	colon      rune = ':'
	space      rune = ' '
	tab        rune = '\t'
	numberSign rune = '#'
	equal      rune = '='
	lt         rune = '<'
	plus       rune = '+'
	minus      rune = '-'
	times      rune = '*'
	over       rune = '/'
	lParen     rune = '('
	rParen     rune = ')'
	semi       rune = ';'
)

type state int

const (
	start state = 1 + iota
	inAssign
	inComment
	inNum
	inId
	done
)

func reservedLookup(s string, reserved []types.ReservedWord) types.TokenType {
	for index := 0; index < len(reserved); index++ {
		if s == reserved[index].Str {
			return reserved[index].TokenType
		}
	}

	return types.ID
}

func Parse(sourceFile string, reserved []types.ReservedWord) []types.Token {
	var tokens []types.Token
	lineno := 0

	source, err := os.Open(sourceFile)
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(source)

	getToken := func(reserved []types.ReservedWord) types.Token {
		var currentToken types.TokenType
		var currentTokenString string
		var currentTokenRunes []rune

		for state := start; state != done; {
			save := true
			r, _, err := reader.ReadRune()
			if err != nil && err != io.EOF {
				panic(err)
			}

			if r == newLine {
				lineno++
			}
			switch state {
			case start:
				if unicode.IsDigit(r) {
					state = inNum
				} else if unicode.IsLetter(r) {
					state = inId
				} else if r == colon {
					state = inAssign
				} else if r == space || r == tab || r == newLine {
					save = false
				} else if r == numberSign {
					save = false
					state = inComment
				} else {
					state = done
					if err == io.EOF {
						save = false
						currentToken = types.ENDFILE
					} else {
						switch r {
						case equal:
							currentToken = types.EQ
						case lt:
							currentToken = types.LT
						case plus:
							currentToken = types.PLUS
						case minus:
							currentToken = types.MINUS
						case times:
							currentToken = types.TIMES
						case over:
							currentToken = types.OVER
						case lParen:
							currentToken = types.LPAREN
						case rParen:
							currentToken = types.RPAREN
						case semi:
							currentToken = types.SEMI
						default:
							currentToken = types.ERROR
						}
					}
				}
			case inComment:
				save = false
				if err == io.EOF {
					state = done
					currentToken = types.ENDFILE
				} else if r == numberSign {
					state = start
				}
			case inAssign:
				state = done
				if r == equal {
					currentToken = types.ASSIGN
				} else {
					err = reader.UnreadRune()
					if err != nil {
						panic(err)
					}
					save = false
					currentToken = types.ERROR
				}
			case inNum:
				if !unicode.IsDigit(r) {
					err = reader.UnreadRune()
					if err != nil {
						panic(err)
					}
					save = false
					state = done
					currentToken = types.NUM
				}
			case inId:
				if !unicode.IsLetter(r) {
					err = reader.UnreadRune()
					if err != nil {
						panic(err)
					}
					save = false
					state = done
					currentToken = types.ID
				}
			case done:
				//Should never happen
				fmt.Printf("Scanner bug: state= %d\n", state)
				state = done
				currentToken = types.ERROR
			default:
				//Should never happen
				fmt.Printf("Scanner bug: state= %d\n", state)
				state = done
				currentToken = types.ERROR
			}

			if save {
				currentTokenRunes = append(currentTokenRunes, r)
			}
			if state == done {
				currentTokenString = string(currentTokenRunes)
				if currentToken == types.ID {
					currentToken = reservedLookup(currentTokenString, reserved)
				}
			}
		}

		return types.Token{currentToken, currentTokenString, lineno}
	}

	for moreTokens := true; moreTokens; {
		token := getToken(reserved)
		tokens = append(tokens, token)

		if token.TokenType == types.ENDFILE {
			moreTokens = false
		}
	}

	defer source.Close()

	return tokens
}
