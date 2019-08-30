package main

import (
  "fmt"
  "io/ioutil"
  "strconv"
  "errors"
)

type Token struct {
  kind string // "intliteral", "punct"
  value string
}

var source []byte
var sourceIndex = 0

func getChar() (byte, error) {
  if sourceIndex == len(source) {
    return 0, errors.New("EOF")
  }
  char := source[sourceIndex]
  sourceIndex++
  return char,nil
}

func ungetChar() {
  sourceIndex--
}

func readNumber(char byte) string {
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

  return string(number)
}

func tokenize() []*Token {
  fmt.Printf("# tokens = ")
  var tokens []*Token
  for {
    char , err := getChar()
    if err != nil {
      break
    }
    switch char {
    case ' ','\n':
      continue
    case '0','1','2','3','4','5','6','7','8','9':
      number := readNumber(char)
      token := &Token{
        kind:  "intliteral",
        value: number,
      }
      tokens = append(tokens, token)
      fmt.Printf(" '%s' ", token.value)
    case '+','-',';':
      token := &Token{
        kind:  "punct",
        value: string([]byte{char}),
      }
      tokens = append(tokens, token)
      fmt.Printf(" '%s' ", token.value)
    default:
      panic(fmt.Sprintf("Invalid char '%c'", char))
    }
  }
  fmt.Printf("\n")
  return tokens
}

// "-1"
// " 2 + 3 "
type Expr struct {
  kind string  // "intliteral", "unary", "binary"
  intval int
  operator string // "+", "-"
  operand *Expr
  left *Expr
  right *Expr
}

var tokens []*Token
var tokenIndex = 0

func getToken() *Token {
  if tokenIndex == len(tokens) {
    return nil
  }
  token := tokens[tokenIndex]
  tokenIndex++
  return token
}

func parseUnaryExpr() *Expr {
  token := getToken()
  switch token.kind {
  case "intliteral":
    number , err := strconv.Atoi(token.value)
    if err != nil {
      panic(err)
    }
    return &Expr{
      kind:   "intliteral",
      intval: number,
    }
  case "punct":
    return &Expr{
      kind:     "unary",
      operator: token.value,
      operand:  parse(),
    }
  default:
    panic("Unexpected token.kind")
  }
}

func parse() *Expr {
  expr := parseUnaryExpr()
  token := getToken()
  if token == nil || token.value == ";" {
    return expr
  }
  switch token.value {
  case "+","-":
    return &Expr{
      kind:     "binary",
      operator: token.value ,
      left:     expr,
      right:    parseUnaryExpr(),
    }
  default:
    panic("Unexpected token.value")
  }
}

func generateExpr(expr *Expr) {
  switch expr.kind {
  case "intliteral":
    fmt.Printf("  mov $%d, %%rax\n", expr.intval)
  case "unary":
    switch expr.operator {
    case "+":
      fmt.Printf("  mov $%d, %%rax\n", expr.operand.intval)
    case "-":
      fmt.Printf("  mov $-%d, %%rax\n", expr.operand.intval)
    }
  case "binary":
    switch expr.operator {
    case "+":
      fmt.Printf("  mov $%d, %%rax\n", expr.left.intval)
      fmt.Printf("  mov $%d, %%rcx\n", expr.right.intval)
      fmt.Printf("  add %%rcx, %%rax\n")
    case "-":
      fmt.Printf("  mov $%d, %%rax\n", expr.left.intval)
      fmt.Printf("  mov $%d, %%rcx\n", expr.right.intval)
      fmt.Printf("  sub %%rcx, %%rax\n")
    }
  default:
    panic("Unexpected expr.kind")
  }
}

func main() {
  source , _ = ioutil.ReadFile("/dev/stdin")
  tokens = tokenize()
  expr := parse()

  fmt.Printf("\n")
  fmt.Printf("  .global main\n")
  fmt.Printf("main:\n")
  generateExpr(expr)
  fmt.Printf("  ret\n")

}
