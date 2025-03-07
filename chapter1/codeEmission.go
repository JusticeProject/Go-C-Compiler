package main

import (
	"fmt"
	"os"
	"strconv"
)

/////////////////////////////////////////////////////////////////////////////////

func doCodeEmission(asm Program_Asm, assemblyFilename string) {
	file, err := os.OpenFile(assemblyFilename, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Could not open/create assembly file", assemblyFilename)
	}

	defer file.Close()

	emitProgram(file, asm)

	// TODO: add comments to the .s file, a # comments out the rest of the line in assembly
}

/////////////////////////////////////////////////////////////////////////////////

func emitProgram(file *os.File, asm Program_Asm) {
	emitFunction(file, asm.fn)
	file.WriteString("\t.section\t.note.GNU-stack,\"\",@progbits\n")
}

/////////////////////////////////////////////////////////////////////////////////

func emitFunction(file *os.File, fn *Function_Asm) {
	file.WriteString("\t.globl " + string(fn.name) + "\n")
	file.WriteString(string(fn.name) + ":\n")
	emitInstructions(file, fn.instructions)
}

/////////////////////////////////////////////////////////////////////////////////

func emitInstructions(file *os.File, instructions []*Instruction) {
	for _, instr := range instructions {
		line := "\t"

		switch instr.typ {
		case MOV_INSTRUCTION:
			line += "movl" + "\t"
			line += getOperandString(instr.src)
			line += ", "
			line += getOperandString(instr.dst)
		case RET_INSTRUCTION:
			line += "ret"
		}

		file.WriteString(line + "\n")
	}
}

/////////////////////////////////////////////////////////////////////////////////

func getOperandString(op *Operand) string {
	switch op.typ {
	case IMMEDIATE_INT_OPERAND:
		return "$" + strconv.FormatInt(int64(op.value), 10)
	case REGISTER_OPERAND:
		// TODO: may need switch statement to determine which register
		return "%" + "eax"
	}

	return ""
}
