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

type Compare_Instruction_Asm struct {
	op1 Operand_Asm
	op2 Operand_Asm
}

/////////////////////////////////////////////////////////////////////////////////

type IDivide_Instruction_Asm struct {
	divisor Operand_Asm
}

/////////////////////////////////////////////////////////////////////////////////

type CDQ_Sign_Extend_Instruction_Asm struct {
}

/////////////////////////////////////////////////////////////////////////////////

type Jump_Instruction_Asm struct {
	target string
}

/////////////////////////////////////////////////////////////////////////////////

type Jump_Conditional_Instruction_Asm struct {
	code   ConditionalCodeAsm
	target string
}

/////////////////////////////////////////////////////////////////////////////////

type Set_Conditional_Instruction_Asm struct {
	code ConditionalCodeAsm
	dst  Operand_Asm
}

/////////////////////////////////////////////////////////////////////////////////

type Label_Instruction_Asm struct {
	name string
}

/////////////////////////////////////////////////////////////////////////////////

type Allocate_Stack_Instruction_Asm struct {
	stackSize Operand_Asm
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
	case NOT_OPERATOR:
		fmt.Println("NOT_OPERATOR not converted directly to Asm")
		os.Exit(1)
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
	// TODO: do I really need to convert? or should I just use the same enum for both applications?
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
	getOperandString(sizeBytes int) string
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

type ConditionalCodeAsm int

const (
	NONE_CODE_ASM ConditionalCodeAsm = iota
	EQUAL_CODE_ASM
	NOT_EQUAL_CODE_ASM
	LESS_THAN_CODE_ASM
	LESS_OR_EQUAL_CODE_ASM
	GREATER_THAN_CODE_ASM
	GREATER_OR_EQUAL_CODE_ASM
)

func convertBinaryOpToCondition(binOp BinaryOperatorType) ConditionalCodeAsm {
	switch binOp {
	case EQUAL_OPERATOR:
		return EQUAL_CODE_ASM
	case NOT_EQUAL_OPERATOR:
		return NOT_EQUAL_CODE_ASM
	case LESS_THAN_OPERATOR:
		return LESS_THAN_CODE_ASM
	case LESS_OR_EQUAL_OPERATOR:
		return LESS_OR_EQUAL_CODE_ASM
	case GREATER_THAN_OPERATOR:
		return GREATER_THAN_CODE_ASM
	case GREATER_OR_EQUAL_OPERATOR:
		return GREATER_OR_EQUAL_CODE_ASM
	default:
		fmt.Println("unknown BinaryOperatorType when converting to code:", binOp)
		os.Exit(1)
	}

	return NONE_CODE_ASM
}

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

//###############################################################################
//###############################################################################
//###############################################################################

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
	if instr.unOp == NOT_OPERATOR {
		cmp := Compare_Instruction_Asm{op1: &Immediate_Int_Operand_Asm{0}, op2: instr.src.valueToAsm()}
		mov := Mov_Instruction_Asm{src: &Immediate_Int_Operand_Asm{0}, dst: instr.dst.valueToAsm()}
		setC := Set_Conditional_Instruction_Asm{code: EQUAL_CODE_ASM, dst: instr.dst.valueToAsm()}

		instructions := []Instruction_Asm{&cmp, &mov, &setC}
		return instructions
	} else {
		src := instr.src.valueToAsm()
		dst := instr.dst.valueToAsm()
		movInstr := Mov_Instruction_Asm{src: src, dst: dst}
		unaryInstr := Unary_Instruction_Asm{unOp: convertUnaryOpToAsm(instr.unOp), src: dst}

		instructions := []Instruction_Asm{&movInstr, &unaryInstr}
		return instructions
	}
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
		firstMov := Mov_Instruction_Asm{src: instr.src1.valueToAsm(), dst: &Register_Operand_Asm{AX_REGISTER_ASM}}
		cdq := CDQ_Sign_Extend_Instruction_Asm{}
		idiv := IDivide_Instruction_Asm{divisor: instr.src2.valueToAsm()}
		secondMov := Mov_Instruction_Asm{src: &Register_Operand_Asm{AX_REGISTER_ASM}, dst: instr.dst.valueToAsm()}
		instructions := []Instruction_Asm{&firstMov, &cdq, &idiv, &secondMov}
		return instructions
	} else if instr.binOp == REMAINDER_OPERATOR {
		firstMov := Mov_Instruction_Asm{src: instr.src1.valueToAsm(), dst: &Register_Operand_Asm{AX_REGISTER_ASM}}
		cdq := CDQ_Sign_Extend_Instruction_Asm{}
		idiv := IDivide_Instruction_Asm{divisor: instr.src2.valueToAsm()}
		secondMov := Mov_Instruction_Asm{src: &Register_Operand_Asm{DX_REGISTER_ASM}, dst: instr.dst.valueToAsm()}
		instructions := []Instruction_Asm{&firstMov, &cdq, &idiv, &secondMov}
		return instructions
	} else if instr.binOp == EQUAL_OPERATOR || instr.binOp == NOT_EQUAL_OPERATOR || instr.binOp == LESS_THAN_OPERATOR ||
		instr.binOp == LESS_OR_EQUAL_OPERATOR || instr.binOp == GREATER_THAN_OPERATOR || instr.binOp == GREATER_OR_EQUAL_OPERATOR {
		cmp := Compare_Instruction_Asm{op1: instr.src2.valueToAsm(), op2: instr.src1.valueToAsm()}
		mov := Mov_Instruction_Asm{src: &Immediate_Int_Operand_Asm{0}, dst: instr.dst.valueToAsm()}
		setC := Set_Conditional_Instruction_Asm{code: convertBinaryOpToCondition(instr.binOp), dst: instr.dst.valueToAsm()}
		instructions := []Instruction_Asm{&cmp, &mov, &setC}
		return instructions
	} else {
		fmt.Println("unknown Binary_Instruction_Tacky")
		os.Exit(1)
	}

	return nil
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Copy_Instruction_Tacky) instructionToAsm() []Instruction_Asm {
	mov := Mov_Instruction_Asm{src: instr.src.valueToAsm(), dst: instr.dst.valueToAsm()}
	return []Instruction_Asm{&mov}
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Jump_Instruction_Tacky) instructionToAsm() []Instruction_Asm {
	jmp := Jump_Instruction_Asm{instr.target}
	return []Instruction_Asm{&jmp}
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Jump_If_Zero_Instruction_Tacky) instructionToAsm() []Instruction_Asm {
	cmp := Compare_Instruction_Asm{op1: &Immediate_Int_Operand_Asm{0}, op2: instr.condition.valueToAsm()}
	jmpC := Jump_Conditional_Instruction_Asm{code: EQUAL_CODE_ASM, target: instr.target}
	return []Instruction_Asm{&cmp, &jmpC}
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Jump_If_Not_Zero_Instruction_Tacky) instructionToAsm() []Instruction_Asm {
	cmp := Compare_Instruction_Asm{op1: &Immediate_Int_Operand_Asm{0}, op2: instr.condition.valueToAsm()}
	jmpC := Jump_Conditional_Instruction_Asm{code: NOT_EQUAL_CODE_ASM, target: instr.target}
	return []Instruction_Asm{&cmp, &jmpC}
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Label_Instruction_Tacky) instructionToAsm() []Instruction_Asm {
	label := Label_Instruction_Asm{instr.name}
	return []Instruction_Asm{&label}
}

//###############################################################################
//###############################################################################
//###############################################################################

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

	// TODO: does this get reset to 0 for each function?
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
		case *Binary_Instruction_Asm:
			convertedInstr.src = replaceIfPseudoregister(convertedInstr.src, &stackOffset, &nameToOffset)
			convertedInstr.dst = replaceIfPseudoregister(convertedInstr.dst, &stackOffset, &nameToOffset)
			pr.fn.instructions[index] = convertedInstr
		case *IDivide_Instruction_Asm:
			convertedInstr.divisor = replaceIfPseudoregister(convertedInstr.divisor, &stackOffset, &nameToOffset)
			pr.fn.instructions[index] = convertedInstr
		case *Compare_Instruction_Asm:
			convertedInstr.op1 = replaceIfPseudoregister(convertedInstr.op1, &stackOffset, &nameToOffset)
			convertedInstr.op2 = replaceIfPseudoregister(convertedInstr.op2, &stackOffset, &nameToOffset)
			pr.fn.instructions[index] = convertedInstr
		case *Set_Conditional_Instruction_Asm:
			convertedInstr.dst = replaceIfPseudoregister(convertedInstr.dst, &stackOffset, &nameToOffset)
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
	firstInstr := Allocate_Stack_Instruction_Asm{stackSize: &op}
	instructions := []Instruction_Asm{&firstInstr}
	pr.fn.instructions = append(instructions, pr.fn.instructions...)

	// rewrite invalid instructions, they can't have both operands be Stack operands
	instructions = []Instruction_Asm{}

	for _, instr := range pr.fn.instructions {

		switch convertedInstr := instr.(type) {
		case *Mov_Instruction_Asm:
			newInstrs := convertedInstr.fixInvalidInstr()
			instructions = append(instructions, newInstrs...)
		case *Binary_Instruction_Asm:
			newInstrs := convertedInstr.fixInvalidInstr()
			instructions = append(instructions, newInstrs...)
		case *IDivide_Instruction_Asm:
			newInstrs := convertedInstr.fixInvalidInstr()
			instructions = append(instructions, newInstrs...)
		case *Compare_Instruction_Asm:
			newInstrs := convertedInstr.fixInvalidInstr()
			instructions = append(instructions, newInstrs...)
		default:
			// don't need to fix it, just add it to the list
			instructions = append(instructions, instr)
		}
	}

	pr.fn.instructions = instructions
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Mov_Instruction_Asm) fixInvalidInstr() []Instruction_Asm {
	_, srcIsStack := instr.src.(*Stack_Operand_Asm)
	_, dstIsStack := instr.dst.(*Stack_Operand_Asm)

	if srcIsStack && dstIsStack {
		intermediateOperand := Register_Operand_Asm{R10_REGISTER_ASM}
		firstInstr := Mov_Instruction_Asm{src: instr.src, dst: &intermediateOperand}
		secondInstr := Mov_Instruction_Asm{src: &intermediateOperand, dst: instr.dst}
		return []Instruction_Asm{&firstInstr, &secondInstr}
	}

	return []Instruction_Asm{instr}
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Binary_Instruction_Asm) fixInvalidInstr() []Instruction_Asm {
	if instr.binOp == ADD_OPERATOR_ASM || instr.binOp == SUB_OPERATOR_ASM {
		_, srcIsStack := instr.src.(*Stack_Operand_Asm)
		_, dstIsStack := instr.dst.(*Stack_Operand_Asm)

		if srcIsStack && dstIsStack {
			intermediateOperand := Register_Operand_Asm{R10_REGISTER_ASM}
			firstInstr := Mov_Instruction_Asm{src: instr.src, dst: &intermediateOperand}
			secondInstr := Binary_Instruction_Asm{binOp: instr.binOp, src: &intermediateOperand, dst: instr.dst}
			return []Instruction_Asm{&firstInstr, &secondInstr}
		}
	} else if instr.binOp == MULT_OPERATOR_ASM {
		_, dstIsStack := instr.dst.(*Stack_Operand_Asm)

		if dstIsStack {
			firstInstr := Mov_Instruction_Asm{src: instr.dst, dst: &Register_Operand_Asm{R11_REGISTER_ASM}}
			secondInstr := Binary_Instruction_Asm{binOp: instr.binOp, src: instr.src, dst: &Register_Operand_Asm{R11_REGISTER_ASM}}
			thirdInstr := Mov_Instruction_Asm{src: &Register_Operand_Asm{R11_REGISTER_ASM}, dst: instr.dst}
			return []Instruction_Asm{&firstInstr, &secondInstr, &thirdInstr}
		}
	}

	return []Instruction_Asm{instr}
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *IDivide_Instruction_Asm) fixInvalidInstr() []Instruction_Asm {
	_, isConstant := instr.divisor.(*Immediate_Int_Operand_Asm)

	if isConstant {
		firstInstr := Mov_Instruction_Asm{src: instr.divisor, dst: &Register_Operand_Asm{R10_REGISTER_ASM}}
		secondInstr := IDivide_Instruction_Asm{divisor: &Register_Operand_Asm{R10_REGISTER_ASM}}
		return []Instruction_Asm{&firstInstr, &secondInstr}
	}

	return []Instruction_Asm{instr}
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Compare_Instruction_Asm) fixInvalidInstr() []Instruction_Asm {
	_, op1IsStack := instr.op1.(*Stack_Operand_Asm)
	_, op2IsStack := instr.op2.(*Stack_Operand_Asm)
	_, op2IsConstant := instr.op2.(*Immediate_Int_Operand_Asm)

	if op1IsStack && op2IsStack {
		intermediateOperand := Register_Operand_Asm{R10_REGISTER_ASM}
		firstInstr := Mov_Instruction_Asm{src: instr.op1, dst: &intermediateOperand}
		secondInstr := Compare_Instruction_Asm{op1: &intermediateOperand, op2: instr.op2}
		return []Instruction_Asm{&firstInstr, &secondInstr}
	} else if op2IsConstant {
		firstInstr := Mov_Instruction_Asm{src: instr.op2, dst: &Register_Operand_Asm{R11_REGISTER_ASM}}
		secondInstr := Compare_Instruction_Asm{op1: instr.op1, op2: &Register_Operand_Asm{R11_REGISTER_ASM}}
		return []Instruction_Asm{&firstInstr, &secondInstr}
	}

	return []Instruction_Asm{instr}
}
