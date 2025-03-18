package main

import (
	"fmt"
	"os"
)

/////////////////////////////////////////////////////////////////////////////////

type Identifier_Info struct {
	uniqueName       string
	fromCurrentScope bool
	hasLinkage       bool
}

/////////////////////////////////////////////////////////////////////////////////

func copyIdentifierMap(input map[string]Identifier_Info) map[string]Identifier_Info {
	output := make(map[string]Identifier_Info)

	for key, value := range input {
		output[key] = Identifier_Info{uniqueName: value.uniqueName,
			fromCurrentScope: false, hasLinkage: value.hasLinkage}
	}

	return output
}

//###############################################################################
//###############################################################################
//###############################################################################

func doSemanticAnalysis(ast Program) Program {
	ast = doIdentifierResolution(ast)
	ast = doTypeChecking(ast)
	ast = doLoopLabeling(ast)
	return ast
}

//###############################################################################
//###############################################################################
//###############################################################################

func doIdentifierResolution(ast Program) Program {
	// key = user-defined variable name
	// value = struct containing globally unique name and bool flag indicating whether it was declared in current scope.
	// maps in Go are passed by reference to a function, so you don't need to pass a map by pointer.
	identifierMap := make(map[string]Identifier_Info)
	for index, _ := range ast.functions {
		ast.functions[index] = resolveFunctionDeclaration(ast.functions[index], identifierMap)
	}

	return ast
}

/////////////////////////////////////////////////////////////////////////////////

func resolveFunctionDeclaration(decl Function_Declaration, identifierMap map[string]Identifier_Info) Function_Declaration {
	prevEntry, funcExists := identifierMap[decl.name]
	if funcExists {
		if prevEntry.fromCurrentScope && !prevEntry.hasLinkage {
			fmt.Println("Semantic error. Duplicate function declaration:", decl.name)
			os.Exit(1)
		}
	}

	identifierMap[decl.name] = Identifier_Info{uniqueName: decl.name, fromCurrentScope: true, hasLinkage: true}

	// the list of function parameters in a declaration starts a new scope, so we need a copy of the map to track them
	innerMap := copyIdentifierMap(identifierMap)
	newParams := []string{}
	for _, param := range decl.params {
		newParam := resolveParam(param, innerMap)
		newParams = append(newParams, newParam)
	}

	var newBody *Block = nil
	if decl.body != nil {
		tempBody := resolveBlock(*decl.body, innerMap)
		newBody = &tempBody
	}
	return Function_Declaration{name: decl.name, params: newParams, body: newBody}
}

/////////////////////////////////////////////////////////////////////////////////

func resolveParam(param string, identifierMap map[string]Identifier_Info) string {
	idInfo, nameExists := identifierMap[param]

	if nameExists && idInfo.fromCurrentScope {
		fmt.Println("Semantic error. Variable", param, "declared more than once in same scope.")
		os.Exit(1)
	}

	uniqueName := makeTempVarName(param)
	identifierMap[param] = Identifier_Info{uniqueName: uniqueName, fromCurrentScope: true, hasLinkage: false}

	return uniqueName
}

/////////////////////////////////////////////////////////////////////////////////

func resolveVariableDeclaration(decl Variable_Declaration, identifierMap map[string]Identifier_Info) Variable_Declaration {
	idInfo, nameExists := identifierMap[decl.name]

	if nameExists && idInfo.fromCurrentScope {
		fmt.Println("Semantic error. Variable", decl.name, "declared more than once in same scope.")
		os.Exit(1)
	}

	uniqueName := makeTempVarName(decl.name)
	identifierMap[decl.name] = Identifier_Info{uniqueName: uniqueName, fromCurrentScope: true, hasLinkage: false}

	var init Expression = nil
	if decl.initializer != nil {
		init = resolveExpression(decl.initializer, identifierMap)
	}

	return Variable_Declaration{name: uniqueName, initializer: init}
}

/////////////////////////////////////////////////////////////////////////////////

func resolveBlock(existingBlock Block, identifierMap map[string]Identifier_Info) Block {
	// keep the existing Block structure but just swap out the Block_Item at each index
	for index, _ := range existingBlock.items {
		existingItem := existingBlock.items[index]
		newItem := resolveBlockItem(existingItem, identifierMap)
		existingBlock.items[index] = newItem
	}

	return existingBlock
}

/////////////////////////////////////////////////////////////////////////////////

func resolveBlockItem(existingItem Block_Item, identifierMap map[string]Identifier_Info) Block_Item {
	switch convertedItem := existingItem.(type) {
	case *Block_Statement:
		newStatement := resolveStatement(convertedItem.st, identifierMap)
		return &Block_Statement{newStatement}
	case *Block_Declaration:
		decl, isVarDecl := convertedItem.decl.(*Variable_Declaration)
		if isVarDecl {
			newDecl := resolveVariableDeclaration(*decl, identifierMap)
			return &Block_Declaration{&newDecl}
		} else {
			funcDecl := convertedItem.decl.(*Function_Declaration)
			if funcDecl.body != nil {
				fmt.Println("Semantic error. Local function declaration can not have a body:", funcDecl.name)
				os.Exit(1)
			}
			newDecl := resolveFunctionDeclaration(*funcDecl, identifierMap)
			return &Block_Declaration{&newDecl}
		}
	default:
		fmt.Println("unknown Block_Item when resolving variables")
		os.Exit(1)
	}

	return nil
}

/////////////////////////////////////////////////////////////////////////////////

func resolveForInit(fi For_Initial_Clause, identifierMap map[string]Identifier_Info) For_Initial_Clause {
	switch convertedInit := fi.(type) {
	case *For_Initial_Declaration:
		// TODO: change to Declaration then back to Variable_Declaration?
		newDecl := resolveVariableDeclaration(convertedInit.decl, identifierMap)
		return &For_Initial_Declaration{decl: newDecl}
	case *For_Initial_Expression:
		newExp := resolveExpression(convertedInit.exp, identifierMap)
		return &For_Initial_Expression{exp: newExp}
	default:
		fmt.Println("unknown For_Initial_Clause when resolving variables.")
		os.Exit(1)
	}

	return nil
}

/////////////////////////////////////////////////////////////////////////////////

func resolveStatement(st Statement, identifierMap map[string]Identifier_Info) Statement {
	if st == nil {
		return nil
	}

	switch convertedSt := st.(type) {
	case *Return_Statement:
		newExp := resolveExpression(convertedSt.exp, identifierMap)
		return &Return_Statement{newExp}
	case *Expression_Statement:
		newExp := resolveExpression(convertedSt.exp, identifierMap)
		return &Expression_Statement{newExp}
	case *If_Statement:
		newCond := resolveExpression(convertedSt.condition, identifierMap)
		newThen := resolveStatement(convertedSt.thenSt, identifierMap)
		newElse := resolveStatement(convertedSt.elseSt, identifierMap)
		return &If_Statement{condition: newCond, thenSt: newThen, elseSt: newElse}
	case *Compound_Statement:
		newIdentifierMap := copyIdentifierMap(identifierMap)
		newBlock := resolveBlock(convertedSt.block, newIdentifierMap)
		return &Compound_Statement{block: newBlock}
	case *Break_Statement:
		return st
	case *Continue_Statement:
		return st
	case *While_Statement:
		newCond := resolveExpression(convertedSt.condition, identifierMap)
		newBody := resolveStatement(convertedSt.body, identifierMap)
		return &While_Statement{condition: newCond, body: newBody}
	case *Do_While_Statement:
		newBody := resolveStatement(convertedSt.body, identifierMap)
		newCond := resolveExpression(convertedSt.condition, identifierMap)
		return &Do_While_Statement{body: newBody, condition: newCond}
	case *For_Statement:
		newIdentifierMap := copyIdentifierMap(identifierMap)
		newInit := resolveForInit(convertedSt.initial, newIdentifierMap)
		newCond := resolveExpression(convertedSt.condition, newIdentifierMap)
		newPost := resolveExpression(convertedSt.post, newIdentifierMap)
		newBody := resolveStatement(convertedSt.body, newIdentifierMap)
		return &For_Statement{initial: newInit, condition: newCond, post: newPost, body: newBody}
	case *Null_Statement:
		return st
	default:
		fmt.Println("unknown Statement type when resolving variables")
		os.Exit(1)
	}

	return nil
}

/////////////////////////////////////////////////////////////////////////////////

func resolveExpression(exp Expression, identifierMap map[string]Identifier_Info) Expression {
	if exp == nil {
		return nil
	}

	switch convertedExp := exp.(type) {
	case *Constant_Int_Expression:
		return exp
	case *Variable_Expression:
		idInfo, varExists := identifierMap[convertedExp.name]
		if varExists {
			return &Variable_Expression{idInfo.uniqueName}
		} else {
			fmt.Println("Semantic error. Undeclared variable:", convertedExp.name)
			os.Exit(1)
		}
	case *Unary_Expression:
		newInner := resolveExpression(convertedExp.innerExp, identifierMap)
		return &Unary_Expression{unOp: convertedExp.unOp, innerExp: newInner}
	case *Binary_Expression:
		newFirst := resolveExpression(convertedExp.firstExp, identifierMap)
		newSecond := resolveExpression(convertedExp.secExp, identifierMap)
		return &Binary_Expression{binOp: convertedExp.binOp, firstExp: newFirst, secExp: newSecond}
	case *Assignment_Expression:
		_, isValidLvalue := convertedExp.lvalue.(*Variable_Expression)
		if !isValidLvalue {
			fmt.Println("Semantic error. Invalid lvalue on left side of assignment.")
			os.Exit(1)
		}
		newLvalue := resolveExpression(convertedExp.lvalue, identifierMap)
		newRightExp := resolveExpression(convertedExp.rightExp, identifierMap)
		return &Assignment_Expression{lvalue: newLvalue, rightExp: newRightExp}
	case *Conditional_Expression:
		newCond := resolveExpression(convertedExp.condition, identifierMap)
		newMiddle := resolveExpression(convertedExp.middleExp, identifierMap)
		newRight := resolveExpression(convertedExp.rightExp, identifierMap)
		return &Conditional_Expression{condition: newCond, middleExp: newMiddle, rightExp: newRight}
	case *Function_Call_Expression:
		idInfo, nameExists := identifierMap[convertedExp.functionName]
		if nameExists {
			newFuncName := idInfo.uniqueName
			newArgs := []Expression{}
			for _, arg := range convertedExp.args {
				newArg := resolveExpression(arg, identifierMap)
				newArgs = append(newArgs, newArg)
			}
			return &Function_Call_Expression{functionName: newFuncName, args: newArgs}
		} else {
			fmt.Println("Semantic error. Trying to use undeclared function:", convertedExp.functionName)
			os.Exit(1)
		}
	default:
		fmt.Println("unknown Expression type when resolving variables")
		os.Exit(1)
	}

	return nil
}

//###############################################################################
//###############################################################################
//###############################################################################

func doTypeChecking(ast Program) Program {
	return ast
}

//###############################################################################
//###############################################################################
//###############################################################################

func doLoopLabeling(ast Program) Program {
	for index, _ := range ast.functions {
		ast.functions[index] = labelFunction(ast.functions[index])
	}

	return ast
}

/////////////////////////////////////////////////////////////////////////////////

func labelFunction(fn Function_Declaration) Function_Declaration {
	if fn.body == nil {
		return fn
	}

	tempBody := labelBlock(*fn.body, "")
	fn.body = &tempBody
	return fn
}

/////////////////////////////////////////////////////////////////////////////////

func labelBlock(bl Block, currentLabel string) Block {
	newItems := []Block_Item{}

	for _, item := range bl.items {
		newItem := labelBlockItem(item, currentLabel)
		newItems = append(newItems, newItem)
	}

	bl.items = newItems
	return bl
}

/////////////////////////////////////////////////////////////////////////////////

func labelBlockItem(bi Block_Item, currentLabel string) Block_Item {
	switch convertedBi := bi.(type) {
	case *Block_Statement:
		newSt := labelStatement(convertedBi.st, currentLabel)
		return &Block_Statement{st: newSt}
	case *Block_Declaration:
		return bi
	}

	return nil
}

/////////////////////////////////////////////////////////////////////////////////

func labelStatement(st Statement, currentLabel string) Statement {
	switch convertedSt := st.(type) {
	case *If_Statement:
		convertedSt.thenSt = labelStatement(convertedSt.thenSt, currentLabel)
		convertedSt.elseSt = labelStatement(convertedSt.elseSt, currentLabel)
		return convertedSt
	case *Compound_Statement:
		convertedSt.block = labelBlock(convertedSt.block, currentLabel)
		return convertedSt
	case *Break_Statement:
		if currentLabel == "" {
			fmt.Println("Semantic error: break statement outside of loop.")
			os.Exit(1)
		}
		convertedSt.label = currentLabel
		return convertedSt
	case *Continue_Statement:
		if currentLabel == "" {
			fmt.Println("Semantic error: continue statement outside of loop.")
			os.Exit(1)
		}
		convertedSt.label = currentLabel
		return convertedSt
	case *While_Statement:
		newLabel := makeLabelName("whileLoop")
		convertedSt.body = labelStatement(convertedSt.body, newLabel)
		convertedSt.label = newLabel
		return convertedSt
	case *Do_While_Statement:
		newLabel := makeLabelName("doWhileLoop")
		convertedSt.body = labelStatement(convertedSt.body, newLabel)
		convertedSt.label = newLabel
		return convertedSt
	case *For_Statement:
		newLabel := makeLabelName("forLoop")
		convertedSt.body = labelStatement(convertedSt.body, newLabel)
		convertedSt.label = newLabel
		return convertedSt
	default:
		return st
	}
}

/////////////////////////////////////////////////////////////////////////////////
