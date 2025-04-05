package main

import (
	"math"
	"os"
	"strconv"
)

/////////////////////////////////////////////////////////////////////////////////

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
	DOUBLE_ASM_TYPE
)

func dataTypeEnumToAssemblyTypeEnum(input DataTypeEnum) AssemblyTypeEnum {
	switch input {
	case INT_TYPE:
		return LONGWORD_ASM_TYPE
	case LONG_TYPE:
		return QUADWORD_ASM_TYPE
	case UNSIGNED_INT_TYPE:
		return LONGWORD_ASM_TYPE
	case UNSIGNED_LONG_TYPE:
		return QUADWORD_ASM_TYPE
	case DOUBLE_TYPE:
		return DOUBLE_ASM_TYPE
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
	case DOUBLE_ASM_TYPE:
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
	case INITIAL_UNSIGNED_INT:
		return 4
	case INITIAL_UNSIGNED_LONG:
		return 8
	case INITIAL_DOUBLE:
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
	// for variables and constant values
	asmTyp     AssemblyTypeEnum
	isStatic   bool
	isConstant bool

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

/////////////////////////////////////////////////////////////////////////////////

// to represent floating point constants
type Static_Constant_Asm struct {
	name         string
	alignment    int32
	initialValue string
	initEnum     InitializerEnum
}

var allStaticConstants = make([]Static_Constant_Asm, 0, 10)

func addStaticConstant(name string, alignment int32, initialValue string, initEnum InitializerEnum) string {
	// check if we already have one that matches, if so then return its name
	for _, st := range allStaticConstants {
		if (st.alignment == alignment) && (st.initialValue == initialValue) && (st.initEnum == initEnum) {
			return st.name
		}
	}

	// doesn't exist yet, so create one and add it
	newName := makeLabelName(name)
	st := Static_Constant_Asm{name: newName, alignment: alignment, initialValue: initialValue, initEnum: initEnum}
	allStaticConstants = append(allStaticConstants, st)
	return newName
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

type Move_Zero_Extend_Instruction_Asm struct {
	src Operand_Asm
	dst Operand_Asm
}

/////////////////////////////////////////////////////////////////////////////////

// Convert scalar double to signed int
type Cvttsd2si_Double_To_Int_Instruction_Asm struct {
	dstAsmType AssemblyTypeEnum
	src        Operand_Asm
	dst        Operand_Asm
}

/////////////////////////////////////////////////////////////////////////////////

// Convert signed int to scalar double
type Cvtsi2sd_Int_To_Double_Instruction_Asm struct {
	srcAsmType AssemblyTypeEnum
	src        Operand_Asm
	dst        Operand_Asm
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

type Divide_Instruction_Asm struct {
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
	SHIFT_RIGHT_OPERATOR_ASM
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
	DIV_DOUBLE_OPERATOR_ASM
	AND_OPERATOR_ASM
	OR_OPERATOR_ASM
	XOR_OPERATOR_ASM
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
	// TODO: ParseUint?
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
	LESS_THAN_CODE_UNSIGNED_ASM
	LESS_OR_EQUAL_CODE_UNSIGNED_ASM
	GREATER_THAN_CODE_UNSIGNED_ASM
	GREATER_OR_EQUAL_CODE_UNSIGNED_ASM
)

func convertBinaryOpToCondition(binOp BinaryOperatorType, signed bool) ConditionalCodeAsm {
	switch binOp {
	case IS_EQUAL_OPERATOR:
		return IS_EQUAL_CODE_ASM
	case NOT_EQUAL_OPERATOR:
		return NOT_EQUAL_CODE_ASM
	case LESS_THAN_OPERATOR:
		if signed {
			return LESS_THAN_CODE_ASM
		} else {
			return LESS_THAN_CODE_UNSIGNED_ASM
		}
	case LESS_OR_EQUAL_OPERATOR:
		if signed {
			return LESS_OR_EQUAL_CODE_ASM
		} else {
			return LESS_OR_EQUAL_CODE_UNSIGNED_ASM
		}
	case GREATER_THAN_OPERATOR:
		if signed {
			return GREATER_THAN_CODE_ASM
		} else {
			return GREATER_THAN_CODE_UNSIGNED_ASM
		}
	case GREATER_OR_EQUAL_OPERATOR:
		if signed {
			return GREATER_OR_EQUAL_CODE_ASM
		} else {
			return GREATER_OR_EQUAL_CODE_UNSIGNED_ASM
		}
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
	XMM0_REGISTER_ASM
	XMM1_REGISTER_ASM
	XMM2_REGISTER_ASM
	XMM3_REGISTER_ASM
	XMM4_REGISTER_ASM
	XMM5_REGISTER_ASM
	XMM6_REGISTER_ASM
	XMM7_REGISTER_ASM
	XMM14_REGISTER_ASM
	XMM15_REGISTER_ASM
)

// the first six arguments when calling a function are placed in these registers
var INT_ARG_REGISTERS = []RegisterTypeAsm{DI_REGISTER_ASM, SI_REGISTER_ASM, DX_REGISTER_ASM,
	CX_REGISTER_ASM, R8_REGISTER_ASM, R9_REGISTER_ASM}

// the first eight double arguments when calling a function are placed in these registers
var DOUBLE_ARG_REGISTERS = []RegisterTypeAsm{XMM0_REGISTER_ASM, XMM1_REGISTER_ASM, XMM2_REGISTER_ASM, XMM3_REGISTER_ASM,
	XMM4_REGISTER_ASM, XMM5_REGISTER_ASM, XMM6_REGISTER_ASM, XMM7_REGISTER_ASM}

//###############################################################################
//###############################################################################
//###############################################################################

func doAssemblyGen(tacky Program_Tacky) Program_Asm {
	asm := tacky.convertToAsm()

	// add all Static_Constant_Asm to the list of top-level constructs and the symbolTableBackend
	for _, stConst := range allStaticConstants {
		asm.topItems = append(asm.topItems, &stConst)
		symAsm := Symbol_Asm{asmTyp: DOUBLE_ASM_TYPE, isStatic: true, isConstant: true, defined: false}
		symbolTableBackend[stConst.name] = symAsm
	}

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
		intRegParams, doubleRegParams, stackParams := classifyParameters(fn.paramNames)

		// copy parameters from general purpose registers
		for index, param := range intRegParams {
			src := Register_Operand_Asm{INT_ARG_REGISTERS[index]}
			mov := Mov_Instruction_Asm{asmTyp: getAsmTypeOfVariable(param), src: &src, dst: &Pseudoregister_Operand_Asm{param}}
			instructions = append(instructions, &mov)
		}

		// copy parameters from floating point registers
		for index, param := range doubleRegParams {
			src := Register_Operand_Asm{DOUBLE_ARG_REGISTERS[index]}
			mov := Mov_Instruction_Asm{asmTyp: DOUBLE_ASM_TYPE, src: &src, dst: &Pseudoregister_Operand_Asm{param}}
			instructions = append(instructions, &mov)
		}

		// copy parameters from the stack
		for index, param := range stackParams {
			// first parameter is at Stack(16), the next is at Stack(24), etc.
			stackOffset := 16 + index*8
			src := Stack_Operand_Asm{int32(stackOffset)}
			mov := Mov_Instruction_Asm{asmTyp: getAsmTypeOfVariable(param), src: &src, dst: &Pseudoregister_Operand_Asm{param}}
			instructions = append(instructions, &mov)
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
	retType := instr.val.getAssemblyType()

	var dst Register_Operand_Asm
	if retType == DOUBLE_ASM_TYPE {
		dst.reg = XMM0_REGISTER_ASM
	} else {
		dst.reg = AX_REGISTER_ASM
	}

	movInstr := Mov_Instruction_Asm{asmTyp: retType, src: instr.val.valueToAsm(), dst: &dst}
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

func (instr *Zero_Extend_Instruction_Tacky) instructionToAsm() []Instruction_Asm {
	mov := Move_Zero_Extend_Instruction_Asm{src: instr.src.valueToAsm(), dst: instr.dst.valueToAsm()}
	return []Instruction_Asm{&mov}
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Double_To_Int_Instruction_Tacky) instructionToAsm() []Instruction_Asm {
	cvt := Cvttsd2si_Double_To_Int_Instruction_Asm{dstAsmType: instr.dst.getAssemblyType(),
		src: instr.src.valueToAsm(), dst: instr.dst.valueToAsm()}
	return []Instruction_Asm{&cvt}
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Double_To_UInt_Instruction_Tacky) instructionToAsm() []Instruction_Asm {
	// TODO: page 328 and page 335, also see errata note on website
	fail("not implemented yet")
	return []Instruction_Asm{}
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Int_To_Double_Instruction_Tacky) instructionToAsm() []Instruction_Asm {
	cvt := Cvtsi2sd_Int_To_Double_Instruction_Asm{srcAsmType: instr.src.getAssemblyType(),
		src: instr.src.valueToAsm(), dst: instr.dst.valueToAsm()}
	return []Instruction_Asm{&cvt}
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *UInt_To_Double_Instruction_Tacky) instructionToAsm() []Instruction_Asm {
	// TODO: page 328 and page 335
	fail("not implemented yet")
	return []Instruction_Asm{}
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Unary_Instruction_Tacky) instructionToAsm() []Instruction_Asm {
	if (instr.unOp == NOT_OPERATOR) && (instr.src.getAssemblyType() == DOUBLE_ASM_TYPE) {
		xmm0 := Register_Operand_Asm{XMM0_REGISTER_ASM}
		bin := Binary_Instruction_Asm{binOp: XOR_OPERATOR_ASM, asmTyp: DOUBLE_ASM_TYPE, src: &xmm0, dst: &xmm0}
		cmp := Compare_Instruction_Asm{asmTyp: DOUBLE_ASM_TYPE, op1: instr.src.valueToAsm(), op2: &xmm0}
		mov := Mov_Instruction_Asm{asmTyp: instr.dst.getAssemblyType(), src: &Immediate_Int_Operand_Asm{"0"}, dst: instr.dst.valueToAsm()}
		setC := Set_Conditional_Instruction_Asm{code: IS_EQUAL_CODE_ASM, dst: instr.dst.valueToAsm()}
		return []Instruction_Asm{&bin, &cmp, &mov, &setC}
	} else if instr.unOp == NOT_OPERATOR {
		cmp := Compare_Instruction_Asm{asmTyp: instr.src.getAssemblyType(), op1: &Immediate_Int_Operand_Asm{"0"}, op2: instr.src.valueToAsm()}
		mov := Mov_Instruction_Asm{asmTyp: instr.dst.getAssemblyType(), src: &Immediate_Int_Operand_Asm{"0"}, dst: instr.dst.valueToAsm()}
		setC := Set_Conditional_Instruction_Asm{code: IS_EQUAL_CODE_ASM, dst: instr.dst.valueToAsm()}
		instructions := []Instruction_Asm{&cmp, &mov, &setC}
		return instructions
	} else if (instr.unOp == NEGATE_OPERATOR) && (instr.src.getAssemblyType() == DOUBLE_ASM_TYPE) {
		minusZero := addStaticConstant("minusZero", 16, "-0.0", INITIAL_DOUBLE)
		src := instr.src.valueToAsm()
		dst := instr.dst.valueToAsm()
		mov := Mov_Instruction_Asm{asmTyp: DOUBLE_ASM_TYPE, src: src, dst: dst}
		bin := Binary_Instruction_Asm{binOp: XOR_OPERATOR_ASM, asmTyp: DOUBLE_ASM_TYPE, src: &Data_Operand_Asm{minusZero}, dst: dst}
		instructions := []Instruction_Asm{&mov, &bin}
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
	} else if (instr.binOp == DIVIDE_OPERATOR) && (instr.src1.getAssemblyType() == DOUBLE_ASM_TYPE || instr.src2.getAssemblyType() == DOUBLE_ASM_TYPE) {
		src1 := instr.src1.valueToAsm()
		dst := instr.dst.valueToAsm()
		movInstr := Mov_Instruction_Asm{asmTyp: DOUBLE_ASM_TYPE, src: src1, dst: dst}

		src2 := instr.src2.valueToAsm()
		binInstr := Binary_Instruction_Asm{binOp: DIV_DOUBLE_OPERATOR_ASM, asmTyp: DOUBLE_ASM_TYPE, src: src2, dst: dst}

		instructions := []Instruction_Asm{&movInstr, &binInstr}
		return instructions
	} else if instr.binOp == DIVIDE_OPERATOR {
		if instr.src1.isSigned() {
			firstMov := Mov_Instruction_Asm{asmTyp: instr.src1.getAssemblyType(), src: instr.src1.valueToAsm(), dst: &Register_Operand_Asm{AX_REGISTER_ASM}}
			cdq := CDQ_Sign_Extend_Instruction_Asm{asmTyp: instr.src1.getAssemblyType()}
			idiv := IDivide_Instruction_Asm{asmTyp: instr.src2.getAssemblyType(), divisor: instr.src2.valueToAsm()}
			secondMov := Mov_Instruction_Asm{asmTyp: instr.dst.getAssemblyType(), src: &Register_Operand_Asm{AX_REGISTER_ASM}, dst: instr.dst.valueToAsm()}
			instructions := []Instruction_Asm{&firstMov, &cdq, &idiv, &secondMov}
			return instructions
		} else {
			firstMov := Mov_Instruction_Asm{asmTyp: instr.src1.getAssemblyType(), src: instr.src1.valueToAsm(), dst: &Register_Operand_Asm{AX_REGISTER_ASM}}
			zero := Mov_Instruction_Asm{asmTyp: instr.src1.getAssemblyType(), src: &Immediate_Int_Operand_Asm{"0"}, dst: &Register_Operand_Asm{DX_REGISTER_ASM}}
			div := Divide_Instruction_Asm{asmTyp: instr.src2.getAssemblyType(), divisor: instr.src2.valueToAsm()}
			secondMov := Mov_Instruction_Asm{asmTyp: instr.dst.getAssemblyType(), src: &Register_Operand_Asm{AX_REGISTER_ASM}, dst: instr.dst.valueToAsm()}
			instructions := []Instruction_Asm{&firstMov, &zero, &div, &secondMov}
			return instructions
		}
	} else if instr.binOp == REMAINDER_OPERATOR {
		if instr.src1.isSigned() {
			firstMov := Mov_Instruction_Asm{asmTyp: instr.src1.getAssemblyType(), src: instr.src1.valueToAsm(), dst: &Register_Operand_Asm{AX_REGISTER_ASM}}
			cdq := CDQ_Sign_Extend_Instruction_Asm{asmTyp: instr.src1.getAssemblyType()}
			idiv := IDivide_Instruction_Asm{asmTyp: instr.src2.getAssemblyType(), divisor: instr.src2.valueToAsm()}
			secondMov := Mov_Instruction_Asm{asmTyp: instr.dst.getAssemblyType(), src: &Register_Operand_Asm{DX_REGISTER_ASM}, dst: instr.dst.valueToAsm()}
			instructions := []Instruction_Asm{&firstMov, &cdq, &idiv, &secondMov}
			return instructions
		} else {
			firstMov := Mov_Instruction_Asm{asmTyp: instr.src1.getAssemblyType(), src: instr.src1.valueToAsm(), dst: &Register_Operand_Asm{AX_REGISTER_ASM}}
			zero := Mov_Instruction_Asm{asmTyp: instr.src1.getAssemblyType(), src: &Immediate_Int_Operand_Asm{"0"}, dst: &Register_Operand_Asm{DX_REGISTER_ASM}}
			div := Divide_Instruction_Asm{asmTyp: instr.src2.getAssemblyType(), divisor: instr.src2.valueToAsm()}
			secondMov := Mov_Instruction_Asm{asmTyp: instr.dst.getAssemblyType(), src: &Register_Operand_Asm{DX_REGISTER_ASM}, dst: instr.dst.valueToAsm()}
			instructions := []Instruction_Asm{&firstMov, &zero, &div, &secondMov}
			return instructions
		}
	} else if instr.binOp == IS_EQUAL_OPERATOR || instr.binOp == NOT_EQUAL_OPERATOR || instr.binOp == LESS_THAN_OPERATOR ||
		instr.binOp == LESS_OR_EQUAL_OPERATOR || instr.binOp == GREATER_THAN_OPERATOR || instr.binOp == GREATER_OR_EQUAL_OPERATOR {
		cmp := Compare_Instruction_Asm{asmTyp: instr.src1.getAssemblyType(), op1: instr.src2.valueToAsm(), op2: instr.src1.valueToAsm()}
		mov := Mov_Instruction_Asm{asmTyp: instr.dst.getAssemblyType(), src: &Immediate_Int_Operand_Asm{"0"}, dst: instr.dst.valueToAsm()}
		signed := instr.src1.isSigned()
		setC := Set_Conditional_Instruction_Asm{code: convertBinaryOpToCondition(instr.binOp, signed), dst: instr.dst.valueToAsm()}
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
	if asmTyp == DOUBLE_ASM_TYPE {
		xmm0 := Register_Operand_Asm{XMM0_REGISTER_ASM}
		bin := Binary_Instruction_Asm{binOp: XOR_OPERATOR_ASM, asmTyp: asmTyp, src: &xmm0, dst: &xmm0}
		cmp := Compare_Instruction_Asm{asmTyp: asmTyp, op1: instr.condition.valueToAsm(), op2: &xmm0}
		jmpC := Jump_Conditional_Instruction_Asm{code: IS_EQUAL_CODE_ASM, target: instr.target}
		return []Instruction_Asm{&bin, &cmp, &jmpC}
	} else {
		cmp := Compare_Instruction_Asm{asmTyp: asmTyp, op1: &Immediate_Int_Operand_Asm{"0"}, op2: instr.condition.valueToAsm()}
		jmpC := Jump_Conditional_Instruction_Asm{code: IS_EQUAL_CODE_ASM, target: instr.target}
		return []Instruction_Asm{&cmp, &jmpC}
	}
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Jump_If_Not_Zero_Instruction_Tacky) instructionToAsm() []Instruction_Asm {
	asmTyp := instr.condition.getAssemblyType()
	if asmTyp == DOUBLE_ASM_TYPE {
		xmm0 := Register_Operand_Asm{XMM0_REGISTER_ASM}
		bin := Binary_Instruction_Asm{binOp: XOR_OPERATOR_ASM, asmTyp: asmTyp, src: &xmm0, dst: &xmm0}
		cmp := Compare_Instruction_Asm{asmTyp: asmTyp, op1: instr.condition.valueToAsm(), op2: &xmm0}
		jmpC := Jump_Conditional_Instruction_Asm{code: NOT_EQUAL_CODE_ASM, target: instr.target}
		return []Instruction_Asm{&bin, &cmp, &jmpC}
	} else {
		cmp := Compare_Instruction_Asm{asmTyp: asmTyp, op1: &Immediate_Int_Operand_Asm{"0"}, op2: instr.condition.valueToAsm()}
		jmpC := Jump_Conditional_Instruction_Asm{code: NOT_EQUAL_CODE_ASM, target: instr.target}
		return []Instruction_Asm{&cmp, &jmpC}
	}
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Label_Instruction_Tacky) instructionToAsm() []Instruction_Asm {
	label := Label_Instruction_Asm{instr.name}
	return []Instruction_Asm{&label}
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Function_Call_Tacky) instructionToAsm() []Instruction_Asm {
	instructions := []Instruction_Asm{}

	// classify the arguments
	intRegArgs, doubleRegArgs, stackArgs := classifyParameters(instr.args)

	// adjust the stack alignment
	var stackPadding int32
	if (len(stackArgs) % 2) == 1 {
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

	// pass some args in general purpose registers
	for index, arg := range intRegArgs {
		src := arg.valueToAsm()
		dst := Register_Operand_Asm{INT_ARG_REGISTERS[index]}
		mov := Mov_Instruction_Asm{asmTyp: arg.getAssemblyType(), src: src, dst: &dst}
		instructions = append(instructions, &mov)
	}

	// pass some args in the floating point registers
	for index, arg := range doubleRegArgs {
		src := arg.valueToAsm()
		dst := Register_Operand_Asm{DOUBLE_ARG_REGISTERS[index]}
		mov := Mov_Instruction_Asm{asmTyp: DOUBLE_ASM_TYPE, src: src, dst: &dst}
		instructions = append(instructions, &mov)
	}

	// pass some args on the stack
	for index := len(stackArgs) - 1; index >= 0; index-- {
		src := stackArgs[index].valueToAsm()
		srcTyp := stackArgs[index].getAssemblyType()
		pushRightAway := canPushToStack(src)
		if pushRightAway || (srcTyp == QUADWORD_ASM_TYPE) || (srcTyp == DOUBLE_ASM_TYPE) {
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
	var src Register_Operand_Asm
	if instr.returnVal.getAssemblyType() == DOUBLE_ASM_TYPE {
		src.reg = XMM0_REGISTER_ASM
	} else {
		src.reg = AX_REGISTER_ASM
	}
	mov := Mov_Instruction_Asm{asmTyp: instr.returnVal.getAssemblyType(), src: &src, dst: instr.returnVal.valueToAsm()}
	instructions = append(instructions, &mov)

	return instructions
}

/////////////////////////////////////////////////////////////////////////////////

func classifyParameters[T any](params []T) ([]T, []T, []T) {
	intRegParams := []T{}
	doubleRegParams := []T{}
	stackParams := []T{}

	for _, p := range params {
		var typ AssemblyTypeEnum
		switch converted := any(p).(type) {
		case string:
			typ = getAsmTypeOfVariable(converted)
		case Value_Tacky:
			typ = converted.getAssemblyType()
		}

		if typ == DOUBLE_ASM_TYPE {
			if len(doubleRegParams) < 8 {
				doubleRegParams = append(doubleRegParams, p)
			} else {
				stackParams = append(stackParams, p)
			}
		} else {
			if len(intRegParams) < 6 {
				intRegParams = append(intRegParams, p)
			} else {
				stackParams = append(stackParams, p)
			}
		}
	}

	return intRegParams, doubleRegParams, stackParams
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
	if val.typ == DOUBLE_TYPE {
		opName := addStaticConstant("staticConst", initToAlignment(INITIAL_DOUBLE), val.value, INITIAL_DOUBLE)
		op := Data_Operand_Asm{opName}
		return &op
	} else {
		return &Immediate_Int_Operand_Asm{value: val.value}
	}
}

/////////////////////////////////////////////////////////////////////////////////

func (val *Constant_Value_Tacky) getAssemblyType() AssemblyTypeEnum {
	return dataTypeEnumToAssemblyTypeEnum(val.typ)
}

/////////////////////////////////////////////////////////////////////////////////

func (val *Constant_Value_Tacky) isSigned() bool {
	return isSigned(val.typ)
}

/////////////////////////////////////////////////////////////////////////////////

func (val *Variable_Value_Tacky) valueToAsm() Operand_Asm {
	return &Pseudoregister_Operand_Asm{name: val.name}
}

/////////////////////////////////////////////////////////////////////////////////

func (val *Variable_Value_Tacky) getAssemblyType() AssemblyTypeEnum {
	return getAsmTypeOfVariable(val.name)
}

/////////////////////////////////////////////////////////////////////////////////

func (val *Variable_Value_Tacky) isSigned() bool {
	return isSigned(symbolTable[val.name].dataTyp.typ)
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

func (st *Static_Constant_Asm) replacePseudoregisters(nameToOffset map[string]int32) {
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
		case *Move_Zero_Extend_Instruction_Asm:
			convertedInstr.src = replaceIfPseudoregister(convertedInstr.src, &fn.stackSize, nameToOffset)
			convertedInstr.dst = replaceIfPseudoregister(convertedInstr.dst, &fn.stackSize, nameToOffset)
			fn.instructions[index] = convertedInstr
		case *Cvttsd2si_Double_To_Int_Instruction_Asm:
			convertedInstr.src = replaceIfPseudoregister(convertedInstr.src, &fn.stackSize, nameToOffset)
			convertedInstr.dst = replaceIfPseudoregister(convertedInstr.dst, &fn.stackSize, nameToOffset)
			fn.instructions[index] = convertedInstr
		case *Cvtsi2sd_Int_To_Double_Instruction_Asm:
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
		case *Divide_Instruction_Asm:
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

		// need to make sure that Quadwords and Doubles are 8-byte aligned on the stack
		if (asmSym.asmTyp == QUADWORD_ASM_TYPE) || (asmSym.asmTyp == DOUBLE_ASM_TYPE) {
			remainder := *stackOffset % 8
			if remainder != 0 {
				// ex: -4 changes to -8, -12 changes to -16
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
		case *Move_Zero_Extend_Instruction_Asm:
			newInstrs := convertedInstr.fixInvalidInstr()
			instructions = append(instructions, newInstrs...)
		case *Cvttsd2si_Double_To_Int_Instruction_Asm:
			newInstrs := convertedInstr.fixInvalidInstr()
			instructions = append(instructions, newInstrs...)
		case *Cvtsi2sd_Int_To_Double_Instruction_Asm:
			newInstrs := convertedInstr.fixInvalidInstr()
			instructions = append(instructions, newInstrs...)
		case *Binary_Instruction_Asm:
			newInstrs := convertedInstr.fixInvalidInstr()
			instructions = append(instructions, newInstrs...)
		case *IDivide_Instruction_Asm:
			newInstrs := convertedInstr.fixInvalidInstr()
			instructions = append(instructions, newInstrs...)
		case *Divide_Instruction_Asm:
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

func (st *Static_Constant_Asm) fixInvalidInstr() {
	// nothing to do here
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Mov_Instruction_Asm) fixInvalidInstr() []Instruction_Asm {
	_, srcIsStack := instr.src.(*Stack_Operand_Asm)
	_, dstIsStack := instr.dst.(*Stack_Operand_Asm)
	_, srcIsStatic := instr.src.(*Data_Operand_Asm)
	_, dstIsStatic := instr.dst.(*Data_Operand_Asm)
	srcIsBigImm := opIsBigImm(instr.src)
	isQuadInstr := instr.asmTyp == QUADWORD_ASM_TYPE
	isDoubleInstr := instr.asmTyp == DOUBLE_ASM_TYPE

	if isDoubleInstr && ((srcIsStack || srcIsStatic) && (dstIsStack || dstIsStatic)) {
		// page 337, mov for doubles can't have both operands be in memory
		xmm14 := Register_Operand_Asm{XMM14_REGISTER_ASM}
		firstInstr := Mov_Instruction_Asm{asmTyp: DOUBLE_ASM_TYPE, src: instr.src, dst: &xmm14}
		secondInstr := Mov_Instruction_Asm{asmTyp: DOUBLE_ASM_TYPE, src: &xmm14, dst: instr.dst}
		return []Instruction_Asm{&firstInstr, &secondInstr}
	} else if ((srcIsStack || srcIsStatic) && (dstIsStack || dstIsStatic)) ||
		(isQuadInstr && srcIsBigImm && dstIsStack) ||
		(isQuadInstr && srcIsBigImm && dstIsStatic) {
		// from page 268 of the book, large Quadwords can't go directly to the stack (memory)
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

func (instr *Move_Zero_Extend_Instruction_Asm) fixInvalidInstr() []Instruction_Asm {
	_, dstIsReg := instr.dst.(*Register_Operand_Asm)
	_, dstIsStack := instr.dst.(*Stack_Operand_Asm)
	_, dstIsStatic := instr.dst.(*Data_Operand_Asm)

	if dstIsReg {
		mov := Mov_Instruction_Asm{asmTyp: LONGWORD_ASM_TYPE, src: instr.src, dst: instr.dst}
		return []Instruction_Asm{&mov}
	} else if dstIsStack || dstIsStatic {
		mov1 := Mov_Instruction_Asm{asmTyp: LONGWORD_ASM_TYPE, src: instr.src, dst: &Register_Operand_Asm{R11_REGISTER_ASM}}
		mov2 := Mov_Instruction_Asm{asmTyp: QUADWORD_ASM_TYPE, src: &Register_Operand_Asm{R11_REGISTER_ASM}, dst: instr.dst}
		return []Instruction_Asm{&mov1, &mov2}
	}

	return []Instruction_Asm{instr}
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Cvttsd2si_Double_To_Int_Instruction_Asm) fixInvalidInstr() []Instruction_Asm {
	_, dstIsStack := instr.dst.(*Stack_Operand_Asm)
	_, dstIsStatic := instr.dst.(*Data_Operand_Asm)

	// the dst of Cvttsd2si must be a register
	if dstIsStack || dstIsStatic {
		r11 := Register_Operand_Asm{R11_REGISTER_ASM}
		cvt := Cvttsd2si_Double_To_Int_Instruction_Asm{dstAsmType: instr.dstAsmType, src: instr.src, dst: &r11}
		mov := Mov_Instruction_Asm{asmTyp: instr.dstAsmType, src: &r11, dst: instr.dst}
		return []Instruction_Asm{&cvt, &mov}
	}

	return []Instruction_Asm{instr}
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Cvtsi2sd_Int_To_Double_Instruction_Asm) fixInvalidInstr() []Instruction_Asm {
	_, srcIsConstant := instr.src.(*Immediate_Int_Operand_Asm)
	_, dstIsStack := instr.dst.(*Stack_Operand_Asm)
	_, dstIsStatic := instr.dst.(*Data_Operand_Asm)

	// The src can't be a constant, the dst must be a register.
	// It's not very optimized to change both when only one or the other occurs, but it simplifies the logic here.
	if srcIsConstant || dstIsStack || dstIsStatic {
		r10 := Register_Operand_Asm{R10_REGISTER_ASM}
		xmm15 := Register_Operand_Asm{XMM15_REGISTER_ASM}
		mov1 := Mov_Instruction_Asm{asmTyp: instr.srcAsmType, src: instr.src, dst: &r10}
		cvt := Cvtsi2sd_Int_To_Double_Instruction_Asm{srcAsmType: instr.srcAsmType, src: &r10, dst: &xmm15}
		mov2 := Mov_Instruction_Asm{asmTyp: DOUBLE_ASM_TYPE, src: &xmm15, dst: instr.dst}
		return []Instruction_Asm{&mov1, &cvt, &mov2}
	}

	return []Instruction_Asm{instr}
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Binary_Instruction_Asm) fixInvalidInstr() []Instruction_Asm {
	if instr.asmTyp == DOUBLE_ASM_TYPE &&
		((instr.binOp == ADD_OPERATOR_ASM) || (instr.binOp == SUB_OPERATOR_ASM) || (instr.binOp == MULT_OPERATOR_ASM) ||
			(instr.binOp == DIV_DOUBLE_OPERATOR_ASM) || (instr.binOp == XOR_OPERATOR_ASM)) {
		_, dstIsStack := instr.dst.(*Stack_Operand_Asm)
		_, dstIsStatic := instr.dst.(*Data_Operand_Asm)

		if dstIsStack || dstIsStatic {
			// page 337, dst must be a register
			xmm15 := Register_Operand_Asm{XMM15_REGISTER_ASM}
			mov1 := Mov_Instruction_Asm{asmTyp: DOUBLE_ASM_TYPE, src: instr.dst, dst: &xmm15}
			bin := Binary_Instruction_Asm{binOp: instr.binOp, asmTyp: DOUBLE_ASM_TYPE, src: instr.src, dst: &xmm15}
			mov2 := Mov_Instruction_Asm{asmTyp: DOUBLE_ASM_TYPE, src: &xmm15, dst: instr.dst}
			return []Instruction_Asm{&mov1, &bin, &mov2}
		}
	} else if instr.binOp == ADD_OPERATOR_ASM || instr.binOp == SUB_OPERATOR_ASM || instr.binOp == AND_OPERATOR_ASM || instr.binOp == OR_OPERATOR_ASM {
		_, srcIsStack := instr.src.(*Stack_Operand_Asm)
		_, dstIsStack := instr.dst.(*Stack_Operand_Asm)
		_, srcIsStatic := instr.src.(*Data_Operand_Asm)
		_, dstIsStatic := instr.dst.(*Data_Operand_Asm)
		srcIsBigImm := opIsBigImm(instr.src)
		isQuadInstr := instr.asmTyp == QUADWORD_ASM_TYPE

		// page 268 of the book, binary instructions can't use large immediate values, so put them in a register first
		if ((srcIsStack || srcIsStatic) && (dstIsStack || dstIsStatic)) || (isQuadInstr && srcIsBigImm) {
			intermediateOperand := Register_Operand_Asm{R10_REGISTER_ASM}
			firstInstr := Mov_Instruction_Asm{asmTyp: instr.asmTyp, src: instr.src, dst: &intermediateOperand}
			secondInstr := Binary_Instruction_Asm{binOp: instr.binOp, asmTyp: instr.asmTyp, src: &intermediateOperand, dst: instr.dst}
			return []Instruction_Asm{&firstInstr, &secondInstr}
		}
	} else if instr.binOp == MULT_OPERATOR_ASM {
		_, dstIsStack := instr.dst.(*Stack_Operand_Asm)
		_, dstIsStatic := instr.dst.(*Data_Operand_Asm)
		srcIsBigImm := opIsBigImm(instr.src)
		isQuadInstr := instr.asmTyp == QUADWORD_ASM_TYPE

		// Moving the original dst to r11 is only when the dst is memory which imul does not allow.
		// Moving the original src to r10 is only when the src is a large immediate value, which imul does not allow (page 268 of book).
		// It's not very optimized to do both changes when only one is necessary, but this simplifies the logic here.
		if dstIsStack || dstIsStatic || (isQuadInstr && srcIsBigImm) {
			r10 := Register_Operand_Asm{R10_REGISTER_ASM}
			r11 := Register_Operand_Asm{R11_REGISTER_ASM}
			firstInstr := Mov_Instruction_Asm{asmTyp: instr.asmTyp, src: instr.src, dst: &r10}
			secInstr := Mov_Instruction_Asm{asmTyp: instr.asmTyp, src: instr.dst, dst: &r11}
			thirdInstr := Binary_Instruction_Asm{binOp: instr.binOp, asmTyp: instr.asmTyp, src: &r10, dst: &r11}
			fourthInstr := Mov_Instruction_Asm{asmTyp: instr.asmTyp, src: &r11, dst: instr.dst}
			return []Instruction_Asm{&firstInstr, &secInstr, &thirdInstr, &fourthInstr}
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

func (instr *Divide_Instruction_Asm) fixInvalidInstr() []Instruction_Asm {
	_, isConstant := instr.divisor.(*Immediate_Int_Operand_Asm)

	if isConstant {
		firstInstr := Mov_Instruction_Asm{asmTyp: instr.asmTyp, src: instr.divisor, dst: &Register_Operand_Asm{R10_REGISTER_ASM}}
		secondInstr := Divide_Instruction_Asm{asmTyp: instr.asmTyp, divisor: &Register_Operand_Asm{R10_REGISTER_ASM}}
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
	op1IsBigImm := opIsBigImm(instr.op1)
	op2IsBimImm := opIsBigImm(instr.op2)
	isQuadInstr := instr.asmTyp == QUADWORD_ASM_TYPE
	isDoubleInstr := instr.asmTyp == DOUBLE_ASM_TYPE

	if isDoubleInstr && (op2IsStack || op2IsStatic || op2IsConstant) {
		xmm15 := Register_Operand_Asm{XMM15_REGISTER_ASM}
		mov := Mov_Instruction_Asm{asmTyp: DOUBLE_ASM_TYPE, src: instr.op2, dst: &xmm15}
		cmp := Compare_Instruction_Asm{asmTyp: DOUBLE_ASM_TYPE, op1: instr.op1, op2: &xmm15}
		return []Instruction_Asm{&mov, &cmp}
	} else if isQuadInstr && (op1IsBigImm || op2IsBimImm) {
		// page 268 of the book, can't use immediate values that are too big
		r10 := Register_Operand_Asm{R10_REGISTER_ASM}
		r11 := Register_Operand_Asm{R11_REGISTER_ASM}
		firstInstr := Mov_Instruction_Asm{asmTyp: instr.asmTyp, src: instr.op1, dst: &r10}
		secInstr := Mov_Instruction_Asm{asmTyp: instr.asmTyp, src: instr.op2, dst: &r11}
		thirdInstr := Compare_Instruction_Asm{asmTyp: instr.asmTyp, op1: &r10, op2: &r11}
		return []Instruction_Asm{&firstInstr, &secInstr, &thirdInstr}
	} else if (op1IsStack || op1IsStatic) && (op2IsStack || op2IsStatic) {
		r10 := Register_Operand_Asm{R10_REGISTER_ASM}
		firstInstr := Mov_Instruction_Asm{asmTyp: instr.asmTyp, src: instr.op1, dst: &r10}
		secondInstr := Compare_Instruction_Asm{asmTyp: instr.asmTyp, op1: &r10, op2: instr.op2}
		return []Instruction_Asm{&firstInstr, &secondInstr}
	} else if op2IsConstant {
		r11 := Register_Operand_Asm{R11_REGISTER_ASM}
		firstInstr := Mov_Instruction_Asm{asmTyp: instr.asmTyp, src: instr.op2, dst: &r11}
		secondInstr := Compare_Instruction_Asm{asmTyp: instr.asmTyp, op1: instr.op1, op2: &r11}
		return []Instruction_Asm{&firstInstr, &secondInstr}
	}

	return []Instruction_Asm{instr}
}

/////////////////////////////////////////////////////////////////////////////////

func (instr *Push_Instruction_Asm) fixInvalidInstr() []Instruction_Asm {
	// page 268 of the book, can't push immediate values that are too big, need to put the values in a register first
	if opIsBigImm(instr.op) {
		r10 := Register_Operand_Asm{R10_REGISTER_ASM}
		firstInstr := Mov_Instruction_Asm{asmTyp: QUADWORD_ASM_TYPE, src: instr.op, dst: &r10}
		secInstr := Push_Instruction_Asm{&r10}
		return []Instruction_Asm{&firstInstr, &secInstr}
	}
	return []Instruction_Asm{instr}
}
