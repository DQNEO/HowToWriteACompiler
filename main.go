package main

import (
	"fmt"
	"strconv"
	"io/ioutil"
	)

func main() {
	var err error
	bytes, _ := ioutil.ReadFile("/dev/stdin")
	var input string = string(bytes)

	number, err := strconv.Atoi(input)
	if err != nil {
		panic(err)
	}
	fmt.Printf("  .global main\n")
	fmt.Printf("main:\n")
	fmt.Printf("  movq $%d, %%rax\n", number)
	fmt.Printf("  ret\n")
}
