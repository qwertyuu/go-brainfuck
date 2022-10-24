package interpreter

type runeProvider func() byte
type outputProcessor func(byte)

var loopJumpPointer int
var loopJumpPositions [30000]int
var mem [30000]byte
var memPointer int
var programPointer int

func RunProgram(program string, inputProvider runeProvider, outputProcessor outputProcessor) {
	loopJumpPointer = 0
	memPointer = 0

	for programPointer = 0; programPointer < len(program); programPointer++ {
		switch program[programPointer] {
		case '+':
			mem[memPointer]++
		case '-':
			mem[memPointer]--
		case '>':
			memPointer++
		case '<':
			memPointer--
		case '.':
			outputProcessor(mem[memPointer])
		case ',':
			mem[memPointer] = inputProvider()
		case '[':
			loopJumpPositions[loopJumpPointer] = programPointer
			loopJumpPointer++
			if mem[memPointer] == 0 {
				programPointer++
				targetLoopingLevel := loopJumpPointer
				// search for the matching ]
				// if you encounter a [, ignore it and its ] counterpart
			L:
				for {
					switch program[programPointer] {
					case '[':
						loopJumpPointer++
					case ']':
						if targetLoopingLevel == loopJumpPointer {
							break L
						}
						loopJumpPointer--
					}
					programPointer++
					if programPointer >= len(program) {
						panic("Could not find matching ] for looping.")
					}
				}
				loopJumpPointer--
			}
		case ']':
			loopJumpPointer--
			programPointer = loopJumpPositions[loopJumpPointer] - 1
		}
	}
}
