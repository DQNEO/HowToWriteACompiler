# How to write a compiler from scratch in 30 mininutes

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
$ echo '30 + 12' | go run main.go
```

This program receives source code from stdin, and emits assembly code to stdout.

If you want to compile and run at once, `asrun` script helps you.

```
$ echo '30 + 12' | go run main.go | ./asrun
```

`asrun` takes assembly code from stdin and executes it while displaying the code and the resulting status code.

```
$ echo  '30 + 12' | go run main.go | ./asrun
-------- a.s ----------------
  .global main
main:
  movq  $30, %rax
  movq  $12, %rbx
  addq %rbx, %rax
  ret
-------- result -------------
42
```

# Design

The compiler has 3 phases.

Source Code -> [Tokenizer] -> Tokens -> [Parser] -> AST -> [Code Generator] -> Assembly

## Tokenizer

Source Code -> [Tokenizer] -> Tokens

Tokenizer analyzes the byte stream of source code, and breaks it down into a list of tokens.

In this compiler, the function `tokenize()` does this task.

## Parser

Tokens -> [Parser] -> AST

Parser analyzes stream of tokens, and composes a tree of nested structs , which represents sytanx structure of source code.

This tree is called AST (Abstract Syntax Tree).

The function `parser()` does this task.

## Code Generator

AST -> [Code Generator] -> Assembly

Code generator converts AST into target language code.

In this compiler, the target language is GAS(GNU Assembly) for X86-64 linux.

The function `generateCode()` does this task.

# How to run unit tests

```
$ ./test.sh
```

# SEE ALSO

This project derives from My Go compiler.

https://github.com/DQNEO/minigo

Actually, the first 7 commits of `minigo` equals to this repo.
