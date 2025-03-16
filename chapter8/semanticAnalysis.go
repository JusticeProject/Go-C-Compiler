package main

import (
	"fmt"
	"os"
)

/////////////////////////////////////////////////////////////////////////////////

type Variable_Info struct {
	uniqueName       string
	fromCurrentBlock bool
}

/////////////////////////////////////////////////////////////////////////////////

func copyVariableMap(input map[string]Variable_Info) map[string]Variable_Info {
	output := make(map[string]Variable_Info)

	for key, value := range input {
		output[key] = Variable_Info{uniqueName: value.uniqueName, fromCurrentBlock: false}
	}

	return output
}

//###############################################################################
//###############################################################################
//###############################################################################

func doSemanticAnalysis(ast Program) Program {
	ast = doVariableResolution(ast)
	ast = doLoopLabeling(ast)
	return ast
}

//###############################################################################
//###############################################################################
//###############################################################################

func doVariableResolution(ast Program) Program {
	// key = user-defined variable name
	// value = struct containing globally unique name and bool flag indicating whether it was declared in current scope.
	// maps in Go are passed by reference to a function, so you don't need to pass a map by pointer.
	// TODO: how to handle multiple functions, does the map get reset?
	variableMap := make(map[string]Variable_Info)
	ast.fn = resolveFunction(ast.fn, variableMap)

	return ast
}

/////////////////////////////////////////////////////////////////////////////////

func resolveFunction(existingFunc Function, variableMap map[string]Variable_Info) Function {
	// TODO: need to check if function name has already been used
	newFunc := Function{name: existingFunc.name}
	newFunc.body = resolveBlock(existingFunc.body, variableMap)
	return newFunc
}

/////////////////////////////////////////////////////////////////////////////////

func resolveBlock(existingBlock Block, variableMap map[string]Variable_Info) Block {
	// keep the existing Block structure but just swap out the Block_Item at each index
	for index, _ := range existingBlock.items {
		existingItem := existingBlock.items[index]
		newItem := resolveBlockItem(existingItem, variableMap)
		existingBlock.items[index] = newItem
	}

	return existingBlock
}

/////////////////////////////////////////////////////////////////////////////////

func resolveBlockItem(existingItem Block_Item, variableMap map[string]Variable_Info) Block_Item {
	switch convertedItem := existingItem.(type) {
	case *Block_Statement:
		newStatement := resolveStatement(convertedItem.st, variableMap)
		return &Block_Statement{newStatement}
	case *Block_Declaration:
		newDecl := resolveDeclaration(convertedItem.decl, variableMap)
		return &Block_Declaration{newDecl}
	default:
		fmt.Println("unknown Block_Item when resolving variables")
		os.Exit(1)
	}

	return nil
}

/////////////////////////////////////////////////////////////////////////////////

func resolveDeclaration(decl Declaration, variableMap map[string]Variable_Info) Declaration {
	varInfo, nameExists := variableMap[decl.name]

	if nameExists && varInfo.fromCurrentBlock {
		fmt.Println("Semantic error. Variable", decl.name, "declared more than once in same scope.")
		os.Exit(1)
	}

	uniqueName := makeTempVarName(decl.name)
	variableMap[decl.name] = Variable_Info{uniqueName: uniqueName, fromCurrentBlock: true}

	var init Expression = nil
	if decl.initializer != nil {
		init = resolveExpression(decl.initializer, variableMap)
	}

	return Declaration{name: uniqueName, initializer: init}
}

/////////////////////////////////////////////////////////////////////////////////

func resolveForInit(fi For_Initial_Clause, variableMap map[string]Variable_Info) For_Initial_Clause {
	switch convertedInit := fi.(type) {
	case *For_Initial_Declaration:
		newDecl := resolveDeclaration(convertedInit.decl, variableMap)
		return &For_Initial_Declaration{decl: newDecl}
	case *For_Initial_Expression:
		newExp := resolveExpression(convertedInit.exp, variableMap)
		return &For_Initial_Expression{exp: newExp}
	default:
		fmt.Println("unknown For_Initial_Clause when resolving variables.")
		os.Exit(1)
	}

	return nil
}

/////////////////////////////////////////////////////////////////////////////////

func resolveStatement(st Statement, variableMap map[string]Variable_Info) Statement {
	if st == nil {
		return nil
	}

	switch convertedSt := st.(type) {
	case *Return_Statement:
		newExp := resolveExpression(convertedSt.exp, variableMap)
		return &Return_Statement{newExp}
	case *Expression_Statement:
		newExp := resolveExpression(convertedSt.exp, variableMap)
		return &Expression_Statement{newExp}
	case *If_Statement:
		newCond := resolveExpression(convertedSt.condition, variableMap)
		newThen := resolveStatement(convertedSt.thenSt, variableMap)
		newElse := resolveStatement(convertedSt.elseSt, variableMap)
		return &If_Statement{condition: newCond, thenSt: newThen, elseSt: newElse}
	case *Compound_Statement:
		newVariableMap := copyVariableMap(variableMap)
		newBlock := resolveBlock(convertedSt.block, newVariableMap)
		return &Compound_Statement{block: newBlock}
	case *Break_Statement:
		return st
	case *Continue_Statement:
		return st
	case *While_Statement:
		newCond := resolveExpression(convertedSt.condition, variableMap)
		newBody := resolveStatement(convertedSt.body, variableMap)
		return &While_Statement{condition: newCond, body: newBody}
	case *Do_While_Statement:
		newBody := resolveStatement(convertedSt.body, variableMap)
		newCond := resolveExpression(convertedSt.condition, variableMap)
		return &Do_While_Statement{body: newBody, condition: newCond}
	case *For_Statement:
		newVariableMap := copyVariableMap(variableMap)
		newInit := resolveForInit(convertedSt.initial, newVariableMap)
		newCond := resolveExpression(convertedSt.condition, newVariableMap)
		newPost := resolveExpression(convertedSt.post, newVariableMap)
		newBody := resolveStatement(convertedSt.body, newVariableMap)
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

func resolveExpression(exp Expression, variableMap map[string]Variable_Info) Expression {
	if exp == nil {
		return nil
	}

	switch convertedExp := exp.(type) {
	case *Constant_Int_Expression:
		return exp
	case *Variable_Expression:
		varInfo, varExists := variableMap[convertedExp.name]
		if varExists {
			return &Variable_Expression{varInfo.uniqueName}
		} else {
			fmt.Println("Semantic error. Undeclared variable:", convertedExp.name)
			os.Exit(1)
		}
	case *Unary_Expression:
		newInner := resolveExpression(convertedExp.innerExp, variableMap)
		return &Unary_Expression{unOp: convertedExp.unOp, innerExp: newInner}
	case *Binary_Expression:
		newFirst := resolveExpression(convertedExp.firstExp, variableMap)
		newSecond := resolveExpression(convertedExp.secExp, variableMap)
		return &Binary_Expression{binOp: convertedExp.binOp, firstExp: newFirst, secExp: newSecond}
	case *Assignment_Expression:
		_, isValidLvalue := convertedExp.lvalue.(*Variable_Expression)
		if !isValidLvalue {
			fmt.Println("Semantic error. Invalid lvalue on left side of assignment.")
			os.Exit(1)
		}
		newLvalue := resolveExpression(convertedExp.lvalue, variableMap)
		newRightExp := resolveExpression(convertedExp.rightExp, variableMap)
		return &Assignment_Expression{lvalue: newLvalue, rightExp: newRightExp}
	case *Conditional_Expression:
		newCond := resolveExpression(convertedExp.condition, variableMap)
		newMiddle := resolveExpression(convertedExp.middleExp, variableMap)
		newRight := resolveExpression(convertedExp.rightExp, variableMap)
		return &Conditional_Expression{condition: newCond, middleExp: newMiddle, rightExp: newRight}
	default:
		fmt.Println("unknown Expression type when resolving variables")
		os.Exit(1)
	}

	return nil
}

//###############################################################################
//###############################################################################
//###############################################################################

func doLoopLabeling(ast Program) Program {
	// TODO: need to handle more than one function
	ast.fn = labelFunction(ast.fn)
	return ast
}

/////////////////////////////////////////////////////////////////////////////////

func labelFunction(fn Function) Function {
	fn.body = labelBlock(fn.body, "")
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
