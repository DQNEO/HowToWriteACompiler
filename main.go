package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
)

type Token struct {
	Type string // "numberliteral"
	Value string
}

var sourceIndex = 0

func getchar() (byte, error) {
	if len(bytes) == sourceIndex {
		return 0, errors.New("EOF")
	}
	char := bytes[sourceIndex]
	sourceIndex++
	return char, nil
}

func ungethar() {
	sourceIndex--
}

func tokenize() []*Token {
	var tokens []*Token

	for {
		char, err := getchar()
		if err != nil {
			break
		}
		switch char {
		case '0','1','2','3','4','5','6','7','8','9':
			var number []byte = []byte{char}
			for {
				char, err := getchar()
				if err != nil {
					break
				}
				if '0' <= char && char <= '9' {
					number = append(number, char)
				} else {
					ungethar()
					break
				}
			}
			token := &Token{
				Type:"numberliteral",
				Value: string(number),
			}
			tokens = append(tokens, token)
		case ' ', '\t','\n':
			continue
		case ';':
			token := &Token{
				Type:"punctuation",
				Value: string([]byte{char}),
			}
			tokens = append(tokens, token)
		}

	}


	return tokens
}

var bytes []byte
func main() {
	var err error
	bytes, _ = ioutil.ReadFile("/dev/stdin")
	tokens := tokenize()

	token := tokens[0]
	number, err := strconv.Atoi(token.Value)
	if err != nil {
		panic(err)
	}
	fmt.Printf(" .global main\n")
	fmt.Printf("main:\n")
	fmt.Printf("  movq $%d, %%rax\n", number)
	fmt.Printf("  ret\n")
}
