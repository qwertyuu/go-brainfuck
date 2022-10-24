package main

import (
	"bufio"
	"fmt"
	"go-brainfuck/v2/internal/interpreter"
	"os"
	"time"
)

var inputPointer int
var inputBuffer string
var inputReader *bufio.Reader

func inputProvider() byte {
	if inputPointer >= len(inputBuffer) {
		fmt.Print("\nAwaiting input\n")
		line, err := inputReader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		inputBuffer = line
		inputPointer = 0
	}
	toReturn := inputBuffer[inputPointer]
	inputPointer++
	return toReturn
}

func outputProcessor(toPrint byte) {
	fmt.Printf("%c", toPrint)
}

func main() {
	t1 := time.Now()
	inputReader = bufio.NewReader(os.Stdin)
	interpreter.RunProgram(os.Args[1], inputProvider, outputProcessor)
	fmt.Printf("\n%s\n", time.Since(t1))
}
