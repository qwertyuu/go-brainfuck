package main

import (
	"fmt"
	"math"
	"strings"
)

type register_setting struct {
	register_index int
	data           byte
	print          bool
}

var final_program strings.Builder
var mem [30000]byte
var register_pointer int

func main() {
	Set_registers([]register_setting{
		{1, 'H', true},  // 72
		{2, 'e', true},  // 101
		{3, ' ', false}, // 32
		{4, 10, false},  // 10
	}, 0)

	// Hello
	Set_register(2, 'l', true)
	Print_register(2)
	Set_register(2, 'o', true)

	// space
	Print_register(3)

	// World
	Set_register(1, 'W', true)
	Print_register(2)
	Set_register(2, 'r', true)
	Set_register(2, 'l', true)
	Set_register(2, 'd', true)
	Set_register(3, '!', true)
	Print_register(4)
	println(final_program.String())
}

func abs(value int) byte {
	if value < 0 {
		return byte(-value)
	}
	return byte(value)
}

func Set_registers(registers []register_setting, buffer_register int) {
	if len(registers) == 1 {
		Set_register(registers[0].register_index, registers[0].data, registers[0].print)
		return
	}
	// TODO: ensure that all register index are unique

	register_value_diffs := make([]int, len(registers))
	for i, register := range registers {
		register_value_diffs[i] = int(register.data) - int(mem[register.register_index])
	}

	min_abs_diff_value := byte(math.MaxUint8)
	for _, diff := range register_value_diffs {
		int_data := abs(diff)
		if int_data < min_abs_diff_value {
			min_abs_diff_value = int_data
		}
	}

	min_divisor_modulo_sum := math.MaxInt
	min_divisor_modulo_factor := byte(0)
	for i := byte(2); i < min_abs_diff_value; i++ {
		divisor_modulo_sum := 0
		for _, diff := range register_value_diffs {
			divisor_modulo_sum += diff % int(i)
			divisor_modulo_sum += diff / int(i)
		}
		if divisor_modulo_sum < min_divisor_modulo_sum {
			min_divisor_modulo_sum = divisor_modulo_sum
			min_divisor_modulo_factor = byte(i)
		}
	}
	fmt.Println(min_divisor_modulo_factor)
	fmt.Println(min_divisor_modulo_sum)

	// set base factor
	Set_register(buffer_register, byte(min_divisor_modulo_factor), false)

	final_program.WriteByte('[')
	final_program.WriteString(get_move_between_registers(buffer_register, registers[0].register_index))
	final_program.WriteString(get_move_between_values(mem[registers[0].register_index], registers[0].data/min_divisor_modulo_factor))
	mem[registers[0].register_index] = registers[0].data / min_divisor_modulo_factor * min_divisor_modulo_factor

	for i := 1; i < len(registers); i++ {
		final_program.WriteString(get_move_between_registers(registers[i-1].register_index, registers[i].register_index))
		final_program.WriteString(get_move_between_values(mem[registers[i].register_index], registers[i].data/min_divisor_modulo_factor))
		mem[registers[i].register_index] = registers[i].data / min_divisor_modulo_factor * min_divisor_modulo_factor
	}
	// move back to the buffer register
	final_program.WriteString(get_move_between_registers(registers[len(registers)-1].register_index, buffer_register))
	final_program.WriteString("-]")

	final_program.WriteString(get_move_between_registers(buffer_register, registers[0].register_index))
	final_program.WriteString(get_move_between_values(mem[registers[0].register_index], registers[0].data))
	mem[registers[0].register_index] = registers[0].data
	if registers[0].print {
		final_program.WriteByte('.')
	}

	for i := 1; i < len(registers); i++ {
		final_program.WriteString(get_move_between_registers(registers[i-1].register_index, registers[i].register_index))
		final_program.WriteString(get_move_between_values(mem[registers[i].register_index], registers[i].data))
		mem[registers[i].register_index] = registers[i].data
		if registers[i].print {
			final_program.WriteByte('.')
		}
	}

	register_pointer = registers[len(registers)-1].register_index
}

func get_move_between_registers(from int, to int) string {
	if from < to {
		return strings.Repeat(">", to-from)
	}
	return strings.Repeat("<", from-to)
}

func get_move_between_values(from byte, to byte) string {
	if from < to {
		return strings.Repeat("+", int(to-from))
	}
	return strings.Repeat("-", int(from-to))
}

func Print_register(buffer_register int) {
	move_to(buffer_register)
	final_program.WriteRune('.')
}

func move_to(destination_register_index int) {
	// move to this register
	final_program.WriteString(get_move_between_registers(register_pointer, destination_register_index))
	register_pointer = destination_register_index
}

func Set_register(destination_register_index int, data byte, print bool) {
	move_to(destination_register_index)
	// check what value this register is at now and apply the diff
	final_program.WriteString(get_move_between_values(mem[destination_register_index], data))
	mem[destination_register_index] = data
	if print {
		final_program.WriteByte('.')
	}
}
