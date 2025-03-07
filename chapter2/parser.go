package main

import (
	"fmt"
	"os"
	"strconv"
)

/////////////////////////////////////////////////////////////////////////////////

type Identifier string

/////////////////////////////////////////////////////////////////////////////////

type Program struct {
	fn *Function
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
	name Identifier
	body *Statement
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
	exp *Expression
}

func (s *Statement) getPrettyPrintLines() []string {
	typeOfStatement := s.getDesc()

	lines := []string{typeOfStatement + "("}
	anotherLine := s.exp.getPrettyPrintLine()
	lines = append(lines, anotherLine)
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
)

type Expression struct {
	typ      ExpressionType
	intValue int32
	innerExp *Expression
	unOp     UnaryOperatorType
}

func (e *Expression) getPrettyPrintLine() string {
	line := ""

	switch e.typ {
	case CONSTANT_INT_EXPRESSION:
		line = "CONSTANT_INT_EXPRESSION" + "(" + strconv.FormatInt(int64(e.intValue), 10) + ")"
	case UNARY_EXPRESSION:
		line = "UNARY_EXPRESSION_" + getPrettyPrintUnary(e.unOp) + "("
		line += e.innerExp.getPrettyPrintLine()
		line += ")"
	}

	return line
}

/////////////////////////////////////////////////////////////////////////////////

type UnaryOperatorType int

const (
	COMPLEMENT_OPERATOR UnaryOperatorType = iota
	NEGATE_OPERATOR
)

func getPrettyPrintUnary(typ UnaryOperatorType) string {
	switch typ {
	case COMPLEMENT_OPERATOR:
		return "COMPLEMENT"
	case NEGATE_OPERATOR:
		return "NEGATE"
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
	pr := Program{fn: &fn}
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
	fn := Function{name: id, body: &st}
	return fn, tokens
}

/////////////////////////////////////////////////////////////////////////////////

func parseStatement(tokens []Token) (Statement, []Token) {
	_, tokens = expect(RETURN_KEYWORD_TOKEN, tokens)
	ex, tokens := parseExpression(tokens)
	_, tokens = expect(SEMICOLON_TOKEN, tokens)
	st := Statement{typ: RETURN_STATEMENT, exp: &ex}
	return st, tokens
}

/////////////////////////////////////////////////////////////////////////////////

func parseExpression(tokens []Token) (Expression, []Token) {
	nextToken := peekToken(tokens)

	if nextToken.tokenType == INT_CONSTANT_TOKEN {
		integer, tokens := parseInteger(tokens)
		ex := Expression{typ: CONSTANT_INT_EXPRESSION, intValue: integer}
		return ex, tokens
	} else if nextToken.tokenType == TILDE_TOKEN || nextToken.tokenType == HYPHEN_TOKEN {
		unopType, tokens := parseUnaryOperator(tokens)
		innerExp, tokens := parseExpression(tokens)
		unExp := Expression{typ: UNARY_EXPRESSION, innerExp: &innerExp, unOp: unopType}
		return unExp, tokens
	} else if nextToken.tokenType == OPEN_PARENTHESIS_TOKEN {
		_, tokens := expect(OPEN_PARENTHESIS_TOKEN, tokens)
		innerExp, tokens := parseExpression(tokens)
		_, tokens = expect(CLOSE_PARENTHESIS_TOKEN, tokens)
		return innerExp, tokens
	} else {
		fmt.Println("Malformed expression.")
		fmt.Println("Unexpected", allRegexp[nextToken.tokenType])
		os.Exit(1)
	}

	// should never reach here, but go compiler complains if there's no return statement
	return Expression{}, []Token{}
}

/////////////////////////////////////////////////////////////////////////////////

func parseUnaryOperator(tokens []Token) (UnaryOperatorType, []Token) {
	unopToken, tokens := takeToken(tokens)

	switch unopToken.tokenType {
	case TILDE_TOKEN:
		return COMPLEMENT_OPERATOR, tokens
	default: // or case HYPHEN_TOKEN:
		return NEGATE_OPERATOR, tokens
	}
}

/////////////////////////////////////////////////////////////////////////////////

func parseIdentifier(tokens []Token) (Identifier, []Token) {
	currentToken, tokens := expect(IDENTIFIER_TOKEN, tokens)
	id := Identifier(currentToken.word)
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
