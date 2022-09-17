# go-brainfuck
 brainfuck interpreter and compiler in go. I did this project in honor of the 10th anniversary of my first brainfuck program which can be seen [here](https://www.youtube.com/shorts/k8ufd-OyN1Y)

## Interpreter
 interpreter.go was the first iteration of my idea, taking only a single argument when running, which is the BF program to interpret. `go run interpreter.go "++++++++++[>+++++++>++++++++++>+++>+<<<<-]>++.>+.+++++++..+++.>++.<<+++++++++++++++.>.+++.------.--------.>+.>."`

## IL Compiler (intermediate language compiler)
 togo/togo.go is the second iteration that makes an BF > Go > Executable IL compiler to run your sweet BF directly in bytecode on your computer.
 It runs differently, you need to specify some arguments to run the program.
 Example:
 `go run togo/togo.go -program="++++++++++[>+++++++>++++++++++>+++>+<<<<-]>++.>+.+++++++..+++.>++.<<+++++++++++++++.>.+++.------.--------.>+.>." -optimize_increments=true -optimize_pointer_increments=true -use_tinygo=false `

 For now it's very dumb. It will fail if you are on linux, but it may end up working anyways, just that it cannot run the final executable. Try forking my repo and messing with togo/togo.go yourself. This code could actually be a python script that generates go and compiles it. Oh well!

## Assembly

I created kind of a small assembly language to generate brainfuck at a higher level. It seems to work as a PoC with typing a "Hello World!" program that prints. I have not tried anything more complex than this however.

## Misc
 Some brainfuck I wrote when I was 16:
 `++++++++++[>+++++<-]>-->,>,<<[>-<-]<++++++++++[>+++++<-]>--[>>l--<<--]<++++++++++[>>>>+>+++>+>+<<<<<<<-]>>>>>>-<+++++>++++++++++[<.>-]<<<+[>>.<.>>>-<<<<-]>>.---<<<[>>>..<<<-]>>>+++.<.>>>[<<.<.>>>-]`
 Takes 2 coordinates as input, produces a 2D plot with a dot drawn at those coordinates.