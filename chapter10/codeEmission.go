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
	for index, _ := range asm.topItems {
		asm.topItems[index].topLevelEmitAsm(file)
	}

	file.WriteString("\t.section\t.note.GNU-stack,\"\",@progbits\n")
}

/////////////////////////////////////////////////////////////////////////////////

func (fn *Function_Asm) topLevelEmitAsm(file *os.File) {
	if fn.global {
		file.WriteString("\t.globl " + fn.name + "\n")
	}
	file.WriteString("\t.text\n")
	file.WriteString(string(fn.name) + ":\n")

	// include the function prologue instructions for preparing the stack
	file.WriteString("\t" + "pushq" + "\t" + "%rbp" + "\n")
	file.WriteString("\t" + "movq" + "\t" + "%rsp, %rbp" + "\n")

	for _, instr := range fn.instructions {
		instr.instrEmitAsm(file)
	}
}

/////////////////////////////////////////////////////////////////////////////////

func (st *Static_Variable_Asm) topLevelEmitAsm(file *os.File) {
	if st.global {
		file.WriteString("\t" + ".globl " + st.name + "\n")
	}

	if st.initialValue == "0" {
		file.WriteString("\t" + ".bss" + "\n")
		file.WriteString("\t" + ".align 4" + "\n")
		file.WriteString(st.name + ":\n")
		file.WriteString("\t" + ".zero 4" + "\n")
	} else {
		file.WriteString("\t" + ".data" + "\n")
		file.WriteString("\t" + ".align 4" + "\n")
		file.WriteString(st.name + ":\n")
		file.WriteString("\t" + ".long " + st.initialValue + "\n")
	}
}

//###############################################################################
//###############################################################################
//###############################################################################

func (instr *Mov_Instruction_Asm) instrEmitAsm(file *os.File) {
	file.WriteString("\t" + "movl" + "\t" + instr.src.getOperandString(4) + ", " + instr.dst.getOperandString(4) + "\n")
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Unary_Instruction_Asm) instrEmitAsm(file *os.File) {
	file.WriteString("\t" + getUnaryOperatorString(instr.unOp) + "\t" + instr.src.getOperandString(4) + "\n")
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Binary_Instruction_Asm) instrEmitAsm(file *os.File) {
	file.WriteString("\t" + getBinaryOperatorString(instr.binOp) + "\t" + instr.src.getOperandString(4) + ", " +
		instr.dst.getOperandString(4) + "\n")
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Compare_Instruction_Asm) instrEmitAsm(file *os.File) {
	file.WriteString("\t" + "cmpl" + "\t" + instr.op1.getOperandString(4) + ", " + instr.op2.getOperandString(4) + "\n")
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *IDivide_Instruction_Asm) instrEmitAsm(file *os.File) {
	file.WriteString("\t" + "idivl" + "\t" + instr.divisor.getOperandString(4) + "\n")
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *CDQ_Sign_Extend_Instruction_Asm) instrEmitAsm(file *os.File) {
	file.WriteString("\t" + "cdq" + "\n")
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Jump_Instruction_Asm) instrEmitAsm(file *os.File) {
	file.WriteString("\t" + "jmp" + "\t\t" + ".L" + instr.target + "\n")
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Jump_Conditional_Instruction_Asm) instrEmitAsm(file *os.File) {
	file.WriteString("\t" + "j" + getConditionalCodeString(instr.code) + "\t\t" + ".L" + instr.target + "\n")
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Set_Conditional_Instruction_Asm) instrEmitAsm(file *os.File) {
	file.WriteString("\t" + "set" + getConditionalCodeString(instr.code) + "\t" + instr.dst.getOperandString(1) + "\n")
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Label_Instruction_Asm) instrEmitAsm(file *os.File) {
	file.WriteString(".L" + instr.name + ":" + "\n")
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Allocate_Stack_Instruction_Asm) instrEmitAsm(file *os.File) {
	file.WriteString("\t" + "subq" + "\t" + instr.stackSize.getOperandString(4) + ", %rsp" + "\n")
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Deallocate_Stack_Instruction_Asm) instrEmitAsm(file *os.File) {
	file.WriteString("\t" + "addq" + "\t" + instr.stackSize.getOperandString(4) + ", %rsp" + "\n")
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Push_Instruction_Asm) instrEmitAsm(file *os.File) {
	file.WriteString("\t" + "pushq" + "\t" + instr.op.getOperandString(8) + "\n")
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Call_Function_Asm) instrEmitAsm(file *os.File) {
	// need to find if the function we are calling is in the current binary object file or somewhere else
	entry, inTable := symbolTable[instr.name]
	if inTable && entry.defined {
		// It must be in the table and have a definition to use this calling method. If it's in the table
		// but not defined then it's just a function declaration so the definition is elsewhere.
		file.WriteString("\t" + "call" + "\t" + instr.name + "\n")
	} else {
		file.WriteString("\t" + "call" + "\t" + instr.name + "@PLT" + "\n")
	}
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
		fail("unknown unary operator")
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
		fail("unknown binary operator")
	}

	return ""
}

//###############################################################################
//###############################################################################
//###############################################################################

func (op *Immediate_Int_Operand_Asm) getOperandString(sizeBytes int) string {
	return "$" + op.value
}

/////////////////////////////////////////////////////////////////////////////////

func (op *Register_Operand_Asm) getOperandString(sizeBytes int) string {
	return "%" + getRegisterString(op.reg, sizeBytes)
}

/////////////////////////////////////////////////////////////////////////////////

func (op *Pseudoregister_Operand_Asm) getOperandString(sizeBytes int) string {
	fail("cannot emit pseudoregister")
	return ""
}

/////////////////////////////////////////////////////////////////////////////////

func (op *Stack_Operand_Asm) getOperandString(sizeBytes int) string {
	return strconv.FormatInt(int64(op.value), 10) + "(%rbp)"
}

/////////////////////////////////////////////////////////////////////////////////

func (op *Data_Operand_Asm) getOperandString(sizeBytes int) string {
	return op.name + "(%rip)"
}

//###############################################################################
//###############################################################################
//###############################################################################

func getConditionalCodeString(code ConditionalCodeAsm) string {
	switch code {
	case IS_EQUAL_CODE_ASM:
		return "e"
	case NOT_EQUAL_CODE_ASM:
		return "ne"
	case LESS_THAN_CODE_ASM:
		return "l"
	case LESS_OR_EQUAL_CODE_ASM:
		return "le"
	case GREATER_THAN_CODE_ASM:
		return "g"
	case GREATER_OR_EQUAL_CODE_ASM:
		return "ge"
	default:
		fail("unknown conditional code")
	}

	return ""
}

//###############################################################################
//###############################################################################
//###############################################################################

func getRegisterString(reg RegisterTypeAsm, sizeBytes int) string {
	switch reg {
	case AX_REGISTER_ASM:
		return getRegisterPrefix(sizeBytes) + "a" + getXRegisterSuffix(sizeBytes)
	case CX_REGISTER_ASM:
		return getRegisterPrefix(sizeBytes) + "c" + getXRegisterSuffix(sizeBytes)
	case DX_REGISTER_ASM:
		return getRegisterPrefix(sizeBytes) + "d" + getXRegisterSuffix(sizeBytes)
	case DI_REGISTER_ASM:
		return getRegisterPrefix(sizeBytes) + "di" + getIRegisterSuffix(sizeBytes)
	case SI_REGISTER_ASM:
		return getRegisterPrefix(sizeBytes) + "si" + getIRegisterSuffix(sizeBytes)
	case R8_REGISTER_ASM:
		return "r8" + getScratchRegisterSuffix(sizeBytes)
	case R9_REGISTER_ASM:
		return "r9" + getScratchRegisterSuffix(sizeBytes)
	case R10_REGISTER_ASM:
		return "r10" + getScratchRegisterSuffix(sizeBytes)
	case R11_REGISTER_ASM:
		return "r11" + getScratchRegisterSuffix(sizeBytes)
	default:
		fail("unknown register")
	}

	return ""
}

/////////////////////////////////////////////////////////////////////////////////

func getRegisterPrefix(sizeBytes int) string {
	switch sizeBytes {
	case 8:
		return "r"
	case 4:
		return "e"
	default:
		return ""
	}
}

/////////////////////////////////////////////////////////////////////////////////

func getXRegisterSuffix(sizeBytes int) string {
	switch sizeBytes {
	case 1:
		return "l"
	default:
		return "x"
	}
}

/////////////////////////////////////////////////////////////////////////////////

func getIRegisterSuffix(sizeBytes int) string {
	switch sizeBytes {
	case 1:
		return "l"
	default:
		return ""
	}
}

/////////////////////////////////////////////////////////////////////////////////

func getScratchRegisterSuffix(sizeBytes int) string {
	switch sizeBytes {
	case 8:
		return ""
	case 4:
		return "d"
	default:
		return "b"
	}
}
