package main

import (
	"math"
	"os"
	"strconv"
)

/////////////////////////////////////////////////////////////////////////////////

// this is needed because gccgo currently doesn't support go 1.24
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

//###############################################################################
//###############################################################################
//###############################################################################

type Program_Asm struct {
	topItems []Top_Level_Asm
}

/////////////////////////////////////////////////////////////////////////////////

type AssemblyTypeEnum int

const (
	NONE_ASM_TYPE AssemblyTypeEnum = iota
	BYTE_ASM_TYPE
	LONGWORD_ASM_TYPE
	QUADWORD_ASM_TYPE
)

func dataTypeEnumToAssemblyTypeEnum(input DataTypeEnum) AssemblyTypeEnum {
	switch input {
	case INT_TYPE:
		return LONGWORD_ASM_TYPE
	case LONG_TYPE:
		return QUADWORD_ASM_TYPE
	case FUNCTION_TYPE:
		return NONE_ASM_TYPE
	}
	fail("Can not convert DataTypeEnum to AssemblyTypeEnum")
	return NONE_ASM_TYPE
}

/////////////////////////////////////////////////////////////////////////////////

func asmTypToAlignment(asmTyp AssemblyTypeEnum) int32 {
	switch asmTyp {
	case LONGWORD_ASM_TYPE:
		return 4
	case QUADWORD_ASM_TYPE:
		return 8
	}
	fail("Can not convert AssemblyTypeEnum to alignment")
	return 0
}

/////////////////////////////////////////////////////////////////////////////////

func initToAlignment(initEnum InitializerEnum) int32 {
	switch initEnum {
	case INITIAL_INT:
		return 4
	case INITIAL_LONG:
		return 8
	}
	fail("Can not convert InitializerEnum to alignment")
	return 0
}

/////////////////////////////////////////////////////////////////////////////////

func getAsmTypeOfVariable(name string) AssemblyTypeEnum {
	typ := symbolTable[name].dataTyp.typ
	return dataTypeEnumToAssemblyTypeEnum(typ)
}

//###############################################################################
//###############################################################################
//###############################################################################

type Symbol_Asm struct {
	// for variables
	asmTyp   AssemblyTypeEnum
	isStatic bool

	// for functions
	defined bool
}

var symbolTableBackend = make(map[string]Symbol_Asm)

//###############################################################################
//###############################################################################
//###############################################################################

type Top_Level_Asm interface {
	replacePseudoregisters(nameToOffset map[string]int32)
	fixInvalidInstr()
	topLevelEmitAsm(file *os.File)
}

/////////////////////////////////////////////////////////////////////////////////

type Function_Asm struct {
	name         string
	global       bool
	instructions []Instruction_Asm
	stackSize    int32
}

/////////////////////////////////////////////////////////////////////////////////

type Static_Variable_Asm struct {
	name         string
	global       bool
	alignment    int32
	initialValue string
	initEnum     InitializerEnum
}

//###############################################################################
//###############################################################################
//###############################################################################

type Instruction_Asm interface {
	instrEmitAsm(file *os.File)
}

/////////////////////////////////////////////////////////////////////////////////

type Mov_Instruction_Asm struct {
	asmTyp AssemblyTypeEnum
	src    Operand_Asm
	dst    Operand_Asm
}

/////////////////////////////////////////////////////////////////////////////////

// Mov with sign extend
type Movsx_Instruction_Asm struct {
	src Operand_Asm
	dst Operand_Asm
}

/////////////////////////////////////////////////////////////////////////////////

type Unary_Instruction_Asm struct {
	unOp   UnaryOperatorTypeAsm
	asmTyp AssemblyTypeEnum
	src    Operand_Asm
}

/////////////////////////////////////////////////////////////////////////////////

type Binary_Instruction_Asm struct {
	binOp  BinaryOperatorTypeAsm
	asmTyp AssemblyTypeEnum
	src    Operand_Asm
	dst    Operand_Asm
}

/////////////////////////////////////////////////////////////////////////////////

type Compare_Instruction_Asm struct {
	asmTyp AssemblyTypeEnum
	op1    Operand_Asm
	op2    Operand_Asm
}

/////////////////////////////////////////////////////////////////////////////////

type IDivide_Instruction_Asm struct {
	asmTyp  AssemblyTypeEnum
	divisor Operand_Asm
}

/////////////////////////////////////////////////////////////////////////////////

type CDQ_Sign_Extend_Instruction_Asm struct {
	asmTyp AssemblyTypeEnum
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

type Push_Instruction_Asm struct {
	op Operand_Asm
}

/////////////////////////////////////////////////////////////////////////////////

type Call_Function_Asm struct {
	name string
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
		fail("NOT_OPERATOR not converted directly to Asm")
	default:
		fail("unknown UnaryOperatorType:")
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
		fail("unknown BinaryOperatorType:")
	}
	return NOP_BINARY_ASM
}

//###############################################################################
//###############################################################################
//###############################################################################

type Operand_Asm interface {
	getOperandString(asmTyp AssemblyTypeEnum) string
}

func opIsBigImm(op Operand_Asm) bool {
	imm, isImm := op.(*Immediate_Int_Operand_Asm)
	if !isImm {
		return false
	}
	integer, _ := strconv.ParseInt(imm.value, 10, 64)
	if integer > math.MaxInt32 {
		return true
	} else {
		return false
	}
}

/////////////////////////////////////////////////////////////////////////////////

type Immediate_Int_Operand_Asm struct {
	value string
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

/////////////////////////////////////////////////////////////////////////////////

// used for static variables
type Data_Operand_Asm struct {
	name string
}

//###############################################################################
//###############################################################################
//###############################################################################

type ConditionalCodeAsm int

const (
	NONE_CODE_ASM ConditionalCodeAsm = iota
	IS_EQUAL_CODE_ASM
	NOT_EQUAL_CODE_ASM
	LESS_THAN_CODE_ASM
	LESS_OR_EQUAL_CODE_ASM
	GREATER_THAN_CODE_ASM
	GREATER_OR_EQUAL_CODE_ASM
)

func convertBinaryOpToCondition(binOp BinaryOperatorType) ConditionalCodeAsm {
	switch binOp {
	case IS_EQUAL_OPERATOR:
		return IS_EQUAL_CODE_ASM
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
		fail("unknown BinaryOperatorType when converting to code")
	}

	return NONE_CODE_ASM
}

/////////////////////////////////////////////////////////////////////////////////

type RegisterTypeAsm int

const (
	AX_REGISTER_ASM RegisterTypeAsm = iota
	CX_REGISTER_ASM
	DX_REGISTER_ASM
	DI_REGISTER_ASM
	SI_REGISTER_ASM
	R8_REGISTER_ASM
	R9_REGISTER_ASM
	R10_REGISTER_ASM
	R11_REGISTER_ASM
	SP_REGISTER_ASM
)

// the first six arguments when calling a function are placed in these registers
var ARG_REGISTERS = []RegisterTypeAsm{DI_REGISTER_ASM, SI_REGISTER_ASM, DX_REGISTER_ASM,
	CX_REGISTER_ASM, R8_REGISTER_ASM, R9_REGISTER_ASM}

//###############################################################################
//###############################################################################
//###############################################################################

func doAssemblyGen(tacky Program_Tacky) Program_Asm {
	asm := tacky.convertToAsm()

	// move symbolTable data to symbolTableBackend
	for name, sym := range symbolTable {
		asmTyp := dataTypeEnumToAssemblyTypeEnum(sym.dataTyp.typ)
		symAsm := Symbol_Asm{asmTyp: asmTyp, isStatic: (sym.attrs == STATIC_ATTRIBUTES), defined: sym.defined}
		symbolTableBackend[name] = symAsm
	}

	asm.replacePseudoregisters()
	asm.instructionFixup()

	return asm
}

/////////////////////////////////////////////////////////////////////////////////

func (pr *Program_Tacky) convertToAsm() Program_Asm {
	topItems := []Top_Level_Asm{}

	for _, item := range pr.topItems {
		itemAsm := item.topLevelToAsm()
		topItems = append(topItems, itemAsm)
	}

	asm := Program_Asm{topItems: topItems}
	return asm
}

/////////////////////////////////////////////////////////////////////////////////

func (fn *Function_Definition_Tacky) topLevelToAsm() Top_Level_Asm {
	fnAsm := Function_Asm{name: fn.name, global: fn.global, stackSize: 0}
	instructions := []Instruction_Asm{}

	// when we call a function that isn't main, at the beginning of that function
	// we move all parameters passed to us onto the stack
	// TODO: eventually we'll support main with parameters (argc, argv)
	if fn.name != "main" {
		for index, param := range fn.paramNames {
			if index < 6 {
				src := Register_Operand_Asm{ARG_REGISTERS[index]}
				mov := Mov_Instruction_Asm{asmTyp: getAsmTypeOfVariable(param), src: &src, dst: &Pseudoregister_Operand_Asm{param}}
				instructions = append(instructions, &mov)
			} else {
				// the seventh parameter is at Stack(16), the eighth is at Stack(24), etc.
				stackOffset := ((index - 4) * 8)
				src := Stack_Operand_Asm{int32(stackOffset)}
				mov := Mov_Instruction_Asm{asmTyp: getAsmTypeOfVariable(param), src: &src, dst: &Pseudoregister_Operand_Asm{param}}
				instructions = append(instructions, &mov)
			}
		}
	}

	for _, instrTacky := range fn.body {
		convertedInstructions := instrTacky.instructionToAsm()
		instructions = append(instructions, convertedInstructions...)
	}

	fnAsm.instructions = instructions
	return &fnAsm
}

/////////////////////////////////////////////////////////////////////////////////

func (st *Static_Variable_Tacky) topLevelToAsm() Top_Level_Asm {
	align := initToAlignment(st.initEnum)
	return &Static_Variable_Asm{name: st.name, global: st.global, alignment: align, initialValue: st.initialValue, initEnum: st.initEnum}
}

//###############################################################################
//###############################################################################
//###############################################################################

func (instr *Return_Instruction_Tacky) instructionToAsm() []Instruction_Asm {
	src := instr.val.valueToAsm()
	retType := instr.val.getAssemblyType()
	dst := Register_Operand_Asm{reg: AX_REGISTER_ASM}
	movInstr := Mov_Instruction_Asm{asmTyp: retType, src: src, dst: &dst}
	retInstr := Ret_Instruction_Asm{}

	instructions := []Instruction_Asm{&movInstr, &retInstr}
	return instructions
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Sign_Extend_Instruction_Tacky) instructionToAsm() []Instruction_Asm {
	mov := Movsx_Instruction_Asm{src: instr.src.valueToAsm(), dst: instr.dst.valueToAsm()}
	return []Instruction_Asm{&mov}
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Truncate_Instruction_Tacky) instructionToAsm() []Instruction_Asm {
	mov := Mov_Instruction_Asm{asmTyp: LONGWORD_ASM_TYPE, src: instr.src.valueToAsm(), dst: instr.dst.valueToAsm()}
	return []Instruction_Asm{&mov}
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Unary_Instruction_Tacky) instructionToAsm() []Instruction_Asm {
	if instr.unOp == NOT_OPERATOR {
		cmp := Compare_Instruction_Asm{asmTyp: instr.src.getAssemblyType(), op1: &Immediate_Int_Operand_Asm{"0"}, op2: instr.src.valueToAsm()}
		mov := Mov_Instruction_Asm{asmTyp: instr.dst.getAssemblyType(), src: &Immediate_Int_Operand_Asm{"0"}, dst: instr.dst.valueToAsm()}
		setC := Set_Conditional_Instruction_Asm{code: IS_EQUAL_CODE_ASM, dst: instr.dst.valueToAsm()}

		instructions := []Instruction_Asm{&cmp, &mov, &setC}
		return instructions
	} else {
		src := instr.src.valueToAsm()
		dst := instr.dst.valueToAsm()
		movInstr := Mov_Instruction_Asm{asmTyp: instr.src.getAssemblyType(), src: src, dst: dst}
		unaryInstr := Unary_Instruction_Asm{unOp: convertUnaryOpToAsm(instr.unOp), asmTyp: instr.dst.getAssemblyType(), src: dst}

		instructions := []Instruction_Asm{&movInstr, &unaryInstr}
		return instructions
	}
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Binary_Instruction_Tacky) instructionToAsm() []Instruction_Asm {
	if instr.binOp == ADD_OPERATOR || instr.binOp == SUBTRACT_OPERATOR || instr.binOp == MULTIPLY_OPERATOR {
		src1 := instr.src1.valueToAsm()
		dst := instr.dst.valueToAsm()
		movInstr := Mov_Instruction_Asm{asmTyp: instr.src1.getAssemblyType(), src: src1, dst: dst}

		src2 := instr.src2.valueToAsm()
		binInstr := Binary_Instruction_Asm{binOp: convertBinaryOpToAsm(instr.binOp), asmTyp: instr.src2.getAssemblyType(), src: src2, dst: dst}

		instructions := []Instruction_Asm{&movInstr, &binInstr}
		return instructions
	} else if instr.binOp == DIVIDE_OPERATOR {
		firstMov := Mov_Instruction_Asm{asmTyp: instr.src1.getAssemblyType(), src: instr.src1.valueToAsm(), dst: &Register_Operand_Asm{AX_REGISTER_ASM}}
		cdq := CDQ_Sign_Extend_Instruction_Asm{asmTyp: instr.src1.getAssemblyType()}
		idiv := IDivide_Instruction_Asm{asmTyp: instr.src2.getAssemblyType(), divisor: instr.src2.valueToAsm()}
		secondMov := Mov_Instruction_Asm{asmTyp: instr.dst.getAssemblyType(), src: &Register_Operand_Asm{AX_REGISTER_ASM}, dst: instr.dst.valueToAsm()}
		instructions := []Instruction_Asm{&firstMov, &cdq, &idiv, &secondMov}
		return instructions
	} else if instr.binOp == REMAINDER_OPERATOR {
		firstMov := Mov_Instruction_Asm{asmTyp: instr.src1.getAssemblyType(), src: instr.src1.valueToAsm(), dst: &Register_Operand_Asm{AX_REGISTER_ASM}}
		cdq := CDQ_Sign_Extend_Instruction_Asm{asmTyp: instr.src1.getAssemblyType()}
		idiv := IDivide_Instruction_Asm{asmTyp: instr.src2.getAssemblyType(), divisor: instr.src2.valueToAsm()}
		secondMov := Mov_Instruction_Asm{asmTyp: instr.dst.getAssemblyType(), src: &Register_Operand_Asm{DX_REGISTER_ASM}, dst: instr.dst.valueToAsm()}
		instructions := []Instruction_Asm{&firstMov, &cdq, &idiv, &secondMov}
		return instructions
	} else if instr.binOp == IS_EQUAL_OPERATOR || instr.binOp == NOT_EQUAL_OPERATOR || instr.binOp == LESS_THAN_OPERATOR ||
		instr.binOp == LESS_OR_EQUAL_OPERATOR || instr.binOp == GREATER_THAN_OPERATOR || instr.binOp == GREATER_OR_EQUAL_OPERATOR {
		cmp := Compare_Instruction_Asm{asmTyp: instr.src1.getAssemblyType(), op1: instr.src2.valueToAsm(), op2: instr.src1.valueToAsm()}
		mov := Mov_Instruction_Asm{asmTyp: instr.dst.getAssemblyType(), src: &Immediate_Int_Operand_Asm{"0"}, dst: instr.dst.valueToAsm()}
		setC := Set_Conditional_Instruction_Asm{code: convertBinaryOpToCondition(instr.binOp), dst: instr.dst.valueToAsm()}
		instructions := []Instruction_Asm{&cmp, &mov, &setC}
		return instructions
	} else {
		fail("unknown Binary_Instruction_Tacky")
	}

	return nil
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Copy_Instruction_Tacky) instructionToAsm() []Instruction_Asm {
	mov := Mov_Instruction_Asm{asmTyp: instr.src.getAssemblyType(), src: instr.src.valueToAsm(), dst: instr.dst.valueToAsm()}
	return []Instruction_Asm{&mov}
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Jump_Instruction_Tacky) instructionToAsm() []Instruction_Asm {
	jmp := Jump_Instruction_Asm{instr.target}
	return []Instruction_Asm{&jmp}
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Jump_If_Zero_Instruction_Tacky) instructionToAsm() []Instruction_Asm {
	asmTyp := instr.condition.getAssemblyType()
	cmp := Compare_Instruction_Asm{asmTyp: asmTyp, op1: &Immediate_Int_Operand_Asm{"0"}, op2: instr.condition.valueToAsm()}
	jmpC := Jump_Conditional_Instruction_Asm{code: IS_EQUAL_CODE_ASM, target: instr.target}
	return []Instruction_Asm{&cmp, &jmpC}
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Jump_If_Not_Zero_Instruction_Tacky) instructionToAsm() []Instruction_Asm {
	asmTyp := instr.condition.getAssemblyType()
	cmp := Compare_Instruction_Asm{asmTyp: asmTyp, op1: &Immediate_Int_Operand_Asm{"0"}, op2: instr.condition.valueToAsm()}
	jmpC := Jump_Conditional_Instruction_Asm{code: NOT_EQUAL_CODE_ASM, target: instr.target}
	return []Instruction_Asm{&cmp, &jmpC}
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Label_Instruction_Tacky) instructionToAsm() []Instruction_Asm {
	label := Label_Instruction_Asm{instr.name}
	return []Instruction_Asm{&label}
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Function_Call_Tacky) instructionToAsm() []Instruction_Asm {
	instructions := []Instruction_Asm{}

	// split the args into two: some go to registers, others go to the stack
	numArgs := len(instr.args)
	numRegisterArgs := minInt(numArgs, 6)
	numStackArgs := numArgs - numRegisterArgs

	registerArgs := instr.args[0:numRegisterArgs]
	stackArgs := instr.args[numRegisterArgs:numArgs]

	// adjust the stack alignment
	var stackPadding int32
	if (numStackArgs % 2) == 1 {
		stackPadding = 8
	} else {
		stackPadding = 0
	}

	if stackPadding != 0 {
		// allocate some space on the stack for the padding
		src := Immediate_Int_Operand_Asm{strconv.FormatInt(int64(stackPadding), 10)}
		dst := Register_Operand_Asm{SP_REGISTER_ASM}
		instr := Binary_Instruction_Asm{binOp: SUB_OPERATOR_ASM, asmTyp: QUADWORD_ASM_TYPE, src: &src, dst: &dst}
		instructions = append(instructions, &instr)
	}

	// pass some args in registers
	for index, arg := range registerArgs {
		src := arg.valueToAsm()
		dst := Register_Operand_Asm{ARG_REGISTERS[index]}
		mov := Mov_Instruction_Asm{asmTyp: arg.getAssemblyType(), src: src, dst: &dst}
		instructions = append(instructions, &mov)
	}

	// pass some args on the stack
	for index := numStackArgs - 1; index >= 0; index-- {
		src := stackArgs[index].valueToAsm()
		srcTyp := stackArgs[index].getAssemblyType()
		pushRightAway := canPushToStack(src)
		if pushRightAway || (srcTyp == QUADWORD_ASM_TYPE) {
			push := Push_Instruction_Asm{src}
			instructions = append(instructions, &push)
		} else {
			mov := Mov_Instruction_Asm{asmTyp: srcTyp, src: src, dst: &Register_Operand_Asm{AX_REGISTER_ASM}}
			push := Push_Instruction_Asm{&Register_Operand_Asm{AX_REGISTER_ASM}}
			instructions = append(instructions, &mov)
			instructions = append(instructions, &push)
		}
	}

	// call the function
	call := Call_Function_Asm{instr.funcName}
	instructions = append(instructions, &call)

	// adjust the stack pointer when we return from the function we just called
	bytesToRemove := int32(8*len(stackArgs)) + stackPadding
	if bytesToRemove != 0 {
		src := Immediate_Int_Operand_Asm{strconv.FormatInt(int64(bytesToRemove), 10)}
		dst := Register_Operand_Asm{SP_REGISTER_ASM}
		instr := Binary_Instruction_Asm{binOp: ADD_OPERATOR_ASM, asmTyp: QUADWORD_ASM_TYPE, src: &src, dst: &dst}
		instructions = append(instructions, &instr)
	}

	// retrieve the return value
	dst := instr.returnVal.valueToAsm()
	mov := Mov_Instruction_Asm{asmTyp: instr.returnVal.getAssemblyType(), src: &Register_Operand_Asm{AX_REGISTER_ASM}, dst: dst}
	instructions = append(instructions, &mov)

	return instructions
}

/////////////////////////////////////////////////////////////////////////////////

func canPushToStack(op Operand_Asm) bool {
	switch op.(type) {
	case *Immediate_Int_Operand_Asm:
		return true
	case *Register_Operand_Asm:
		return true
	default:
		return false
	}
}

//###############################################################################
//###############################################################################
//###############################################################################

func (val *Constant_Value_Tacky) valueToAsm() Operand_Asm {
	return &Immediate_Int_Operand_Asm{value: val.value}
}

/////////////////////////////////////////////////////////////////////////////////

func (val *Constant_Value_Tacky) getAssemblyType() AssemblyTypeEnum {
	return dataTypeEnumToAssemblyTypeEnum(val.typ)
}

/////////////////////////////////////////////////////////////////////////////////

func (val *Variable_Value_Tacky) valueToAsm() Operand_Asm {
	return &Pseudoregister_Operand_Asm{name: val.name}
}

/////////////////////////////////////////////////////////////////////////////////

func (val *Variable_Value_Tacky) getAssemblyType() AssemblyTypeEnum {
	return getAsmTypeOfVariable(val.name)
}

//###############################################################################
//###############################################################################
//###############################################################################

func (pr *Program_Asm) replacePseudoregisters() {
	for index, _ := range pr.topItems {
		nameToOffset := make(map[string]int32)
		// store the stack size for the function
		pr.topItems[index].replacePseudoregisters(nameToOffset)
	}
}

/////////////////////////////////////////////////////////////////////////////////

func (st *Static_Variable_Asm) replacePseudoregisters(nameToOffset map[string]int32) {
	// nothing to do here
}

/////////////////////////////////////////////////////////////////////////////////

func (fn *Function_Asm) replacePseudoregisters(nameToOffset map[string]int32) {
	for index, _ := range fn.instructions {
		switch convertedInstr := fn.instructions[index].(type) {
		case *Mov_Instruction_Asm:
			convertedInstr.src = replaceIfPseudoregister(convertedInstr.src, &fn.stackSize, nameToOffset)
			convertedInstr.dst = replaceIfPseudoregister(convertedInstr.dst, &fn.stackSize, nameToOffset)
			fn.instructions[index] = convertedInstr
		case *Movsx_Instruction_Asm:
			convertedInstr.src = replaceIfPseudoregister(convertedInstr.src, &fn.stackSize, nameToOffset)
			convertedInstr.dst = replaceIfPseudoregister(convertedInstr.dst, &fn.stackSize, nameToOffset)
			fn.instructions[index] = convertedInstr
		case *Unary_Instruction_Asm:
			convertedInstr.src = replaceIfPseudoregister(convertedInstr.src, &fn.stackSize, nameToOffset)
			fn.instructions[index] = convertedInstr
		case *Binary_Instruction_Asm:
			convertedInstr.src = replaceIfPseudoregister(convertedInstr.src, &fn.stackSize, nameToOffset)
			convertedInstr.dst = replaceIfPseudoregister(convertedInstr.dst, &fn.stackSize, nameToOffset)
			fn.instructions[index] = convertedInstr
		case *IDivide_Instruction_Asm:
			convertedInstr.divisor = replaceIfPseudoregister(convertedInstr.divisor, &fn.stackSize, nameToOffset)
			fn.instructions[index] = convertedInstr
		case *Compare_Instruction_Asm:
			convertedInstr.op1 = replaceIfPseudoregister(convertedInstr.op1, &fn.stackSize, nameToOffset)
			convertedInstr.op2 = replaceIfPseudoregister(convertedInstr.op2, &fn.stackSize, nameToOffset)
			fn.instructions[index] = convertedInstr
		case *Set_Conditional_Instruction_Asm:
			convertedInstr.dst = replaceIfPseudoregister(convertedInstr.dst, &fn.stackSize, nameToOffset)
			fn.instructions[index] = convertedInstr
		case *Push_Instruction_Asm:
			convertedInstr.op = replaceIfPseudoregister(convertedInstr.op, &fn.stackSize, nameToOffset)
			fn.instructions[index] = convertedInstr
		}
	}
}

/////////////////////////////////////////////////////////////////////////////////

func replaceIfPseudoregister(op Operand_Asm, stackOffset *int32, nameToOffset map[string]int32) Operand_Asm {
	if op == nil {
		return nil
	}

	convertedOp, isPseudo := op.(*Pseudoregister_Operand_Asm)

	if !isPseudo {
		return op
	}

	existingOffset, alreadyExists := nameToOffset[convertedOp.name]
	if alreadyExists {
		return &Stack_Operand_Asm{value: existingOffset}
	} else {
		asmSym, inSymTable := symbolTableBackend[convertedOp.name]
		if inSymTable && asmSym.isStatic {
			return &Data_Operand_Asm{name: convertedOp.name}
		}

		*stackOffset = *stackOffset - asmTypToAlignment(asmSym.asmTyp)

		// need to make sure that Quadwords are 8-byte aligned on the stack
		if asmSym.asmTyp == QUADWORD_ASM_TYPE {
			remainder := *stackOffset % 8
			if remainder != 0 {
				*stackOffset = (*stackOffset/8)*8 - 8
			}
		}

		nameToOffset[convertedOp.name] = *stackOffset
		return &Stack_Operand_Asm{value: *stackOffset}
	}
}

//###############################################################################
//###############################################################################
//###############################################################################

func (pr *Program_Asm) instructionFixup() {
	for index, _ := range pr.topItems {
		pr.topItems[index].fixInvalidInstr()
	}
}

/////////////////////////////////////////////////////////////////////////////////

func (fn *Function_Asm) fixInvalidInstr() {
	// round up the stack size to the nearest multiple of 16, although we're actually rounding down since it's negative...
	newStackSize := fn.stackSize
	remainder := newStackSize % 16
	if remainder != 0 {
		newStackSize = (newStackSize/16)*16 - 16
	}
	fn.stackSize = newStackSize

	// insert instruction to allocate space on the stack
	src := Immediate_Int_Operand_Asm{strconv.FormatInt(int64(-fn.stackSize), 10)}
	dst := Register_Operand_Asm{SP_REGISTER_ASM}
	firstInstr := Binary_Instruction_Asm{binOp: SUB_OPERATOR_ASM, asmTyp: QUADWORD_ASM_TYPE, src: &src, dst: &dst}
	instructions := []Instruction_Asm{&firstInstr}
	fn.instructions = append(instructions, fn.instructions...)

	// rewrite invalid instructions, they can't have both operands be memory addresses
	instructions = []Instruction_Asm{}

	for _, instr := range fn.instructions {

		switch convertedInstr := instr.(type) {
		case *Mov_Instruction_Asm:
			newInstrs := convertedInstr.fixInvalidInstr()
			instructions = append(instructions, newInstrs...)
		case *Movsx_Instruction_Asm:
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
		case *Push_Instruction_Asm:
			newInstrs := convertedInstr.fixInvalidInstr()
			instructions = append(instructions, newInstrs...)
		default:
			// don't need to fix it, just add it to the list
			instructions = append(instructions, instr)
		}
	}

	fn.instructions = instructions
}

/////////////////////////////////////////////////////////////////////////////////

func (st *Static_Variable_Asm) fixInvalidInstr() {
	// nothing to do here
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Mov_Instruction_Asm) fixInvalidInstr() []Instruction_Asm {
	_, srcIsStack := instr.src.(*Stack_Operand_Asm)
	_, dstIsStack := instr.dst.(*Stack_Operand_Asm)
	_, srcIsStatic := instr.src.(*Data_Operand_Asm)
	_, dstIsStatic := instr.dst.(*Data_Operand_Asm)
	// TODO: page 268 of the book??

	if (srcIsStack || srcIsStatic) && (dstIsStack || dstIsStatic) {
		intermediateOperand := Register_Operand_Asm{R10_REGISTER_ASM}
		firstInstr := Mov_Instruction_Asm{asmTyp: instr.asmTyp, src: instr.src, dst: &intermediateOperand}
		secondInstr := Mov_Instruction_Asm{asmTyp: instr.asmTyp, src: &intermediateOperand, dst: instr.dst}
		return []Instruction_Asm{&firstInstr, &secondInstr}
	}

	return []Instruction_Asm{instr}
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Movsx_Instruction_Asm) fixInvalidInstr() []Instruction_Asm {
	_, srcIsConst := instr.src.(*Immediate_Int_Operand_Asm)
	_, dstIsStack := instr.dst.(*Stack_Operand_Asm)

	r10 := Register_Operand_Asm{R10_REGISTER_ASM}
	r11 := Register_Operand_Asm{R11_REGISTER_ASM}

	if srcIsConst && dstIsStack {
		firstInstr := Mov_Instruction_Asm{asmTyp: LONGWORD_ASM_TYPE, src: instr.src, dst: &r10}
		secInstr := Movsx_Instruction_Asm{src: &r10, dst: &r11}
		thirdInstr := Mov_Instruction_Asm{asmTyp: QUADWORD_ASM_TYPE, src: &r11, dst: instr.dst}
		return []Instruction_Asm{&firstInstr, &secInstr, &thirdInstr}
	} else if srcIsConst {
		firstInstr := Mov_Instruction_Asm{asmTyp: LONGWORD_ASM_TYPE, src: instr.src, dst: &r10}
		secInstr := Movsx_Instruction_Asm{src: &r10, dst: instr.dst}
		return []Instruction_Asm{&firstInstr, &secInstr}
	} else if dstIsStack {
		firstInstr := Movsx_Instruction_Asm{src: instr.src, dst: &r11}
		secInstr := Mov_Instruction_Asm{asmTyp: QUADWORD_ASM_TYPE, src: &r11, dst: instr.dst}
		return []Instruction_Asm{&firstInstr, &secInstr}
	}

	return []Instruction_Asm{instr}
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Binary_Instruction_Asm) fixInvalidInstr() []Instruction_Asm {
	if instr.binOp == ADD_OPERATOR_ASM || instr.binOp == SUB_OPERATOR_ASM {
		_, srcIsStack := instr.src.(*Stack_Operand_Asm)
		_, dstIsStack := instr.dst.(*Stack_Operand_Asm)
		_, srcIsStatic := instr.src.(*Data_Operand_Asm)
		_, dstIsStatic := instr.dst.(*Data_Operand_Asm)
		// TODO: page 268 of the book
		//srcIsBigImm := opIsBigImm(instr.src)
		//dstIsBigImm := opIsBigImm(instr.dst)
		//isQuadInstr := instr.asmTyp == QUADWORD_ASM_TYPE

		if (srcIsStack || srcIsStatic) && (dstIsStack || dstIsStatic) {
			intermediateOperand := Register_Operand_Asm{R10_REGISTER_ASM}
			firstInstr := Mov_Instruction_Asm{asmTyp: instr.asmTyp, src: instr.src, dst: &intermediateOperand}
			secondInstr := Binary_Instruction_Asm{binOp: instr.binOp, asmTyp: instr.asmTyp, src: &intermediateOperand, dst: instr.dst}
			return []Instruction_Asm{&firstInstr, &secondInstr}
		} /*else if isQuadInstr {

		}*/
	} else if instr.binOp == MULT_OPERATOR_ASM {
		_, dstIsStack := instr.dst.(*Stack_Operand_Asm)
		_, dstIsStatic := instr.dst.(*Data_Operand_Asm)
		// TODO: page 268 of the book

		if dstIsStack || dstIsStatic {
			firstInstr := Mov_Instruction_Asm{asmTyp: instr.asmTyp, src: instr.dst, dst: &Register_Operand_Asm{R11_REGISTER_ASM}}
			secondInstr := Binary_Instruction_Asm{binOp: instr.binOp, asmTyp: instr.asmTyp, src: instr.src, dst: &Register_Operand_Asm{R11_REGISTER_ASM}}
			thirdInstr := Mov_Instruction_Asm{asmTyp: instr.asmTyp, src: &Register_Operand_Asm{R11_REGISTER_ASM}, dst: instr.dst}
			return []Instruction_Asm{&firstInstr, &secondInstr, &thirdInstr}
		}
	}

	return []Instruction_Asm{instr}
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *IDivide_Instruction_Asm) fixInvalidInstr() []Instruction_Asm {
	_, isConstant := instr.divisor.(*Immediate_Int_Operand_Asm)

	if isConstant {
		firstInstr := Mov_Instruction_Asm{asmTyp: instr.asmTyp, src: instr.divisor, dst: &Register_Operand_Asm{R10_REGISTER_ASM}}
		secondInstr := IDivide_Instruction_Asm{asmTyp: instr.asmTyp, divisor: &Register_Operand_Asm{R10_REGISTER_ASM}}
		return []Instruction_Asm{&firstInstr, &secondInstr}
	}

	return []Instruction_Asm{instr}
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Compare_Instruction_Asm) fixInvalidInstr() []Instruction_Asm {
	_, op1IsStack := instr.op1.(*Stack_Operand_Asm)
	_, op2IsStack := instr.op2.(*Stack_Operand_Asm)
	_, op1IsStatic := instr.op1.(*Data_Operand_Asm)
	_, op2IsStatic := instr.op2.(*Data_Operand_Asm)
	_, op2IsConstant := instr.op2.(*Immediate_Int_Operand_Asm)
	// TODO: page 268 of the book

	if (op1IsStack || op1IsStatic) && (op2IsStack || op2IsStatic) {
		intermediateOperand := Register_Operand_Asm{R10_REGISTER_ASM}
		firstInstr := Mov_Instruction_Asm{asmTyp: instr.asmTyp, src: instr.op1, dst: &intermediateOperand}
		secondInstr := Compare_Instruction_Asm{asmTyp: instr.asmTyp, op1: &intermediateOperand, op2: instr.op2}
		return []Instruction_Asm{&firstInstr, &secondInstr}
	} else if op2IsConstant {
		firstInstr := Mov_Instruction_Asm{asmTyp: instr.asmTyp, src: instr.op2, dst: &Register_Operand_Asm{R11_REGISTER_ASM}}
		secondInstr := Compare_Instruction_Asm{asmTyp: instr.asmTyp, op1: instr.op1, op2: &Register_Operand_Asm{R11_REGISTER_ASM}}
		return []Instruction_Asm{&firstInstr, &secondInstr}
	}

	return []Instruction_Asm{instr}
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Push_Instruction_Asm) fixInvalidInstr() []Instruction_Asm {
	// TODO: page 268 of the book
	return []Instruction_Asm{instr}
}
