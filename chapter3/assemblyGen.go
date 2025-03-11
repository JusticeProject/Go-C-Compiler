package main

import (
	"fmt"
	"os"
)

//###############################################################################
//###############################################################################
//###############################################################################

type Program_Asm struct {
	fn Function_Asm
}

//###############################################################################
//###############################################################################
//###############################################################################

type Function_Asm struct {
	name         string
	instructions []Instruction_Asm
}

//###############################################################################
//###############################################################################
//###############################################################################

type Instruction_Asm interface {
	instrEmitAsm(file *os.File)
}

/////////////////////////////////////////////////////////////////////////////////

type Mov_Instruction_Asm struct {
	src Operand_Asm
	dst Operand_Asm
}

/////////////////////////////////////////////////////////////////////////////////

type Unary_Instruction_Asm struct {
	unOp UnaryOperatorTypeAsm
	src  Operand_Asm
}

/////////////////////////////////////////////////////////////////////////////////

type Binary_Instruction_Asm struct {
	binOp BinaryOperatorTypeAsm
	src   Operand_Asm
	dst   Operand_Asm
}

/////////////////////////////////////////////////////////////////////////////////

type IDivide_Instruction_Asm struct {
	src1 Operand_Asm
	src2 Operand_Asm
	dst  Operand_Asm
}

/////////////////////////////////////////////////////////////////////////////////

type CDQ_Sign_Extend_Instruction_Asm struct {
}

/////////////////////////////////////////////////////////////////////////////////

type Allocate_Stack_Instruction_Asm struct {
	op Operand_Asm
}

/////////////////////////////////////////////////////////////////////////////////

type Ret_Instruction_Asm struct {
}

//###############################################################################
//###############################################################################
//###############################################################################

type UnaryOperatorTypeAsm int

const (
	NOP_UNARY_ASM UnaryOperatorTypeAsm = iota
	NEGATE_OPERATOR_ASM
	NOT_OPERATOR_ASM
)

func convertUnaryOpToAsm(unOp UnaryOperatorType) UnaryOperatorTypeAsm {
	switch unOp {
	case COMPLEMENT_OPERATOR:
		return NOT_OPERATOR_ASM
	case NEGATE_OPERATOR:
		return NEGATE_OPERATOR_ASM
	default:
		fmt.Println("unknown UnaryOperatorType:", unOp)
		os.Exit(1)
	}
	return NOP_UNARY_ASM
}

//###############################################################################
//###############################################################################
//###############################################################################

type BinaryOperatorTypeAsm int

const (
	NOP_BINARY_ASM BinaryOperatorTypeAsm = iota
	ADD_OPERATOR_ASM
	SUB_OPERATOR_ASM
	MULT_OPERATOR_ASM
)

func convertBinaryOpToAsm(binOp BinaryOperatorType) BinaryOperatorTypeAsm {
	switch binOp {
	case ADD_OPERATOR:
		return ADD_OPERATOR_ASM
	case SUBTRACT_OPERATOR:
		return SUB_OPERATOR_ASM
	case MULTIPLY_OPERATOR:
		return MULT_OPERATOR_ASM
	default:
		fmt.Println("unknown BinaryOperatorType:", binOp)
		os.Exit(1)
	}
	return NOP_BINARY_ASM
}

//###############################################################################
//###############################################################################
//###############################################################################

type Operand_Asm interface {
	getOperandString() string
}

/////////////////////////////////////////////////////////////////////////////////

type Immediate_Int_Operand_Asm struct {
	value int32
}

/////////////////////////////////////////////////////////////////////////////////

type Register_Operand_Asm struct {
	reg RegisterTypeAsm
}

/////////////////////////////////////////////////////////////////////////////////

type Pseudoregister_Operand_Asm struct {
	name string
}

/////////////////////////////////////////////////////////////////////////////////

type Stack_Operand_Asm struct {
	value int32
}

//###############################################################################
//###############################################################################
//###############################################################################

type RegisterTypeAsm int

const (
	AX_REGISTER_ASM RegisterTypeAsm = iota
	DX_REGISTER_ASM
	R10_REGISTER_ASM
	R11_REGISTER_ASM
)

//###############################################################################
//###############################################################################
//###############################################################################

func doAssemblyGen(tacky Program_Tacky) Program_Asm {
	asm := tacky.convertToAsm()
	stackOffset := asm.replacePseudoregisters()
	asm.instructionFixup(stackOffset)

	return asm
}

/////////////////////////////////////////////////////////////////////////////////

func (pr *Program_Tacky) convertToAsm() Program_Asm {
	// TODO: need to handle more than one function in the program
	fnAsm := pr.fn.convertToAsm()
	asm := Program_Asm{fn: fnAsm}
	return asm
}

/////////////////////////////////////////////////////////////////////////////////

func (fn *Function_Tacky) convertToAsm() Function_Asm {
	fnAsm := Function_Asm{name: fn.name}

	for _, instrTacky := range fn.body {
		convertedInstructions := instrTacky.instructionToAsm()
		fnAsm.instructions = append(fnAsm.instructions, convertedInstructions...)
	}

	return fnAsm
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Return_Instruction_Tacky) instructionToAsm() []Instruction_Asm {
	src := instr.val.valueToAsm()
	dst := Register_Operand_Asm{reg: AX_REGISTER_ASM}
	movInstr := Mov_Instruction_Asm{src: src, dst: &dst}
	retInstr := Ret_Instruction_Asm{}

	instructions := []Instruction_Asm{&movInstr, &retInstr}
	return instructions
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Unary_Instruction_Tacky) instructionToAsm() []Instruction_Asm {
	src := instr.src.valueToAsm()
	dst := instr.dst.valueToAsm()
	movInstr := Mov_Instruction_Asm{src: src, dst: dst}
	unaryInstr := Unary_Instruction_Asm{unOp: convertUnaryOpToAsm(instr.unOp), src: dst}

	instructions := []Instruction_Asm{&movInstr, &unaryInstr}
	return instructions
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Binary_Instruction_Tacky) instructionToAsm() []Instruction_Asm {
	if instr.binOp == ADD_OPERATOR || instr.binOp == SUBTRACT_OPERATOR || instr.binOp == MULTIPLY_OPERATOR {
		src1 := instr.src1.valueToAsm()
		dst := instr.dst.valueToAsm()
		movInstr := Mov_Instruction_Asm{src: src1, dst: dst}

		src2 := instr.src2.valueToAsm()
		binInstr := Binary_Instruction_Asm{binOp: convertBinaryOpToAsm(instr.binOp), src: src2, dst: dst}

		instructions := []Instruction_Asm{&movInstr, &binInstr}
		return instructions
	} else if instr.binOp == DIVIDE_OPERATOR {
		// TODO:
	} else if instr.binOp == REMAINDER_OPERATOR {
		// TODO:
	}

	return nil
}

/////////////////////////////////////////////////////////////////////////////////

func (val *Constant_Value_Tacky) valueToAsm() Operand_Asm {
	return &Immediate_Int_Operand_Asm{value: val.value}
}

/////////////////////////////////////////////////////////////////////////////////

func (val *Variable_Value_Tacky) valueToAsm() Operand_Asm {
	return &Pseudoregister_Operand_Asm{name: val.name}
}

//###############################################################################
//###############################################################################
//###############################################################################

func (pr *Program_Asm) replacePseudoregisters() int32 {
	// TODO: need to handle more than one function

	var stackOffset int32 = 0
	nameToOffset := make(map[string]int32)

	for index, _ := range pr.fn.instructions {
		switch convertedInstr := pr.fn.instructions[index].(type) {
		case *Mov_Instruction_Asm:
			convertedInstr.src = replaceIfPseudoregister(convertedInstr.src, &stackOffset, &nameToOffset)
			convertedInstr.dst = replaceIfPseudoregister(convertedInstr.dst, &stackOffset, &nameToOffset)
			pr.fn.instructions[index] = convertedInstr
		case *Unary_Instruction_Asm:
			convertedInstr.src = replaceIfPseudoregister(convertedInstr.src, &stackOffset, &nameToOffset)
			pr.fn.instructions[index] = convertedInstr
		}

	}

	return stackOffset
}

/////////////////////////////////////////////////////////////////////////////////

func replaceIfPseudoregister(op Operand_Asm, stackOffset *int32, nameToOffset *map[string]int32) Operand_Asm {
	if op == nil {
		return nil
	}

	convertedOp, isPseudo := op.(*Pseudoregister_Operand_Asm)

	if !isPseudo {
		return op
	}

	existingOffset, alreadyExists := (*nameToOffset)[convertedOp.name]
	if alreadyExists {
		return &Stack_Operand_Asm{value: existingOffset}
	} else {
		*stackOffset = *stackOffset - 4
		(*nameToOffset)[convertedOp.name] = *stackOffset
		return &Stack_Operand_Asm{value: *stackOffset}
	}
}

//###############################################################################
//###############################################################################
//###############################################################################

func (pr *Program_Asm) instructionFixup(stackOffset int32) {
	// TODO: need to handle more than one function

	// insert instruction to allocate space on the stack
	op := Immediate_Int_Operand_Asm{value: -stackOffset}
	firstInstr := Allocate_Stack_Instruction_Asm{op: &op}
	instructions := []Instruction_Asm{&firstInstr}
	pr.fn.instructions = append(instructions, pr.fn.instructions...)

	// rewrite invalid instructions, they can't have both operands be Stack operands
	instructions = []Instruction_Asm{}

	for index, _ := range pr.fn.instructions {
		newInstrs := fixInvalidInstr(pr.fn.instructions[index])
		instructions = append(instructions, newInstrs...)
	}

	pr.fn.instructions = instructions
}

/////////////////////////////////////////////////////////////////////////////////

func fixInvalidInstr(instr Instruction_Asm) []Instruction_Asm {
	switch convertedInstr := instr.(type) {
	case *Mov_Instruction_Asm:
		_, srcIsStack := convertedInstr.src.(*Stack_Operand_Asm)
		_, dstIsStack := convertedInstr.dst.(*Stack_Operand_Asm)

		if srcIsStack && dstIsStack {
			intermediateOperand := Register_Operand_Asm{reg: R10_REGISTER_ASM}
			firstInstr := Mov_Instruction_Asm{src: convertedInstr.src, dst: &intermediateOperand}
			secondInstr := Mov_Instruction_Asm{src: &intermediateOperand, dst: convertedInstr.dst}
			return []Instruction_Asm{&firstInstr, &secondInstr}
		}
	}

	return []Instruction_Asm{instr}
}
