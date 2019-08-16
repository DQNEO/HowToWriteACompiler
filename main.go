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
	fmt.Printf("# Tokens : ")

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
		case ';', '+', '-', '*':
			token = &Token{
				kind:  "punct",
				value: string([]byte{char}),
			}
		default:
			panic(fmt.Sprintf("Invalid char: '%c'", char))
		}

		fmt.Printf(" '%s'", token.value)
		tokens = append(tokens, token)
	}

	fmt.Printf("\n")
	return tokens
}

var tokens []*Token
var tokenIndex int = 0

type Expr struct {
	kind     string // "intliteral", "unary"
	intval   int    // for intliteral
	operator string // "-", "+"
	operand  *Expr  // for unary expr
	left     *Expr  // for binary expr
	right    *Expr  // for binary expr
}

func getToken() *Token {
	if tokenIndex >= len(tokens) {
		return nil
	}
	token := tokens[tokenIndex]
	tokenIndex++
	return token
}

func ungetToken() {
	tokenIndex--
}

func parseUnaryExpr() *Expr {
	token := getToken()
	if token.kind == "intliteral" {
		intval, err := strconv.Atoi(token.value)
		if err != nil {
			panic(err)
		}
		return &Expr{
			kind:   "intliteral",
			intval: intval,
		}
	} else if token.kind == "punct" {
		operand := parseUnaryExpr()
		return &Expr{
			kind:     "unary",
			operator: token.value,
			operand:  operand,
		}
	}

	return nil
}

func parse() *Expr {
	expr := parseUnaryExpr()

	for {
		tok := getToken()
		if tok == nil || tok.value == ";" {
			return expr
		}

		if tok.value == "+" || tok.value == "-"  || tok.value == "*" {
			left := expr
			right := parseUnaryExpr()
			return &Expr{
				kind:     "binary",
				operator: tok.value,
				left:     left,
				right:    right,
			}
		}
	}

	return expr
}

func generateExpr(expr *Expr) {
	switch expr.kind {
	case "intliteral":
		fmt.Printf("  movq $%d, %%rax # %s\n", expr.intval, expr.kind)
	case "unary":
		if expr.operator == "-" {
			fmt.Printf("  movq $-%d, %%rax # %s\n", expr.operand.intval, expr.operand.kind)
		} else {
			fmt.Printf("  movq $%d, %%rax # %s\n", expr.operand.intval, expr.operand.kind)
		}
	case "binary":
		fmt.Printf("  movq $%d, %%rax # %s\n", expr.left.intval, expr.left.kind)
		fmt.Printf("  movq $%d, %%rbx # %s\n", expr.right.intval, expr.right.kind)
		switch expr.operator {
		case "+":
			fmt.Printf("  addq %%rbx, %%rax\n")
		case "-":
			fmt.Printf("  subq %%rbx, %%rax\n")
		case "*":
			fmt.Printf("  imulq %%rbx, %%rax\n")
		default:
			panic("generator: unknown operator:" + expr.operator)
		}
	default:
		panic("generator: unknown expr type:" + expr.kind)
	}
}

func generateCode(expr *Expr) {
	fmt.Printf("  .global main\n")
	fmt.Printf("main:\n")
	generateExpr(expr)
	fmt.Printf("  ret\n")
}

func main() {
	source, _ = ioutil.ReadFile("/dev/stdin")
	tokens = tokenize()
	expr := parse()
	generateCode(expr)
}
