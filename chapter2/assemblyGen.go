package main

import (
	"fmt"
	"os"
)

/////////////////////////////////////////////////////////////////////////////////

type Program_Asm struct {
	fn *Function_Asm
}

/////////////////////////////////////////////////////////////////////////////////

type Function_Asm struct {
	name         string
	instructions []*Instruction_Asm
}

/////////////////////////////////////////////////////////////////////////////////

type InstructionTypeAsm int

const (
	MOV_INSTRUCTION_ASM InstructionTypeAsm = iota
	UNARY_INSTRUCTION_ASM
	ALLOCATE_STACK_INSTRUCTION_ASM
	RET_INSTRUCTION_ASM
)

type UnaryOperatorTypeAsm int

const (
	NOP_ASM UnaryOperatorTypeAsm = iota
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
	return NOP_ASM
}

type Instruction_Asm struct {
	typ  InstructionTypeAsm
	unOp UnaryOperatorTypeAsm
	src  *Operand_Asm
	dst  *Operand_Asm
}

/////////////////////////////////////////////////////////////////////////////////

type OperandTypeAsm int

const (
	IMMEDIATE_INT_OPERAND_ASM OperandTypeAsm = iota
	REGISTER_OPERAND_ASM
	PSEUDOREGISTER_OPERAND_ASM
	STACK_OPERAND_ASM
)

type RegisterTypeAsm int

const (
	AX_REGISTER_ASM RegisterTypeAsm = iota
	R10_REGISTER_ASM
)

type Operand_Asm struct {
	typ   OperandTypeAsm
	reg   RegisterTypeAsm
	value int32
	name  string
}

/////////////////////////////////////////////////////////////////////////////////

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

func (fn *Function_Tacky) convertToAsm() *Function_Asm {
	fnAsm := Function_Asm{name: fn.name}

	for _, instrTacky := range fn.body {
		convertedInstructions := instrTacky.convertToAsm()
		fnAsm.instructions = append(fnAsm.instructions, convertedInstructions...)
	}

	return &fnAsm
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Instruction_Tacky) convertToAsm() []*Instruction_Asm {
	instructions := []*Instruction_Asm{}

	switch instr.typ {
	case RETURN_INSTRUCTION_TACKY:
		src := instr.src.convertToAsm()
		dst := Operand_Asm{typ: REGISTER_OPERAND_ASM, reg: AX_REGISTER_ASM}
		movInstr := Instruction_Asm{typ: MOV_INSTRUCTION_ASM, src: src, dst: &dst}
		instructions = append(instructions, &movInstr)
		retInstr := Instruction_Asm{typ: RET_INSTRUCTION_ASM}
		instructions = append(instructions, &retInstr)
	case UNARY_INSTRUCTION_TACKY:
		src := instr.src.convertToAsm()
		dst := instr.dst.convertToAsm()
		movInstr := Instruction_Asm{typ: MOV_INSTRUCTION_ASM, src: src, dst: dst}
		instructions = append(instructions, &movInstr)
		// there's only one operand for the unary instruction so below we set both the src and dst
		// to the same thing in case we pick the wrong one in the next stage
		unaryInstr := Instruction_Asm{typ: UNARY_INSTRUCTION_ASM, unOp: convertUnaryOpToAsm(instr.unOp), src: dst, dst: dst}
		instructions = append(instructions, &unaryInstr)
	}

	return instructions
}

/////////////////////////////////////////////////////////////////////////////////

func (val *Value_Tacky) convertToAsm() *Operand_Asm {
	switch val.typ {
	case CONSTANT_VALUE_TACKY:
		return &Operand_Asm{typ: IMMEDIATE_INT_OPERAND_ASM, value: val.value}
	case VARIABLE_VALUE_TACKY:
		return &Operand_Asm{typ: PSEUDOREGISTER_OPERAND_ASM, name: val.name}
	}

	return nil
}

/////////////////////////////////////////////////////////////////////////////////

func (pr *Program_Asm) replacePseudoregisters() int32 {
	// TODO: need to handle more than one function

	var stackOffset int32 = 0
	nameToOffset := make(map[string]int32)

	for _, instr := range pr.fn.instructions {
		instr.src.replaceIfPseudoregister(&stackOffset, &nameToOffset)
		instr.dst.replaceIfPseudoregister(&stackOffset, &nameToOffset)
	}

	return stackOffset
}

/////////////////////////////////////////////////////////////////////////////////

func (op *Operand_Asm) replaceIfPseudoregister(stackOffset *int32, nameToOffset *map[string]int32) {
	if op == nil {
		return
	}

	if op.typ == PSEUDOREGISTER_OPERAND_ASM {
		op.typ = STACK_OPERAND_ASM

		existingOffset, alreadyExists := (*nameToOffset)[op.name]
		if alreadyExists {
			op.value = existingOffset
		} else {
			*stackOffset = *stackOffset - 4
			op.value = *stackOffset
			(*nameToOffset)[op.name] = *stackOffset
		}
	}
}

/////////////////////////////////////////////////////////////////////////////////

func (pr *Program_Asm) instructionFixup(stackOffset int32) {
	// TODO: need to handle more than one function

	// insert instruction to allocate space on the stack
	op := Operand_Asm{typ: IMMEDIATE_INT_OPERAND_ASM, value: -stackOffset}
	firstInstr := Instruction_Asm{typ: ALLOCATE_STACK_INSTRUCTION_ASM, src: &op}
	instructions := []*Instruction_Asm{&firstInstr}
	pr.fn.instructions = append(instructions, pr.fn.instructions...)

	// rewrite invalid Mov instructions, they can't have both operands be Stack operands
	instructions = []*Instruction_Asm{}

	for _, instr := range pr.fn.instructions {
		newInstrs := instr.fixInvalidInstr()
		instructions = append(instructions, newInstrs...)
	}

	pr.fn.instructions = instructions
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Instruction_Asm) fixInvalidInstr() []*Instruction_Asm {
	if instr.typ == MOV_INSTRUCTION_ASM && instr.src.typ == STACK_OPERAND_ASM && instr.dst.typ == STACK_OPERAND_ASM {
		intermediateOperand := Operand_Asm{typ: REGISTER_OPERAND_ASM, reg: R10_REGISTER_ASM}
		firstInstr := Instruction_Asm{typ: MOV_INSTRUCTION_ASM, src: instr.src, dst: &intermediateOperand}
		secondInstr := Instruction_Asm{typ: MOV_INSTRUCTION_ASM, src: &intermediateOperand, dst: instr.dst}
		return []*Instruction_Asm{&firstInstr, &secondInstr}
	} else {
		return []*Instruction_Asm{instr}
	}
}
