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

	asm.emitAssembly(file)

	// TODO: add comments to the .s file, a # comments out the rest of the line in assembly
}

/////////////////////////////////////////////////////////////////////////////////

func (asm *Program_Asm) emitAssembly(file *os.File) {
	// TODO: need to handle more than one function
	asm.fn.emitAssembly(file)
	file.WriteString("\t.section\t.note.GNU-stack,\"\",@progbits\n")
}

/////////////////////////////////////////////////////////////////////////////////

func (fn *Function_Asm) emitAssembly(file *os.File) {
	file.WriteString("\t.globl " + string(fn.name) + "\n")
	file.WriteString(string(fn.name) + ":\n")

	// include the function prologue instructions for preparing the stack
	file.WriteString("\t" + "pushq" + "\t" + "%rbp" + "\n")
	file.WriteString("\t" + "movq" + "\t" + "%rsp, %rbp" + "\n")

	for _, instr := range fn.instructions {
		instr.emitAssembly(file)
	}
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Instruction_Asm) emitAssembly(file *os.File) {
	switch instr.typ {
	case MOV_INSTRUCTION_ASM:
		file.WriteString("\t" + "movl" + "\t" + getOperandString(instr.src) + ", " + getOperandString(instr.dst) + "\n")
	case RET_INSTRUCTION_ASM:
		// include the function epilogue instructions for restoring the stack
		file.WriteString("\t" + "movq" + "\t" + "%rbp, %rsp" + "\n")
		file.WriteString("\t" + "popq" + "\t" + "%rbp" + "\n")
		file.WriteString("\t" + "ret" + "\n")
	case UNARY_INSTRUCTION_ASM:
		file.WriteString("\t" + getUnaryOperatorString(instr.unOp) + "\t" + getOperandString(instr.dst) + "\n")
	case ALLOCATE_STACK_INSTRUCTION_ASM:
		file.WriteString("\t" + "subq" + "\t" + getOperandString(instr.src) + ", %rsp" + "\n")
	default:
		fmt.Println("unknown instruction:", instr.typ)
		os.Exit(1)
	}
}

/////////////////////////////////////////////////////////////////////////////////

func getUnaryOperatorString(unOp UnaryOperatorTypeAsm) string {
	switch unOp {
	case NEGATE_OPERATOR_ASM:
		return "negl"
	case NOT_OPERATOR_ASM:
		return "notl"
	default:
		fmt.Println("unknown unary operator:", unOp)
		os.Exit(1)
	}

	return ""
}

/////////////////////////////////////////////////////////////////////////////////

func getOperandString(op *Operand_Asm) string {
	switch op.typ {
	case IMMEDIATE_INT_OPERAND_ASM:
		return "$" + strconv.FormatInt(int64(op.value), 10)
	case REGISTER_OPERAND_ASM:
		return "%" + getRegisterString(op.reg)
	case STACK_OPERAND_ASM:
		return strconv.FormatInt(int64(op.value), 10) + "(%rbp)"
	default:
		fmt.Println("unknown operand type:", op.typ)
		os.Exit(1)
	}

	return ""
}

/////////////////////////////////////////////////////////////////////////////////

func getRegisterString(reg RegisterTypeAsm) string {
	switch reg {
	case AX_REGISTER_ASM:
		return "eax"
	case R10_REGISTER_ASM:
		return "r10d"
	default:
		fmt.Println("unknown register:", reg)
		os.Exit(1)
	}

	return ""
}
