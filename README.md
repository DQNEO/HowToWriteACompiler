# How to write a compiler from scratch

This repository explains how to write a compiler from scratch by Go.
The compiler has some constraints

* Can compile only arithmetic operations.
* Runs only on Linux
* Outputs X86-64 assembly (GAS)

# Usage

First you need to run a docker container and get in it.

```
./docker-run
```

And then you can use the compiler.


```
echo '30 + 12' | go run main.go
```

This program receives source code from stdin, and emit compiled code to stdout.

If you want to compile and run it immediately, you can use `asrun` script.

```
echo '30 + 12' | go run main.go | ./asrun
```

`asrun` takes assembly code from stdin, execute it and shows the code with the resulting status code.

```
$ echo  '30 + 12' | go run main.go | ./asrun
-------- a.s ----------------
  .global main
main:
  movq $30, %rax # intliteral
  movq $12, %rbx # intliteral
  addq %rbx, %rax
  ret
-------- result -------------
42
```

# Design

There are 3 phases in this compiler.

Source Code -> [Tokenizer] -> Tokens -> [Parser] -> AST -> [Code Generator] -> Assembly

## Tokenizer

Source Code -> [Tokenizer] -> Tokens

Tokenizer analyzes byte stream of source code, and break it down into a list of tokens.

In this compiler, `tokenize()` does this task.

## Parser

Tokens -> [Parser] -> AST

Parser analyzes stream of tokens, and composer a tree of nested structs , which represents sytanx structure of source code.

This tree is called a AST (Abstract Syntax Tree).

`parser()` does this task.

## Code Generator

AST -> [Code Generator] -> Assembly

Code generator converts ASTs into code of the target language.

In this compiler, the target language is GAS(GNU Assembly) for X86-64 linux.

`generateCode()` does this task.

# Test

see [test.sh](test.sh)
