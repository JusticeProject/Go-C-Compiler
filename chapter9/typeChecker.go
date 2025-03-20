package main

import (
	"fmt"
	"os"
)

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

type Symbol struct {
	typ     Data_Type
	defined bool
}

var symbolTable = make(map[string]Symbol)

//###############################################################################
//###############################################################################
//###############################################################################

func doTypeChecking(ast Program) Program {
	for index, _ := range ast.functions {
		typeCheckFuncDecl(ast.functions[index])
	}

	return ast
}

/////////////////////////////////////////////////////////////////////////////////

func typeCheckDeclaration(decl Declaration) {
	switch convertedDecl := decl.(type) {
	case *Function_Declaration:
		typeCheckFuncDecl(*convertedDecl)
	case *Variable_Declaration:
		typeCheckVarDecl(*convertedDecl)
	}
}

/////////////////////////////////////////////////////////////////////////////////

func typeCheckFuncDecl(decl Function_Declaration) {
	newTyp := Function_Type{paramCount: len(decl.params)}
	hasBody := (decl.body != nil)
	alreadyDefined := false

	oldDecl, inSymbolTable := symbolTable[decl.name]
	if inSymbolTable {
		if !oldDecl.typ.isEqual(&newTyp) {
			fmt.Println("Incompatible function declarations:", decl.name, newTyp.paramCount, "params")
			fmt.Println("does not match previous declaration of", decl.name)
			os.Exit(1)
		}
		alreadyDefined = oldDecl.defined
		if alreadyDefined && hasBody {
			fmt.Println("Function", decl.name, "has two definitions.")
			os.Exit(1)
		}
	}

	symbolTable[decl.name] = Symbol{typ: &newTyp, defined: (alreadyDefined || hasBody)}

	if hasBody {
		for _, param := range decl.params {
			symbolTable[param] = Symbol{typ: &Int_Type{}}
		}
		typeCheckBlock(*decl.body)
	}
}

/////////////////////////////////////////////////////////////////////////////////

func typeCheckVarDecl(decl Variable_Declaration) {
	// every variable should have a unique name at this point, so it won't conflict with any existing entry
	// TODO: need to handle other data types
	symbolTable[decl.name] = Symbol{typ: &Int_Type{}}

	if decl.initializer != nil {
		typeCheckExpression(decl.initializer)
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
		typeCheckDeclaration(convertedItem.decl)
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
		typeCheckVarDecl(convertedInit.decl)
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
			fmt.Println("Function name", convertedExp.name, "used as variable")
			os.Exit(1)
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
			fmt.Println("Function call to", convertedExp.functionName, "does not match any known function declaration.")
			os.Exit(1)
		}
		for _, arg := range convertedExp.args {
			typeCheckExpression(arg)
		}
	}
}

/////////////////////////////////////////////////////////////////////////////////
