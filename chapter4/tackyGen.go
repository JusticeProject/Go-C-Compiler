package main

import "strconv"

//###############################################################################
//###############################################################################
//###############################################################################

var tempVarCounter int64 = -1

func makeTempVarName() string {
	// TODO: I could pass in the name of the function so the variables would be named something
	// like main.0, main.1, addFunction.2, etc.
	tempVarCounter++
	return "tmp." + strconv.FormatInt(tempVarCounter, 10)
}

//###############################################################################
//###############################################################################
//###############################################################################

type Program_Tacky struct {
	fn Function_Tacky
}

//###############################################################################
//###############################################################################
//###############################################################################

type Function_Tacky struct {
	name string
	body []Instruction_Tacky
}

//###############################################################################
//###############################################################################
//###############################################################################

type Instruction_Tacky interface {
	instructionToAsm() []Instruction_Asm
}

/////////////////////////////////////////////////////////////////////////////////

type Return_Instruction_Tacky struct {
	val Value_Tacky
}

/////////////////////////////////////////////////////////////////////////////////

type Unary_Instruction_Tacky struct {
	unOp UnaryOperatorType
	src  Value_Tacky
	dst  Value_Tacky
}

/////////////////////////////////////////////////////////////////////////////////

type Binary_Instruction_Tacky struct {
	binOp BinaryOperatorType
	src1  Value_Tacky
	src2  Value_Tacky
	dst   Value_Tacky
}

//###############################################################################
//###############################################################################
//###############################################################################

type Value_Tacky interface {
	valueToAsm() Operand_Asm
}

/////////////////////////////////////////////////////////////////////////////////

type Constant_Value_Tacky struct {
	value int32
}

/////////////////////////////////////////////////////////////////////////////////

type Variable_Value_Tacky struct {
	name string
}

//###############################################################################
//###############################################################################
//###############################################################################

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

func (fn *Function) genTacky() Function_Tacky {
	bodyTac := fn.body.statementToTacky()
	fnTac := Function_Tacky{name: fn.name, body: bodyTac}
	return fnTac
}

//###############################################################################
//###############################################################################
//###############################################################################

func (st *Return_Statement) statementToTacky() []Instruction_Tacky {
	instructions := []Instruction_Tacky{}
	val, instructions := st.exp.expToTacky(instructions)
	instr := Return_Instruction_Tacky{val: val}
	instructions = append(instructions, &instr)
	return instructions
}

//###############################################################################
//###############################################################################
//###############################################################################

func (exp *Constant_Int_Expression) expToTacky(instructions []Instruction_Tacky) (Value_Tacky, []Instruction_Tacky) {
	val := Constant_Value_Tacky{value: exp.intValue}
	return &val, instructions
}

/////////////////////////////////////////////////////////////////////////////////

func (exp *Unary_Expression) expToTacky(instructions []Instruction_Tacky) (Value_Tacky, []Instruction_Tacky) {
	src, instructions := exp.innerExp.expToTacky(instructions)
	dstName := makeTempVarName()
	dst := Variable_Value_Tacky{name: dstName}
	// TODO: will I need a helper function to convert the Unary Operator type to its TACKY equivalent?
	instr := Unary_Instruction_Tacky{unOp: exp.unOp, src: src, dst: &dst}
	instructions = append(instructions, &instr)
	return &dst, instructions
}

/////////////////////////////////////////////////////////////////////////////////

func (exp *Binary_Expression) expToTacky(instructions []Instruction_Tacky) (Value_Tacky, []Instruction_Tacky) {
	src1, instructions := exp.firstExp.expToTacky(instructions)
	src2, instructions := exp.secExp.expToTacky(instructions)
	dstName := makeTempVarName()
	dst := Variable_Value_Tacky{name: dstName}
	instr := Binary_Instruction_Tacky{binOp: exp.binOp, src1: src1, src2: src2, dst: &dst}
	instructions = append(instructions, &instr)
	return &dst, instructions
}
