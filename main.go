package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
)

type Token struct {
	kind  string // "intliteral", "punct"
	value string
}

var source []byte
var sourceIndex int = 0

func getChar() (byte, error) {
	if sourceIndex == len(source) {
		return 0, errors.New("EOF")
	}
	char := source[sourceIndex]
	sourceIndex++
	return char, nil
}

func ungetChar() {
	sourceIndex--
}

func tokenize() []*Token {
	var tokens []*Token
	for {
		char, err := getChar()
		if err != nil {
			break
		}
		var token *Token
		switch char {
		case ' ', '\t', '\n':
			continue
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			var number []byte = []byte{char}
			for {
				char, err := getChar()
				if err != nil {
					break
				}
				if '0' <= char && char <= '9' {
					number = append(number, char)
				} else {
					ungetChar()
					break
				}
			}
			token = &Token{
				kind:  "intliteral",
				value: string(number),
			}
		case ';':
			token = &Token{
				kind: "punct",
				value: string([]byte{char}),
			}
		default:
			panic(fmt.Sprintf("tokenizer: Invalid char: '%c'", char))
		}

		tokens = append(tokens, token)
	}

	return tokens
}

func main() {
	var err error
	source, _ = ioutil.ReadFile("/dev/stdin")
	tokens := tokenize()
	token0 := tokens[0]
	number, err := strconv.Atoi(token0.value)
	if err != nil {
		panic(err)
	}
	fmt.Printf("  .global main\n")
	fmt.Printf("main:\n")
	fmt.Printf("  movq $%d, %%rax\n", number)
	fmt.Printf("  ret\n")
}
