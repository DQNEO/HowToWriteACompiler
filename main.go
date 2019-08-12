package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
)

func main() {
	var source []byte
	source, _ = ioutil.ReadFile("/dev/stdin")
	var input string = string(source)

	number, err := strconv.Atoi(input)
	if err != nil {
		panic(err)
	}
	fmt.Printf("  .global main\n")
	fmt.Printf("main:\n")
	fmt.Printf("  movq $%d, %%rax\n", number)
	fmt.Printf("  ret\n")
}
