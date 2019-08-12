package main

import (
	"fmt"
	"strconv"
	"io/ioutil"
	)

type Token struct {
	Type string // "numberliteral"
	Value string
}

func tokenize(bytes []byte) []*Token {
	var tokens []*Token

	var token *Token
	char := bytes[0]
	if '0' <= char && char <= '9' {
		token = &Token{
			Type:"numberliteral",
			Value: string(bytes),
		}
	}

	tokens = append(tokens, token)
	return tokens
}

func main() {
	var err error
	bytes, _ := ioutil.ReadFile("/dev/stdin")
	tokens := tokenize(bytes)

	token := tokens[0]
	number, err := strconv.Atoi(token.Value)
	if err != nil {
		panic(err)
	}
	fmt.Printf("  .global main\n")
	fmt.Printf("main:\n")
	fmt.Printf("  movq $%d, %%rax\n", number)
	fmt.Printf("  ret\n")
}
