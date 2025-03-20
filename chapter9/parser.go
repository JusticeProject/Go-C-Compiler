package main

import (
	"fmt"
	"os"
	"strconv"
)

//###############################################################################
//###############################################################################
//###############################################################################

type Program struct {
	functions []Function_Declaration
}

//###############################################################################
//###############################################################################
//###############################################################################

type Declaration interface {
	declToTacky() []Instruction_Tacky
	getPrettyPrintLines() []string
}

type Variable_Declaration struct {
	name        string
	initializer Expression
}

type Function_Declaration struct {
	name string
	// TODO: will need to change string to Identifier which will hold a string and a type?
	params []string
	body   *Block
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

type Constant_Int_Expression struct {
	intValue int32
}

type Variable_Expression struct {
	name string
}

type Unary_Expression struct {
	unOp     UnaryOperatorType
	innerExp Expression
}

type Binary_Expression struct {
	binOp    BinaryOperatorType
	firstExp Expression
	secExp   Expression
}

type Assignment_Expression struct {
	lvalue   Expression
	rightExp Expression
}

// example: a == 3 ? 1 : 2
type Conditional_Expression struct {
	condition Expression
	middleExp Expression
	rightExp  Expression
}

type Function_Call_Expression struct {
	functionName string
	args         []Expression
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
		fmt.Println("Sytnax Error. Tokens remaining after parsing program:")
		fmt.Println(tokens)
		os.Exit(1)
	}

	// print the ast in a well-formatted way
	prettyPrint(ast)

	return ast
}

/////////////////////////////////////////////////////////////////////////////////

func parseProgram(tokens []Token) (Program, []Token) {
	functions := []Function_Declaration{}

	// TODO: change to nextToken isDataType()
	for peekToken(tokens).tokenType == INT_KEYWORD_TOKEN {
		var decl Declaration
		decl, tokens = parseDeclaration(tokens)
		fnDecl, ok := decl.(*Function_Declaration)
		if ok {
			functions = append(functions, *fnDecl)
		}

	}

	pr := Program{functions: functions}
	return pr, tokens
}

/////////////////////////////////////////////////////////////////////////////////

func parseDeclaration(tokens []Token) (Declaration, []Token) {
	// TODO: need to handle other data types, expectDataType() helper function?
	_, tokens = expect(INT_KEYWORD_TOKEN, tokens)
	name, tokens := parseIdentifier(tokens)

	if peekToken(tokens).tokenType == OPEN_PARENTHESIS_TOKEN {
		// it's a function declaration or definition
		_, tokens = expect(OPEN_PARENTHESIS_TOKEN, tokens)
		params, tokens := parseParamList(tokens)
		_, tokens = expect(CLOSE_PARENTHESIS_TOKEN, tokens)

		if peekToken(tokens).tokenType == SEMICOLON_TOKEN {
			// it's a function declaration
			_, tokens = expect(SEMICOLON_TOKEN, tokens)
			fn := Function_Declaration{name: name, params: params, body: nil}
			return &fn, tokens
		} else {
			// it's a function definition
			block, tokens := parseBlock(tokens)
			fn := Function_Declaration{name: name, params: params, body: &block}
			return &fn, tokens
		}
	} else {
		// it's a variable declaration
		decl := Variable_Declaration{name: name}

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

func parseParamList(tokens []Token) ([]string, []Token) {
	params := []string{}

	if peekToken(tokens).tokenType == VOID_KEYWORD_TOKEN {
		_, tokens = expect(VOID_KEYWORD_TOKEN, tokens)
	} else {
		foundComma := false
		for (peekToken(tokens).tokenType != CLOSE_PARENTHESIS_TOKEN) || foundComma {
			// TODO: need to handle other data types
			_, tokens = expect(INT_KEYWORD_TOKEN, tokens)
			var id string
			id, tokens = parseIdentifier(tokens)
			params = append(params, id)
			if peekToken(tokens).tokenType == COMMA_TOKEN {
				_, tokens = expect(COMMA_TOKEN, tokens)
				foundComma = true
			} else {
				break
			}
		}
	}

	return params, tokens
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
	firstToken := peekToken(tokens)
	// TODO: need to handle other data types
	if firstToken.tokenType == INT_KEYWORD_TOKEN {
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
	// TODO: use helper function isDataType(nextToken) to check if next token is one of the keywords int, bool, float, etc.
	nextToken := peekToken(tokens)

	if nextToken.tokenType == INT_KEYWORD_TOKEN {
		decl, tokens := parseDeclaration(tokens)
		// TODO: don't like this
		varDecl, ok := decl.(*Variable_Declaration)
		if ok {
			return &For_Initial_Declaration{decl: *varDecl}, tokens
		} else {
			fmt.Println("Expected Variable Declaration at beginning of for loop")
			os.Exit(1)
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
		fmt.Println("unknown token type:", token.tokenType)
		os.Exit(1)
	}

	return 0
}

/////////////////////////////////////////////////////////////////////////////////

func parseFactor(tokens []Token) (Expression, []Token) {
	nextToken := peekToken(tokens)

	if nextToken.tokenType == INT_CONSTANT_TOKEN {
		integer, tokens := parseInteger(tokens)
		ex := Constant_Int_Expression{intValue: integer}
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
			v := Variable_Expression{name}
			return &v, tokens
		}
	} else if nextToken.tokenType == TILDE_TOKEN || nextToken.tokenType == HYPHEN_TOKEN || nextToken.tokenType == EXCLAMATION_TOKEN {
		unopType, tokens := parseUnaryOperator(tokens)
		innerExp, tokens := parseFactor(tokens)
		unExp := Unary_Expression{innerExp: innerExp, unOp: unopType}
		return &unExp, tokens
	} else if nextToken.tokenType == OPEN_PARENTHESIS_TOKEN {
		_, tokens := expect(OPEN_PARENTHESIS_TOKEN, tokens)
		innerExp, tokens := parseExpression(tokens, 0)
		_, tokens = expect(CLOSE_PARENTHESIS_TOKEN, tokens)
		return innerExp, tokens
	} else {
		fmt.Println("Malformed expression.")
		fmt.Println("Unexpected", allRegexp[nextToken.tokenType])
		os.Exit(1)
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

func parseInteger(tokens []Token) (int32, []Token) {
	currentToken, tokens := expect(INT_CONSTANT_TOKEN, tokens)
	// TODO: what about 8, 16, 64-bit integers?
	integer, _ := strconv.ParseInt(currentToken.word, 10, 64)
	return int32(integer), tokens
}

/////////////////////////////////////////////////////////////////////////////////

func expect(expected TokenEnum, tokens []Token) (Token, []Token) {
	actual, tokens := takeToken(tokens)

	if actual.tokenType != expected {
		// TODO: make this error msg more human readable, need function to convert TokenEnum to string
		fmt.Println("Syntax error.")
		fmt.Println("Expected", allRegexp[expected])
		fmt.Println("but found", allRegexp[actual.tokenType])
		os.Exit(1)
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
