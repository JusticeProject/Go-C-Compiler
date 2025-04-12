package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

//###############################################################################
//###############################################################################
//###############################################################################

type Program struct {
	decls []Declaration
}

//###############################################################################
//###############################################################################
//###############################################################################

type Declaration interface {
	declToTacky() []Instruction_Tacky
	getPrettyPrintLines() []string
}

type Variable_Declaration struct {
	name         string
	initializer  Expression
	dTyp         Data_Type
	storageClass StorageClassEnum
}

type Function_Declaration struct {
	name         string
	paramNames   []string
	body         *Block
	dTyp         Data_Type
	storageClass StorageClassEnum
}

/////////////////////////////////////////////////////////////////////////////////

type StorageClassEnum int

const (
	NONE_STORAGE_CLASS StorageClassEnum = iota
	STATIC_STORAGE_CLASS
	EXTERN_STORAGE_CLASS
)

func getStorageClass(token TokenEnum) StorageClassEnum {
	switch token {
	case STATIC_KEYWORD_TOKEN:
		return STATIC_STORAGE_CLASS
	case EXTERN_KEYWORD_TOKEN:
		return EXTERN_STORAGE_CLASS
	default:
		return NONE_STORAGE_CLASS
	}
}

//###############################################################################
//###############################################################################
//###############################################################################

type DataTypeEnum int

const (
	NONE_TYPE DataTypeEnum = iota
	INT_TYPE
	LONG_TYPE
	UNSIGNED_INT_TYPE
	UNSIGNED_LONG_TYPE
	DOUBLE_TYPE
	FUNCTION_TYPE
	POINTER_TYPE
)

type Data_Type struct {
	typ DataTypeEnum

	// for FUNCTION_TYPE
	paramTypes []*Data_Type
	returnType *Data_Type

	// for POINTER_TYPE
	refType *Data_Type

	// TODO: if this struct changes, update isEqualType() also
}

func (dt *Data_Type) isEqualType(input *Data_Type) bool {
	if dt == nil && input == nil {
		return true
	} else if dt == nil && input != nil {
		return false
	} else if dt != nil && input == nil {
		return false
	}
	// else both are not nil, so do a further comparison

	if dt.typ != input.typ {
		return false
	}
	if len(dt.paramTypes) != len(input.paramTypes) {
		return false
	}
	for i := 0; i < len(dt.paramTypes); i++ {
		if !dt.paramTypes[i].isEqualType(input.paramTypes[i]) {
			return false
		}
	}
	if !dt.returnType.isEqualType(input.returnType) {
		return false
	}
	return dt.refType.isEqualType(input.refType)
}

//###############################################################################
//###############################################################################
//###############################################################################

type Block struct {
	items []Block_Item
}

//###############################################################################
//###############################################################################
//###############################################################################

type Block_Item interface {
	blockItemToTacky() []Instruction_Tacky
	getPrettyPrintLines() []string
}

type Block_Statement struct {
	st Statement
}

type Block_Declaration struct {
	decl Declaration
}

//###############################################################################
//###############################################################################
//###############################################################################

type For_Initial_Clause interface {
	forInitialToTacky() []Instruction_Tacky
	getPrettyPrintLines() []string
}

type For_Initial_Declaration struct {
	// this should always be a Variable_Declaration, so don't use the Declaration interface
	decl Variable_Declaration
}

type For_Initial_Expression struct {
	exp Expression
}

//###############################################################################
//###############################################################################
//###############################################################################

type Statement interface {
	statementToTacky() []Instruction_Tacky
	getPrettyPrintLines() []string
}

type Return_Statement struct {
	exp Expression
}

type Expression_Statement struct {
	exp Expression
}

type If_Statement struct {
	condition Expression
	thenSt    Statement
	elseSt    Statement
}

type Compound_Statement struct {
	block Block
}

type Break_Statement struct {
	label string
}

type Continue_Statement struct {
	label string
}

type While_Statement struct {
	condition Expression
	body      Statement
	label     string
}

type Do_While_Statement struct {
	body      Statement
	condition Expression
	label     string
}

type For_Statement struct {
	initial   For_Initial_Clause
	condition Expression
	post      Expression
	body      Statement
	label     string
}

// example: while (true) {;}
// the ; is a null statement
type Null_Statement struct {
}

//###############################################################################
//###############################################################################
//###############################################################################

type Expression interface {
	expToTacky(instructions []Instruction_Tacky) (Expression_Result_Tacky, []Instruction_Tacky)
	getPrettyPrintLines() []string
}

type Constant_Value_Expression struct {
	dTyp      Data_Type
	value     string
	resultTyp Data_Type
}

type Variable_Expression struct {
	name      string
	resultTyp Data_Type
}

type Cast_Expression struct {
	targetType Data_Type
	innerExp   Expression
	resultTyp  Data_Type
}

type Unary_Expression struct {
	unOp      UnaryOperatorType
	innerExp  Expression
	resultTyp Data_Type
}

type Binary_Expression struct {
	binOp     BinaryOperatorType
	firstExp  Expression
	secExp    Expression
	resultTyp Data_Type
}

type Assignment_Expression struct {
	lvalue    Expression
	rightExp  Expression
	resultTyp Data_Type
}

// example: a == 3 ? 1 : 2
type Conditional_Expression struct {
	condition Expression
	middleExp Expression
	rightExp  Expression
	resultTyp Data_Type
}

type Function_Call_Expression struct {
	functionName string
	args         []Expression
	resultTyp    Data_Type
}

type Dereference_Expression struct {
	innerExp  Expression
	resultTyp Data_Type
}

type Address_Of_Expression struct {
	innerExp  Expression
	resultTyp Data_Type
}

//###############################################################################
//###############################################################################
//###############################################################################

type Declarator interface {
	// returns the identifier (name), the derived type, and list of param names (if any)
	processDeclarator(baseTyp Data_Type) (string, Data_Type, []string)
}

type Identifier_Declarator struct {
	name string
}

type Pointer_Declarator struct {
	innerDec Declarator
}

type Function_Declarator struct {
	paramInfos []Param_Info
	innerDec   Declarator
}

/////////////////////////////////////////////////////////////////////////////////

type Param_Info struct {
	dTyp Data_Type
	dec  Declarator
}

/////////////////////////////////////////////////////////////////////////////////

type Abstract_Declarator interface {
	processAbstractDeclarator(baseTyp Data_Type) Data_Type
}

type Abstract_Pointer_Declarator struct {
	innerDec Abstract_Declarator
}

type Abstract_Base_Declarator struct {
}

//###############################################################################
//###############################################################################
//###############################################################################

type UnaryOperatorType int

const (
	NONE_UNARY_OPERATOR UnaryOperatorType = iota
	COMPLEMENT_OPERATOR
	NEGATE_OPERATOR
	NOT_OPERATOR
	DEREFERENCE_OPERATOR
	ADDRESS_OF_OPERATOR
)

func getUnaryOperator(token Token) UnaryOperatorType {
	switch token.tokenType {
	case TILDE_TOKEN:
		return COMPLEMENT_OPERATOR
	case HYPHEN_TOKEN:
		return NEGATE_OPERATOR
	case EXCLAMATION_TOKEN:
		return NOT_OPERATOR
	case ASTERISK_TOKEN:
		return DEREFERENCE_OPERATOR
	case AMPERSAND_TOKEN:
		return ADDRESS_OF_OPERATOR
	}

	return NONE_UNARY_OPERATOR
}

func isUnaryOperator(token Token) bool {
	unOp := getUnaryOperator(token)
	if unOp == NONE_UNARY_OPERATOR {
		return false
	} else {
		return true
	}
}

/////////////////////////////////////////////////////////////////////////////////

type BinaryOperatorType int

const (
	NONE_BINARY_OPERATOR BinaryOperatorType = iota
	ADD_OPERATOR
	SUBTRACT_OPERATOR
	MULTIPLY_OPERATOR
	DIVIDE_OPERATOR
	REMAINDER_OPERATOR
	AND_OPERATOR
	OR_OPERATOR
	IS_EQUAL_OPERATOR
	NOT_EQUAL_OPERATOR
	LESS_THAN_OPERATOR
	LESS_OR_EQUAL_OPERATOR
	GREATER_THAN_OPERATOR
	GREATER_OR_EQUAL_OPERATOR
	ASSIGNMENT_OPERATOR
	CONDITIONAL_OPERATOR
)

func getBinaryOperator(token Token) BinaryOperatorType {
	switch token.tokenType {
	case HYPHEN_TOKEN:
		return SUBTRACT_OPERATOR
	case PLUS_TOKEN:
		return ADD_OPERATOR
	case ASTERISK_TOKEN:
		return MULTIPLY_OPERATOR
	case FORWARD_SLASH_TOKEN:
		return DIVIDE_OPERATOR
	case PERCENT_TOKEN:
		return REMAINDER_OPERATOR
	case TWO_AMPERSANDS_TOKEN:
		return AND_OPERATOR
	case TWO_VERTICAL_BARS_TOKEN:
		return OR_OPERATOR
	case TWO_EQUAL_SIGNS_TOKEN:
		return IS_EQUAL_OPERATOR
	case EXCLAMATION_EQUAL_TOKEN:
		return NOT_EQUAL_OPERATOR
	case LESS_THAN_TOKEN:
		return LESS_THAN_OPERATOR
	case LESS_OR_EQUAL_TOKEN:
		return LESS_OR_EQUAL_OPERATOR
	case GREATER_THAN_TOKEN:
		return GREATER_THAN_OPERATOR
	case GREATER_OR_EQUAL_TOKEN:
		return GREATER_OR_EQUAL_OPERATOR
	case EQUAL_TOKEN:
		return ASSIGNMENT_OPERATOR
	case QUESTION_TOKEN:
		return CONDITIONAL_OPERATOR
	}

	return NONE_BINARY_OPERATOR
}

func isBinaryOperator(token Token) bool {
	binOp := getBinaryOperator(token)
	if binOp == NONE_BINARY_OPERATOR {
		return false
	} else {
		return true
	}
}

//###############################################################################
//###############################################################################
//###############################################################################

func doParser(tokens []Token) Program {
	// get the Abstract Syntax Tree of the entire program
	ast, tokens := parseProgram(tokens)

	// if there are any remaining tokens, generate syntax error
	if len(tokens) > 0 {
		fail("Sytnax Error. Tokens remaining after parsing program.")
	}

	// print the ast in a well-formatted way
	prettyPrint(ast)

	return ast
}

/////////////////////////////////////////////////////////////////////////////////

func parseProgram(tokens []Token) (Program, []Token) {
	decls := []Declaration{}

	for {
		var decl Declaration
		decl, tokens = parseDeclaration(tokens)
		if decl == nil {
			break
		} else {
			decls = append(decls, decl)
		}
	}

	pr := Program{decls: decls}
	return pr, tokens
}

/////////////////////////////////////////////////////////////////////////////////

func parseDeclaration(tokens []Token) (Declaration, []Token) {
	var specifiers []TokenEnum
	specifiers, tokens = parseSpecifiers(tokens, true)
	if len(specifiers) == 0 {
		return nil, tokens
	}
	baseType, storageClass := analyzeTypeAndStorageClass(specifiers)
	dec, tokens := parseDeclarator(tokens)
	name, decType, paramNames := dec.processDeclarator(baseType)

	if decType.typ == FUNCTION_TYPE {
		if peekToken(tokens).tokenType == SEMICOLON_TOKEN {
			// it's a function declaration
			_, tokens = expect(SEMICOLON_TOKEN, tokens)
			fn := Function_Declaration{name: name, paramNames: paramNames, body: nil, dTyp: decType, storageClass: storageClass}
			return &fn, tokens
		} else {
			// it's a function definition
			block, tokens := parseBlock(tokens)
			fn := Function_Declaration{name: name, paramNames: paramNames, body: &block, dTyp: decType, storageClass: storageClass}
			return &fn, tokens
		}
	} else {
		// it's a variable declaration
		decl := Variable_Declaration{name: name, dTyp: decType, storageClass: storageClass}

		if peekToken(tokens).tokenType == EQUAL_TOKEN {
			// it has an initializer
			_, tokens = expect(EQUAL_TOKEN, tokens)
			var exp Expression
			exp, tokens = parseExpression(tokens, 0)
			decl.initializer = exp
		}
		_, tokens = expect(SEMICOLON_TOKEN, tokens)
		return &decl, tokens
	}
}

/////////////////////////////////////////////////////////////////////////////////

func parseDeclarator(tokens []Token) (Declarator, []Token) {
	if peekToken(tokens).tokenType == ASTERISK_TOKEN {
		// it's a pointer
		_, tokens = expect(ASTERISK_TOKEN, tokens)
		innerDec, tokens := parseDeclarator(tokens)
		pDec := Pointer_Declarator{innerDec: innerDec}
		return &pDec, tokens
	} else {
		// it's a direct declarator
		dec, tokens := parseDirectDeclarator(tokens)
		return dec, tokens
	}
}

/////////////////////////////////////////////////////////////////////////////////

func parseDirectDeclarator(tokens []Token) (Declarator, []Token) {
	simpleDec, tokens := parseSimpleDeclarator(tokens)

	if peekToken(tokens).tokenType == OPEN_PARENTHESIS_TOKEN {
		paramInfos, tokens := parseParamList(tokens)
		funDec := Function_Declarator{paramInfos: paramInfos, innerDec: simpleDec}
		return &funDec, tokens
	} else {
		return simpleDec, tokens
	}
}

/////////////////////////////////////////////////////////////////////////////////

func parseSimpleDeclarator(tokens []Token) (Declarator, []Token) {
	if peekToken(tokens).tokenType == OPEN_PARENTHESIS_TOKEN {
		_, tokens = expect(OPEN_PARENTHESIS_TOKEN, tokens)
		dec, tokens := parseDeclarator(tokens)
		_, tokens = expect(CLOSE_PARENTHESIS_TOKEN, tokens)
		return dec, tokens
	} else {
		name, tokens := parseIdentifier(tokens)
		identDec := Identifier_Declarator{name}
		return &identDec, tokens
	}
}

/////////////////////////////////////////////////////////////////////////////////

func (dec *Identifier_Declarator) processDeclarator(baseTyp Data_Type) (string, Data_Type, []string) {
	return dec.name, baseTyp, []string{}
}

/////////////////////////////////////////////////////////////////////////////////

func (dec *Pointer_Declarator) processDeclarator(baseTyp Data_Type) (string, Data_Type, []string) {
	derivedType := Data_Type{typ: POINTER_TYPE, refType: &baseTyp}
	return dec.innerDec.processDeclarator(derivedType)
}

/////////////////////////////////////////////////////////////////////////////////

func (dec *Function_Declarator) processDeclarator(baseTyp Data_Type) (string, Data_Type, []string) {
	ident, isIdent := dec.innerDec.(*Identifier_Declarator)
	if !isIdent {
		fail("Can't apply additional type derivations to a function type")
	}
	paramNames := []string{}
	paramTypes := []*Data_Type{}
	for _, paramInfo := range dec.paramInfos {
		paramName, paramType, _ := paramInfo.dec.processDeclarator(paramInfo.dTyp)
		if paramType.typ == FUNCTION_TYPE {
			fail("Function pointers in parameters aren't supported")
		}
		paramNames = append(paramNames, paramName)
		paramTypes = append(paramTypes, &paramType)
	}
	derivedType := Data_Type{typ: FUNCTION_TYPE, paramTypes: paramTypes, returnType: &baseTyp}
	return ident.name, derivedType, paramNames
}

/////////////////////////////////////////////////////////////////////////////////

func parseAbstractDeclarator(tokens []Token) (Abstract_Declarator, []Token) {
	if peekToken(tokens).tokenType == ASTERISK_TOKEN {
		_, tokens = expect(ASTERISK_TOKEN, tokens)
		innerDec, tokens := parseAbstractDeclarator(tokens)
		return &Abstract_Pointer_Declarator{innerDec}, tokens
	} else if peekToken(tokens).tokenType == OPEN_PARENTHESIS_TOKEN {
		_, tokens = expect(OPEN_PARENTHESIS_TOKEN, tokens)
		innerDec, tokens := parseAbstractDeclarator(tokens)
		_, tokens = expect(CLOSE_PARENTHESIS_TOKEN, tokens)
		return &Abstract_Pointer_Declarator{innerDec}, tokens
	} else {
		return &Abstract_Base_Declarator{}, tokens
	}
}

/////////////////////////////////////////////////////////////////////////////////

func (absDec *Abstract_Pointer_Declarator) processAbstractDeclarator(baseTyp Data_Type) Data_Type {
	derivedType := Data_Type{typ: POINTER_TYPE, refType: &baseTyp}
	return absDec.innerDec.processAbstractDeclarator(derivedType)
}

/////////////////////////////////////////////////////////////////////////////////

func (absDec *Abstract_Base_Declarator) processAbstractDeclarator(baseTyp Data_Type) Data_Type {
	return baseTyp
}

/////////////////////////////////////////////////////////////////////////////////

func parseSpecifiers(tokens []Token, storageClassAllowed bool) ([]TokenEnum, []Token) {
	specifiers := []TokenEnum{}

	for isSpecifier(peekToken(tokens)) {
		var spec Token
		spec, tokens = takeToken(tokens)
		specifiers = append(specifiers, spec.tokenType)
	}

	if !storageClassAllowed {
		if isSpecifierInList(STATIC_KEYWORD_TOKEN, specifiers) || isSpecifierInList(EXTERN_KEYWORD_TOKEN, specifiers) {
			fail("Storage class specifier not allowed in parameter lists and cast expressions")
		}
	}

	return specifiers, tokens
}

/////////////////////////////////////////////////////////////////////////////////

func isSpecifier(token Token) bool {
	switch token.tokenType {
	case INT_KEYWORD_TOKEN:
		return true
	case LONG_KEYWORD_TOKEN:
		return true
	case SIGNED_KEYWORD_TOKEN:
		return true
	case UNSIGNED_KEYWORD_TOKEN:
		return true
	case DOUBLE_KEYWORD_TOKEN:
		return true
	case STATIC_KEYWORD_TOKEN:
		return true
	case EXTERN_KEYWORD_TOKEN:
		return true
	default:
		return false
	}
}

/////////////////////////////////////////////////////////////////////////////////

func isSpecifierInList(tokenType TokenEnum, specifiers []TokenEnum) bool {
	for index, _ := range specifiers {
		if specifiers[index] == tokenType {
			return true
		}
	}
	return false
}

/////////////////////////////////////////////////////////////////////////////////

func removeDuplicateSpecifier(specifiers []TokenEnum) []TokenEnum {
	set := make(map[TokenEnum]bool)
	result := []TokenEnum{}

	for _, spec := range specifiers {
		if !set[spec] {
			set[spec] = true
			result = append(result, spec)
		}
	}
	return result
}

func hasDuplicateSpecifier(specifiers []TokenEnum) bool {
	specsNoDupes := removeDuplicateSpecifier(specifiers)
	if len(specifiers) == len(specsNoDupes) {
		return false
	} else {
		return true
	}
}

/////////////////////////////////////////////////////////////////////////////////

func analyzeType(specifiers []TokenEnum) Data_Type {
	if len(specifiers) == 0 {
		fail("Missing type specifier")
	}
	if hasDuplicateSpecifier(specifiers) {
		fail("Can't use same specifier more than once")
	}
	if isSpecifierInList(SIGNED_KEYWORD_TOKEN, specifiers) && isSpecifierInList(UNSIGNED_KEYWORD_TOKEN, specifiers) {
		fail("Can't use both signed and unsigned specifiers")
	}

	if isSpecifierInList(DOUBLE_KEYWORD_TOKEN, specifiers) {
		if len(specifiers) == 1 {
			return Data_Type{typ: DOUBLE_TYPE}
		} else {
			fail("Can't combine 'double' with other type specifiers")
		}
	}
	if isSpecifierInList(UNSIGNED_KEYWORD_TOKEN, specifiers) && isSpecifierInList(LONG_KEYWORD_TOKEN, specifiers) {
		return Data_Type{typ: UNSIGNED_LONG_TYPE}
	}
	if isSpecifierInList(UNSIGNED_KEYWORD_TOKEN, specifiers) {
		return Data_Type{typ: UNSIGNED_INT_TYPE}
	}
	if isSpecifierInList(LONG_KEYWORD_TOKEN, specifiers) {
		return Data_Type{typ: LONG_TYPE}
	}

	return Data_Type{typ: INT_TYPE}
}

/////////////////////////////////////////////////////////////////////////////////

func isDataTypeKeyword(token TokenEnum) bool {
	switch token {
	case INT_KEYWORD_TOKEN:
		return true
	case LONG_KEYWORD_TOKEN:
		return true
	case SIGNED_KEYWORD_TOKEN:
		return true
	case UNSIGNED_KEYWORD_TOKEN:
		return true
	case DOUBLE_KEYWORD_TOKEN:
		return true
	default:
		return false
	}
}

/////////////////////////////////////////////////////////////////////////////////

func analyzeTypeAndStorageClass(specifiers []TokenEnum) (Data_Type, StorageClassEnum) {
	types := []TokenEnum{}
	storageClasses := []TokenEnum{}
	for _, spec := range specifiers {
		if isDataTypeKeyword(spec) {
			types = append(types, spec)
		} else {
			storageClasses = append(storageClasses, spec)
		}
	}

	dTyp := analyzeType(types)

	if len(storageClasses) > 1 {
		fail("Invalid storage class")
	}

	storageClass := NONE_STORAGE_CLASS
	if len(storageClasses) == 1 {
		storageClass = getStorageClass(storageClasses[0])
	}

	return dTyp, storageClass
}

/////////////////////////////////////////////////////////////////////////////////

func parseParamList(tokens []Token) ([]Param_Info, []Token) {
	paramInfos := []Param_Info{}

	_, tokens = expect(OPEN_PARENTHESIS_TOKEN, tokens)

	if peekToken(tokens).tokenType == VOID_KEYWORD_TOKEN {
		// there are no params
		_, tokens = expect(VOID_KEYWORD_TOKEN, tokens)
	} else {
		foundComma := false
		for (peekToken(tokens).tokenType != CLOSE_PARENTHESIS_TOKEN) || foundComma {
			// get the type, static and extern are not allowed for params
			var specifiers []TokenEnum
			specifiers, tokens = parseSpecifiers(tokens, false)
			baseType := analyzeType(specifiers)

			var dec Declarator
			dec, tokens = parseDeclarator(tokens)
			name, decType, _ := dec.processDeclarator(baseType)

			// add it to the list
			paramInfo := Param_Info{dTyp: decType, dec: &Identifier_Declarator{name}}
			paramInfos = append(paramInfos, paramInfo)

			if peekToken(tokens).tokenType == COMMA_TOKEN {
				_, tokens = expect(COMMA_TOKEN, tokens)
				foundComma = true
			} else {
				break
			}
		}
	}

	_, tokens = expect(CLOSE_PARENTHESIS_TOKEN, tokens)

	return paramInfos, tokens
}

/////////////////////////////////////////////////////////////////////////////////

func parseArgList(tokens []Token) ([]Expression, []Token) {
	args := []Expression{}
	foundComma := false

	for (peekToken(tokens).tokenType != CLOSE_PARENTHESIS_TOKEN) || foundComma {
		var exp Expression
		exp, tokens = parseExpression(tokens, 0)
		args = append(args, exp)

		if peekToken(tokens).tokenType == COMMA_TOKEN {
			_, tokens = expect(COMMA_TOKEN, tokens)
			foundComma = true
		} else {
			break
		}
	}

	return args, tokens
}

/////////////////////////////////////////////////////////////////////////////////

func parseBlock(tokens []Token) (Block, []Token) {
	_, tokens = expect(OPEN_BRACE_TOKEN, tokens)

	items := []Block_Item{}
	for peekToken(tokens).tokenType != CLOSE_BRACE_TOKEN {
		var bItem Block_Item
		bItem, tokens = parseBlockItem(tokens)
		items = append(items, bItem)
	}

	_, tokens = expect(CLOSE_BRACE_TOKEN, tokens)
	bl := Block{items: items}
	return bl, tokens
}

/////////////////////////////////////////////////////////////////////////////////

func parseBlockItem(tokens []Token) (Block_Item, []Token) {
	if isSpecifier(peekToken(tokens)) {
		// it's a declaration
		decl, tokens := parseDeclaration(tokens)
		declBlock := Block_Declaration{decl}
		return &declBlock, tokens
	} else {
		// it's a statement
		st, tokens := parseStatement(tokens)
		stBlock := Block_Statement{st}
		return &stBlock, tokens
	}
}

/////////////////////////////////////////////////////////////////////////////////

func parseForInitial(tokens []Token) (For_Initial_Clause, []Token) {
	nextToken := peekToken(tokens)

	if isDataTypeKeyword(nextToken.tokenType) {
		decl, tokens := parseDeclaration(tokens)
		varDecl, ok := decl.(*Variable_Declaration)
		if ok {
			return &For_Initial_Declaration{decl: *varDecl}, tokens
		} else {
			fail("Expected Variable Declaration at beginning of for loop")
		}
	} else {
		// must be an (optional) expression
		exp, tokens := parseOptionalExpression(tokens, SEMICOLON_TOKEN)
		return &For_Initial_Expression{exp: exp}, tokens
	}

	return nil, tokens
}

/////////////////////////////////////////////////////////////////////////////////

func parseStatement(tokens []Token) (Statement, []Token) {
	nextToken := peekToken(tokens)

	if nextToken.tokenType == RETURN_KEYWORD_TOKEN {
		_, tokens = expect(RETURN_KEYWORD_TOKEN, tokens)
		ex, tokens := parseExpression(tokens, 0)
		_, tokens = expect(SEMICOLON_TOKEN, tokens)
		st := Return_Statement{exp: ex}
		return &st, tokens
	} else if nextToken.tokenType == SEMICOLON_TOKEN {
		_, tokens = expect(SEMICOLON_TOKEN, tokens)
		return &Null_Statement{}, tokens
	} else if nextToken.tokenType == IF_KEYWORD_TOKEN {
		_, tokens = expect(IF_KEYWORD_TOKEN, tokens)
		_, tokens = expect(OPEN_PARENTHESIS_TOKEN, tokens)
		cond, tokens := parseExpression(tokens, 0)
		_, tokens = expect(CLOSE_PARENTHESIS_TOKEN, tokens)
		thenSt, tokens := parseStatement(tokens)
		var elseSt Statement = nil
		if peekToken(tokens).tokenType == ELSE_KEYWORD_TOKEN {
			_, tokens = expect(ELSE_KEYWORD_TOKEN, tokens)
			elseSt, tokens = parseStatement(tokens)
		}
		return &If_Statement{condition: cond, thenSt: thenSt, elseSt: elseSt}, tokens
	} else if nextToken.tokenType == OPEN_BRACE_TOKEN {
		block, tokens := parseBlock(tokens)
		return &Compound_Statement{block: block}, tokens
	} else if nextToken.tokenType == BREAK_KEYWORD_TOKEN {
		_, tokens = expect(BREAK_KEYWORD_TOKEN, tokens)
		_, tokens = expect(SEMICOLON_TOKEN, tokens)
		return &Break_Statement{}, tokens
	} else if nextToken.tokenType == CONTINUE_KEYWORD_TOKEN {
		_, tokens = expect(CONTINUE_KEYWORD_TOKEN, tokens)
		_, tokens = expect(SEMICOLON_TOKEN, tokens)
		return &Continue_Statement{}, tokens
	} else if nextToken.tokenType == WHILE_KEYWORD_TOKEN {
		_, tokens = expect(WHILE_KEYWORD_TOKEN, tokens)
		_, tokens = expect(OPEN_PARENTHESIS_TOKEN, tokens)
		condition, tokens := parseExpression(tokens, 0)
		_, tokens = expect(CLOSE_PARENTHESIS_TOKEN, tokens)
		body, tokens := parseStatement(tokens)
		return &While_Statement{condition: condition, body: body}, tokens
	} else if nextToken.tokenType == DO_KEYWORD_TOKEN {
		_, tokens = expect(DO_KEYWORD_TOKEN, tokens)
		body, tokens := parseStatement(tokens)
		_, tokens = expect(WHILE_KEYWORD_TOKEN, tokens)
		_, tokens = expect(OPEN_PARENTHESIS_TOKEN, tokens)
		condition, tokens := parseExpression(tokens, 0)
		_, tokens = expect(CLOSE_PARENTHESIS_TOKEN, tokens)
		_, tokens = expect(SEMICOLON_TOKEN, tokens)
		return &Do_While_Statement{body: body, condition: condition}, tokens
	} else if nextToken.tokenType == FOR_KEYWORD_TOKEN {
		_, tokens = expect(FOR_KEYWORD_TOKEN, tokens)
		_, tokens = expect(OPEN_PARENTHESIS_TOKEN, tokens)
		forInit, tokens := parseForInitial(tokens)
		condition, tokens := parseOptionalExpression(tokens, SEMICOLON_TOKEN)
		post, tokens := parseOptionalExpression(tokens, CLOSE_PARENTHESIS_TOKEN)
		body, tokens := parseStatement(tokens)
		return &For_Statement{initial: forInit, condition: condition, post: post, body: body}, tokens
	} else {
		var exp Expression
		exp, tokens = parseExpression(tokens, 0)
		_, tokens = expect(SEMICOLON_TOKEN, tokens)
		return &Expression_Statement{exp: exp}, tokens
	}
}

/////////////////////////////////////////////////////////////////////////////////

func parseOptionalExpression(tokens []Token, expectedEndToken TokenEnum) (Expression, []Token) {
	nextToken := peekToken(tokens)

	if nextToken.tokenType == expectedEndToken {
		_, tokens = expect(expectedEndToken, tokens)
		return nil, tokens
	} else {
		exp, tokens := parseExpression(tokens, 0)
		_, tokens = expect(expectedEndToken, tokens)
		return exp, tokens
	}
}

/////////////////////////////////////////////////////////////////////////////////

func parseExpression(tokens []Token, minPrecedence int) (Expression, []Token) {
	left, tokens := parseFactor(tokens)
	nextToken := peekToken(tokens)

	for isBinaryOperator(nextToken) && getPrecedence(nextToken) >= minPrecedence {
		if nextToken.tokenType == EQUAL_TOKEN {
			_, tokens = expect(EQUAL_TOKEN, tokens)
			var right Expression
			right, tokens = parseExpression(tokens, getPrecedence(nextToken))
			left = &Assignment_Expression{lvalue: left, rightExp: right}
		} else if nextToken.tokenType == QUESTION_TOKEN {
			_, tokens = expect(QUESTION_TOKEN, tokens)
			var middleExp Expression
			middleExp, tokens = parseExpression(tokens, 0)
			_, tokens = expect(COLON_TOKEN, tokens)
			var rightExp Expression
			rightExp, tokens = parseExpression(tokens, getPrecedence(nextToken))
			left = &Conditional_Expression{condition: left, middleExp: middleExp, rightExp: rightExp}
		} else {
			var binOpType BinaryOperatorType
			binOpType, tokens = parseBinaryOperator(tokens)
			var right Expression
			right, tokens = parseExpression(tokens, getPrecedence(nextToken)+1)
			left = &Binary_Expression{binOp: binOpType, firstExp: left, secExp: right}
		}
		nextToken = peekToken(tokens)
	}

	return left, tokens
}

/////////////////////////////////////////////////////////////////////////////////

func getPrecedence(token Token) int {
	switch token.tokenType {
	case ASTERISK_TOKEN:
		return 50
	case FORWARD_SLASH_TOKEN:
		return 50
	case PERCENT_TOKEN:
		return 50
	case PLUS_TOKEN:
		return 45
	case HYPHEN_TOKEN:
		return 45
	case LESS_THAN_TOKEN:
		return 35
	case LESS_OR_EQUAL_TOKEN:
		return 35
	case GREATER_THAN_TOKEN:
		return 35
	case GREATER_OR_EQUAL_TOKEN:
		return 35
	case TWO_EQUAL_SIGNS_TOKEN:
		return 30
	case EXCLAMATION_EQUAL_TOKEN:
		return 30
	case TWO_AMPERSANDS_TOKEN:
		return 10
	case TWO_VERTICAL_BARS_TOKEN:
		return 5
	case QUESTION_TOKEN:
		return 3
	case EQUAL_TOKEN:
		return 1
	default:
		fail("unknown token type")
	}

	return 0
}

/////////////////////////////////////////////////////////////////////////////////

func parseFactor(tokens []Token) (Expression, []Token) {
	nextToken := peekToken(tokens)

	if constantTokenToDataType(nextToken) != NONE_TYPE {
		value, typ, tokens := parseConstantValue(tokens)
		ex := Constant_Value_Expression{dTyp: Data_Type{typ: typ}, value: value}
		return &ex, tokens
	} else if nextToken.tokenType == IDENTIFIER_TOKEN {
		name, tokens := parseIdentifier(tokens)
		if peekToken(tokens).tokenType == OPEN_PARENTHESIS_TOKEN {
			// it's a function call
			_, tokens = expect(OPEN_PARENTHESIS_TOKEN, tokens)
			args, tokens := parseArgList(tokens)
			_, tokens = expect(CLOSE_PARENTHESIS_TOKEN, tokens)
			f := Function_Call_Expression{functionName: name, args: args}
			return &f, tokens
		} else {
			// it's just a variable expression
			v := Variable_Expression{name: name}
			return &v, tokens
		}
	} else if isUnaryOperator(nextToken) {
		unopType, tokens := parseUnaryOperator(tokens)
		innerExp, tokens := parseFactor(tokens)
		if unopType == DEREFERENCE_OPERATOR {
			return &Dereference_Expression{innerExp: innerExp}, tokens
		} else if unopType == ADDRESS_OF_OPERATOR {
			return &Address_Of_Expression{innerExp: innerExp}, tokens
		} else {
			unExp := Unary_Expression{innerExp: innerExp, unOp: unopType}
			return &unExp, tokens
		}
	} else if nextToken.tokenType == OPEN_PARENTHESIS_TOKEN {
		_, tokens = expect(OPEN_PARENTHESIS_TOKEN, tokens)
		if isDataTypeKeyword(peekToken(tokens).tokenType) {
			// must be a cast expression
			specifiers, tokens := parseSpecifiers(tokens, false)
			baseTyp := analyzeType(specifiers)
			absDec, tokens := parseAbstractDeclarator(tokens)
			derivedTyp := absDec.processAbstractDeclarator(baseTyp)
			_, tokens = expect(CLOSE_PARENTHESIS_TOKEN, tokens)
			exp, tokens := parseFactor(tokens)
			cast := Cast_Expression{targetType: derivedTyp, innerExp: exp}
			return &cast, tokens
		} else {
			// must be another expression within parentheses
			innerExp, tokens := parseExpression(tokens, 0)
			_, tokens = expect(CLOSE_PARENTHESIS_TOKEN, tokens)
			return innerExp, tokens
		}
	} else {
		fail("Malformed expression. Unexpected", allRegexp[nextToken.tokenType].String())
	}

	// should never reach here, but go compiler complains if there's no return statement
	return nil, []Token{}
}

/////////////////////////////////////////////////////////////////////////////////

func parseUnaryOperator(tokens []Token) (UnaryOperatorType, []Token) {
	unopToken, tokens := takeToken(tokens)
	unOpTyp := getUnaryOperator(unopToken)
	return unOpTyp, tokens
}

/////////////////////////////////////////////////////////////////////////////////

func parseBinaryOperator(tokens []Token) (BinaryOperatorType, []Token) {
	binopToken, tokens := takeToken(tokens)
	binOpTyp := getBinaryOperator(binopToken)
	return binOpTyp, tokens
}

/////////////////////////////////////////////////////////////////////////////////

func parseIdentifier(tokens []Token) (string, []Token) {
	currentToken, tokens := expect(IDENTIFIER_TOKEN, tokens)
	id := currentToken.word
	return id, tokens
}

/////////////////////////////////////////////////////////////////////////////////

func constantTokenToDataType(token Token) DataTypeEnum {
	switch token.tokenType {
	case INT_CONSTANT_TOKEN:
		return INT_TYPE
	case LONG_CONSTANT_TOKEN:
		return LONG_TYPE
	case UNSIGNED_INT_CONSTANT_TOKEN:
		return UNSIGNED_INT_TYPE
	case UNSIGNED_LONG_CONSTANT_TOKEN:
		return UNSIGNED_LONG_TYPE
	case DOUBLE_CONSTANT_TOKEN:
		return DOUBLE_TYPE
	default:
		return NONE_TYPE
	}
}

/////////////////////////////////////////////////////////////////////////////////

func parseConstantValue(tokens []Token) (string, DataTypeEnum, []Token) {
	allowedTokens := []TokenEnum{INT_CONSTANT_TOKEN, LONG_CONSTANT_TOKEN, UNSIGNED_INT_CONSTANT_TOKEN,
		UNSIGNED_LONG_CONSTANT_TOKEN, DOUBLE_CONSTANT_TOKEN}
	currentToken, tokens := expectMultiple(allowedTokens, tokens)

	dataTyp := constantTokenToDataType(currentToken)

	// remove the trailing characters that some constants have, like the L in 100L, Go's parser doesn't like them
	currentToken.word = strings.TrimRight(currentToken.word, "lLuU")

	if (dataTyp == INT_TYPE) || (dataTyp == LONG_TYPE) {
		integer, err := strconv.ParseInt(currentToken.word, 10, 64)
		if err != nil {
			fail("Could not parse integer:", err.Error())
		}
		if (dataTyp == INT_TYPE) && (integer > math.MaxInt32) {
			dataTyp = LONG_TYPE
		}
	} else if (dataTyp == UNSIGNED_INT_TYPE) || (dataTyp == UNSIGNED_LONG_TYPE) {
		integer, err := strconv.ParseUint(currentToken.word, 10, 64)
		if err != nil {
			fail("Could not parse unsigned integer:", err.Error())
		}
		if (dataTyp == UNSIGNED_INT_TYPE) && (integer > math.MaxUint32) {
			dataTyp = UNSIGNED_LONG_TYPE
		}
	} else if dataTyp == DOUBLE_TYPE {
		currentToken.word = roundDouble(currentToken.word)
	}

	return currentToken.word, dataTyp, tokens
}

/////////////////////////////////////////////////////////////////////////////////

func expect(expected TokenEnum, tokens []Token) (Token, []Token) {
	actual, tokens := takeToken(tokens)

	if actual.tokenType != expected {
		fail("Syntax error. Expected", allRegexp[expected].String(), "but found", allRegexp[actual.tokenType].String())
	}

	return actual, tokens
}

/////////////////////////////////////////////////////////////////////////////////

func expectMultiple(expected []TokenEnum, tokens []Token) (Token, []Token) {
	actual, tokens := takeToken(tokens)

	found := false
	for index, _ := range expected {
		if expected[index] == actual.tokenType {
			found = true
		}
	}
	if !found {
		fail("Syntax error. Unexpected", actual.word)
	}

	return actual, tokens
}

/////////////////////////////////////////////////////////////////////////////////

func peekToken(tokens []Token) Token {
	if len(tokens) == 0 {
		fmt.Println("peekToken(): Ran out of tokens.")
		return Token{word: "", tokenType: NONE_TOKEN}
	}

	firstToken := tokens[0]
	return firstToken
}

/////////////////////////////////////////////////////////////////////////////////

func takeToken(tokens []Token) (Token, []Token) {
	// check for no more tokens, don't call os.Exit here, we don't have enough information to print a useful error message
	if len(tokens) == 0 {
		fmt.Println("takeToken(): Ran out of tokens.")
		return Token{word: "", tokenType: NONE_TOKEN}, tokens
	}

	firstToken := tokens[0]
	tokens = tokens[1:]

	return firstToken, tokens
}
