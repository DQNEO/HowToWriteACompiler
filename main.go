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
		case ';','+','-':
			token := &Token{
				Type:"punctuation",
				Value: string([]byte{char}),
			}
			tokens = append(tokens, token)
		case ' ', '\t','\n':
			continue
		default:
			panic(fmt.Sprintf("Invalid char: %c", char))
		}

	}


	return tokens
}

var bytes []byte
var tokens []*Token
var tokenIndex int = 0

// Node is an expression
type Node struct {
	Type string // "intliteral", "unary"
	intval int
	operator string
	operand *Node
}

func getToken() *Token {
	if tokenIndex >= len(tokens ) {
		return nil
	}
	token := tokens[tokenIndex]
	tokenIndex++
	return token
}

func parseUnaryExpr() *Node {
	token := getToken()
	if token.Type == "numberliteral" {
		intval, _ := strconv.Atoi(token.Value)
		return &Node{
			Type:   "intliteral",
			intval: intval,
		}
	} else if token.Type == "punctuation" {
		operand := parseUnaryExpr()
		return &Node{
			Type: "unary",
			operator:token.Value,
			operand: operand,
		}
	}

	return nil
}

func parseExpr() *Node {
	node := parseUnaryExpr()
	return node
}

func generateExpression(node *Node) {
	if node.Type == "intliteral" {
		fmt.Printf("  movq $%d, %%rax # %s\n", node.intval, node.Type)
	} else if node.Type == "unary" {
		if node.operator == "-" {
			fmt.Printf("  movq $-%d, %%rax # %s\n", node.operand.intval, node.operand.Type)
		} else {
			fmt.Printf("  movq $%d, %%rax # %s\n", node.operand.intval, node.operand.Type)
		}
	}
}

func generateCode(node *Node) {
	fmt.Printf("  .global main\n")
	fmt.Printf("main:\n")
	generateExpression(node)
	fmt.Printf("  ret\n")
}

func main() {
	bytes, _ = ioutil.ReadFile("/dev/stdin")
	tokens = tokenize()
	node := parseExpr()
	generateCode(node)
}
