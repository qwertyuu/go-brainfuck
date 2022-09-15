package main

import (
	"fmt"
	"os"
)

func main() {
	mem := [255]int{}
	mem_pointer := 0
	program := os.Args[1]
	program_pointer := 0
	fmt.Println(mem)

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
		}
		fmt.Println(mem)
		program_pointer += 1
		if program_pointer == len(program) {
			break
		}
	}

}
