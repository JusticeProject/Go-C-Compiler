package main

import "strconv"

//###############################################################################
//###############################################################################
//###############################################################################

var tempVarCounter int64 = -1

func makeTempVarName(prefix string) string {
	// TODO: I could pass in the name of the function so the variables would be named something
	// like main.0, main.1, addFunction.2, etc.
	tempVarCounter++
	if len(prefix) > 0 {
		return prefix + "." + strconv.FormatInt(tempVarCounter, 10)
	} else {
		return "tmp." + strconv.FormatInt(tempVarCounter, 10)
	}
}

/////////////////////////////////////////////////////////////////////////////////

var labelCounter int64 = -1

func makeLabelName(name string) string {
	labelCounter++
	return name + strconv.FormatInt(labelCounter, 10)
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

/////////////////////////////////////////////////////////////////////////////////

type Copy_Instruction_Tacky struct {
	src Value_Tacky
	dst Value_Tacky
}

/////////////////////////////////////////////////////////////////////////////////

type Jump_Instruction_Tacky struct {
	target string
}

/////////////////////////////////////////////////////////////////////////////////

type Jump_If_Zero_Instruction_Tacky struct {
	condition Value_Tacky
	target    string
}

/////////////////////////////////////////////////////////////////////////////////

type Jump_If_Not_Zero_Instruction_Tacky struct {
	condition Value_Tacky
	target    string
}

/////////////////////////////////////////////////////////////////////////////////

type Label_Instruction_Tacky struct {
	name string
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
	//bodyTac := fn.body.statementToTacky()
	//fnTac := Function_Tacky{name: fn.name, body: bodyTac}
	// TODO:
	return Function_Tacky{}
}

//###############################################################################
//###############################################################################
//###############################################################################

func (b *Block_Statement) blockToTacky() []Instruction_Tacky {
	// TODO:
	return []Instruction_Tacky{}
}

/////////////////////////////////////////////////////////////////////////////////

func (b *Block_Declaration) blockToTacky() []Instruction_Tacky {
	// TODO:
	return []Instruction_Tacky{}
}

//###############################################################################
//###############################################################################
//###############################################################################

func (d *Declaration) genTacky() []Instruction_Tacky {
	// TODO:
	return []Instruction_Tacky{}
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

/////////////////////////////////////////////////////////////////////////////////

func (st *Expression_Statement) statementToTacky() []Instruction_Tacky {
	// TODO:
	return []Instruction_Tacky{}
}

/////////////////////////////////////////////////////////////////////////////////

func (st *Null_Statement) statementToTacky() []Instruction_Tacky {
	// TODO:
	return []Instruction_Tacky{}
}

//###############################################################################
//###############################################################################
//###############################################################################

func (exp *Constant_Int_Expression) expToTacky(instructions []Instruction_Tacky) (Value_Tacky, []Instruction_Tacky) {
	val := Constant_Value_Tacky{value: exp.intValue}
	return &val, instructions
}

/////////////////////////////////////////////////////////////////////////////////

func (exp *Variable_Expression) expToTacky(instructions []Instruction_Tacky) (Value_Tacky, []Instruction_Tacky) {
	// TODO:
	return nil, []Instruction_Tacky{}
}

/////////////////////////////////////////////////////////////////////////////////

func (exp *Unary_Expression) expToTacky(instructions []Instruction_Tacky) (Value_Tacky, []Instruction_Tacky) {
	src, instructions := exp.innerExp.expToTacky(instructions)
	dstName := makeTempVarName("")
	dst := Variable_Value_Tacky{name: dstName}
	// TODO: will I need a helper function to convert the Unary Operator type to its TACKY equivalent?
	instr := Unary_Instruction_Tacky{unOp: exp.unOp, src: src, dst: &dst}
	instructions = append(instructions, &instr)
	return &dst, instructions
}

/////////////////////////////////////////////////////////////////////////////////

func (exp *Binary_Expression) expToTacky(instructions []Instruction_Tacky) (Value_Tacky, []Instruction_Tacky) {
	// some operators can short-circuit on the first expression, so we handle them differently
	if exp.binOp == AND_OPERATOR {
		v1, instructions := exp.firstExp.expToTacky(instructions)
		false_label := makeLabelName("and_false")
		j1 := Jump_If_Zero_Instruction_Tacky{condition: v1, target: false_label}
		instructions = append(instructions, &j1)
		v2, instructions := exp.secExp.expToTacky(instructions)
		j2 := Jump_If_Zero_Instruction_Tacky{condition: v2, target: false_label}
		instructions = append(instructions, &j2)
		result := Variable_Value_Tacky{makeTempVarName("")}
		cp1 := Copy_Instruction_Tacky{src: &Constant_Value_Tacky{1}, dst: &result}
		instructions = append(instructions, &cp1)
		end := makeLabelName("end")
		j3 := Jump_Instruction_Tacky{end}
		instructions = append(instructions, &j3)
		lb1 := Label_Instruction_Tacky{false_label}
		instructions = append(instructions, &lb1)
		cp2 := Copy_Instruction_Tacky{src: &Constant_Value_Tacky{0}, dst: &result}
		instructions = append(instructions, &cp2)
		lb2 := Label_Instruction_Tacky{end}
		instructions = append(instructions, &lb2)
		return &result, instructions
	} else if exp.binOp == OR_OPERATOR {
		v1, instructions := exp.firstExp.expToTacky(instructions)
		true_label := makeLabelName("or_true")
		j1 := Jump_If_Not_Zero_Instruction_Tacky{condition: v1, target: true_label}
		instructions = append(instructions, &j1)
		v2, instructions := exp.secExp.expToTacky(instructions)
		j2 := Jump_If_Not_Zero_Instruction_Tacky{condition: v2, target: true_label}
		instructions = append(instructions, &j2)
		result := Variable_Value_Tacky{makeTempVarName("")}
		cp1 := Copy_Instruction_Tacky{src: &Constant_Value_Tacky{0}, dst: &result}
		instructions = append(instructions, &cp1)
		end := makeLabelName("end")
		j3 := Jump_Instruction_Tacky{end}
		instructions = append(instructions, &j3)
		lb1 := Label_Instruction_Tacky{true_label}
		instructions = append(instructions, &lb1)
		cp2 := Copy_Instruction_Tacky{src: &Constant_Value_Tacky{1}, dst: &result}
		instructions = append(instructions, &cp2)
		lb2 := Label_Instruction_Tacky{end}
		instructions = append(instructions, &lb2)
		return &result, instructions
	} else {
		src1, instructions := exp.firstExp.expToTacky(instructions)
		src2, instructions := exp.secExp.expToTacky(instructions)
		dstName := makeTempVarName("")
		dst := Variable_Value_Tacky{dstName}
		instr := Binary_Instruction_Tacky{binOp: exp.binOp, src1: src1, src2: src2, dst: &dst}
		instructions = append(instructions, &instr)
		return &dst, instructions
	}
}

/////////////////////////////////////////////////////////////////////////////////

func (exp *Assignment_Expression) expToTacky(instructions []Instruction_Tacky) (Value_Tacky, []Instruction_Tacky) {
	// TODO:
	return nil, []Instruction_Tacky{}
}
