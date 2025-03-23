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
	topItems []Top_Level_Tacky
}

//###############################################################################
//###############################################################################
//###############################################################################

type Top_Level_Tacky interface {
	topLevelToAsm()
}

/////////////////////////////////////////////////////////////////////////////////

type Function_Definition_Tacky struct {
	name   string
	global bool
	params []string
	body   []Instruction_Tacky
}

/////////////////////////////////////////////////////////////////////////////////

type Static_Variable_Tacky struct {
	name         string
	global       bool
	initialValue int32
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

/////////////////////////////////////////////////////////////////////////////////

type Function_Call_Tacky struct {
	funcName  string
	args      []Value_Tacky
	returnVal Value_Tacky
}

//###############################################################################
//###############################################################################
//###############################################################################

type Value_Tacky interface {
	valueToAsm() Operand_Asm
}

/////////////////////////////////////////////////////////////////////////////////

// TODO: could switch this to using enums for the data type and a string to hold the actual value
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
	topItems := []Top_Level_Tacky{}

	for _, decl := range pr.decls {
		fnDecl, isFunc := decl.(*Function_Declaration)
		if !isFunc {
			// we won't deal with file scope variable declarations now, we'll do that later
			continue
		}

		instrs := fnDecl.declToTacky()
		if len(instrs) > 0 {
			// function definitions will have at least one instruction, function declarations won't have any instructions,
			// we will only keep the function definitions
			global := symbolTable[fnDecl.name].attrs.isGlobalAttribute()
			tacFunc := Function_Definition_Tacky{name: fnDecl.name, global: global, params: fnDecl.params, body: instrs}
			topItems = append(topItems, &tacFunc)
		}
	}

	// add the top level items that are in the symbol table
	moreItems := convertSymbolsToTacky()
	topItems = append(topItems, moreItems...)

	tacky := Program_Tacky{topItems: topItems}
	return tacky
}

//###############################################################################
//###############################################################################
//###############################################################################

func convertSymbolsToTacky() []Top_Level_Tacky {
	topItems := []Top_Level_Tacky{}

	for name, sym := range symbolTable {
		switch convertedAttr := sym.attrs.(type) {
		case *Static_Attributes:
			switch convertedInit := convertedAttr.init.(type) {
			case *Initial_Int:
				v := Static_Variable_Tacky{name: name, global: convertedAttr.global, initialValue: convertedInit.value}
				topItems = append(topItems, &v)
			case *Tentative:
				v := Static_Variable_Tacky{name: name, global: convertedAttr.global, initialValue: 0}
				topItems = append(topItems, &v)
			default:
				continue
			}
		default:
			continue
		}
	}

	return topItems
}

/////////////////////////////////////////////////////////////////////////////////

func (d *Variable_Declaration) declToTacky() []Instruction_Tacky {
	if d.initializer == nil {
		// no instructions needed
		return []Instruction_Tacky{}
	} else if (d.storageClass == STATIC_STORAGE_CLASS) || (d.storageClass == EXTERN_STORAGE_CLASS) {
		// don't emit tacky for local variable declarations with static or extern specifiers,
		// we handle that at the top level
		return []Instruction_Tacky{}
	} else {
		// get the instructions for the initializer
		instructions := []Instruction_Tacky{}
		result, instructions := d.initializer.expToTacky(instructions)

		// assign the value from the initializer to the declared variable
		v := Variable_Value_Tacky{d.name}
		cp := Copy_Instruction_Tacky{result, &v}
		instructions = append(instructions, &cp)
		return instructions
	}
}

/////////////////////////////////////////////////////////////////////////////////

func (fn *Function_Declaration) declToTacky() []Instruction_Tacky {
	if fn.body == nil {
		// no instructions needed
		return []Instruction_Tacky{}
	}

	bodyTac := fn.body.blockToTacky()

	// Add a return statement to the end of every function just in case the original source didn't have one.
	// If it already had a return statement then no big deal becuase this new ret instruction will never run.
	ret := Return_Instruction_Tacky{&Constant_Value_Tacky{0}}
	bodyTac = append(bodyTac, &ret)

	return bodyTac
}

//###############################################################################
//###############################################################################
//###############################################################################

func (b *Block) blockToTacky() []Instruction_Tacky {
	instructions := []Instruction_Tacky{}

	for _, bItem := range b.items {
		moreInstr := bItem.blockItemToTacky()
		instructions = append(instructions, moreInstr...)
	}

	return instructions
}

//###############################################################################
//###############################################################################
//###############################################################################

func (bi *Block_Statement) blockItemToTacky() []Instruction_Tacky {
	return bi.st.statementToTacky()
}

/////////////////////////////////////////////////////////////////////////////////

func (bi *Block_Declaration) blockItemToTacky() []Instruction_Tacky {
	return bi.decl.declToTacky()
}

//###############################################################################
//###############################################################################
//###############################################################################

func (fid *For_Initial_Declaration) forInitialToTacky() []Instruction_Tacky {
	return fid.decl.declToTacky()
}

/////////////////////////////////////////////////////////////////////////////////

func (fie *For_Initial_Expression) forInitialToTacky() []Instruction_Tacky {
	if fie.exp == nil {
		return []Instruction_Tacky{}
	} else {
		_, instructions := fie.exp.expToTacky([]Instruction_Tacky{})
		return instructions
	}
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
	instructions := []Instruction_Tacky{}
	_, instructions = st.exp.expToTacky(instructions)
	return instructions
}

/////////////////////////////////////////////////////////////////////////////////

func (st *If_Statement) statementToTacky() []Instruction_Tacky {
	if st.elseSt == nil {
		c, instructions := st.condition.expToTacky([]Instruction_Tacky{})
		endLabel := makeLabelName("end")
		jmp := Jump_If_Zero_Instruction_Tacky{condition: c, target: endLabel}
		instructions = append(instructions, &jmp)
		moreInstr := st.thenSt.statementToTacky()
		instructions = append(instructions, moreInstr...)
		lblInstr := Label_Instruction_Tacky{endLabel}
		instructions = append(instructions, &lblInstr)
		return instructions
	} else {
		c, instructions := st.condition.expToTacky([]Instruction_Tacky{})
		elseLabel := makeLabelName("else")
		jmpElse := Jump_If_Zero_Instruction_Tacky{condition: c, target: elseLabel}
		instructions = append(instructions, &jmpElse)
		moreInstr := st.thenSt.statementToTacky()
		instructions = append(instructions, moreInstr...)
		endLabel := makeLabelName("end")
		jmpEnd := Jump_Instruction_Tacky{endLabel}
		instructions = append(instructions, &jmpEnd)
		elseLabelInstr := Label_Instruction_Tacky{elseLabel}
		instructions = append(instructions, &elseLabelInstr)
		moreInstr = st.elseSt.statementToTacky()
		instructions = append(instructions, moreInstr...)
		endLabelInstr := Label_Instruction_Tacky{endLabel}
		instructions = append(instructions, &endLabelInstr)
		return instructions
	}
}

/////////////////////////////////////////////////////////////////////////////////

func (st *Compound_Statement) statementToTacky() []Instruction_Tacky {
	return st.block.blockToTacky()
}

/////////////////////////////////////////////////////////////////////////////////

func (st *Break_Statement) statementToTacky() []Instruction_Tacky {
	jmp := Jump_Instruction_Tacky{"break_" + st.label}
	return []Instruction_Tacky{&jmp}
}

/////////////////////////////////////////////////////////////////////////////////

func (st *Continue_Statement) statementToTacky() []Instruction_Tacky {
	jmp := Jump_Instruction_Tacky{"continue_" + st.label}
	return []Instruction_Tacky{&jmp}
}

/////////////////////////////////////////////////////////////////////////////////

func (st *While_Statement) statementToTacky() []Instruction_Tacky {
	instructions := []Instruction_Tacky{}

	continueLabel := Label_Instruction_Tacky{"continue_" + st.label}
	instructions = append(instructions, &continueLabel)

	v, instructions := st.condition.expToTacky(instructions)

	jmpBreak := Jump_If_Zero_Instruction_Tacky{condition: v, target: "break_" + st.label}
	instructions = append(instructions, &jmpBreak)

	moreInstr := st.body.statementToTacky()
	instructions = append(instructions, moreInstr...)

	jmpContinue := Jump_Instruction_Tacky{"continue_" + st.label}
	instructions = append(instructions, &jmpContinue)

	breakLabel := Label_Instruction_Tacky{"break_" + st.label}
	instructions = append(instructions, &breakLabel)

	return instructions
}

/////////////////////////////////////////////////////////////////////////////////

func (st *Do_While_Statement) statementToTacky() []Instruction_Tacky {
	instructions := []Instruction_Tacky{}

	startLabel := Label_Instruction_Tacky{"start_" + st.label}
	instructions = append(instructions, &startLabel)

	moreInstr := st.body.statementToTacky()
	instructions = append(instructions, moreInstr...)

	continueLabel := Label_Instruction_Tacky{"continue_" + st.label}
	instructions = append(instructions, &continueLabel)

	v, instructions := st.condition.expToTacky(instructions)

	jmp := Jump_If_Not_Zero_Instruction_Tacky{condition: v, target: "start_" + st.label}
	instructions = append(instructions, &jmp)

	breakLabel := Label_Instruction_Tacky{"break_" + st.label}
	instructions = append(instructions, &breakLabel)

	return instructions
}

/////////////////////////////////////////////////////////////////////////////////

func (st *For_Statement) statementToTacky() []Instruction_Tacky {
	instructions := []Instruction_Tacky{}

	moreInstr := st.initial.forInitialToTacky()
	instructions = append(instructions, moreInstr...)

	startLabel := Label_Instruction_Tacky{"start_" + st.label}
	instructions = append(instructions, &startLabel)

	if st.condition != nil {
		var v Value_Tacky
		v, instructions = st.condition.expToTacky(instructions)
		jmpBreak := Jump_If_Zero_Instruction_Tacky{condition: v, target: "break_" + st.label}
		instructions = append(instructions, &jmpBreak)
	}

	moreInstr = st.body.statementToTacky()
	instructions = append(instructions, moreInstr...)

	continueLabel := Label_Instruction_Tacky{"continue_" + st.label}
	instructions = append(instructions, &continueLabel)

	if st.post != nil {
		_, instructions = st.post.expToTacky(instructions)
	}

	jmp := Jump_Instruction_Tacky{target: "start_" + st.label}
	instructions = append(instructions, &jmp)

	breakLabel := Label_Instruction_Tacky{"break_" + st.label}
	instructions = append(instructions, &breakLabel)

	return instructions
}

/////////////////////////////////////////////////////////////////////////////////

func (st *Null_Statement) statementToTacky() []Instruction_Tacky {
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
	return &Variable_Value_Tacky{exp.name}, instructions
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
	// evaluate the right side of the expression
	result, instructions := exp.rightExp.expToTacky(instructions)

	varExp, _ := exp.lvalue.(*Variable_Expression)
	v := Variable_Value_Tacky{varExp.name}

	// store the right side of the expression in the lvalue
	cp := Copy_Instruction_Tacky{result, &v}
	instructions = append(instructions, &cp)

	return &v, instructions
}

/////////////////////////////////////////////////////////////////////////////////

func (exp *Conditional_Expression) expToTacky(instructions []Instruction_Tacky) (Value_Tacky, []Instruction_Tacky) {
	c, instructions := exp.condition.expToTacky(instructions)
	rightLabel := makeLabelName("rightExp")
	jmp := Jump_If_Zero_Instruction_Tacky{c, rightLabel}
	instructions = append(instructions, &jmp)
	v1, instructions := exp.middleExp.expToTacky(instructions)
	result := Variable_Value_Tacky{makeTempVarName("")}
	cp1 := Copy_Instruction_Tacky{v1, &result}
	instructions = append(instructions, &cp1)
	endLabel := makeLabelName("end")
	jmpEnd := Jump_Instruction_Tacky{endLabel}
	instructions = append(instructions, &jmpEnd)
	rightLabelInstr := Label_Instruction_Tacky{rightLabel}
	instructions = append(instructions, &rightLabelInstr)
	v2, instructions := exp.rightExp.expToTacky(instructions)
	cp2 := Copy_Instruction_Tacky{v2, &result}
	instructions = append(instructions, &cp2)
	endLabelInstr := Label_Instruction_Tacky{endLabel}
	instructions = append(instructions, &endLabelInstr)
	return &result, instructions
}

/////////////////////////////////////////////////////////////////////////////////

func (e *Function_Call_Expression) expToTacky(instructions []Instruction_Tacky) (Value_Tacky, []Instruction_Tacky) {
	argsTacky := []Value_Tacky{}

	for _, argExp := range e.args {
		var argTac Value_Tacky
		argTac, instructions = argExp.expToTacky(instructions)
		argsTacky = append(argsTacky, argTac)
	}

	retVal := Variable_Value_Tacky{makeTempVarName("")}
	fn := Function_Call_Tacky{funcName: e.functionName, args: argsTacky, returnVal: &retVal}
	instructions = append(instructions, &fn)

	return &retVal, instructions
}
