package main

//###############################################################################
//###############################################################################
//###############################################################################

func doLoopLabeling(ast Program) Program {
	for index, _ := range ast.decls {
		fnDecl, isFunc := ast.decls[index].(*Function_Declaration)
		if isFunc {
			ast.decls[index] = labelFunction(*fnDecl)
		}
	}

	return ast
}

/////////////////////////////////////////////////////////////////////////////////

func labelFunction(fn Function_Declaration) Declaration {
	if fn.body == nil {
		return &fn
	}

	tempBody := labelBlock(*fn.body, "")
	fn.body = &tempBody
	return &fn
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
			fail("Semantic error: break statement outside of loop.")
		}
		convertedSt.label = currentLabel
		return convertedSt
	case *Continue_Statement:
		if currentLabel == "" {
			fail("Semantic error: continue statement outside of loop.")
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
