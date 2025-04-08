package main

/////////////////////////////////////////////////////////////////////////////////

type InitializerEnum int

const (
	NO_INITIALIZER InitializerEnum = iota
	TENTATIVE_INIT
	INITIAL_INT
	INITIAL_LONG
	INITIAL_UNSIGNED_INT
	INITIAL_UNSIGNED_LONG
	INITIAL_DOUBLE
)

func dataTypeEnumToInitEnum(input DataTypeEnum) InitializerEnum {
	switch input {
	case INT_TYPE:
		return INITIAL_INT
	case LONG_TYPE:
		return INITIAL_LONG
	case UNSIGNED_INT_TYPE:
		return INITIAL_UNSIGNED_INT
	case UNSIGNED_LONG_TYPE:
		return INITIAL_UNSIGNED_LONG
	case DOUBLE_TYPE:
		return INITIAL_DOUBLE
	}

	fail("Can't convert DataTypeEnum to InitializerEnum")
	return NO_INITIALIZER
}

/////////////////////////////////////////////////////////////////////////////////

type AttributeEnum int

const (
	NONE_ATTRIBUTES AttributeEnum = iota
	FUNCTION_ATTRIBUTES
	STATIC_ATTRIBUTES
	LOCAL_ATTRIBUTES
)

/////////////////////////////////////////////////////////////////////////////////

type Symbol struct {
	dataTyp      Data_Type
	attrs        AttributeEnum
	defined      bool
	global       bool
	initEnum     InitializerEnum
	initialValue string
}

var symbolTable = make(map[string]Symbol)

//###############################################################################
//###############################################################################
//###############################################################################

func setResultType(exp Expression, typ DataTypeEnum) Expression {
	switch convertedExp := exp.(type) {
	case *Constant_Value_Expression:
		convertedExp.resultTyp = typ
		return convertedExp
	case *Variable_Expression:
		convertedExp.resultTyp = typ
		return convertedExp
	case *Cast_Expression:
		convertedExp.resultTyp = typ
		return convertedExp
	case *Unary_Expression:
		convertedExp.resultTyp = typ
		return convertedExp
	case *Binary_Expression:
		convertedExp.resultTyp = typ
		return convertedExp
	case *Assignment_Expression:
		convertedExp.resultTyp = typ
		return convertedExp
	case *Conditional_Expression:
		convertedExp.resultTyp = typ
		return convertedExp
	case *Function_Call_Expression:
		convertedExp.resultTyp = typ
		return convertedExp
	default:
		fail("Unknown Expression in setResultType")
	}
	return nil
}

/////////////////////////////////////////////////////////////////////////////////

func getResultType(exp Expression) DataTypeEnum {
	switch convertedExp := exp.(type) {
	case *Constant_Value_Expression:
		return convertedExp.resultTyp
	case *Variable_Expression:
		return convertedExp.resultTyp
	case *Cast_Expression:
		return convertedExp.resultTyp
	case *Unary_Expression:
		return convertedExp.resultTyp
	case *Binary_Expression:
		return convertedExp.resultTyp
	case *Assignment_Expression:
		return convertedExp.resultTyp
	case *Conditional_Expression:
		return convertedExp.resultTyp
	case *Function_Call_Expression:
		return convertedExp.resultTyp
	default:
		fail("Unknown Expression in getResultType")
	}
	return NONE_TYPE
}

/////////////////////////////////////////////////////////////////////////////////

func convertToType(exp Expression, newTyp DataTypeEnum) Expression {
	if getResultType(exp) == newTyp {
		return exp
	}
	castExp := Cast_Expression{targetType: newTyp, innerExp: exp}
	return setResultType(&castExp, newTyp)
}

/////////////////////////////////////////////////////////////////////////////////

func size(typ DataTypeEnum) int32 {
	return asmTypToAlignment(dataTypeEnumToAssemblyTypeEnum(typ))
}

/////////////////////////////////////////////////////////////////////////////////

func isSigned(typ DataTypeEnum) bool {
	switch typ {
	case INT_TYPE:
		return true
	case LONG_TYPE:
		return true
	case UNSIGNED_INT_TYPE:
		return false
	case UNSIGNED_LONG_TYPE:
		return false
	case DOUBLE_TYPE:
		// TODO:
		return false
	}
	fail("Can't determine signedness")
	return false
}

/////////////////////////////////////////////////////////////////////////////////

func getCommonType(typ1 DataTypeEnum, typ2 DataTypeEnum) DataTypeEnum {
	if typ1 == typ2 {
		return typ1
	}

	if (typ1 == DOUBLE_TYPE) || (typ2 == DOUBLE_TYPE) {
		return DOUBLE_TYPE
	}

	if size(typ1) == size(typ2) {
		if isSigned(typ1) {
			return typ2
		} else {
			return typ1
		}
	}

	if size(typ1) > size(typ2) {
		return typ1
	} else {
		return typ2
	}
}

//###############################################################################
//###############################################################################
//###############################################################################

func doTypeChecking(ast Program) Program {
	for index, _ := range ast.decls {
		ast.decls[index] = typeCheckFileScopeDeclaration(ast.decls[index])
	}

	return ast
}

/////////////////////////////////////////////////////////////////////////////////

func typeCheckFileScopeDeclaration(decl Declaration) Declaration {
	switch convertedDecl := decl.(type) {
	case *Function_Declaration:
		newDecl := typeCheckFuncDecl(*convertedDecl)
		return &newDecl
	case *Variable_Declaration:
		newDecl := typeCheckFileScopeVarDecl(*convertedDecl)
		return &newDecl
	}
	return nil
}

/////////////////////////////////////////////////////////////////////////////////

func typeCheckFuncDecl(decl Function_Declaration) Function_Declaration {
	newTyp := decl.dTyp
	hasBody := (decl.body != nil)
	alreadyDefined := false
	global := (decl.storageClass != STATIC_STORAGE_CLASS)

	oldDecl, inSymbolTable := symbolTable[decl.name]
	if inSymbolTable {
		if !oldDecl.dataTyp.isEqualType(&newTyp) {
			fail("Incompatible function declarations for function", decl.name)
		}
		alreadyDefined = oldDecl.defined
		if alreadyDefined && hasBody {
			fail("Function", decl.name, "has two definitions.")
		}

		if oldDecl.global && decl.storageClass == STATIC_STORAGE_CLASS {
			fail("Static function declaration follows non-static declaration of", decl.name)
		}
		global = oldDecl.global
	}

	symbolTable[decl.name] = Symbol{dataTyp: newTyp, attrs: FUNCTION_ATTRIBUTES, defined: (alreadyDefined || hasBody), global: global}

	if hasBody {
		for index, paramName := range decl.paramNames {
			paramType := decl.dTyp.paramTypes[index].typ
			// every variable should have a unique name at this point, so it won't conflict with any existing entry
			symbolTable[paramName] = Symbol{dataTyp: Data_Type{typ: paramType}}
		}
		*decl.body = typeCheckBlock(*decl.body, decl.name)
	}

	return decl
}

/////////////////////////////////////////////////////////////////////////////////

func typeCheckFileScopeVarDecl(decl Variable_Declaration) Variable_Declaration {
	// every variable should have a unique name at this point, so it won't conflict with any existing entry
	var initEnum InitializerEnum = NO_INITIALIZER
	var initialValue string = ""

	constValExp, isConst := decl.initializer.(*Constant_Value_Expression)
	if isConst {
		initEnum = dataTypeEnumToInitEnum(decl.dTyp.typ)
		decl.initializer = convertToType(decl.initializer, decl.dTyp.typ)
		// TODO: if the constant value is a long that doesn't fit into an int (2147483650L) then
		// strconv.ParseInt(value, 10, 64), then cast int64 to int32 (for example), then back to string
		// I tested this and the assembler will truncate it for me.
		initialValue = constValExp.value
	} else if decl.initializer == nil {
		if decl.storageClass == EXTERN_STORAGE_CLASS {
			initEnum = NO_INITIALIZER
		} else {
			initEnum = TENTATIVE_INIT
		}
	} else {
		fail("Non-constant initializer for variable", decl.name)
	}

	global := (decl.storageClass != STATIC_STORAGE_CLASS)

	oldDecl, alreadyExists := symbolTable[decl.name]
	if alreadyExists {
		if !oldDecl.dataTyp.isEqualType(&decl.dTyp) {
			fail("Data types don't match for variable", decl.name)
		}
		if decl.storageClass == EXTERN_STORAGE_CLASS {
			global = oldDecl.global
		} else if oldDecl.global != global {
			fail("Conflicting variable linkage")
		}

		// TODO: update this when more types are available, and the else if below.
		// We don't want to initialize a variable twice because the two values could be conflicting,
		// so if both decl's initialize then throw an error.
		if (oldDecl.initEnum == INITIAL_INT) || (oldDecl.initEnum == INITIAL_LONG) || (oldDecl.initEnum == INITIAL_UNSIGNED_INT) ||
			(oldDecl.initEnum == INITIAL_UNSIGNED_LONG) || (oldDecl.initEnum == INITIAL_DOUBLE) {
			if initEnum == oldDecl.initEnum {
				fail("Conflicting file scope variable declarations")
			} else {
				initEnum = oldDecl.initEnum
				initialValue = oldDecl.initialValue
			}
		} else if (initEnum != INITIAL_INT) && (initEnum != INITIAL_LONG) && (initEnum != INITIAL_UNSIGNED_INT) &&
			(initEnum != INITIAL_UNSIGNED_LONG) && (initEnum != INITIAL_DOUBLE) && (oldDecl.initEnum == TENTATIVE_INIT) {
			initEnum = TENTATIVE_INIT
		}
	}

	symbolTable[decl.name] = Symbol{dataTyp: decl.dTyp, attrs: STATIC_ATTRIBUTES, global: global,
		initEnum: initEnum, initialValue: initialValue}

	return decl
}

/////////////////////////////////////////////////////////////////////////////////

func typeCheckLocalVarDecl(decl Variable_Declaration) Variable_Declaration {
	// every variable should have a unique name at this point, so it won't conflict with any existing entry
	if decl.storageClass == EXTERN_STORAGE_CLASS {
		if decl.initializer != nil {
			fail("Initializer on local extern variable declaration")
		}
		oldDecl, alreadyExists := symbolTable[decl.name]
		if alreadyExists {
			if !oldDecl.dataTyp.isEqualType(&decl.dTyp) {
				fail("Data types don't match for variable", decl.name)
			}
		} else {
			symbolTable[decl.name] = Symbol{dataTyp: decl.dTyp, attrs: STATIC_ATTRIBUTES, global: true, initEnum: NO_INITIALIZER}
		}
	} else if decl.storageClass == STATIC_STORAGE_CLASS {
		var initEnum InitializerEnum = NO_INITIALIZER
		var initialValue string = ""
		constValExp, isConstVal := decl.initializer.(*Constant_Value_Expression)
		if isConstVal {
			initEnum = dataTypeEnumToInitEnum(decl.dTyp.typ)
			decl.initializer = convertToType(decl.initializer, decl.dTyp.typ)
			// TODO:
			// if the constant value is a long that doesn't fit into an int (2147483650L) then
			// strconv.ParseInt(value, 10, 64), then cast int64 to int32 (for example), then back to string
			initialValue = constValExp.value
		} else if decl.initializer == nil {
			initEnum = dataTypeEnumToInitEnum(decl.dTyp.typ)
			initialValue = "0"
		} else {
			fail("Non-constant initializer on local static variable")
		}
		symbolTable[decl.name] = Symbol{dataTyp: decl.dTyp, attrs: STATIC_ATTRIBUTES, global: false,
			initEnum: initEnum, initialValue: initialValue}
	} else {
		// it's an automatic variable
		symbolTable[decl.name] = Symbol{dataTyp: decl.dTyp, attrs: LOCAL_ATTRIBUTES}
		if decl.initializer != nil {
			decl.initializer = typeCheckExpression(decl.initializer)
			decl.initializer = convertToType(decl.initializer, decl.dTyp.typ)
		}
	}

	return decl
}

/////////////////////////////////////////////////////////////////////////////////

func typeCheckBlock(b Block, funcName string) Block {
	for index, _ := range b.items {
		b.items[index] = typeCheckBlockItem(b.items[index], funcName)
	}
	return b
}

/////////////////////////////////////////////////////////////////////////////////

func typeCheckBlockItem(bi Block_Item, funcName string) Block_Item {
	switch convertedItem := bi.(type) {
	case *Block_Statement:
		convertedItem.st = typeCheckStatement(convertedItem.st, funcName)
		return convertedItem
	case *Block_Declaration:
		decl, isVarDecl := convertedItem.decl.(*Variable_Declaration)
		if isVarDecl {
			newDecl := typeCheckLocalVarDecl(*decl)
			convertedItem.decl = &newDecl
			return convertedItem
		} else {
			funcDecl := convertedItem.decl.(*Function_Declaration)
			newFuncDecl := typeCheckFuncDecl(*funcDecl)
			convertedItem.decl = &newFuncDecl
			return convertedItem
		}
	}
	return nil
}

/////////////////////////////////////////////////////////////////////////////////

func typeCheckStatement(st Statement, funcName string) Statement {
	if st == nil {
		return nil
	}

	switch convertedSt := st.(type) {
	case *Return_Statement:
		convertedSt.exp = typeCheckExpression(convertedSt.exp)
		retType := symbolTable[funcName].dataTyp.returnType.typ
		convertedSt.exp = convertToType(convertedSt.exp, retType)
		return convertedSt
	case *Expression_Statement:
		convertedSt.exp = typeCheckExpression(convertedSt.exp)
		return convertedSt
	case *If_Statement:
		convertedSt.condition = typeCheckExpression(convertedSt.condition)
		convertedSt.thenSt = typeCheckStatement(convertedSt.thenSt, funcName)
		if convertedSt.elseSt != nil {
			convertedSt.elseSt = typeCheckStatement(convertedSt.elseSt, funcName)
		}
		return convertedSt
	case *Compound_Statement:
		convertedSt.block = typeCheckBlock(convertedSt.block, funcName)
		return convertedSt
	case *Break_Statement:
		return st
	case *Continue_Statement:
		return st
	case *While_Statement:
		convertedSt.condition = typeCheckExpression(convertedSt.condition)
		convertedSt.body = typeCheckStatement(convertedSt.body, funcName)
		return convertedSt
	case *Do_While_Statement:
		convertedSt.body = typeCheckStatement(convertedSt.body, funcName)
		convertedSt.condition = typeCheckExpression(convertedSt.condition)
		return convertedSt
	case *For_Statement:
		convertedSt.initial = typeCheckForInitial(convertedSt.initial)
		if convertedSt.condition != nil {
			convertedSt.condition = typeCheckExpression(convertedSt.condition)
		}
		if convertedSt.post != nil {
			convertedSt.post = typeCheckExpression(convertedSt.post)
		}
		convertedSt.body = typeCheckStatement(convertedSt.body, funcName)
		return convertedSt
	case *Null_Statement:
		return st
	}

	fail("Unknown Statement type in typeCheckStatement")
	return nil
}

/////////////////////////////////////////////////////////////////////////////////

func typeCheckForInitial(initial For_Initial_Clause) For_Initial_Clause {
	switch convertedInit := initial.(type) {
	case *For_Initial_Declaration:
		if convertedInit.decl.storageClass != NONE_STORAGE_CLASS {
			fail("For loop initializer can not have storage-class specifier")
		}
		convertedInit.decl = typeCheckLocalVarDecl(convertedInit.decl)
		return convertedInit
	case *For_Initial_Expression:
		convertedInit.exp = typeCheckExpression(convertedInit.exp)
		return convertedInit
	}
	return nil
}

/////////////////////////////////////////////////////////////////////////////////

func typeCheckExpression(exp Expression) Expression {
	if exp == nil {
		return nil
	}

	switch convertedExp := exp.(type) {
	case *Constant_Value_Expression:
		return setResultType(convertedExp, convertedExp.typ)
	case *Variable_Expression:
		typ := symbolTable[convertedExp.name].dataTyp.typ
		if typ == FUNCTION_TYPE {
			fail("Function name", convertedExp.name, "used as variable")
		}
		return setResultType(convertedExp, typ)
	case *Cast_Expression:
		newInner := typeCheckExpression(convertedExp.innerExp)
		newCast := Cast_Expression{targetType: convertedExp.targetType, innerExp: newInner}
		return setResultType(&newCast, convertedExp.targetType)
	case *Unary_Expression:
		newInner := typeCheckExpression(convertedExp.innerExp)
		if (convertedExp.unOp == COMPLEMENT_OPERATOR) && (getResultType(newInner) == DOUBLE_TYPE) {
			fail("Can't take the bitwise complement of a double")
		}
		newUnary := Unary_Expression{unOp: convertedExp.unOp, innerExp: newInner}
		if convertedExp.unOp == NOT_OPERATOR {
			return setResultType(&newUnary, INT_TYPE)
		} else {
			return setResultType(&newUnary, getResultType(newInner))
		}
	case *Binary_Expression:
		newFirstExp := typeCheckExpression(convertedExp.firstExp)
		newSecExp := typeCheckExpression(convertedExp.secExp)
		if convertedExp.binOp == REMAINDER_OPERATOR {
			if (getResultType(newFirstExp) == DOUBLE_TYPE) || (getResultType(newSecExp) == DOUBLE_TYPE) {
				fail("Can't take the remainder using doubles")
			}
		}
		if (convertedExp.binOp == AND_OPERATOR) || (convertedExp.binOp == OR_OPERATOR) {
			newBinExp := Binary_Expression{binOp: convertedExp.binOp, firstExp: newFirstExp, secExp: newSecExp}
			return setResultType(&newBinExp, INT_TYPE)
		}
		typ1 := getResultType(newFirstExp)
		typ2 := getResultType(newSecExp)
		commonTyp := getCommonType(typ1, typ2)
		newFirstExp = convertToType(newFirstExp, commonTyp)
		newSecExp = convertToType(newSecExp, commonTyp)
		newBinExp := Binary_Expression{binOp: convertedExp.binOp, firstExp: newFirstExp, secExp: newSecExp}
		if (convertedExp.binOp == ADD_OPERATOR) || (convertedExp.binOp == SUBTRACT_OPERATOR) || (convertedExp.binOp == MULTIPLY_OPERATOR) ||
			(convertedExp.binOp == DIVIDE_OPERATOR) || (convertedExp.binOp == REMAINDER_OPERATOR) {
			return setResultType(&newBinExp, commonTyp)
		} else {
			// comparisons (less than, equal, etc. have a type of int)
			return setResultType(&newBinExp, INT_TYPE)
		}
	case *Assignment_Expression:
		newLvalue := typeCheckExpression(convertedExp.lvalue)
		newRightExp := typeCheckExpression(convertedExp.rightExp)
		leftTyp := getResultType(newLvalue)
		newRightExp = convertToType(newRightExp, leftTyp)
		assignExp := Assignment_Expression{lvalue: newLvalue, rightExp: newRightExp}
		return setResultType(&assignExp, leftTyp)
	case *Conditional_Expression:
		newMiddle := typeCheckExpression(convertedExp.middleExp)
		newRight := typeCheckExpression(convertedExp.rightExp)
		middleTyp := getResultType(newMiddle)
		rightTyp := getResultType(newRight)
		commonTyp := getCommonType(middleTyp, rightTyp)
		newMiddle = convertToType(newMiddle, commonTyp)
		newRight = convertToType(newRight, commonTyp)
		newCond := typeCheckExpression(convertedExp.condition)
		newExp := Conditional_Expression{condition: newCond, middleExp: newMiddle, rightExp: newRight}
		return setResultType(&newExp, commonTyp)
	case *Function_Call_Expression:
		existingSym, inTable := symbolTable[convertedExp.functionName]

		if !inTable {
			fail("Calling a function that's not in the symbol table:", convertedExp.functionName)
		}

		existingTyp := existingSym.dataTyp
		if existingTyp.typ != FUNCTION_TYPE {
			fail("Variable used as function name:", convertedExp.functionName)
		}

		if len(existingTyp.paramTypes) != len(convertedExp.args) {
			fail("Function called with the wrong number of arguments:", convertedExp.functionName)
		}

		newArgs := []Expression{}
		for index, _ := range convertedExp.args {
			newArg := typeCheckExpression(convertedExp.args[index])
			newArg = convertToType(newArg, existingTyp.paramTypes[index].typ)
			newArgs = append(newArgs, newArg)
		}

		callExp := Function_Call_Expression{functionName: convertedExp.functionName, args: newArgs}
		return setResultType(&callExp, existingTyp.returnType.typ)
	}

	fail("Unknown Expression type in typeCheckExpression")
	return nil
}
