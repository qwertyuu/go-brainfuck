package interpreter

import (
	"testing"
)


func TestOutput(t *testing.T) {
	RunExpectedOutput("+++++++++++++++++++++++++++++++++++++++++++++++++++++.", "", "5", t)
}

func TestInputAndOutput(t *testing.T) {
	program := ",>,+.<."
	RunExpectedOutput(program, "44", "54", t)
}

func TestHelloWorld(t *testing.T) {
	program := "++++++++++[>+++++++>++++++++++>+++>+<<<<-]>++.>+.+++++++..+++.>++.<<+++++++++++++++.>.+++.------.--------.>+.>."
	RunExpectedOutput(program, "", "Hello World!\n", t)
}

func TestBFInterpreter(t *testing.T) {
	// @see http://www.hevanet.com/cristofd/dbfi.b
	program := ">>>+[[-]>>[-]++>+>+++++++[<++++>>++<-]++>>+>+>+++++[>++>++++++<<-]+>>>,<++[[>[->>]<[>>]<<-]<[<]<+>>[>]>[<+>-[[<+>-]>]<[[[-]<]++<-[<+++++++++>[<->-]>>]>>]]<<]<]<[[<]>[[>]>>[>>]+[<<]<[<]<+>>-]>[>]+[->>]<<<<[[<<]<[<]+<<[+>+<<-[>-->+<<-[>+<[>>+<<-]]]>[<+>-]<]++>>-->[>]>>[>>]]<<[>>+<[[<]<]>[[<<]<[<]+[-<+>>-[<<+>++>-[<->[<<+>>-]]]<[>+<-]>]>[>]>]>[>>]>>]<<[>>+>>+>>]<<[->>>>>>>>]<<[>.>>>>>>>]<<[>->>>>>]<<[>,>>>]<<[>+>]<<[+<<]<]"
	helloWorldProgram := "++++++++++[>+++++++>++++++++++>+++>+<<<<-]>++.>+.+++++++..+++.>++.<<+++++++++++++++.>.+++.------.--------.>+.>.!"

	RunExpectedOutput(program, helloWorldProgram, "Hello World!\n", t)
}

func RunExpectedOutput(program string, input string, expectedOutput string, t *testing.T) {
	outputs := ""
	outputProcessor := func(output byte) {
		outputs += string(output)
	}
	char := 0
	provider := func() byte {
		toRet := input[char]
		char++
		return toRet
	}
	RunProgram(program, provider, outputProcessor)
	if outputs != expectedOutput {
		t.Errorf(`Output should be '%s' but got %q`, expectedOutput, outputs)
	}
}