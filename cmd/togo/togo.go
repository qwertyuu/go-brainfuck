package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func commitIncrement(program *strings.Builder, counter *int, toIncrement string) {
	if *counter > 0 {
		if *counter == 1 {
			program.WriteString(toIncrement + "++\n")
		} else {
			program.WriteString(toIncrement + " += " + strconv.Itoa(*counter) + "\n")
		}
	} else if *counter < 0 {
		if *counter == -1 {
			program.WriteString(toIncrement + "--\n")
		} else {
			program.WriteString(toIncrement + " -= " + strconv.Itoa(-*counter) + "\n")
		}
	}
	*counter = 0
}

var optimizeIncrements = flag.Bool("optimize_increments", false, "Optimize increments by grouping them. might make weird stuff if your program is based around wrapping byte values (overflow or underflow) while this is active")
var optimizePointerIncrements = flag.Bool("optimize_pointer_increments", true, "Optimize pointer increments. I see no reasons to disable this.")

var useTinygo = flag.Bool("use_tinygo", false, "Whether to use the go or tinygo compiler. tinygo needs to be installed and creates very small executables that run as fast or faster than go executables")
var formatCode = flag.Bool("format_code", true, "Format the IL code. Easier for a human to read but longer to run this script.")
var compile = flag.Bool("compile", true, "Compile code after IL translation")
var timeProgram = flag.Bool("time_program", true, "Add runtime stopwatch to the output")

var program = flag.String("program", "", "BF program to convert to Go")

func main() {
	flag.Parse() // after declaring flags we need to parse them
	var finalProgram strings.Builder
	incrementCounter := 0
	memIncrementCounter := 0
	needsInput := strings.Contains(*program, ",")
	imports := []string{
		"fmt",
	}

	if *timeProgram {
		imports = append(imports, "time")
	}
	if needsInput {
		imports = append(imports, "bufio")
		imports = append(imports, "os")
	}
	fmt.Println("Translating to Go code...")

	finalProgram.WriteString(`
	package main

	import (`)
	finalProgram.WriteString("\"" + strings.Join(imports, "\";\"") + "\"")
	finalProgram.WriteRune(')')
	if needsInput {
		finalProgram.WriteString(`
		func getInput() {
			if inputPointer >= len(inputBuffer) {
				fmt.Print("\nAwaiting input\n")
				line, err := in.ReadString('\n')
				if err != nil {
					panic(err)
				}
				inputBuffer = line
				inputPointer = 0
			}
		}
		
		var inputBuffer string
		var inputPointer int
		var in *bufio.Reader
		`)
	}
	finalProgram.WriteString(`
	func main() {
`)
	if *timeProgram {
		finalProgram.WriteString(`
		t1 := time.Now()
		`)
	}
	finalProgram.WriteString(`mem := [30000]byte{}
	memPointer := 0
	`)
	if needsInput {
		finalProgram.WriteString("in = bufio.NewReader(os.Stdin)\n")
	}

	for _, instruction := range *program {
		if *optimizeIncrements && !(instruction == '+' || instruction == '-') {
			commitIncrement(&finalProgram, &incrementCounter, "mem[memPointer]")
		}
		if *optimizePointerIncrements && !(instruction == '>' || instruction == '<') {
			commitIncrement(&finalProgram, &memIncrementCounter, "memPointer")
		}
		switch instruction {
		case '+':
			if *optimizeIncrements {
				incrementCounter++
			} else {
				finalProgram.WriteString("mem[memPointer]++\n")
			}
		case '-':
			if *optimizeIncrements {
				incrementCounter--
			} else {
				finalProgram.WriteString("mem[memPointer]--\n")
			}
		case '>':
			if *optimizePointerIncrements {
				memIncrementCounter++
			} else {
				finalProgram.WriteString("memPointer++\n")
			}
		case '<':
			if *optimizePointerIncrements {
				memIncrementCounter--
			} else {
				finalProgram.WriteString("memPointer--\n")
			}
		case '.':
			finalProgram.WriteString("fmt.Printf(\"%c\", mem[memPointer])\n")
		case ',':
			finalProgram.WriteString("getInput()\n")
			finalProgram.WriteString("mem[memPointer] = inputBuffer[inputPointer]\n")
			finalProgram.WriteString("inputPointer++\n")
		case '[':
			finalProgram.WriteString("for mem[memPointer] > 0 {\n")
		case ']':
			finalProgram.WriteString("}\n")
		}
	}
	if *timeProgram {
		finalProgram.WriteString("fmt.Printf(\"\\n%s\\n\", time.Since(t1))\n")
	}
	finalProgram.WriteByte('}')

	fmt.Println("Writing Go code to disk...")
	err := os.WriteFile("asgo/bf.go", []byte(finalProgram.String()), 0644)
	if err != nil {
		panic(err)
	}
	if *formatCode {
		fmt.Println("Formatting Go code... (why not)")
		cmd := exec.Command("gofmt", "-w", "asgo/bf.go")
		err = cmd.Run()
		if err != nil {
			panic(err)
		}
	}
	if *compile {
		fmt.Println("Compiling Go code...")
		var cmd *exec.Cmd
		if *useTinygo {
			cmd = exec.Command("tinygo", "build", "-opt=2", "-o", "bin/", "asgo/bf.go")
		} else {
			cmd = exec.Command("go", "build", "-ldflags=-s -w", "-o", "bin/", "asgo/bf.go")
		}
		err = cmd.Run()
		if err != nil {
			panic(err)
		}
	}

}
