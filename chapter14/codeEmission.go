package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

/////////////////////////////////////////////////////////////////////////////////

func truncateDoubleToInteger(input string) string {
	if strings.Contains(input, ".") || strings.Contains(input, "e") || strings.Contains(input, "E") {
		double, err := strconv.ParseFloat(input, 64)
		if err != nil {
			fail("Failed to convert double to integer:", err.Error())
		}
		integer := int64(double)
		return strconv.FormatInt(integer, 10)
	} else {
		// it wasn't in floating point format, so just return it as-is
		return input
	}
}

/////////////////////////////////////////////////////////////////////////////////

func roundDouble(input string) string {
	value, err := strconv.ParseFloat(input, 64)
	if err != nil {
		if (value == math.Inf(1)) || (value == math.Inf(-1)) {
			fmt.Println("Warning:", input, "rounded to", value)
		} else {
			fail("Could not parse double:", err.Error())
		}
	}
	// Go does the rounding for us when we call ParseFloat and then FormatFloat
	result := strconv.FormatFloat(value, 'G', 24, 64)
	return result
}

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

	alignStr := strconv.FormatInt(int64(st.alignment), 10)
	typStr := ""
	if (st.initEnum == INITIAL_INT) || (st.initEnum == INITIAL_UNSIGNED_INT) {
		typStr = ".long "
		st.initialValue = truncateDoubleToInteger(st.initialValue)
	} else if (st.initEnum == INITIAL_LONG) || (st.initEnum == INITIAL_UNSIGNED_LONG) {
		typStr = ".quad "
		st.initialValue = truncateDoubleToInteger(st.initialValue)
	} else if st.initEnum == INITIAL_DOUBLE {
		typStr = ".double "
		st.initialValue = roundDouble(st.initialValue)
	}

	if (st.initialValue == "0") && (st.initEnum != INITIAL_DOUBLE) {
		file.WriteString("\t" + ".bss" + "\n")
		file.WriteString("\t" + ".align " + alignStr + "\n")
		file.WriteString(st.name + ":\n")
		file.WriteString("\t" + ".zero " + alignStr + "\n")
	} else {
		file.WriteString("\t" + ".data" + "\n")
		file.WriteString("\t" + ".align " + alignStr + "\n")
		file.WriteString(st.name + ":\n")
		file.WriteString("\t" + typStr + st.initialValue + "\n")
	}
}

/////////////////////////////////////////////////////////////////////////////////

func (st *Static_Constant_Asm) topLevelEmitAsm(file *os.File) {
	alignStr := strconv.FormatInt(int64(st.alignment), 10)

	if st.initEnum == INITIAL_DOUBLE {
		file.WriteString("\t" + ".section\t.rodata" + "\n")
		file.WriteString("\t" + ".align " + alignStr + "\n")
		file.WriteString(st.name + ":\n")
		st.initialValue = roundDouble(st.initialValue)
		file.WriteString("\t" + ".double " + st.initialValue + "\n")
	} else {
		fail("Static_Constant_Asm currently only supports doubles")
	}
}

//###############################################################################
//###############################################################################
//###############################################################################

func (instr *Mov_Instruction_Asm) instrEmitAsm(file *os.File) {
	file.WriteString("\t" + "mov" + getInstructionSuffix(instr.asmTyp) + "\t" +
		instr.src.getOperandString(instr.asmTyp) + ", " + instr.dst.getOperandString(instr.asmTyp) + "\n")
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Movsx_Instruction_Asm) instrEmitAsm(file *os.File) {
	file.WriteString("\t" + "movslq" + "\t" + instr.src.getOperandString(LONGWORD_ASM_TYPE) + ", " +
		instr.dst.getOperandString(QUADWORD_ASM_TYPE) + "\n")
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Move_Zero_Extend_Instruction_Asm) instrEmitAsm(file *os.File) {
	fail("Move_Zero_Extend_Instruction_Asm should have been rewritten in the previous step")
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Cvttsd2si_Double_To_Int_Instruction_Asm) instrEmitAsm(file *os.File) {
	file.WriteString("\t" + "cvttsd2si" + getInstructionSuffix(instr.dstAsmType) + "\t" +
		instr.src.getOperandString(instr.dstAsmType) + ", " + instr.dst.getOperandString(instr.dstAsmType) + "\n")
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Cvtsi2sd_Int_To_Double_Instruction_Asm) instrEmitAsm(file *os.File) {
	file.WriteString("\t" + "cvtsi2sd" + getInstructionSuffix(instr.srcAsmType) + "\t" +
		instr.src.getOperandString(instr.srcAsmType) + ", " + instr.dst.getOperandString(instr.srcAsmType) + "\n")
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Unary_Instruction_Asm) instrEmitAsm(file *os.File) {
	file.WriteString("\t" + getUnaryOperatorString(instr.unOp) + getInstructionSuffix(instr.asmTyp) + "\t" +
		instr.src.getOperandString(instr.asmTyp) + "\n")
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Binary_Instruction_Asm) instrEmitAsm(file *os.File) {
	file.WriteString("\t" + getBinaryOperatorString(instr.binOp, instr.asmTyp) + "\t" +
		instr.src.getOperandString(instr.asmTyp) + ", " + instr.dst.getOperandString(instr.asmTyp) + "\n")
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Compare_Instruction_Asm) instrEmitAsm(file *os.File) {
	var cmpStr string
	if instr.asmTyp == DOUBLE_ASM_TYPE {
		cmpStr = "comi"
	} else {
		cmpStr = "cmp"
	}
	file.WriteString("\t" + cmpStr + getInstructionSuffix(instr.asmTyp) + "\t" +
		instr.op1.getOperandString(instr.asmTyp) + ", " + instr.op2.getOperandString(instr.asmTyp) + "\n")
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *IDivide_Instruction_Asm) instrEmitAsm(file *os.File) {
	file.WriteString("\t" + "idiv" + getInstructionSuffix(instr.asmTyp) + "\t" + instr.divisor.getOperandString(instr.asmTyp) + "\n")
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Divide_Instruction_Asm) instrEmitAsm(file *os.File) {
	file.WriteString("\t" + "div" + getInstructionSuffix(instr.asmTyp) + "\t" + instr.divisor.getOperandString(instr.asmTyp) + "\n")
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *CDQ_Sign_Extend_Instruction_Asm) instrEmitAsm(file *os.File) {
	if instr.asmTyp == QUADWORD_ASM_TYPE {
		file.WriteString("\t" + "cqo" + "\n")
	} else {
		file.WriteString("\t" + "cdq" + "\n")
	}
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
	file.WriteString("\t" + "set" + getConditionalCodeString(instr.code) + "\t" + instr.dst.getOperandString(BYTE_ASM_TYPE) + "\n")
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Label_Instruction_Asm) instrEmitAsm(file *os.File) {
	file.WriteString(".L" + instr.name + ":" + "\n")
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Push_Instruction_Asm) instrEmitAsm(file *os.File) {
	file.WriteString("\t" + "pushq" + "\t" + instr.op.getOperandString(QUADWORD_ASM_TYPE) + "\n")
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Call_Function_Asm) instrEmitAsm(file *os.File) {
	// need to find if the function we are calling is in the current binary object file or somewhere else
	entry, inTable := symbolTableBackend[instr.name]
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
		return "neg"
	case NOT_OPERATOR_ASM:
		return "not"
	case SHIFT_RIGHT_OPERATOR_ASM:
		return "shr"
	default:
		fail("unknown unary operator")
	}

	return ""
}

//###############################################################################
//###############################################################################
//###############################################################################

func getBinaryOperatorString(binOp BinaryOperatorTypeAsm, asmTyp AssemblyTypeEnum) string {
	switch binOp {
	case ADD_OPERATOR_ASM:
		return "add" + getInstructionSuffix(asmTyp)
	case SUB_OPERATOR_ASM:
		return "sub" + getInstructionSuffix(asmTyp)
	case MULT_OPERATOR_ASM:
		if asmTyp == DOUBLE_ASM_TYPE {
			return "mul" + getInstructionSuffix(asmTyp)
		} else {
			return "imul" + getInstructionSuffix(asmTyp)
		}
	case DIV_DOUBLE_OPERATOR_ASM:
		return "div" + getInstructionSuffix(asmTyp)
	case AND_OPERATOR_ASM:
		return "and" + getInstructionSuffix(asmTyp)
	case OR_OPERATOR_ASM:
		return "or" + getInstructionSuffix(asmTyp)
	case XOR_OPERATOR_ASM:
		if asmTyp == DOUBLE_ASM_TYPE {
			return "xorpd"
		} else {
			return "xor" + getInstructionSuffix(asmTyp)
		}
	default:
		fail("unknown binary operator")
	}

	return ""
}

//###############################################################################
//###############################################################################
//###############################################################################

func (op *Immediate_Int_Operand_Asm) getOperandString(asmTyp AssemblyTypeEnum) string {
	// TODO: truncate the value since it could be in floating point format??
	// use helper function truncateDoubleToInteger()??
	return "$" + op.value
}

/////////////////////////////////////////////////////////////////////////////////

func (op *Register_Operand_Asm) getOperandString(asmTyp AssemblyTypeEnum) string {
	return "%" + getRegisterString(op.reg, asmTyp)
}

/////////////////////////////////////////////////////////////////////////////////

func (op *Pseudoregister_Operand_Asm) getOperandString(asmTyp AssemblyTypeEnum) string {
	fail("cannot emit pseudoregister")
	return ""
}

/////////////////////////////////////////////////////////////////////////////////

func (op *Stack_Operand_Asm) getOperandString(asmTyp AssemblyTypeEnum) string {
	return strconv.FormatInt(int64(op.value), 10) + "(%rbp)"
}

/////////////////////////////////////////////////////////////////////////////////

func (op *Data_Operand_Asm) getOperandString(asmTyp AssemblyTypeEnum) string {
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
	case LESS_THAN_CODE_UNSIGNED_ASM:
		return "b"
	case LESS_OR_EQUAL_CODE_UNSIGNED_ASM:
		return "be"
	case GREATER_THAN_CODE_UNSIGNED_ASM:
		return "a"
	case GREATER_OR_EQUAL_CODE_UNSIGNED_ASM:
		return "ae"
	default:
		fail("unknown conditional code")
	}

	return ""
}

//###############################################################################
//###############################################################################
//###############################################################################

func getInstructionSuffix(asmTyp AssemblyTypeEnum) string {
	switch asmTyp {
	case QUADWORD_ASM_TYPE:
		return "q"
	case LONGWORD_ASM_TYPE:
		return "l"
	case DOUBLE_ASM_TYPE:
		return "sd"
	default:
		fail("unknown AssemblyTypeEnum")
	}
	return ""
}

/////////////////////////////////////////////////////////////////////////////////

func getRegisterString(reg RegisterTypeAsm, asmTyp AssemblyTypeEnum) string {
	switch reg {
	case AX_REGISTER_ASM:
		return getRegisterPrefix(asmTyp) + "a" + getXRegisterSuffix(asmTyp)
	case CX_REGISTER_ASM:
		return getRegisterPrefix(asmTyp) + "c" + getXRegisterSuffix(asmTyp)
	case DX_REGISTER_ASM:
		return getRegisterPrefix(asmTyp) + "d" + getXRegisterSuffix(asmTyp)
	case DI_REGISTER_ASM:
		return getRegisterPrefix(asmTyp) + "di" + getIRegisterSuffix(asmTyp)
	case SI_REGISTER_ASM:
		return getRegisterPrefix(asmTyp) + "si" + getIRegisterSuffix(asmTyp)
	case R8_REGISTER_ASM:
		return "r8" + getScratchRegisterSuffix(asmTyp)
	case R9_REGISTER_ASM:
		return "r9" + getScratchRegisterSuffix(asmTyp)
	case R10_REGISTER_ASM:
		return "r10" + getScratchRegisterSuffix(asmTyp)
	case R11_REGISTER_ASM:
		return "r11" + getScratchRegisterSuffix(asmTyp)
	case SP_REGISTER_ASM:
		return "rsp"
	case XMM0_REGISTER_ASM:
		return "xmm0"
	case XMM1_REGISTER_ASM:
		return "xmm1"
	case XMM2_REGISTER_ASM:
		return "xmm2"
	case XMM3_REGISTER_ASM:
		return "xmm3"
	case XMM4_REGISTER_ASM:
		return "xmm4"
	case XMM5_REGISTER_ASM:
		return "xmm5"
	case XMM6_REGISTER_ASM:
		return "xmm6"
	case XMM7_REGISTER_ASM:
		return "xmm7"
	case XMM14_REGISTER_ASM:
		return "xmm14"
	case XMM15_REGISTER_ASM:
		return "xmm15"
	default:
		fail("unknown register")
	}

	return ""
}

/////////////////////////////////////////////////////////////////////////////////

func getRegisterPrefix(asmTyp AssemblyTypeEnum) string {
	switch asmTyp {
	case QUADWORD_ASM_TYPE:
		return "r"
	case LONGWORD_ASM_TYPE:
		return "e"
	default:
		return ""
	}
}

/////////////////////////////////////////////////////////////////////////////////

func getXRegisterSuffix(asmTyp AssemblyTypeEnum) string {
	switch asmTyp {
	case BYTE_ASM_TYPE:
		return "l"
	default:
		return "x"
	}
}

/////////////////////////////////////////////////////////////////////////////////

func getIRegisterSuffix(asmTyp AssemblyTypeEnum) string {
	switch asmTyp {
	case BYTE_ASM_TYPE:
		return "l"
	default:
		return ""
	}
}

/////////////////////////////////////////////////////////////////////////////////

func getScratchRegisterSuffix(asmTyp AssemblyTypeEnum) string {
	switch asmTyp {
	case QUADWORD_ASM_TYPE:
		return ""
	case LONGWORD_ASM_TYPE:
		return "d"
	default:
		return "b"
	}
}
