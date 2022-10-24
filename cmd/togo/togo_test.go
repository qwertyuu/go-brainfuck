package main

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"
)


func TestOutput(t *testing.T) {
	RunExpectedOutput("+++++++++++++++++++++++++++++++++++++++++++++++++++++.", "", "5", t)
}

func TestInputOutput(t *testing.T) {
	RunExpectedOutput(",>,+.<.", "44\n", "\nAwaiting input\n54", t)
}

func TestHelloWorld(t *testing.T) {
	program := "++++++++++[>+++++++>++++++++++>+++>+<<<<-]>++.>+.+++++++..+++.>++.<<+++++++++++++++.>.+++.------.--------.>+.>."
	RunExpectedOutput(program, "", "Hello World!\n", t)
}

func TestBFInterpreter(t *testing.T) {
	// @see http://www.hevanet.com/cristofd/dbfi.b
	program := ">>>+[[-]>>[-]++>+>+++++++[<++++>>++<-]++>>+>+>+++++[>++>++++++<<-]+>>>,<++[[>[->>]<[>>]<<-]<[<]<+>>[>]>[<+>-[[<+>-]>]<[[[-]<]++<-[<+++++++++>[<->-]>>]>>]]<<]<]<[[<]>[[>]>>[>>]+[<<]<[<]<+>>-]>[>]+[->>]<<<<[[<<]<[<]+<<[+>+<<-[>-->+<<-[>+<[>>+<<-]]]>[<+>-]<]++>>-->[>]>>[>>]]<<[>>+<[[<]<]>[[<<]<[<]+[-<+>>-[<<+>++>-[<->[<<+>>-]]]<[>+<-]>]>[>]>]>[>>]>>]<<[>>+>>+>>]<<[->>>>>>>>]<<[>.>>>>>>>]<<[>->>>>>]<<[>,>>>]<<[>+>]<<[+<<]<]"
	helloWorldProgram := "++++++++++[>+++++++>++++++++++>+++>+<<<<-]>++.>+.+++++++..+++.>++.<<+++++++++++++++.>.+++.------.--------.>+.>.!\n"

	RunExpectedOutput(program, helloWorldProgram, "\nAwaiting input\nHello World!\n", t)
}

func RunExpectedOutput(program string, input string, expectedOutput string, t *testing.T) {
	os.Chdir("../..")
	defer os.Chdir("cmd/togo")
	os.Args[1] = "-program=\"" + program + "\""
	os.Args[2] = "-optimize_increments=true"
	os.Args[3] = "-time_program=false"
	main()

	cmd := exec.Command("go", "run", "asgo/bf.go")
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
	if out.String() != expectedOutput {
		t.Fatalf("Expected '%s', got '%s'", expectedOutput, out.String())
	}
}