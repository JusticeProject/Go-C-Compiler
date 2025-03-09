package main

import "strconv"

/////////////////////////////////////////////////////////////////////////////////

var tempVarCounter int64 = -1

func makeTempVarName() string {
	// TODO: I could pass in the name of the function so the variables would be named something
	// like main.0, main.1, addFunction.2, etc.
	tempVarCounter++
	return "tmp." + strconv.FormatInt(tempVarCounter, 10)
}

/////////////////////////////////////////////////////////////////////////////////

type Program_Tacky struct {
	fn *Function_Tacky
}

/////////////////////////////////////////////////////////////////////////////////

type Function_Tacky struct {
	name string
	body []*Instruction_Tacky
}

/////////////////////////////////////////////////////////////////////////////////

type InstructionTypeTacky int

const (
	RETURN_INSTRUCTION_TACKY InstructionTypeTacky = iota
	UNARY_INSTRUCTION_TACKY
)

type Instruction_Tacky struct {
	typ  InstructionTypeTacky
	unOp UnaryOperatorType
	src  *Value_Tacky
	dst  *Value_Tacky
}

/////////////////////////////////////////////////////////////////////////////////

type ValueTypeTacky int

const (
	CONSTANT_VALUE_TACKY ValueTypeTacky = iota
	VARIABLE_VALUE_TACKY
)

type Value_Tacky struct {
	typ   ValueTypeTacky
	value int32
	name  string
}

/////////////////////////////////////////////////////////////////////////////////

// TAC = three-address code
func doTackyGen(ast Program) Program_Tacky {
	return ast.genTacky()
}

/////////////////////////////////////////////////////////////////////////////////

func (pr *Program) genTacky() Program_Tacky {
	// TODO: need to handle more than one function
	fnTac := pr.fn.genTacky()
	tacky := Program_Tacky{fn: fnTac}
	return tacky
}

/////////////////////////////////////////////////////////////////////////////////

func (fn *Function) genTacky() *Function_Tacky {
	bodyTac := fn.body.genTacky()
	fnTac := Function_Tacky{name: fn.name, body: bodyTac}
	return &fnTac
}

/////////////////////////////////////////////////////////////////////////////////

func (st *Statement) genTacky() []*Instruction_Tacky {
	switch st.typ {
	case RETURN_STATEMENT:
		instructions := []*Instruction_Tacky{}
		val, instructions := st.exp.genTacky(instructions)
		instr := Instruction_Tacky{typ: RETURN_INSTRUCTION_TACKY, src: val, dst: val}
		instructions = append(instructions, &instr)
		return instructions
	}

	return nil
}

/////////////////////////////////////////////////////////////////////////////////

func (exp *Expression) genTacky(instructions []*Instruction_Tacky) (*Value_Tacky, []*Instruction_Tacky) {
	switch exp.typ {
	case CONSTANT_INT_EXPRESSION:
		val := Value_Tacky{typ: CONSTANT_VALUE_TACKY, value: exp.intValue}
		return &val, instructions
	case UNARY_EXPRESSION:
		src, instructions := exp.innerExp.genTacky(instructions)
		dstName := makeTempVarName()
		dst := Value_Tacky{typ: VARIABLE_VALUE_TACKY, name: dstName}
		// TODO: will I need a helper function to convert the Unary Operator type to its TACKY equivalent?
		instr := Instruction_Tacky{typ: UNARY_INSTRUCTION_TACKY, unOp: exp.unOp, src: src, dst: &dst}
		instructions = append(instructions, &instr)
		return &dst, instructions
	}

	return nil, instructions
}
