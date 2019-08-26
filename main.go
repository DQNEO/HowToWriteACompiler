package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
)

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

type Token struct {
	kind  string // "intliteral", "punct"
	value string
}

func readNumber(char byte) string {
	number := []byte{char}
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

	return string(number)
}

func tokenize() []*Token {
	var tokens []*Token
	fmt.Printf("# Tokens : ")

	for {
		char, err := getChar()
		if err != nil {
			break
		}
		switch char {
		case ' ', '\t', '\n':
			continue
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			intliteral := readNumber(char)
			token := &Token{
				kind:  "intliteral",
				value: intliteral,
			}
			tokens = append(tokens, token)
			fmt.Printf(" '%s'", token.value)
		case ';', '+', '-', '*', '/':
			token := &Token{
				kind:  "punct",
				value: string([]byte{char}),
			}
			tokens = append(tokens, token)
			fmt.Printf(" '%s'", token.value)
		default:
			panic(fmt.Sprintf("tokenizer: Invalid char: '%c'", char))
		}
	}

	fmt.Printf("\n")
	return tokens
}

var tokens []*Token
var tokenIndex int = 0

func getToken() *Token {
	if tokenIndex == len(tokens) {
		return nil
	}
	token := tokens[tokenIndex]
	tokenIndex++
	return token
}

type Expr struct {
	kind     string // "intliteral", "unary"
	intval   int    // for intliteral
	operator string // "-", "+", ...
	operand  *Expr  // for unary expr
	left     *Expr  // for binary expr
	right    *Expr  // for binary expr
}

func parseUnaryExpr() *Expr {
	token := getToken()
	switch token.kind {
	case "intliteral":
		intval, err := strconv.Atoi(token.value)
		if err != nil {
			panic(err)
		}
		return &Expr{
			kind:   "intliteral",
			intval: intval,
		}
	case "punct":
		operator := token.value
		operand := parseUnaryExpr()
		return &Expr{
			kind:     "unary",
			operator: operator,
			operand:  operand,
		}
	default:
		return nil
	}
}

func parse() *Expr {
	expr := parseUnaryExpr()

	for {
		token := getToken()
		if token == nil || token.value == ";" {
			return expr
		}

		switch token.value {
		case "+", "-", "*", "/":
			left := expr
			right := parseUnaryExpr()
			return &Expr{
				kind:     "binary",
				operator: token.value,
				left:     left,
				right:    right,
			}
		default:
			panic("unexpected token:" + token.value)
		}
	}
}

func generateExpr(expr *Expr) {
	switch expr.kind {
	case "intliteral":
		fmt.Printf("  movq $%d, %%rax\n", expr.intval)
	case "unary":
		switch expr.operator {
		case "-":
			fmt.Printf("  movq $-%d, %%rax\n", expr.operand.intval)
		case "+":
			fmt.Printf("  movq $%d, %%rax\n", expr.operand.intval)
		default:
			panic("generator: Unknown unary operator:" + expr.operator)
		}
	case "binary":
		fmt.Printf("  movq $%d, %%rax\n", expr.left.intval)
		fmt.Printf("  movq $%d, %%rcx\n", expr.right.intval)

		switch expr.operator {
		case "+":
			fmt.Printf("  addq %%rcx, %%rax\n")
		case "-":
			fmt.Printf("  subq %%rcx, %%rax\n")
		case "*":
			fmt.Printf("  imulq %%rcx, %%rax\n")
		case "/":
			fmt.Printf("  movq $0, %%rdx\n")
			fmt.Printf("  idiv %%rcx\n")
		default:
			panic("generator: Unknown binary operator:" + expr.operator)
		}
	default:
		panic("generator: Unknown expr.kind:" + expr.kind)
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
