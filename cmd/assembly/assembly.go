package main

import (
	"math"
	"strings"
)

type registerSetting struct {
	registerIndex int
	data           byte
	print          bool
}

var finalProgram strings.Builder
var mem [30000]byte
var registerPointer int

func main() {
	SetRegisters([]registerSetting{
		{1, 'H', true},  // 72
		{2, 'e', true},  // 101
		{3, ' ', false}, // 32
		{4, 10, false},  // 10
	}, 0)

	// Hello
	SetRegister(2, 'l', true)
	PrintRegister(2)
	SetRegister(2, 'o', true)

	// space
	PrintRegister(3)

	// World
	SetRegister(1, 'W', true)
	PrintRegister(2)
	SetRegister(2, 'r', true)
	SetRegister(2, 'l', true)
	SetRegister(2, 'd', true)
	SetRegister(3, '!', true)
	PrintRegister(4)
	println(finalProgram.String())
}

func abs(value int) byte {
	if value < 0 {
		return byte(-value)
	}
	return byte(value)
}

func SetRegisters(registers []registerSetting, bufferRegister int) {
	if len(registers) == 1 {
		SetRegister(registers[0].registerIndex, registers[0].data, registers[0].print)
		return
	}
	// TODO: ensure that all register index are unique

	registerValueDiffs := make([]int, len(registers))
	for i, register := range registers {
		registerValueDiffs[i] = int(register.data) - int(mem[register.registerIndex])
	}

	minAbsDiffValue := byte(math.MaxUint8)
	for _, diff := range registerValueDiffs {
		intData := abs(diff)
		if intData < minAbsDiffValue {
			minAbsDiffValue = intData
		}
	}

	minDivisorModuloSum := math.MaxInt
	minDivisorModuloFactor := byte(0)
	for i := byte(2); i < minAbsDiffValue; i++ {
		divisorModuloSum := 0
		for _, diff := range registerValueDiffs {
			divisorModuloSum += diff % int(i)
			divisorModuloSum += diff / int(i)
		}
		if divisorModuloSum < minDivisorModuloSum {
			minDivisorModuloSum = divisorModuloSum
			minDivisorModuloFactor = byte(i)
		}
	}
	if minDivisorModuloFactor == 0 {
		for i := 0; i < len(registers); i++ {
			SetRegister(registers[i].registerIndex, registers[i].data, registers[i].print)
		}
		return
	}

	// set base factor
	SetRegister(bufferRegister, byte(minDivisorModuloFactor), false)

	finalProgram.WriteByte('[')
	finalProgram.WriteString(getMoveBetweenRegisters(bufferRegister, registers[0].registerIndex))
	finalProgram.WriteString(getMoveBetweenValues(mem[registers[0].registerIndex], registers[0].data/minDivisorModuloFactor))
	mem[registers[0].registerIndex] = registers[0].data / minDivisorModuloFactor * minDivisorModuloFactor

	for i := 1; i < len(registers); i++ {
		finalProgram.WriteString(getMoveBetweenRegisters(registers[i-1].registerIndex, registers[i].registerIndex))
		finalProgram.WriteString(getMoveBetweenValues(mem[registers[i].registerIndex], registers[i].data/minDivisorModuloFactor))
		mem[registers[i].registerIndex] = registers[i].data / minDivisorModuloFactor * minDivisorModuloFactor
	}
	// move back to the buffer register
	finalProgram.WriteString(getMoveBetweenRegisters(registers[len(registers)-1].registerIndex, bufferRegister))
	finalProgram.WriteString("-]")

	finalProgram.WriteString(getMoveBetweenRegisters(bufferRegister, registers[0].registerIndex))
	finalProgram.WriteString(getMoveBetweenValues(mem[registers[0].registerIndex], registers[0].data))
	mem[registers[0].registerIndex] = registers[0].data
	if registers[0].print {
		finalProgram.WriteByte('.')
	}

	for i := 1; i < len(registers); i++ {
		finalProgram.WriteString(getMoveBetweenRegisters(registers[i-1].registerIndex, registers[i].registerIndex))
		finalProgram.WriteString(getMoveBetweenValues(mem[registers[i].registerIndex], registers[i].data))
		mem[registers[i].registerIndex] = registers[i].data
		if registers[i].print {
			finalProgram.WriteByte('.')
		}
	}

	registerPointer = registers[len(registers)-1].registerIndex
}

func getMoveBetweenRegisters(from int, to int) string {
	if from < to {
		return strings.Repeat(">", to-from)
	}
	return strings.Repeat("<", from-to)
}

func getMoveBetweenValues(from byte, to byte) string {
	if from < to {
		return strings.Repeat("+", int(to-from))
	}
	return strings.Repeat("-", int(from-to))
}

func PrintRegister(bufferRegister int) {
	moveTo(bufferRegister)
	finalProgram.WriteRune('.')
}

func moveTo(destinationRegisterIndex int) {
	// move to this register
	finalProgram.WriteString(getMoveBetweenRegisters(registerPointer, destinationRegisterIndex))
	registerPointer = destinationRegisterIndex
}

func SetRegister(destinationRegisterIndex int, data byte, print bool) {
	moveTo(destinationRegisterIndex)
	// check what value this register is at now and apply the diff
	finalProgram.WriteString(getMoveBetweenValues(mem[destinationRegisterIndex], data))
	mem[destinationRegisterIndex] = data
	if print {
		finalProgram.WriteByte('.')
	}
}
