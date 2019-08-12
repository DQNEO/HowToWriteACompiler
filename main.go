package main

import "fmt"

func main() {
	fmt.Printf("  .global main\n")
	fmt.Printf("main:\n")
	fmt.Printf("  movq $42, %%rax\n")
	fmt.Printf("  ret\n")
}
