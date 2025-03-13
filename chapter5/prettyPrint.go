package main

import (
	"fmt"
	"os"
	"strconv"
)

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

func (p *Program) getPrettyPrintLines() []string {
	lines := []string{"Program("}
	newLines := p.fn.getPrettyPrintLines()
	lines = append(lines, newLines...)
	lines = append(lines, ")")
	return lines
}

/////////////////////////////////////////////////////////////////////////////////

func (f *Function) getPrettyPrintLines() []string {
	lines := []string{"Function("}
	lines = append(lines, "name="+string(f.name)+",")
	lines = append(lines, "body=")

	for _, block := range f.body {
		moreLines := block.getPrettyPrintLines()
		lines = append(lines, moreLines...)
	}

	lines = append(lines, ")")
	return lines
}

//###############################################################################
//###############################################################################
//###############################################################################

func (b *Block_Statement) getPrettyPrintLines() []string {
	return b.st.getPrettyPrintLines()
}

/////////////////////////////////////////////////////////////////////////////////

func (b *Block_Declaration) getPrettyPrintLines() []string {
	return b.decl.getPrettyPrintLines()
}

//###############################################################################
//###############################################################################
//###############################################################################

func (d *Declaration) getPrettyPrintLines() []string {
	lines := []string{"DECLARATION("}
	lines = append(lines, "name="+string(d.name)+",")
	lines = append(lines, "initializer=")
	if d.initializer != nil {
		moreLines := d.initializer.getPrettyPrintLines()
		lines = append(lines, moreLines...)
	}

	lines = append(lines, "),")
	return lines
}

//###############################################################################
//###############################################################################
//###############################################################################

func (s *Return_Statement) getPrettyPrintLines() []string {
	lines := []string{"RETURN_STATEMENT("}
	moreLines := s.exp.getPrettyPrintLines()
	lines = append(lines, moreLines...)
	lines = append(lines, "),")
	return lines
}

/////////////////////////////////////////////////////////////////////////////////

func (s *Expression_Statement) getPrettyPrintLines() []string {
	lines := []string{"EXPRESSION_STATEMENT("}
	moreLines := s.exp.getPrettyPrintLines()
	lines = append(lines, moreLines...)
	lines = append(lines, "),")
	return lines
}

/////////////////////////////////////////////////////////////////////////////////

func (s *Null_Statement) getPrettyPrintLines() []string {
	return []string{"NULL_STATEMENT(),"}
}

//###############################################################################
//###############################################################################
//###############################################################################

func (e *Constant_Int_Expression) getPrettyPrintLines() []string {
	line := "CONSTANT_INT_EXPRESSION" + "(" + strconv.FormatInt(int64(e.intValue), 10) + ")"
	return []string{line}
}

/////////////////////////////////////////////////////////////////////////////////

func (e *Variable_Expression) getPrettyPrintLines() []string {
	return []string{"VARIABLE_EXPRESSION(" + e.name + ")"}
}

/////////////////////////////////////////////////////////////////////////////////

func (e *Unary_Expression) getPrettyPrintLines() []string {
	lines := []string{}
	line := "UNARY_EXPRESSION_" + getPrettyPrintUnary(e.unOp) + "("
	lines = append(lines, line)
	moreLines := e.innerExp.getPrettyPrintLines()
	lines = append(lines, moreLines...)
	lines = append(lines, ")")
	return lines
}

/////////////////////////////////////////////////////////////////////////////////

func (e *Binary_Expression) getPrettyPrintLines() []string {
	lines := []string{}
	line := "BINARY_EXPRESSION_" + getPrettyPrintBinary(e.binOp) + "("
	lines = append(lines, line)
	moreLines := e.firstExp.getPrettyPrintLines()
	moreLines[len(moreLines)-1] = moreLines[len(moreLines)-1] + ","
	lines = append(lines, moreLines...)
	moreLines = e.secExp.getPrettyPrintLines()
	lines = append(lines, moreLines...)
	lines = append(lines, ")")
	return lines
}

/////////////////////////////////////////////////////////////////////////////////

func (e *Assignment_Expression) getPrettyPrintLines() []string {
	lines := []string{"ASSIGNMENT_EXPRESSION("}
	lines = append(lines, "lvalue=")
	moreLines := e.lvalue.getPrettyPrintLines()
	lines = append(lines, moreLines...)
	lines = append(lines, "rightExp=")
	moreLines = e.rightExp.getPrettyPrintLines()
	lines = append(lines, moreLines...)
	lines = append(lines, ")")
	return lines
}

//###############################################################################
//###############################################################################
//###############################################################################

func getPrettyPrintUnary(typ UnaryOperatorType) string {
	switch typ {
	case COMPLEMENT_OPERATOR:
		return "COMPLEMENT"
	case NEGATE_OPERATOR:
		return "NEGATE"
	case NOT_OPERATOR:
		return "NOT"
	default:
		fmt.Println("unknown Unary operator:", typ)
		os.Exit(1)
	}

	return ""
}

/////////////////////////////////////////////////////////////////////////////////

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
	case AND_OPERATOR:
		return "AND"
	case OR_OPERATOR:
		return "OR"
	case IS_EQUAL_OPERATOR:
		return "IS_EQUAL"
	case NOT_EQUAL_OPERATOR:
		return "NOT_EQUAL"
	case LESS_THAN_OPERATOR:
		return "LESS_THAN"
	case LESS_OR_EQUAL_OPERATOR:
		return "LESS_OR_EQUAL"
	case GREATER_THAN_OPERATOR:
		return "GREATER_THAN"
	case GREATER_OR_EQUAL_OPERATOR:
		return "GREATER_OR_EQUAL"
	default:
		fmt.Println("Unknown binary operator:", typ)
		os.Exit(1)
	}

	return ""
}
