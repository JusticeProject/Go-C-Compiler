package main

import "strconv"

//###############################################################################
//###############################################################################
//###############################################################################

var tempVarCounter int64 = -1

func makeTempVarName(prefix string) string {
	tempVarCounter++
	if len(prefix) > 0 {
		return prefix + "." + strconv.FormatInt(tempVarCounter, 10)
	} else {
		return "tmp." + strconv.FormatInt(tempVarCounter, 10)
	}
}

/////////////////////////////////////////////////////////////////////////////////

func makeLabelName(name string) string {
	tempVarCounter++
	return name + strconv.FormatInt(tempVarCounter, 10)
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
	topLevelToAsm() Top_Level_Asm
}

/////////////////////////////////////////////////////////////////////////////////

type Function_Definition_Tacky struct {
	name       string
	global     bool
	paramNames []string
	body       []Instruction_Tacky
}

/////////////////////////////////////////////////////////////////////////////////

type Static_Variable_Tacky struct {
	name         string
	global       bool
	initialValue string
	initEnum     InitializerEnum
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

type Sign_Extend_Instruction_Tacky struct {
	src Value_Tacky
	dst Value_Tacky
}

/////////////////////////////////////////////////////////////////////////////////

type Truncate_Instruction_Tacky struct {
	src Value_Tacky
	dst Value_Tacky
}

/////////////////////////////////////////////////////////////////////////////////

type Zero_Extend_Instruction_Tacky struct {
	src Value_Tacky
	dst Value_Tacky
}

/////////////////////////////////////////////////////////////////////////////////

type Double_To_Int_Instruction_Tacky struct {
	src Value_Tacky
	dst Value_Tacky
}

/////////////////////////////////////////////////////////////////////////////////

type Double_To_UInt_Instruction_Tacky struct {
	src Value_Tacky
	dst Value_Tacky
}

/////////////////////////////////////////////////////////////////////////////////

type Int_To_Double_Instruction_Tacky struct {
	src Value_Tacky
	dst Value_Tacky
}

/////////////////////////////////////////////////////////////////////////////////

type UInt_To_Double_Instruction_Tacky struct {
	src Value_Tacky
	dst Value_Tacky
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

// get address of variable src and store it in dst
type Get_Address_Instruction_Tacky struct {
	src Value_Tacky
	dst Value_Tacky
}

/////////////////////////////////////////////////////////////////////////////////

type Load_Instruction_Tacky struct {
	srcPtr Value_Tacky
	dst    Value_Tacky
}

/////////////////////////////////////////////////////////////////////////////////

type Store_Instruction_Tacky struct {
	src    Value_Tacky
	dstPtr Value_Tacky
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

type Expression_Result_Tacky interface {
	isExpResult()
}

/////////////////////////////////////////////////////////////////////////////////

type Plain_Operand_Tacky struct {
	val Value_Tacky
}

func (er *Plain_Operand_Tacky) isExpResult() {}

/////////////////////////////////////////////////////////////////////////////////

type Dereferenced_Pointer_Tacky struct {
	ptr Value_Tacky
}

func (er *Dereferenced_Pointer_Tacky) isExpResult() {}

//###############################################################################
//###############################################################################
//###############################################################################

type Value_Tacky interface {
	valueToAsm() Operand_Asm
	getDataType() DataTypeEnum
	getAssemblyType() AssemblyTypeEnum
	isSigned() bool
}

/////////////////////////////////////////////////////////////////////////////////

type Constant_Value_Tacky struct {
	typ   DataTypeEnum
	value string
}

/////////////////////////////////////////////////////////////////////////////////

type Variable_Value_Tacky struct {
	name string
}

func makeTackyVariable(typ DataTypeEnum) Variable_Value_Tacky {
	varName := makeTempVarName("")
	symbolTable[varName] = Symbol{dataTyp: Data_Type{typ: typ}, attrs: LOCAL_ATTRIBUTES}
	return Variable_Value_Tacky{varName}
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
			global := symbolTable[fnDecl.name].global
			tacFunc := Function_Definition_Tacky{name: fnDecl.name, global: global, paramNames: fnDecl.paramNames, body: instrs}
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
		switch sym.attrs {
		case STATIC_ATTRIBUTES:
			switch sym.initEnum {
			case NO_INITIALIZER:
				continue
			case TENTATIVE_INIT:
				v := Static_Variable_Tacky{name: name, global: sym.global, initialValue: "0",
					initEnum: dataTypeEnumToInitEnum(sym.dataTyp.typ)}
				topItems = append(topItems, &v)
			default:
				// it has an initializer with an int, long, float, etc.
				v := Static_Variable_Tacky{name: name, global: sym.global, initialValue: sym.initialValue, initEnum: sym.initEnum}
				topItems = append(topItems, &v)
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
		result, instructions := expToTackyAndConvert(d.initializer, instructions)

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
	ret := Return_Instruction_Tacky{&Constant_Value_Tacky{typ: INT_TYPE, value: "0"}}
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
		_, instructions := expToTackyAndConvert(fie.exp, []Instruction_Tacky{})
		return instructions
	}
}

//###############################################################################
//###############################################################################
//###############################################################################

func (st *Return_Statement) statementToTacky() []Instruction_Tacky {
	instructions := []Instruction_Tacky{}
	val, instructions := expToTackyAndConvert(st.exp, instructions)
	instr := Return_Instruction_Tacky{val: val}
	instructions = append(instructions, &instr)
	return instructions
}

/////////////////////////////////////////////////////////////////////////////////

func (st *Expression_Statement) statementToTacky() []Instruction_Tacky {
	instructions := []Instruction_Tacky{}
	_, instructions = expToTackyAndConvert(st.exp, instructions)
	return instructions
}

/////////////////////////////////////////////////////////////////////////////////

func (st *If_Statement) statementToTacky() []Instruction_Tacky {
	if st.elseSt == nil {
		c, instructions := expToTackyAndConvert(st.condition, []Instruction_Tacky{})
		endLabel := makeLabelName("end")
		jmp := Jump_If_Zero_Instruction_Tacky{condition: c, target: endLabel}
		instructions = append(instructions, &jmp)
		moreInstr := st.thenSt.statementToTacky()
		instructions = append(instructions, moreInstr...)
		lblInstr := Label_Instruction_Tacky{endLabel}
		instructions = append(instructions, &lblInstr)
		return instructions
	} else {
		c, instructions := expToTackyAndConvert(st.condition, []Instruction_Tacky{})
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

	v, instructions := expToTackyAndConvert(st.condition, instructions)

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

	v, instructions := expToTackyAndConvert(st.condition, instructions)

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
		v, instructions = expToTackyAndConvert(st.condition, instructions)
		jmpBreak := Jump_If_Zero_Instruction_Tacky{condition: v, target: "break_" + st.label}
		instructions = append(instructions, &jmpBreak)
	}

	moreInstr = st.body.statementToTacky()
	instructions = append(instructions, moreInstr...)

	continueLabel := Label_Instruction_Tacky{"continue_" + st.label}
	instructions = append(instructions, &continueLabel)

	if st.post != nil {
		_, instructions = expToTackyAndConvert(st.post, instructions)
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

func expToTackyAndConvert(exp Expression, instructions []Instruction_Tacky) (Value_Tacky, []Instruction_Tacky) {
	result, instructions := exp.expToTacky(instructions)

	switch convertedRes := result.(type) {
	case *Plain_Operand_Tacky:
		return convertedRes.val, instructions
	case *Dereferenced_Pointer_Tacky:
		dst := makeTackyVariable(getResultType(exp).typ)
		load := Load_Instruction_Tacky{srcPtr: convertedRes.ptr, dst: &dst}
		instructions = append(instructions, &load)
		return &dst, instructions
	}
	return nil, []Instruction_Tacky{}
}

/////////////////////////////////////////////////////////////////////////////////

func (exp *Constant_Value_Expression) expToTacky(instructions []Instruction_Tacky) (Expression_Result_Tacky, []Instruction_Tacky) {
	val := Constant_Value_Tacky{typ: exp.dTyp.typ, value: exp.value}
	return &Plain_Operand_Tacky{&val}, instructions
}

/////////////////////////////////////////////////////////////////////////////////

func (exp *Variable_Expression) expToTacky(instructions []Instruction_Tacky) (Expression_Result_Tacky, []Instruction_Tacky) {
	v := Variable_Value_Tacky{exp.name}
	return &Plain_Operand_Tacky{&v}, instructions
}

/////////////////////////////////////////////////////////////////////////////////

func (exp *Cast_Expression) expToTacky(instructions []Instruction_Tacky) (Expression_Result_Tacky, []Instruction_Tacky) {
	innerResult, instructions := expToTackyAndConvert(exp.innerExp, instructions)
	innerType := getResultType(exp.innerExp)

	// if they are both the same type then nothing more to do
	if exp.targetType.isEqualType(&innerType) {
		return &Plain_Operand_Tacky{innerResult}, instructions
	}

	dst := makeTackyVariable(exp.targetType.typ)
	// TODO: update as we add more data types
	if exp.targetType.typ == DOUBLE_TYPE {
		if (innerType.typ == INT_TYPE) || (innerType.typ == LONG_TYPE) {
			newInstr := Int_To_Double_Instruction_Tacky{src: innerResult, dst: &dst}
			instructions = append(instructions, &newInstr)
		} else if (innerType.typ == UNSIGNED_INT_TYPE) || (innerType.typ == UNSIGNED_LONG_TYPE) {
			newInstr := UInt_To_Double_Instruction_Tacky{src: innerResult, dst: &dst}
			instructions = append(instructions, &newInstr)
		} else {
			fail("Cast not supported")
		}
	} else if innerType.typ == DOUBLE_TYPE {
		if (exp.targetType.typ == INT_TYPE) || (exp.targetType.typ == LONG_TYPE) {
			newInstr := Double_To_Int_Instruction_Tacky{src: innerResult, dst: &dst}
			instructions = append(instructions, &newInstr)
		} else if (exp.targetType.typ == UNSIGNED_INT_TYPE) || (exp.targetType.typ == UNSIGNED_LONG_TYPE) {
			newInstr := Double_To_UInt_Instruction_Tacky{src: innerResult, dst: &dst}
			instructions = append(instructions, &newInstr)
		} else {
			fail("Cast not supported")
		}
	} else if size(exp.targetType.typ) == size(innerType.typ) {
		newInstr := Copy_Instruction_Tacky{src: innerResult, dst: &dst}
		instructions = append(instructions, &newInstr)
	} else if size(exp.targetType.typ) < size(innerType.typ) {
		newInstr := Truncate_Instruction_Tacky{src: innerResult, dst: &dst}
		instructions = append(instructions, &newInstr)
	} else if isSigned(innerType.typ) {
		// the target type is bigger, do a sign extend since the inner type is signed
		newInstr := Sign_Extend_Instruction_Tacky{src: innerResult, dst: &dst}
		instructions = append(instructions, &newInstr)
	} else {
		// the target type is bigger, do a zero extend since the inner type is unsigned
		newInstr := Zero_Extend_Instruction_Tacky{src: innerResult, dst: &dst}
		instructions = append(instructions, &newInstr)
	}

	return &Plain_Operand_Tacky{&dst}, instructions
}

/////////////////////////////////////////////////////////////////////////////////

func (exp *Unary_Expression) expToTacky(instructions []Instruction_Tacky) (Expression_Result_Tacky, []Instruction_Tacky) {
	src, instructions := expToTackyAndConvert(exp.innerExp, instructions)
	dst := makeTackyVariable(getResultType(exp).typ)
	instr := Unary_Instruction_Tacky{unOp: exp.unOp, src: src, dst: &dst}
	instructions = append(instructions, &instr)
	return &Plain_Operand_Tacky{&dst}, instructions
}

/////////////////////////////////////////////////////////////////////////////////

func (exp *Binary_Expression) expToTacky(instructions []Instruction_Tacky) (Expression_Result_Tacky, []Instruction_Tacky) {
	// some operators can short-circuit on the first expression, so we handle them differently
	if exp.binOp == AND_OPERATOR {
		v1, instructions := expToTackyAndConvert(exp.firstExp, instructions)
		false_label := makeLabelName("and_false")
		j1 := Jump_If_Zero_Instruction_Tacky{condition: v1, target: false_label}
		instructions = append(instructions, &j1)
		v2, instructions := expToTackyAndConvert(exp.secExp, instructions)
		j2 := Jump_If_Zero_Instruction_Tacky{condition: v2, target: false_label}
		instructions = append(instructions, &j2)
		result := makeTackyVariable(getResultType(exp).typ)
		cp1 := Copy_Instruction_Tacky{src: &Constant_Value_Tacky{typ: INT_TYPE, value: "1"}, dst: &result}
		instructions = append(instructions, &cp1)
		end := makeLabelName("end")
		j3 := Jump_Instruction_Tacky{end}
		instructions = append(instructions, &j3)
		lb1 := Label_Instruction_Tacky{false_label}
		instructions = append(instructions, &lb1)
		cp2 := Copy_Instruction_Tacky{src: &Constant_Value_Tacky{typ: INT_TYPE, value: "0"}, dst: &result}
		instructions = append(instructions, &cp2)
		lb2 := Label_Instruction_Tacky{end}
		instructions = append(instructions, &lb2)
		return &Plain_Operand_Tacky{&result}, instructions
	} else if exp.binOp == OR_OPERATOR {
		v1, instructions := expToTackyAndConvert(exp.firstExp, instructions)
		true_label := makeLabelName("or_true")
		j1 := Jump_If_Not_Zero_Instruction_Tacky{condition: v1, target: true_label}
		instructions = append(instructions, &j1)
		v2, instructions := expToTackyAndConvert(exp.secExp, instructions)
		j2 := Jump_If_Not_Zero_Instruction_Tacky{condition: v2, target: true_label}
		instructions = append(instructions, &j2)
		result := makeTackyVariable(getResultType(exp).typ)
		cp1 := Copy_Instruction_Tacky{src: &Constant_Value_Tacky{typ: INT_TYPE, value: "0"}, dst: &result}
		instructions = append(instructions, &cp1)
		end := makeLabelName("end")
		j3 := Jump_Instruction_Tacky{end}
		instructions = append(instructions, &j3)
		lb1 := Label_Instruction_Tacky{true_label}
		instructions = append(instructions, &lb1)
		cp2 := Copy_Instruction_Tacky{src: &Constant_Value_Tacky{typ: INT_TYPE, value: "1"}, dst: &result}
		instructions = append(instructions, &cp2)
		lb2 := Label_Instruction_Tacky{end}
		instructions = append(instructions, &lb2)
		return &Plain_Operand_Tacky{&result}, instructions
	} else {
		src1, instructions := expToTackyAndConvert(exp.firstExp, instructions)
		src2, instructions := expToTackyAndConvert(exp.secExp, instructions)
		dst := makeTackyVariable(getResultType(exp).typ)
		instr := Binary_Instruction_Tacky{binOp: exp.binOp, src1: src1, src2: src2, dst: &dst}
		instructions = append(instructions, &instr)
		return &Plain_Operand_Tacky{&dst}, instructions
	}
}

/////////////////////////////////////////////////////////////////////////////////

func (exp *Assignment_Expression) expToTacky(instructions []Instruction_Tacky) (Expression_Result_Tacky, []Instruction_Tacky) {
	lval, instructions := exp.lvalue.expToTacky(instructions)
	rval, instructions := expToTackyAndConvert(exp.rightExp, instructions)

	switch convertedLval := lval.(type) {
	case *Plain_Operand_Tacky:
		cp := Copy_Instruction_Tacky{src: rval, dst: convertedLval.val}
		instructions = append(instructions, &cp)
		return lval, instructions
	case *Dereferenced_Pointer_Tacky:
		store := Store_Instruction_Tacky{src: rval, dstPtr: convertedLval.ptr}
		instructions = append(instructions, &store)
		return &Plain_Operand_Tacky{rval}, instructions
	}

	return nil, []Instruction_Tacky{}
}

/////////////////////////////////////////////////////////////////////////////////

func (exp *Conditional_Expression) expToTacky(instructions []Instruction_Tacky) (Expression_Result_Tacky, []Instruction_Tacky) {
	c, instructions := expToTackyAndConvert(exp.condition, instructions)
	rightLabel := makeLabelName("rightExp")
	jmp := Jump_If_Zero_Instruction_Tacky{c, rightLabel}
	instructions = append(instructions, &jmp)
	v1, instructions := expToTackyAndConvert(exp.middleExp, instructions)
	result := makeTackyVariable(getResultType(exp).typ)
	cp1 := Copy_Instruction_Tacky{v1, &result}
	instructions = append(instructions, &cp1)
	endLabel := makeLabelName("end")
	jmpEnd := Jump_Instruction_Tacky{endLabel}
	instructions = append(instructions, &jmpEnd)
	rightLabelInstr := Label_Instruction_Tacky{rightLabel}
	instructions = append(instructions, &rightLabelInstr)
	v2, instructions := expToTackyAndConvert(exp.rightExp, instructions)
	cp2 := Copy_Instruction_Tacky{v2, &result}
	instructions = append(instructions, &cp2)
	endLabelInstr := Label_Instruction_Tacky{endLabel}
	instructions = append(instructions, &endLabelInstr)
	return &Plain_Operand_Tacky{&result}, instructions
}

/////////////////////////////////////////////////////////////////////////////////

func (e *Function_Call_Expression) expToTacky(instructions []Instruction_Tacky) (Expression_Result_Tacky, []Instruction_Tacky) {
	argsTacky := []Value_Tacky{}

	for _, argExp := range e.args {
		var argTac Value_Tacky
		argTac, instructions = expToTackyAndConvert(argExp, instructions)
		argsTacky = append(argsTacky, argTac)
	}

	retVal := makeTackyVariable(getResultType(e).typ)
	fn := Function_Call_Tacky{funcName: e.functionName, args: argsTacky, returnVal: &retVal}
	instructions = append(instructions, &fn)

	return &Plain_Operand_Tacky{&retVal}, instructions
}

/////////////////////////////////////////////////////////////////////////////////

func (e *Dereference_Expression) expToTacky(instructions []Instruction_Tacky) (Expression_Result_Tacky, []Instruction_Tacky) {
	result, instructions := expToTackyAndConvert(e.innerExp, instructions)
	return &Dereferenced_Pointer_Tacky{result}, instructions
}

/////////////////////////////////////////////////////////////////////////////////

func (e *Address_Of_Expression) expToTacky(instructions []Instruction_Tacky) (Expression_Result_Tacky, []Instruction_Tacky) {
	v, instructions := e.innerExp.expToTacky(instructions)

	switch convertedV := v.(type) {
	case *Plain_Operand_Tacky:
		dst := makeTackyVariable(getResultType(e).typ)
		getAddr := Get_Address_Instruction_Tacky{src: convertedV.val, dst: &dst}
		instructions = append(instructions, &getAddr)
		return &Plain_Operand_Tacky{&dst}, instructions
	case *Dereferenced_Pointer_Tacky:
		return &Plain_Operand_Tacky{convertedV.ptr}, instructions
	}

	return nil, []Instruction_Tacky{}
}
