package main

/////////////////////////////////////////////////////////////////////////////////

type Data_Type interface {
	isEqual(input Data_Type) bool
}

/////////////////////////////////////////////////////////////////////////////////

type Int_Type struct {
}

func (t *Int_Type) isEqual(input Data_Type) bool {
	_, isIntType := input.(*Int_Type)
	return isIntType
}

// TODO: could maybe switch these to enums in a Data_Type struct with paramCount

/////////////////////////////////////////////////////////////////////////////////

type Function_Type struct {
	paramCount int
}

func (t *Function_Type) isEqual(input Data_Type) bool {
	converted, isFuncType := input.(*Function_Type)
	if isFuncType {
		return (t.paramCount == converted.paramCount)
	}
	return false
}

//###############################################################################
//###############################################################################
//###############################################################################

type Initial_Value interface {
	isConstantValue() bool
	isTentative() bool
}

/////////////////////////////////////////////////////////////////////////////////

type Tentative struct{}

func (t *Tentative) isConstantValue() bool { return false }
func (t *Tentative) isTentative() bool     { return true }

/////////////////////////////////////////////////////////////////////////////////

type Initial_Int struct {
	value int32
}

// TODO: could maybe switch these to enums

func (i *Initial_Int) isConstantValue() bool { return true }
func (i *Initial_Int) isTentative() bool     { return false }

/////////////////////////////////////////////////////////////////////////////////

type No_Initializer struct{}

func (n *No_Initializer) isConstantValue() bool { return false }
func (n *No_Initializer) isTentative() bool     { return false }

//###############################################################################
//###############################################################################
//###############################################################################

type Identifier_Attributes interface {
	isGlobalAttribute() bool
	isConstant() bool
	isTentative() bool
}

/////////////////////////////////////////////////////////////////////////////////

type Function_Attributes struct {
	defined bool
	global  bool
}

func (f *Function_Attributes) isGlobalAttribute() bool { return f.global }
func (f *Function_Attributes) isConstant() bool        { return false }
func (f *Function_Attributes) isTentative() bool       { return false }

/////////////////////////////////////////////////////////////////////////////////

type Static_Attributes struct {
	init   Initial_Value
	global bool
}

func (s *Static_Attributes) isGlobalAttribute() bool { return s.global }
func (s *Static_Attributes) isConstant() bool        { return s.init.isConstantValue() }
func (s *Static_Attributes) isTentative() bool       { return s.init.isTentative() }

/////////////////////////////////////////////////////////////////////////////////

type Local_Attributes struct{}

func (l *Local_Attributes) isGlobalAttribute() bool { return false }
func (l *Local_Attributes) isConstant() bool        { return false }
func (l *Local_Attributes) isTentative() bool       { return false }

//###############################################################################
//###############################################################################
//###############################################################################

type Symbol struct {
	typ   Data_Type
	attrs Identifier_Attributes
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
	newTyp := Function_Type{paramCount: len(decl.params)}
	hasBody := (decl.body != nil)
	alreadyDefined := false
	global := (decl.storageClass != STATIC_STORAGE_CLASS)

	oldDecl, inSymbolTable := symbolTable[decl.name]
	if inSymbolTable {
		if !oldDecl.typ.isEqual(&newTyp) {
			fail("Incompatible function declarations for function", decl.name)
		}
		oldAttrs, _ := oldDecl.attrs.(*Function_Attributes)
		alreadyDefined = oldAttrs.defined
		if alreadyDefined && hasBody {
			fail("Function", decl.name, "has two definitions.")
		}

		if oldAttrs.global && decl.storageClass == STATIC_STORAGE_CLASS {
			fail("Static function declaration follows non-static declaration of", decl.name)
		}
		global = oldAttrs.global
	}

	attrs := Function_Attributes{defined: (alreadyDefined || hasBody), global: global}
	symbolTable[decl.name] = Symbol{typ: &newTyp, attrs: &attrs}

	if hasBody {
		for _, param := range decl.params {
			// every variable should have a unique name at this point, so it won't conflict with any existing entry
			symbolTable[param] = Symbol{typ: &Int_Type{}}
		}
		typeCheckBlock(*decl.body)
	}
}

/////////////////////////////////////////////////////////////////////////////////

func typeCheckFileScopeVarDecl(decl Variable_Declaration) {
	// every variable should have a unique name at this point, so it won't conflict with any existing entry
	// TODO: need to handle other data types, could pass the initializer to a function and get a value back
	var init Initial_Value = nil
	constIntExp, isConst := decl.initializer.(*Constant_Int_Expression)
	if isConst {
		init = &Initial_Int{value: constIntExp.intValue}
	} else if decl.initializer == nil {
		if decl.storageClass == EXTERN_STORAGE_CLASS {
			init = &No_Initializer{}
		} else {
			init = &Tentative{}
		}
	} else {
		fail("Non-constant initializer for variable", decl.name)
	}

	global := (decl.storageClass != STATIC_STORAGE_CLASS)

	oldDecl, alreadyExists := symbolTable[decl.name]
	if alreadyExists {
		_, isFunc := oldDecl.typ.(*Function_Type)
		if isFunc {
			fail("Function redeclared as variable", decl.name)
		}
		if decl.storageClass == EXTERN_STORAGE_CLASS {
			global = oldDecl.attrs.isGlobalAttribute()
		} else if oldDecl.attrs.isGlobalAttribute() != global {
			fail("Conflicting variable linkage")
		}

		if oldDecl.attrs.isConstant() {
			if init.isConstantValue() {
				fail("Conflicting file scope variable declarations")
			} else {
				init = oldDecl.attrs.(*Static_Attributes).init
			}
		} else if !init.isConstantValue() && oldDecl.attrs.isTentative() {
			init = &Tentative{}
		}
	}

	attrs := Static_Attributes{init: init, global: global}
	symbolTable[decl.name] = Symbol{typ: &Int_Type{}, attrs: &attrs}
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
			_, isFunc := oldDecl.typ.(*Function_Type)
			if isFunc {
				fail("Function redeclared as variable")
			}
		} else {
			symbolTable[decl.name] = Symbol{typ: &Int_Type{}, attrs: &Static_Attributes{init: &No_Initializer{}, global: true}}
		}
	} else if decl.storageClass == STATIC_STORAGE_CLASS {
		// TODO: need to handle other data types
		var init Initial_Value
		constIntExp, isConstInt := decl.initializer.(*Constant_Int_Expression)
		if isConstInt {
			init = &Initial_Int{value: constIntExp.intValue}
		} else if decl.initializer == nil {
			init = &Initial_Int{value: 0}
		} else {
			fail("Non-constant initializer on local static variable")
		}
		symbolTable[decl.name] = Symbol{typ: &Int_Type{}, attrs: &Static_Attributes{init: init, global: false}}
	} else {
		symbolTable[decl.name] = Symbol{typ: &Int_Type{}, attrs: &Local_Attributes{}}
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
		_, isFuncType := entry.typ.(*Function_Type)
		if isFuncType {
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
		callType := Function_Type{paramCount: len(convertedExp.args)}
		if !existingSymbol.typ.isEqual(&callType) {
			fail("Function call to", convertedExp.functionName, "does not match any known function declaration.")
		}
		for _, arg := range convertedExp.args {
			typeCheckExpression(arg)
		}
	}
}

/////////////////////////////////////////////////////////////////////////////////
