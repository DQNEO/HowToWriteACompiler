package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
)

type Token struct {
	kind  string // "intliteral"
	value string
}

var source []byte
var sourceIndex = 0

func getchar() (byte, error) {
	if len(source) == sourceIndex {
		return 0, errors.New("EOF")
	}
	char := source[sourceIndex]
	sourceIndex++
	return char, nil
}

func ungethar() {
	sourceIndex--
}

func tokenize() []*Token {
	var tokens []*Token
	fmt.Printf("# Tokens:")

	for {
		char, err := getchar()
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
			token = &Token{
				kind:  "intliteral",
				value: string(number),
			}
		case ';', '+', '-', '*':
			token = &Token{
				kind:  "punctuation",
				value: string([]byte{char}),
			}
		default:
			panic(fmt.Sprintf("Invalid char: %c", char))
		}

		fmt.Printf(" \"%s\"", token.value)
		tokens = append(tokens, token)

	}

	fmt.Printf("\n")
	return tokens
}

var tokens []*Token
var tokenIndex int = 0

// Node is an expression
type Node struct {
	kind     string // "intliteral", "unary"
	intval   int
	operator string
	operand  *Node // for unary
	left     *Node // for binary
	right    *Node // for binary
}

func getToken() *Token {
	if tokenIndex >= len(tokens) {
		return nil
	}
	token := tokens[tokenIndex]
	tokenIndex++
	return token
}

func parseUnaryExpr() *Node {
	token := getToken()
	if token.kind == "intliteral" {
		intval, _ := strconv.Atoi(token.value)
		return &Node{
			kind:   "intliteral",
			intval: intval,
		}
	} else if token.kind == "punctuation" {
		operand := parseUnaryExpr()
		return &Node{
			kind:     "unary",
			operator: token.value,
			operand:  operand,
		}
	}

	return nil
}

func parseExpr() *Node {
	node := parseUnaryExpr()

	for {
		tok := getToken()
		if tok == nil || tok.value == ";" {
			return node
		}

		if tok.value == "+" || tok.value == "-"  || tok.value == "*" {
			left := node
			right := parseUnaryExpr()
			return &Node{
				kind:     "binary",
				operator: tok.value,
				left:     left,
				right:    right,
			}
		}
	}

	return node
}

func generateExpression(node *Node) {
	switch node.kind {
	case "intliteral":
		fmt.Printf("  movq $%d, %%rax # %s\n", node.intval, node.kind)
	case "unary":
		if node.operator == "-" {
			fmt.Printf("  movq $-%d, %%rax # %s\n", node.operand.intval, node.operand.kind)
		} else {
			fmt.Printf("  movq $%d, %%rax # %s\n", node.operand.intval, node.operand.kind)
		}
	case "binary":
		fmt.Printf("  movq $%d, %%rax # %s\n", node.left.intval, node.left.kind)
		fmt.Printf("  movq $%d, %%rbx # %s\n", node.right.intval, node.right.kind)
		switch node.operator {
		case "+":
			fmt.Printf("  addq %%rbx, %%rax\n")
		case "-":
			fmt.Printf("  subq %%rbx, %%rax\n")
		case "*":
			fmt.Printf("  imulq %%rbx, %%rax\n")
		default:
			panic("generator: unknown operator:" + node.operator)
		}
	default:
		panic("generator: unknown node type:" + node.kind)
	}
}

func generateCode(node *Node) {
	fmt.Printf("  .global main\n")
	fmt.Printf("main:\n")
	generateExpression(node)
	fmt.Printf("  ret\n")
}

func main() {
	source, _ = ioutil.ReadFile("/dev/stdin")
	tokens = tokenize()
	node := parseExpr()
	generateCode(node)
}
