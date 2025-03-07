package main

/////////////////////////////////////////////////////////////////////////////////

type Program_Asm struct {
	fn *Function_Asm
}

/////////////////////////////////////////////////////////////////////////////////

type Function_Asm struct {
	name         Identifier
	instructions []*Instruction
}

/////////////////////////////////////////////////////////////////////////////////

type InstructionType int

const (
	MOV_INSTRUCTION InstructionType = iota
	RET_INSTRUCTION
)

type Instruction struct {
	typ InstructionType
	src *Operand
	dst *Operand
}

/////////////////////////////////////////////////////////////////////////////////

type OperandType int

const (
	IMMEDIATE_INT_OPERAND OperandType = iota
	REGISTER_OPERAND
)

type RegisterType int

const (
	EAX_REGISTER RegisterType = iota
)

type Operand struct {
	typ   OperandType
	reg   RegisterType
	value int32
}

/////////////////////////////////////////////////////////////////////////////////

func doAssemblyGen(ast Program) Program_Asm {
	asm := Program_Asm{}
	asm.fn = genFunction(ast.fn)

	return asm
}

/////////////////////////////////////////////////////////////////////////////////

func genFunction(fn *Function) *Function_Asm {
	retFn := Function_Asm{}
	retFn.name = fn.name
	retFn.instructions = genInstructions(fn.body)

	return &retFn
}

/////////////////////////////////////////////////////////////////////////////////

func genInstructions(st *Statement) []*Instruction {
	instructions := []*Instruction{}

	switch st.typ {
	case RETURN_STATEMENT:
		// TODO: need to determine src for more complicated Expressions, maybe genExpressionInstructions(Expression or ExpressionType)
		src := Operand{typ: IMMEDIATE_INT_OPERAND, value: st.exp.intValue}
		dst := Operand{typ: REGISTER_OPERAND, reg: EAX_REGISTER}
		movInstr := Instruction{typ: MOV_INSTRUCTION, src: &src, dst: &dst}
		instructions = append(instructions, &movInstr)
		retInstr := Instruction{typ: RET_INSTRUCTION}
		instructions = append(instructions, &retInstr)
	case IF_STATEMENT:

	}

	return instructions
}
