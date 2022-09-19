package interpreter

type rune_provider func() byte
type output_processor func(byte)

var loop_jump_pointer int
var loop_jump_positions [30000]int
var mem [30000]byte
var mem_pointer int
var program_pointer int

func RunProgram(program string, input_provider rune_provider, output_processor output_processor) {
	loop_jump_pointer = 0
	mem_pointer = 0

	for program_pointer = 0; program_pointer < len(program); program_pointer++ {
		switch program[program_pointer] {
		case '+':
			mem[mem_pointer]++
		case '-':
			mem[mem_pointer]--
		case '>':
			mem_pointer++
		case '<':
			mem_pointer--
		case '.':
			output_processor(mem[mem_pointer])
		case ',':
			mem[mem_pointer] = input_provider()
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
					switch program[program_pointer] {
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
	}
}
