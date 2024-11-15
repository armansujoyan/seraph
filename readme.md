[![Review Assignment Due Date](https://classroom.github.com/assets/deadline-readme-button-22041afd0340ce965d47ae6ef1cefeee28c7c493a6346c4f15d667ab976d596c.svg)](https://classroom.github.com/a/6UCCtKXu)


in this assignment you need to extend the already written scanner with the parser and code generator.

language syntax is described in the following EBNF:

```
prg              = prgHeader varDefs "begin" statementSeq "end" "."
prgHeader        = "program" identifier ";"
varDefs          = [ "var" varSeq { varSeq } ]
varSeq           = identifier { "," identifier } : type ";"
type             = "integer" | "string"
statementSeq     = { ( simpleAssignment | complexAssignment) ";" }
simpleAssignment = identifier ":=" operand ";"
complexAssignment= identifier ":=" operand mathOP operand ";"
operand          = identifier | number
mathOp           = "+" | "-"
```

identifiers, numbers you already identify in the scanner.

your parser also needs to conduct semantic analisys, you will need to create a symbol table, which will be used in the following checks:

variable name cannot be the same as program name.

the variable can be declared only once.

the result of the work will be a minimal compiler.

generated code should be for the real cpu of our server: x86_64 ISA.

the result of the last calculation should be available as an exit code of the program.

if you can instead link to libc or use a kernel call write(), and introduce a function in the language or in your library with a function that is able to print, you'll get extra points.

tip: you can do that by copy pasting code from lab work and always linking the objcect file that does printf and calling printf.

then you shoudl have the function 'write'.

writeCall = write "(" identifier ")" ";"

after the assembly code is generated, your compiler should call assembler and linker to produce the resulting binary.

we should be able to compile your compiler source by just typing 'make'.

we won't consider checking your work otherwise.

your compiler should accept a commandline argument - source code to compile.

we will use this by calling your compiler and giving it different source files to compile.

./yourcompiler ourcode.pas

we expect, that if we give your compiler file 'source.pas' then it'll produce the file 'source.s', which will be later assembled to 'source.o', and linked to the program named 'source'.


