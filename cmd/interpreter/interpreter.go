package main

import (
	"bufio"
	"fmt"
	"go-brainfuck/v2/internal/interpreter"
	"os"
	"time"
)

var input_pointer int
var input_buffer string
var input_reader *bufio.Reader

func input_provider() byte {
	if input_pointer >= len(input_buffer) {
		fmt.Print("\nAwaiting input\n")
		line, err := input_reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		input_buffer = line
		input_pointer = 0
	}
	to_return := input_buffer[input_pointer]
	input_pointer++
	return to_return
}

func output_processor(to_print byte) {
	fmt.Printf("%c", to_print)
}

func main() {
	t1 := time.Now()
	input_reader = bufio.NewReader(os.Stdin)
	interpreter.RunProgram(os.Args[1], input_provider, output_processor)
	fmt.Printf("\n%s\n", time.Since(t1))
}
