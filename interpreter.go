package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func main() {
	t1 := time.Now()
	mem := [30000]byte{}
	loop_jump_positions := [30000]int{}
	loop_jump_pointer := 0
	input_buffer := ""
	input_pointer := 0
	mem_pointer := 0
	program := os.Args[1]
	program_pointer := 0
	in := bufio.NewReader(os.Stdin)

	for {
		instruction := program[program_pointer]
		switch instruction {
		case '+':
			mem[mem_pointer]++
		case '-':
			mem[mem_pointer]--
		case '>':
			mem_pointer++
		case '<':
			mem_pointer--
		case '.':
			fmt.Printf("%c", mem[mem_pointer])
		case ',':
			if input_pointer >= len(input_buffer) {
				fmt.Print("\nAwaiting input\n")
				line, err := in.ReadString('\r')
				if err != nil {
					panic(err)
				}
				input_buffer = line
				input_pointer = 0
			}
			mem[mem_pointer] = input_buffer[input_pointer]
			input_pointer++
		case '[':
			loop_jump_positions[loop_jump_pointer] = program_pointer
			loop_jump_pointer++
			if mem[mem_pointer] == 0 {
				program_pointer++
				target_looping_level := loop_jump_pointer
				// search for the matching ]
				// if you encounter a [, ignore it and its ] counterpart
			L:
				for {
					instruction := program[program_pointer]
					switch instruction {
					case '[':
						loop_jump_pointer++
					case ']':
						if target_looping_level == loop_jump_pointer {
							break L
						}
						loop_jump_pointer--
					}
					program_pointer++
					if program_pointer >= len(program) {
						panic("Could not find matching ] for looping.")
					}
				}
				loop_jump_pointer--
			}
		case ']':
			loop_jump_pointer--
			program_pointer = loop_jump_positions[loop_jump_pointer] - 1
		}
		program_pointer++
		if program_pointer >= len(program) {
			break
		}
	}
	fmt.Printf("\n%s\n", time.Since(t1))
}
