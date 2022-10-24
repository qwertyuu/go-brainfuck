package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func commit_increment(program *strings.Builder, counter *int, to_increment string) {
	if *counter > 0 {
		if *counter == 1 {
			program.WriteString(to_increment + "++\n")
		} else {
			program.WriteString(to_increment + " += " + strconv.Itoa(*counter) + "\n")
		}
	} else if *counter < 0 {
		if *counter == -1 {
			program.WriteString(to_increment + "--\n")
		} else {
			program.WriteString(to_increment + " -= " + strconv.Itoa(-*counter) + "\n")
		}
	}
	*counter = 0
}

var optimize_increments bool
var optimize_pointer_increments bool
var use_tinygo bool
var format_code bool
var compile bool
var time_program bool

func main() {
	program := ""
	flag.StringVar(&program, "program", "", "BF program to convert to Go")
	flag.BoolVar(&optimize_increments, "optimize_increments", false, "Optimize increments by grouping them. might make weird stuff if your program is based around wrapping byte values (overflow or underflow) while this is active")
	flag.BoolVar(&optimize_pointer_increments, "optimize_pointer_increments", true, "Optimize pointer increments. I see no reasons to disable this.")
	flag.BoolVar(&use_tinygo, "use_tinygo", false, "Whether to use the go or tinygo compiler. tinygo needs to be installed and creates very small executables that run as fast or faster than go executables")
	flag.BoolVar(&format_code, "format_code", true, "Format the IL code. Easier for a human to read but longer to run this script.")
	flag.BoolVar(&compile, "compile", true, "Compile code after IL translation")
	flag.BoolVar(&time_program, "time_program", true, "Add runtime stopwatch to the output")

	flag.Parse() // after declaring flags we need to call it
	var final_program strings.Builder
	increment_counter := 0
	mem_increment_counter := 0
	needs_input := strings.Contains(program, ",")
	imports := []string{
		"fmt",
	}

	if time_program {
		imports = append(imports, "time")
	}
	if needs_input {
		imports = append(imports, "bufio")
		imports = append(imports, "os")
	}
	fmt.Println("Translating to Go code...")

	final_program.WriteString(`
	package main

	import (`)
	final_program.WriteString("\"" + strings.Join(imports, "\";\"") + "\"")
	final_program.WriteRune(')')
	if needs_input {
		final_program.WriteString(`
		func get_input() {
			if input_pointer >= len(input_buffer) {
				fmt.Print("\nAwaiting input\n")
				line, err := in.ReadString('\n')
				if err != nil {
					panic(err)
				}
				input_buffer = line
				input_pointer = 0
			}
		}
		
		var input_buffer string
		var input_pointer int
		var in *bufio.Reader
		`)
	}
	final_program.WriteString(`
	func main() {
`)
	if time_program {
		final_program.WriteString(`
		t1 := time.Now()
		`)
	}
	final_program.WriteString(`mem := [30000]byte{}
	mem_pointer := 0
	`)
	if needs_input {
		final_program.WriteString("in = bufio.NewReader(os.Stdin)\n")
	}

	for _, instruction := range program {
		if optimize_increments && !(instruction == '+' || instruction == '-') {
			commit_increment(&final_program, &increment_counter, "mem[mem_pointer]")
		}
		if optimize_pointer_increments && !(instruction == '>' || instruction == '<') {
			commit_increment(&final_program, &mem_increment_counter, "mem_pointer")
		}
		switch instruction {
		case '+':
			if optimize_increments {
				increment_counter++
			} else {
				final_program.WriteString("mem[mem_pointer]++\n")
			}
		case '-':
			if optimize_increments {
				increment_counter--
			} else {
				final_program.WriteString("mem[mem_pointer]--\n")
			}
		case '>':
			if optimize_pointer_increments {
				mem_increment_counter++
			} else {
				final_program.WriteString("mem_pointer++\n")
			}
		case '<':
			if optimize_pointer_increments {
				mem_increment_counter--
			} else {
				final_program.WriteString("mem_pointer--\n")
			}
		case '.':
			final_program.WriteString("fmt.Printf(\"%c\", mem[mem_pointer])\n")
		case ',':
			final_program.WriteString("get_input()\n")
			final_program.WriteString("mem[mem_pointer] = input_buffer[input_pointer]\n")
			final_program.WriteString("input_pointer++\n")
		case '[':
			final_program.WriteString("for mem[mem_pointer] > 0 {\n")
		case ']':
			final_program.WriteString("}\n")
		}
	}
	if time_program {
		final_program.WriteString("fmt.Printf(\"\\n%s\\n\", time.Since(t1))\n")
	}
	final_program.WriteByte('}')

	fmt.Println("Writing Go code to disk...")
	err := os.WriteFile("asgo/bf.go", []byte(final_program.String()), 0644)
	if err != nil {
		panic(err)
	}
	if format_code {
		fmt.Println("Formatting Go code... (why not)")
		cmd := exec.Command("gofmt", "-w", "asgo/bf.go")
		err = cmd.Run()
		if err != nil {
			panic(err)
		}
	}
	if compile {
		fmt.Println("Compiling Go code...")
		var cmd *exec.Cmd
		if use_tinygo {
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
