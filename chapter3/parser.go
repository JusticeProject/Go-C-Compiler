package main

import (
	"fmt"
	"os"
	"strconv"
)

/////////////////////////////////////////////////////////////////////////////////

type Program struct {
	fn Function
}

func (p *Program) getPrettyPrintLines() []string {
	lines := []string{"Program("}
	newLines := p.fn.getPrettyPrintLines()
	lines = append(lines, newLines...)
	lines = append(lines, ")")
	return lines
}

/////////////////////////////////////////////////////////////////////////////////

type Function struct {
	name string
	body Statement
}

func (f *Function) getPrettyPrintLines() []string {
	lines := []string{"Function("}
	lines = append(lines, "name="+string(f.name)+",")
	lines = append(lines, "body=")
	moreLines := f.body.getPrettyPrintLines()
	lines = append(lines, moreLines...)
	lines = append(lines, ")")
	return lines
}

/////////////////////////////////////////////////////////////////////////////////

type StatementType int

const (
	RETURN_STATEMENT StatementType = iota
	IF_STATEMENT
)

type Statement struct {
	typ StatementType
	exp Expression
}

func (s *Statement) getPrettyPrintLines() []string {
	typeOfStatement := s.getDesc()

	lines := []string{typeOfStatement + "("}
	moreLines := s.exp.getPrettyPrintLines()
	lines = append(lines, moreLines...)
	lines = append(lines, ")")
	return lines
}

func (s *Statement) getDesc() string {
	switch s.typ {
	case RETURN_STATEMENT:
		return "RETURN_STATEMENT"
	case IF_STATEMENT:
		return "IF_STATEMENT"
	default:
		return "UNKNOWN_STATEMENT"
	}
}

/////////////////////////////////////////////////////////////////////////////////

type ExpressionType int

const (
	CONSTANT_INT_EXPRESSION ExpressionType = iota
	UNARY_EXPRESSION
	BINARY_EXPRESSION
)

type Expression struct {
	typ      ExpressionType
	intValue int32
	firstExp *Expression
	secExp   *Expression
	unOp     UnaryOperatorType
	binOp    BinaryOperatorType
}

func (e *Expression) getPrettyPrintLines() []string {
	lines := []string{}

	switch e.typ {
	case CONSTANT_INT_EXPRESSION:
		line := "CONSTANT_INT_EXPRESSION" + "(" + strconv.FormatInt(int64(e.intValue), 10) + ")"
		lines = append(lines, line)
	case UNARY_EXPRESSION:
		line := "UNARY_EXPRESSION_" + getPrettyPrintUnary(e.unOp) + "("
		lines = append(lines, line)
		moreLines := e.firstExp.getPrettyPrintLines()
		lines = append(lines, moreLines...)
		lines = append(lines, ")")
	case BINARY_EXPRESSION:
		line := "BINARY_EXPRESSION_" + getPrettyPrintBinary(e.binOp) + "("
		lines = append(lines, line)
		moreLines := e.firstExp.getPrettyPrintLines()
		moreLines[len(moreLines)-1] = moreLines[len(moreLines)-1] + ","
		lines = append(lines, moreLines...)
		moreLines = e.secExp.getPrettyPrintLines()
		lines = append(lines, moreLines...)
		lines = append(lines, ")")
	}

	return lines
}

/////////////////////////////////////////////////////////////////////////////////

type UnaryOperatorType int

const (
	NONE_UNARY_OPERATOR UnaryOperatorType = iota
	COMPLEMENT_OPERATOR
	NEGATE_OPERATOR
)

func getPrettyPrintUnary(typ UnaryOperatorType) string {
	switch typ {
	case COMPLEMENT_OPERATOR:
		return "COMPLEMENT"
	case NEGATE_OPERATOR:
		return "NEGATE"
	default:
		fmt.Println("unknown Unary operator:", typ)
		os.Exit(1)
	}

	return ""
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
	}

	return NONE_BINARY_OPERATOR
}

func getPrettyPrintBinary(typ BinaryOperatorType) string {
	switch typ {
	case ADD_OPERATOR:
		return "ADD"
	case SUBTRACT_OPERATOR:
		return "SUBTRACT"
	case MULTIPLY_OPERATOR:
		return "MULTIPLY"
	case DIVIDE_OPERATOR:
		return "DIVIDE"
	case REMAINDER_OPERATOR:
		return "REMAINDER"
	default:
		fmt.Println("Unknown binary operator:", typ)
		os.Exit(1)
	}

	return ""
}

/////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////

func doParser(tokens []Token) Program {
	// get the Abstract Syntax Tree of the entire program
	ast, tokens := parseProgram(tokens)

	// if there are any remaining tokens, generate syntax error
	if len(tokens) > 0 {
		fmt.Println("Sytnax Error. Tokens remaining after parsing program:")
		fmt.Println(tokens)
		os.Exit(1)
	}

	// pretty-print the ast
	prettyPrint(ast)

	return ast
}

/////////////////////////////////////////////////////////////////////////////////

func parseProgram(tokens []Token) (Program, []Token) {
	fn, tokens := parseFunction(tokens)
	pr := Program{fn: fn}
	return pr, tokens
}

/////////////////////////////////////////////////////////////////////////////////

func parseFunction(tokens []Token) (Function, []Token) {
	_, tokens = expect(INT_KEYWORD_TOKEN, tokens)
	id, tokens := parseIdentifier(tokens)
	_, tokens = expect(OPEN_PARENTHESIS_TOKEN, tokens)
	_, tokens = expect(VOID_KEYWORD_TOKEN, tokens)
	_, tokens = expect(CLOSE_PARENTHESIS_TOKEN, tokens)
	_, tokens = expect(OPEN_BRACE_TOKEN, tokens)
	st, tokens := parseStatement(tokens)
	_, tokens = expect(CLOSE_BRACE_TOKEN, tokens)
	fn := Function{name: id, body: st}
	return fn, tokens
}

/////////////////////////////////////////////////////////////////////////////////

func parseStatement(tokens []Token) (Statement, []Token) {
	_, tokens = expect(RETURN_KEYWORD_TOKEN, tokens)
	ex, tokens := parseExpression(tokens, 0)
	_, tokens = expect(SEMICOLON_TOKEN, tokens)
	st := Statement{typ: RETURN_STATEMENT, exp: *ex}
	return st, tokens
}

/////////////////////////////////////////////////////////////////////////////////

func parseExpression(tokens []Token, minPrecedence int) (*Expression, []Token) {
	left, tokens := parseFactor(tokens)
	nextToken := peekToken(tokens)

	for getBinaryOperator(nextToken) != NONE_BINARY_OPERATOR && getPrecedence(nextToken) >= minPrecedence {
		var binOpType BinaryOperatorType
		binOpType, tokens = parseBinaryOperator(tokens)
		var right *Expression
		right, tokens = parseExpression(tokens, getPrecedence(nextToken)+1)
		left = &Expression{typ: BINARY_EXPRESSION, firstExp: left, secExp: right, binOp: binOpType}
		nextToken = peekToken(tokens)
	}

	return left, tokens
}

/////////////////////////////////////////////////////////////////////////////////

func getPrecedence(token Token) int {
	// TODO: should this switch on BinaryOperatorType or TokenEnum?
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
	default:
		fmt.Println("unknown token type:", token.tokenType)
		os.Exit(1)
	}

	return 0
}

/////////////////////////////////////////////////////////////////////////////////

func parseFactor(tokens []Token) (*Expression, []Token) {
	nextToken := peekToken(tokens)

	if nextToken.tokenType == INT_CONSTANT_TOKEN {
		integer, tokens := parseInteger(tokens)
		ex := Expression{typ: CONSTANT_INT_EXPRESSION, intValue: integer}
		return &ex, tokens
	} else if nextToken.tokenType == TILDE_TOKEN || nextToken.tokenType == HYPHEN_TOKEN {
		unopType, tokens := parseUnaryOperator(tokens)
		firstExp, tokens := parseFactor(tokens)
		unExp := Expression{typ: UNARY_EXPRESSION, firstExp: firstExp, unOp: unopType}
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

	switch unopToken.tokenType {
	case TILDE_TOKEN:
		return COMPLEMENT_OPERATOR, tokens
	case HYPHEN_TOKEN:
		return NEGATE_OPERATOR, tokens
	default:
		fmt.Println("unknown unary operator:", unopToken.tokenType)
		os.Exit(1)
	}

	return NONE_UNARY_OPERATOR, tokens
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

/////////////////////////////////////////////////////////////////////////////////

func prettyPrint(pr Program) {
	lines := pr.getPrettyPrintLines()

	numTabs := 0
	prevEndOfLine := true
	for _, line := range lines {
		if prevEndOfLine {
			prevEndOfLine = printWithTabs(line, numTabs)
		} else {
			prevEndOfLine = printWithTabs(line, 0)
		}

		lastChar := line[len(line)-1]

		if lastChar == '(' {
			numTabs++
		} else if lastChar == ')' {
			numTabs--
		}
	}
}

/////////////////////////////////////////////////////////////////////////////////

func printWithTabs(text string, numTabs int) bool {
	tabString := ""
	for i := 0; i < numTabs; i++ {
		tabString += "    "
	}

	fmt.Printf("%v%v", tabString, text)

	// return true if it ended with a newline
	lastChar := text[len(text)-1]
	if lastChar == '(' || lastChar == ')' || lastChar == ',' {
		fmt.Printf("\n")
		return true
	} else {
		return false
	}
}

/////////////////////////////////////////////////////////////////////////////////
