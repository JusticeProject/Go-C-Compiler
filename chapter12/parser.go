package main

import (
	"fmt"
	"math"
	"strconv"
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

func getStorageClass(token Token) StorageClassEnum {
	switch token.tokenType {
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
	FUNCTION_TYPE
)

type Data_Type struct {
	typ DataTypeEnum

	// for FUNCTION_TYPE
	paramTypes []*Data_Type
	returnType *Data_Type

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
	return dt.returnType.isEqualType(input.returnType)
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
	expToTacky(instructions []Instruction_Tacky) (Value_Tacky, []Instruction_Tacky)
	getPrettyPrintLines() []string
}

type Constant_Value_Expression struct {
	typ       DataTypeEnum
	value     string
	resultTyp DataTypeEnum
}

type Variable_Expression struct {
	name      string
	resultTyp DataTypeEnum
}

type Cast_Expression struct {
	targetType DataTypeEnum
	innerExp   Expression
	resultTyp  DataTypeEnum
}

type Unary_Expression struct {
	unOp      UnaryOperatorType
	innerExp  Expression
	resultTyp DataTypeEnum
}

type Binary_Expression struct {
	binOp     BinaryOperatorType
	firstExp  Expression
	secExp    Expression
	resultTyp DataTypeEnum
}

type Assignment_Expression struct {
	lvalue    Expression
	rightExp  Expression
	resultTyp DataTypeEnum
}

// example: a == 3 ? 1 : 2
type Conditional_Expression struct {
	condition Expression
	middleExp Expression
	rightExp  Expression
	resultTyp DataTypeEnum
}

type Function_Call_Expression struct {
	functionName string
	args         []Expression
	resultTyp    DataTypeEnum
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
)

func getUnaryOperator(token Token) UnaryOperatorType {
	switch token.tokenType {
	case TILDE_TOKEN:
		return COMPLEMENT_OPERATOR
	case HYPHEN_TOKEN:
		return NEGATE_OPERATOR
	case EXCLAMATION_TOKEN:
		return NOT_OPERATOR
	}

	return NONE_UNARY_OPERATOR
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
	var specifiers []Token
	specifiers, tokens = parseSpecifiers(tokens)
	if len(specifiers) == 0 {
		return nil, tokens
	}
	typ, storageClass := analyzeTypeAndStorageClass(specifiers)
	name, tokens := parseIdentifier(tokens)

	if peekToken(tokens).tokenType == OPEN_PARENTHESIS_TOKEN {
		// it's a function declaration or definition
		_, tokens = expect(OPEN_PARENTHESIS_TOKEN, tokens)
		paramNames, paramTypes, tokens := parseParamList(tokens)
		_, tokens = expect(CLOSE_PARENTHESIS_TOKEN, tokens)
		funcType := Data_Type{typ: FUNCTION_TYPE, paramTypes: paramTypes, returnType: &Data_Type{typ: typ}}

		if peekToken(tokens).tokenType == SEMICOLON_TOKEN {
			// it's a function declaration
			_, tokens = expect(SEMICOLON_TOKEN, tokens)
			fn := Function_Declaration{name: name, paramNames: paramNames, body: nil, dTyp: funcType, storageClass: storageClass}
			return &fn, tokens
		} else {
			// it's a function definition
			block, tokens := parseBlock(tokens)
			fn := Function_Declaration{name: name, paramNames: paramNames, body: &block, dTyp: funcType, storageClass: storageClass}
			return &fn, tokens
		}
	} else {
		// it's a variable declaration
		decl := Variable_Declaration{name: name, dTyp: Data_Type{typ: typ}, storageClass: storageClass}

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

func parseSpecifiers(tokens []Token) ([]Token, []Token) {
	specifiers := []Token{}

	for isSpecifier(peekToken(tokens)) {
		var spec Token
		spec, tokens = takeToken(tokens)
		specifiers = append(specifiers, spec)
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
	case STATIC_KEYWORD_TOKEN:
		return true
	case EXTERN_KEYWORD_TOKEN:
		return true
	default:
		return false
	}
}

/////////////////////////////////////////////////////////////////////////////////

func isSpecifierInList(tokenType TokenEnum, specifiers []Token) bool {
	for index, _ := range specifiers {
		if specifiers[index].tokenType == tokenType {
			return true
		}
	}
	return false
}

/////////////////////////////////////////////////////////////////////////////////

func analyzeType(specifiers []Token) DataTypeEnum {
	if len(specifiers) == 1 && isSpecifierInList(INT_KEYWORD_TOKEN, specifiers) {
		return INT_TYPE
	}

	if isSpecifierInList(LONG_KEYWORD_TOKEN, specifiers) {
		if len(specifiers) == 1 {
			return LONG_TYPE
		}
		if len(specifiers) == 2 && isSpecifierInList(INT_KEYWORD_TOKEN, specifiers) {
			return LONG_TYPE
		}
	}

	fail("Invalid type specifier")
	return NONE_TYPE
}

/////////////////////////////////////////////////////////////////////////////////

func isDataTypeKeyword(token Token) bool {
	switch token.tokenType {
	case INT_KEYWORD_TOKEN:
		return true
	case LONG_KEYWORD_TOKEN:
		return true
	default:
		return false
	}
}

/////////////////////////////////////////////////////////////////////////////////

func analyzeTypeAndStorageClass(specifiers []Token) (DataTypeEnum, StorageClassEnum) {
	types := []Token{}
	storageClasses := []Token{}
	for _, spec := range specifiers {
		if isDataTypeKeyword(spec) {
			types = append(types, spec)
		} else {
			storageClasses = append(storageClasses, spec)
		}
	}

	typ := analyzeType(types)

	if len(storageClasses) > 1 {
		fail("Invalid storage class")
	}

	storageClass := NONE_STORAGE_CLASS
	if len(storageClasses) == 1 {
		storageClass = getStorageClass(storageClasses[0])
	}

	return typ, storageClass
}

/////////////////////////////////////////////////////////////////////////////////

func parseParamList(tokens []Token) ([]string, []*Data_Type, []Token) {
	paramNames := []string{}
	paramTypes := []*Data_Type{}

	if peekToken(tokens).tokenType == VOID_KEYWORD_TOKEN {
		// there are no params
		_, tokens = expect(VOID_KEYWORD_TOKEN, tokens)
	} else {
		foundComma := false
		for (peekToken(tokens).tokenType != CLOSE_PARENTHESIS_TOKEN) || foundComma {
			// get the type, static and extern are not allowed for params
			var specifiers []Token
			specifiers, tokens = parseSpecifiers(tokens)
			typ := analyzeType(specifiers)
			paramTypes = append(paramTypes, &Data_Type{typ: typ})

			// get the name
			var id string
			id, tokens = parseIdentifier(tokens)
			paramNames = append(paramNames, id)
			if peekToken(tokens).tokenType == COMMA_TOKEN {
				_, tokens = expect(COMMA_TOKEN, tokens)
				foundComma = true
			} else {
				break
			}
		}
	}

	return paramNames, paramTypes, tokens
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

	if isDataTypeKeyword(nextToken) {
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
		ex := Constant_Value_Expression{typ: typ, value: value}
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
	} else if nextToken.tokenType == TILDE_TOKEN || nextToken.tokenType == HYPHEN_TOKEN || nextToken.tokenType == EXCLAMATION_TOKEN {
		unopType, tokens := parseUnaryOperator(tokens)
		innerExp, tokens := parseFactor(tokens)
		unExp := Unary_Expression{innerExp: innerExp, unOp: unopType}
		return &unExp, tokens
	} else if nextToken.tokenType == OPEN_PARENTHESIS_TOKEN {
		_, tokens = expect(OPEN_PARENTHESIS_TOKEN, tokens)
		if isDataTypeKeyword(peekToken(tokens)) {
			// must be a cast expression
			specifiers, tokens := parseSpecifiers(tokens)
			typ := analyzeType(specifiers)
			_, tokens = expect(CLOSE_PARENTHESIS_TOKEN, tokens)
			exp, tokens := parseFactor(tokens)
			cast := Cast_Expression{targetType: typ, innerExp: exp}
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
	default:
		return NONE_TYPE
	}
}

/////////////////////////////////////////////////////////////////////////////////

func parseConstantValue(tokens []Token) (string, DataTypeEnum, []Token) {
	currentToken, tokens := expectMultiple([]TokenEnum{INT_CONSTANT_TOKEN, LONG_CONSTANT_TOKEN}, tokens)

	dataTyp := constantTokenToDataType(currentToken)

	if (dataTyp == INT_TYPE) || (dataTyp == LONG_TYPE) {
		integer, err := strconv.ParseInt(currentToken.word, 10, 64)
		if err != nil {
			fail("Could not parse integer:", err.Error())
		}

		if (dataTyp == INT_TYPE) && (integer > math.MaxInt32) {
			dataTyp = LONG_TYPE
		}
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
