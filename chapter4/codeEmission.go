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
		instr.instrEmitAsm(file)
	}
}

//###############################################################################
//###############################################################################
//###############################################################################

func (instr *Mov_Instruction_Asm) instrEmitAsm(file *os.File) {
	file.WriteString("\t" + "movl" + "\t" + instr.src.getOperandString() + ", " + instr.dst.getOperandString() + "\n")
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Unary_Instruction_Asm) instrEmitAsm(file *os.File) {
	file.WriteString("\t" + getUnaryOperatorString(instr.unOp) + "\t" + instr.src.getOperandString() + "\n")
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Binary_Instruction_Asm) instrEmitAsm(file *os.File) {
	file.WriteString("\t" + getBinaryOperatorString(instr.binOp) + "\t" + instr.src.getOperandString() + ", " +
		instr.dst.getOperandString() + "\n")
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *IDivide_Instruction_Asm) instrEmitAsm(file *os.File) {
	file.WriteString("\t" + "idivl" + "\t" + instr.divisor.getOperandString() + "\n")
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *CDQ_Sign_Extend_Instruction_Asm) instrEmitAsm(file *os.File) {
	file.WriteString("\t" + "cdq" + "\n")
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Allocate_Stack_Instruction_Asm) instrEmitAsm(file *os.File) {
	file.WriteString("\t" + "subq" + "\t" + instr.stackSize.getOperandString() + ", %rsp" + "\n")
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Ret_Instruction_Asm) instrEmitAsm(file *os.File) {
	// include the function epilogue instructions for restoring the stack
	file.WriteString("\t" + "movq" + "\t" + "%rbp, %rsp" + "\n")
	file.WriteString("\t" + "popq" + "\t" + "%rbp" + "\n")
	file.WriteString("\t" + "ret" + "\n")
}

//###############################################################################
//###############################################################################
//###############################################################################

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

//###############################################################################
//###############################################################################
//###############################################################################

func getBinaryOperatorString(binOp BinaryOperatorTypeAsm) string {
	switch binOp {
	case ADD_OPERATOR_ASM:
		return "addl"
	case SUB_OPERATOR_ASM:
		return "subl"
	case MULT_OPERATOR_ASM:
		return "imull"
	default:
		fmt.Println("unknown binary operator:", binOp)
		os.Exit(1)
	}

	return ""
}

//###############################################################################
//###############################################################################
//###############################################################################

func (op *Immediate_Int_Operand_Asm) getOperandString() string {
	return "$" + strconv.FormatInt(int64(op.value), 10)
}

/////////////////////////////////////////////////////////////////////////////////

func (op *Register_Operand_Asm) getOperandString() string {
	return "%" + getRegisterString(op.reg)
}

/////////////////////////////////////////////////////////////////////////////////

func (op *Pseudoregister_Operand_Asm) getOperandString() string {
	fmt.Println("cannot emit pseudoregister")
	os.Exit(1)
	return ""
}

/////////////////////////////////////////////////////////////////////////////////

func (op *Stack_Operand_Asm) getOperandString() string {
	return strconv.FormatInt(int64(op.value), 10) + "(%rbp)"
}

//###############################################################################
//###############################################################################
//###############################################################################

func getRegisterString(reg RegisterTypeAsm) string {
	switch reg {
	case AX_REGISTER_ASM:
		return "eax"
	case DX_REGISTER_ASM:
		return "edx"
	case R10_REGISTER_ASM:
		return "r10d"
	case R11_REGISTER_ASM:
		return "r11d"
	default:
		fmt.Println("unknown register:", reg)
		os.Exit(1)
	}

	return ""
}
