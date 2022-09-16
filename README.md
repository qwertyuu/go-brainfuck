# go-brainfuck
 brainfuck interpreter and compiler in go.

## Interpreter
 interpreter.go was the first iteration of my idea, taking only a single argument when running, which is the BF program to interpret. `go run interpreter.go "++++++++++[>+++++++>++++++++++>+++>+<<<<-]>++.>+.+++++++..+++.>++.<<+++++++++++++++.>.+++.------.--------.>+.>."`

## IL Compiler (intermediate language compiler)
 togo/togo.go is the second iteration that makes an BF > Go > Executable IL compiler to run your sweet BF directly in bytecode on your computer.
 It runs differently, you need to specify some arguments to run the program.
 Example:
 `go run togo/togo.go -program="++++++++++[>+++++++>++++++++++>+++>+<<<<-]>++.>+.+++++++..+++.>++.<<+++++++++++++++.>.+++.------.--------.>+.>." -optimize_increments=true -optimize_pointer_increments=true -use_tinygo=false `

## Misc
 Some brainfuck I wrote when I was 16:
 `++++++++++[>+++++<-]>-->,>,<<[>-<-]<++++++++++[>+++++<-]>--[>>l--<<--]<++++++++++[>>>>+>+++>+>+<<<<<<<-]>>>>>>-<+++++>++++++++++[<.>-]<<<+[>>.<.>>>-<<<<-]>>.---<<<[>>>..<<<-]>>>+++.<.>>>[<<.<.>>>-]`
 Takes 2 coordinates as input, produces a 2D plot with a dot drawn at those coordinates.