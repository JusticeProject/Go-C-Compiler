package main

import (
	"fmt"
	"os"
)

/////////////////////////////////////////////////////////////////////////////////

func doSemanticAnalysis(ast Program) Program {
	ast = doVariableResolution(ast)
	return ast
}

/////////////////////////////////////////////////////////////////////////////////

func doVariableResolution(ast Program) Program {
	// key = user-defined variable name
	// value = globally unique name
	// TODO: how to handle multiple functions, does the map get reset?
	variableMap := make(map[string]string)

	for index, _ := range ast.fn.body {
		existingBlock := ast.fn.body[index]
		newBlock := resolveBlock(existingBlock, &variableMap)
		ast.fn.body[index] = newBlock
	}

	return ast
}

/////////////////////////////////////////////////////////////////////////////////

func resolveBlock(existingBlock Block_Item, variableMap *map[string]string) Block_Item {
	switch convertedBlock := existingBlock.(type) {
	case *Block_Statement:
		newStatement := resolveStatement(convertedBlock.st, variableMap)
		return &Block_Statement{newStatement}
	case *Block_Declaration:
		newDecl := resolveDeclaration(convertedBlock.decl, variableMap)
		return &Block_Declaration{newDecl}
	default:
		fmt.Println("unknown Block_Item when resolving variables")
		os.Exit(1)
	}

	return nil
}

/////////////////////////////////////////////////////////////////////////////////

func resolveDeclaration(decl Declaration, variableMap *map[string]string) Declaration {
	_, nameExists := (*variableMap)[decl.name]

	if nameExists {
		fmt.Println("Semantic error. Variable", decl.name, "already exists.")
		os.Exit(1)
	}

	uniqueName := makeTempVarName(decl.name)
	(*variableMap)[decl.name] = uniqueName

	var init Expression = nil
	if decl.initializer != nil {
		init = resolveExpression(decl.initializer, variableMap)
	}

	return Declaration{name: uniqueName, initializer: init}
}

/////////////////////////////////////////////////////////////////////////////////

func resolveStatement(st Statement, variableMap *map[string]string) Statement {
	switch convertedSt := st.(type) {
	case *Return_Statement:
		newExp := resolveExpression(convertedSt.exp, variableMap)
		return &Return_Statement{newExp}
	case *Expression_Statement:
		newExp := resolveExpression(convertedSt.exp, variableMap)
		return &Expression_Statement{newExp}
	case *Null_Statement:
		return st
	default:
		fmt.Println("unknown Statement type when resolving variables")
		os.Exit(1)
	}

	return nil
}

/////////////////////////////////////////////////////////////////////////////////

func resolveExpression(exp Expression, variableMap *map[string]string) Expression {
	switch convertedExp := exp.(type) {
	case *Constant_Int_Expression:
		return exp
	case *Variable_Expression:
		uniqueName, varExists := (*variableMap)[convertedExp.name]
		if varExists {
			return &Variable_Expression{uniqueName}
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
	default:
		fmt.Println("unknown Expression type when resolving variables")
		os.Exit(1)
	}

	return nil
}
