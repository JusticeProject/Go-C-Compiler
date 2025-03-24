package main

/////////////////////////////////////////////////////////////////////////////////

type InitializerEnum int

const (
	NO_INITIALIZER InitializerEnum = iota
	INITIAL_INT
	TENTATIVE_INIT
)

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
	initializer  InitializerEnum
	initialValue string
}

var symbolTable = make(map[string]Symbol)

//###############################################################################
//###############################################################################
//###############################################################################

func doTypeChecking(ast Program) Program {
	for index, _ := range ast.decls {
		typeCheckFileScopeDeclaration(ast.decls[index])
	}

	return ast
}

/////////////////////////////////////////////////////////////////////////////////

func typeCheckFileScopeDeclaration(decl Declaration) {
	switch convertedDecl := decl.(type) {
	case *Function_Declaration:
		typeCheckFuncDecl(*convertedDecl)
	case *Variable_Declaration:
		typeCheckFileScopeVarDecl(*convertedDecl)
	}
}

/////////////////////////////////////////////////////////////////////////////////

func typeCheckFuncDecl(decl Function_Declaration) {
	// TODO: add return type to newTyp
	newTyp := Data_Type{typ: FUNCTION_TYPE, paramCount: len(decl.params)}
	hasBody := (decl.body != nil)
	alreadyDefined := false
	global := (decl.storageClass != STATIC_STORAGE_CLASS)

	oldDecl, inSymbolTable := symbolTable[decl.name]
	if inSymbolTable {
		if !oldDecl.dataTyp.isEqualType(newTyp) {
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
		for _, param := range decl.params {
			// every variable should have a unique name at this point, so it won't conflict with any existing entry
			symbolTable[param] = Symbol{dataTyp: Data_Type{typ: INT_TYPE}}
		}
		typeCheckBlock(*decl.body)
	}
}

/////////////////////////////////////////////////////////////////////////////////

func typeCheckFileScopeVarDecl(decl Variable_Declaration) {
	// every variable should have a unique name at this point, so it won't conflict with any existing entry
	// TODO: need to handle other data types, could pass the initializer to a function and get a value back
	var initializer InitializerEnum = NO_INITIALIZER
	var initialValue string = ""

	constIntExp, isConst := decl.initializer.(*Constant_Int_Expression)
	if isConst {
		initializer = INITIAL_INT
		initialValue = constIntExp.intValue
	} else if decl.initializer == nil {
		if decl.storageClass == EXTERN_STORAGE_CLASS {
			initializer = NO_INITIALIZER
		} else {
			initializer = TENTATIVE_INIT
		}
	} else {
		fail("Non-constant initializer for variable", decl.name)
	}

	global := (decl.storageClass != STATIC_STORAGE_CLASS)

	oldDecl, alreadyExists := symbolTable[decl.name]
	if alreadyExists {
		if oldDecl.dataTyp.typ == FUNCTION_TYPE {
			fail("Function redeclared as variable", decl.name)
		}
		if decl.storageClass == EXTERN_STORAGE_CLASS {
			global = oldDecl.global
		} else if oldDecl.global != global {
			fail("Conflicting variable linkage")
		}

		if oldDecl.initializer == INITIAL_INT {
			if initializer == INITIAL_INT {
				fail("Conflicting file scope variable declarations")
			} else {
				initializer = oldDecl.initializer
				initialValue = oldDecl.initialValue
			}
		} else if (initializer != INITIAL_INT) && (oldDecl.initializer == TENTATIVE_INIT) {
			initializer = TENTATIVE_INIT
		}
	}

	symbolTable[decl.name] = Symbol{dataTyp: Data_Type{typ: INT_TYPE}, attrs: STATIC_ATTRIBUTES, global: global,
		initializer: initializer, initialValue: initialValue}
}

/////////////////////////////////////////////////////////////////////////////////

func typeCheckLocalVarDecl(decl Variable_Declaration) {
	// every variable should have a unique name at this point, so it won't conflict with any existing entry
	// TODO: need to handle other data types
	if decl.storageClass == EXTERN_STORAGE_CLASS {
		if decl.initializer != nil {
			fail("Initializer on local extern variable declaration")
		}
		oldDecl, alreadyExists := symbolTable[decl.name]
		if alreadyExists {
			if oldDecl.dataTyp.typ == FUNCTION_TYPE {
				fail("Function redeclared as variable")
			}
		} else {
			symbolTable[decl.name] = Symbol{dataTyp: Data_Type{typ: INT_TYPE}, attrs: STATIC_ATTRIBUTES, global: true, initializer: NO_INITIALIZER}
		}
	} else if decl.storageClass == STATIC_STORAGE_CLASS {
		// TODO: need to handle other data types
		var initializer InitializerEnum = NO_INITIALIZER
		var initialValue string = ""
		constIntExp, isConstInt := decl.initializer.(*Constant_Int_Expression)
		if isConstInt {
			initializer = INITIAL_INT
			initialValue = constIntExp.intValue
		} else if decl.initializer == nil {
			initializer = INITIAL_INT
			initialValue = "0"
		} else {
			fail("Non-constant initializer on local static variable")
		}
		symbolTable[decl.name] = Symbol{dataTyp: Data_Type{typ: INT_TYPE}, attrs: STATIC_ATTRIBUTES, global: false,
			initializer: initializer, initialValue: initialValue}
	} else {
		symbolTable[decl.name] = Symbol{dataTyp: Data_Type{typ: INT_TYPE}, attrs: LOCAL_ATTRIBUTES}
		if decl.initializer != nil {
			typeCheckExpression(decl.initializer)
		}
	}
}

/////////////////////////////////////////////////////////////////////////////////

func typeCheckBlock(b Block) {
	for _, item := range b.items {
		typeCheckBlockItem(item)
	}
}

/////////////////////////////////////////////////////////////////////////////////

func typeCheckBlockItem(bi Block_Item) {
	switch convertedItem := bi.(type) {
	case *Block_Statement:
		typeCheckStatement(convertedItem.st)
	case *Block_Declaration:
		decl, isVarDecl := convertedItem.decl.(*Variable_Declaration)
		if isVarDecl {
			typeCheckLocalVarDecl(*decl)
		} else {
			funcDecl := convertedItem.decl.(*Function_Declaration)
			typeCheckFuncDecl(*funcDecl)
		}
	}
}

/////////////////////////////////////////////////////////////////////////////////

func typeCheckStatement(st Statement) {
	switch convertedSt := st.(type) {
	case *Return_Statement:
		typeCheckExpression(convertedSt.exp)
	case *Expression_Statement:
		typeCheckExpression(convertedSt.exp)
	case *If_Statement:
		typeCheckExpression(convertedSt.condition)
		typeCheckStatement(convertedSt.thenSt)
		if convertedSt.elseSt != nil {
			typeCheckStatement(convertedSt.elseSt)
		}
	case *Compound_Statement:
		typeCheckBlock(convertedSt.block)
	case *While_Statement:
		typeCheckExpression(convertedSt.condition)
		typeCheckStatement(convertedSt.body)
	case *Do_While_Statement:
		typeCheckStatement(convertedSt.body)
		typeCheckExpression(convertedSt.condition)
	case *For_Statement:
		typeCheckForInitial(convertedSt.initial)
		if convertedSt.condition != nil {
			typeCheckExpression(convertedSt.condition)
		}
		if convertedSt.post != nil {
			typeCheckExpression(convertedSt.post)
		}
		typeCheckStatement(convertedSt.body)
	}
}

/////////////////////////////////////////////////////////////////////////////////

func typeCheckForInitial(initial For_Initial_Clause) {
	switch convertedInit := initial.(type) {
	case *For_Initial_Declaration:
		if convertedInit.decl.storageClass != NONE_STORAGE_CLASS {
			fail("For loop initializer can not have storage-class specifier")
		}
		typeCheckLocalVarDecl(convertedInit.decl)
	case *For_Initial_Expression:
		typeCheckExpression(convertedInit.exp)
	}
}

/////////////////////////////////////////////////////////////////////////////////

func typeCheckExpression(exp Expression) {
	switch convertedExp := exp.(type) {
	case *Variable_Expression:
		entry := symbolTable[convertedExp.name]
		if entry.dataTyp.typ == FUNCTION_TYPE {
			fail("Function name", convertedExp.name, "used as variable")
		}
	case *Unary_Expression:
		typeCheckExpression(convertedExp.innerExp)
	case *Binary_Expression:
		typeCheckExpression(convertedExp.firstExp)
		typeCheckExpression(convertedExp.secExp)
	case *Assignment_Expression:
		typeCheckExpression(convertedExp.lvalue)
		typeCheckExpression(convertedExp.rightExp)
	case *Conditional_Expression:
		typeCheckExpression(convertedExp.condition)
		typeCheckExpression(convertedExp.middleExp)
		typeCheckExpression(convertedExp.rightExp)
	case *Function_Call_Expression:
		existingSymbol := symbolTable[convertedExp.functionName]
		// TODO: add return type to callType
		callType := Data_Type{typ: FUNCTION_TYPE, paramCount: len(convertedExp.args)}
		if !existingSymbol.dataTyp.isEqualType(callType) {
			fail("Function call to", convertedExp.functionName, "does not match any known function declaration.")
		}
		for _, arg := range convertedExp.args {
			typeCheckExpression(arg)
		}
	}
}

/////////////////////////////////////////////////////////////////////////////////
