package main

import (
	"fmt"
	"strings"
)

/////////////////////////////////////////////////////////////////////////////////

func prettyPrint(pr Program) {
	lines := pr.getPrettyPrintLines()

	numTabs := 0
	prevEndOfLine := true
	for _, line := range lines {
		if isRightIndent(line) {
			numTabs++
			continue
		} else if isLeftIndent(line) {
			numTabs--
			continue
		}

		if prevEndOfLine {
			prevEndOfLine = printWithTabs(line, numTabs)
		} else {
			prevEndOfLine = printWithTabs(line, 0)
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
	if len(text) > 0 {
		lastChar := text[len(text)-1]
		if lastChar == '(' || lastChar == ')' || lastChar == ',' {
			fmt.Printf("\n")
			return true
		}
	}
	return false
}

/////////////////////////////////////////////////////////////////////////////////

func doRightIndent() string {
	return "->"
}

/////////////////////////////////////////////////////////////////////////////////

func doLeftIndent() string {
	return "<-"
}

/////////////////////////////////////////////////////////////////////////////////

func isRightIndent(text string) bool {
	if text == "->" {
		return true
	} else {
		return false
	}
}

/////////////////////////////////////////////////////////////////////////////////

func isLeftIndent(text string) bool {
	if text == "<-" {
		return true
	} else {
		return false
	}
}

//###############################################################################
//###############################################################################
//###############################################################################

func (p *Program) getPrettyPrintLines() []string {
	lines := []string{"Program(", doRightIndent()}

	for _, decl := range p.decls {
		newLines := decl.getPrettyPrintLines()
		lines = append(lines, newLines...)
	}

	lines = append(lines, doLeftIndent())
	lines = append(lines, ")")
	return lines
}

//###############################################################################
//###############################################################################
//###############################################################################

func (d *Variable_Declaration) getPrettyPrintLines() []string {
	lines := []string{"VARIABLE_DECLARATION(", doRightIndent()}
	lines = append(lines, "name="+string(d.name)+",")
	lines = append(lines, "initializer=")
	if d.initializer == nil {
		lines = append(lines, "NONE")
	} else {
		moreLines := d.initializer.getPrettyPrintLines()
		lines = append(lines, moreLines...)
	}
	lines = append(lines, ",")
	lines = append(lines, "storageClass="+getPrettyPrintStorageClass(d.storageClass))

	lines = append(lines, doLeftIndent())
	lines = append(lines, ")")
	return lines
}

/////////////////////////////////////////////////////////////////////////////////

func (f *Function_Declaration) getPrettyPrintLines() []string {

	// TODO: print out the param types and return type
	if f.body == nil {
		// function declaration
		lines := []string{"Function_Declaration(", doRightIndent(), "name=" + f.name + ","}
		lines = append(lines, "paramNames=")
		lines = append(lines, strings.Join(f.paramNames, ","))
		lines = append(lines, ",")
		lines = append(lines, "storageClass="+getPrettyPrintStorageClass(f.storageClass))
		lines = append(lines, doLeftIndent())
		lines = append(lines, ")")
		return lines
	} else {
		// function definition
		lines := []string{"Function_Definition(", doRightIndent(), "name=" + f.name + ","}
		lines = append(lines, "paramNames=")
		lines = append(lines, strings.Join(f.paramNames, ","))
		lines = append(lines, ",")
		lines = append(lines, "body=")
		moreLines := f.body.getPrettyPrintLines()
		lines = append(lines, moreLines...)
		lines = append(lines, ",")
		lines = append(lines, "storageClass="+getPrettyPrintStorageClass(f.storageClass))
		lines = append(lines, doLeftIndent())
		lines = append(lines, ")")
		return lines
	}
}

/////////////////////////////////////////////////////////////////////////////////

func (bl *Block) getPrettyPrintLines() []string {
	lines := []string{}

	for _, bItem := range bl.items {
		moreLines := bItem.getPrettyPrintLines()
		lines = append(lines, moreLines...)
	}

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

func (fid *For_Initial_Declaration) getPrettyPrintLines() []string {
	return fid.decl.getPrettyPrintLines()
}

/////////////////////////////////////////////////////////////////////////////////

func (fie *For_Initial_Expression) getPrettyPrintLines() []string {
	if fie.exp == nil {
		return []string{}
	} else {
		return fie.exp.getPrettyPrintLines()
	}
}

//###############################################################################
//###############################################################################
//###############################################################################

func (s *Return_Statement) getPrettyPrintLines() []string {
	lines := []string{"RETURN_STATEMENT(", doRightIndent()}
	moreLines := s.exp.getPrettyPrintLines()
	lines = append(lines, moreLines...)
	lines = append(lines, doLeftIndent())
	lines = append(lines, ")")
	return lines
}

/////////////////////////////////////////////////////////////////////////////////

func (s *Expression_Statement) getPrettyPrintLines() []string {
	lines := []string{"EXPRESSION_STATEMENT(", doRightIndent()}
	moreLines := s.exp.getPrettyPrintLines()
	lines = append(lines, moreLines...)
	lines = append(lines, doLeftIndent())
	lines = append(lines, ")")
	return lines
}

/////////////////////////////////////////////////////////////////////////////////

func (s *If_Statement) getPrettyPrintLines() []string {
	lines := []string{"IF_STATEMENT(", doRightIndent(), "condition="}
	moreLines := s.condition.getPrettyPrintLines()
	moreLines[len(moreLines)-1] = moreLines[len(moreLines)-1] + ","
	lines = append(lines, moreLines...)
	lines = append(lines, "then=")
	moreLines = s.thenSt.getPrettyPrintLines()
	if len(moreLines) > 0 {
		moreLines[len(moreLines)-1] = moreLines[len(moreLines)-1] + ","
		lines = append(lines, moreLines...)
	}
	lines = append(lines, "else=")
	if s.elseSt != nil {
		moreLines = s.elseSt.getPrettyPrintLines()
		lines = append(lines, moreLines...)
	}
	lines = append(lines, doLeftIndent())
	lines = append(lines, ")")
	return lines
}

/////////////////////////////////////////////////////////////////////////////////

func (s *Compound_Statement) getPrettyPrintLines() []string {
	return s.block.getPrettyPrintLines()
}

/////////////////////////////////////////////////////////////////////////////////

func (s *Break_Statement) getPrettyPrintLines() []string {
	return []string{"BREAK_STATEMENT()"}
}

/////////////////////////////////////////////////////////////////////////////////

func (s *Continue_Statement) getPrettyPrintLines() []string {
	return []string{"CONTINUE_STATEMENT()"}
}

/////////////////////////////////////////////////////////////////////////////////

func (s *While_Statement) getPrettyPrintLines() []string {
	lines := []string{"WHILE(", doRightIndent(), "condition="}
	moreLines := s.condition.getPrettyPrintLines()
	moreLines[len(moreLines)-1] = moreLines[len(moreLines)-1] + ","
	lines = append(lines, moreLines...)
	lines = append(lines, "body=")
	moreLines = s.body.getPrettyPrintLines()
	lines = append(lines, moreLines...)
	lines = append(lines, doLeftIndent())
	lines = append(lines, ")")
	return lines
}

/////////////////////////////////////////////////////////////////////////////////

func (s *Do_While_Statement) getPrettyPrintLines() []string {
	lines := []string{"DO_WHILE(", doRightIndent(), "body="}
	moreLines := s.body.getPrettyPrintLines()
	moreLines[len(moreLines)-1] = moreLines[len(moreLines)-1] + ","
	lines = append(lines, moreLines...)
	lines = append(lines, "condition=")
	moreLines = s.condition.getPrettyPrintLines()
	lines = append(lines, moreLines...)
	lines = append(lines, doLeftIndent())
	lines = append(lines, ")")
	return lines
}

/////////////////////////////////////////////////////////////////////////////////

func (s *For_Statement) getPrettyPrintLines() []string {
	lines := []string{"FOR(", doRightIndent(), "initial="}

	moreLines := s.initial.getPrettyPrintLines()
	if len(moreLines) > 0 {
		moreLines[len(moreLines)-1] = moreLines[len(moreLines)-1] + ","
	} else {
		moreLines = []string{","}
	}
	lines = append(lines, moreLines...)

	lines = append(lines, "condition=")
	if s.condition != nil {
		moreLines = s.condition.getPrettyPrintLines()
		moreLines[len(moreLines)-1] = moreLines[len(moreLines)-1] + ","
		lines = append(lines, moreLines...)
	} else {
		lines = append(lines, ",")
	}

	lines = append(lines, "post=")
	if s.post != nil {
		moreLines = s.post.getPrettyPrintLines()
		moreLines[len(moreLines)-1] = moreLines[len(moreLines)-1] + ","
		lines = append(lines, moreLines...)
	} else {
		lines = append(lines, ",")
	}

	lines = append(lines, "body=")
	moreLines = s.body.getPrettyPrintLines()
	lines = append(lines, moreLines...)

	lines = append(lines, doLeftIndent())
	lines = append(lines, ")")
	return lines
}

/////////////////////////////////////////////////////////////////////////////////

func (s *Null_Statement) getPrettyPrintLines() []string {
	return []string{"NULL_STATEMENT()"}
}

//###############################################################################
//###############################################################################
//###############################################################################

func (e *Constant_Value_Expression) getPrettyPrintLines() []string {
	line := "CONSTANT_VALUE_EXPRESSION_" + getPrettyPrintDataType(e.dTyp.typ) + "(" + e.value + ")"
	return []string{line}
}

/////////////////////////////////////////////////////////////////////////////////

func (e *Variable_Expression) getPrettyPrintLines() []string {
	return []string{"VARIABLE_EXPRESSION(" + e.name + ")"}
}

/////////////////////////////////////////////////////////////////////////////////

func (e *Cast_Expression) getPrettyPrintLines() []string {
	lines := []string{"CAST_EXPRESSION_" + getPrettyPrintDataType(e.targetType.typ) + "(", doRightIndent()}
	moreLines := e.innerExp.getPrettyPrintLines()
	lines = append(lines, moreLines...)
	lines = append(lines, doLeftIndent())
	lines = append(lines, ")")
	return lines
}

/////////////////////////////////////////////////////////////////////////////////

func (e *Unary_Expression) getPrettyPrintLines() []string {
	lines := []string{"UNARY_EXPRESSION_" + getPrettyPrintUnary(e.unOp) + "(", doRightIndent()}
	moreLines := e.innerExp.getPrettyPrintLines()
	lines = append(lines, moreLines...)
	lines = append(lines, doLeftIndent())
	lines = append(lines, ")")
	return lines
}

/////////////////////////////////////////////////////////////////////////////////

func (e *Binary_Expression) getPrettyPrintLines() []string {
	lines := []string{"BINARY_EXPRESSION_" + getPrettyPrintBinary(e.binOp) + "(", doRightIndent()}
	moreLines := e.firstExp.getPrettyPrintLines()
	moreLines[len(moreLines)-1] = moreLines[len(moreLines)-1] + ","
	lines = append(lines, moreLines...)
	moreLines = e.secExp.getPrettyPrintLines()
	lines = append(lines, moreLines...)
	lines = append(lines, doLeftIndent())
	lines = append(lines, ")")
	return lines
}

/////////////////////////////////////////////////////////////////////////////////

func (e *Assignment_Expression) getPrettyPrintLines() []string {
	lines := []string{"ASSIGNMENT_EXPRESSION(", doRightIndent()}
	lines = append(lines, "lvalue=")
	moreLines := e.lvalue.getPrettyPrintLines()
	moreLines[len(moreLines)-1] = moreLines[len(moreLines)-1] + ","
	lines = append(lines, moreLines...)
	lines = append(lines, "rightExp=")
	moreLines = e.rightExp.getPrettyPrintLines()
	lines = append(lines, moreLines...)
	lines = append(lines, doLeftIndent())
	lines = append(lines, ")")
	return lines
}

/////////////////////////////////////////////////////////////////////////////////

func (e *Conditional_Expression) getPrettyPrintLines() []string {
	lines := []string{"CONDITIONAL_EXPRESSION(", doRightIndent(), "condition="}
	moreLines := e.condition.getPrettyPrintLines()
	moreLines[len(moreLines)-1] = moreLines[len(moreLines)-1] + ","
	lines = append(lines, moreLines...)
	lines = append(lines, "middleExp=")
	moreLines = e.middleExp.getPrettyPrintLines()
	moreLines[len(moreLines)-1] = moreLines[len(moreLines)-1] + ","
	lines = append(lines, moreLines...)
	lines = append(lines, "rightExp=")
	moreLines = e.rightExp.getPrettyPrintLines()
	lines = append(lines, moreLines...)
	lines = append(lines, doLeftIndent())
	lines = append(lines, ")")
	return lines
}

/////////////////////////////////////////////////////////////////////////////////

func (e *Function_Call_Expression) getPrettyPrintLines() []string {
	lines := []string{"FUNCTION_CALL(", doRightIndent(), "functionName=", e.functionName, ",", "args="}

	for _, arg := range e.args {
		moreLines := arg.getPrettyPrintLines()
		moreLines[len(moreLines)-1] = moreLines[len(moreLines)-1] + ","
		lines = append(lines, moreLines...)
	}

	lines = append(lines, doLeftIndent())
	lines = append(lines, ")")

	return lines
}

/////////////////////////////////////////////////////////////////////////////////

func (e *Dereference_Expression) getPrettyPrintLines() []string {
	// TODO:
	return []string{}
}

/////////////////////////////////////////////////////////////////////////////////

func (e *Address_Of_Expression) getPrettyPrintLines() []string {
	// TODO:
	return []string{}
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
		fail("unknown Unary operator")
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
		fail("Unknown binary operator")
	}

	return ""
}

/////////////////////////////////////////////////////////////////////////////////

func getPrettyPrintStorageClass(stCl StorageClassEnum) string {
	switch stCl {
	case NONE_STORAGE_CLASS:
		return "NONE"
	case STATIC_STORAGE_CLASS:
		return "STATIC"
	case EXTERN_STORAGE_CLASS:
		return "EXTERN"
	}
	return ""
}

/////////////////////////////////////////////////////////////////////////////////

func getPrettyPrintDataType(typ DataTypeEnum) string {
	switch typ {
	case INT_TYPE:
		return "INT"
	case LONG_TYPE:
		return "LONG"
	default:
		return ""
	}
}
