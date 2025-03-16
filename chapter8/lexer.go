package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

/////////////////////////////////////////////////////////////////////////////////

// * is 0 or more, \b is word boundary, use raw strings to make it easier
// + is 1 or more
var regexp_identifier *regexp.Regexp = regexp.MustCompile(`[a-zA-Z_][0-9A-Za-z_]*\b`)
var regexp_int_constant *regexp.Regexp = regexp.MustCompile(`[0-9]+\b`)
var regexp_int_keyword *regexp.Regexp = regexp.MustCompile(`int\b`)
var regexp_void_keyword *regexp.Regexp = regexp.MustCompile(`void\b`)
var regexp_return_keyword *regexp.Regexp = regexp.MustCompile(`return\b`)
var regexp_open_parenthesis *regexp.Regexp = regexp.MustCompile(`\(`)
var regexp_close_parenthesis *regexp.Regexp = regexp.MustCompile(`\)`)
var regexp_open_brace *regexp.Regexp = regexp.MustCompile(`{`)
var regexp_close_brace *regexp.Regexp = regexp.MustCompile(`}`)
var regexp_semicolon *regexp.Regexp = regexp.MustCompile(`;`)
var regexp_tilde *regexp.Regexp = regexp.MustCompile(`~`)
var regexp_hyphen *regexp.Regexp = regexp.MustCompile(`-`)
var regexp_two_hyphens *regexp.Regexp = regexp.MustCompile(`--`)
var regexp_plus *regexp.Regexp = regexp.MustCompile(`\+`)
var regexp_asterisk *regexp.Regexp = regexp.MustCompile(`\*`)
var regexp_forward_slash *regexp.Regexp = regexp.MustCompile(`/`)
var regexp_percent *regexp.Regexp = regexp.MustCompile(`%`)
var regexp_exclamation *regexp.Regexp = regexp.MustCompile(`!`)
var regexp_two_ampersands *regexp.Regexp = regexp.MustCompile(`&&`)
var regexp_two_vertical_bars *regexp.Regexp = regexp.MustCompile(`\|\|`)
var regexp_two_equal_signs *regexp.Regexp = regexp.MustCompile(`==`)
var regexp_exclamation_equal *regexp.Regexp = regexp.MustCompile(`!=`)
var regexp_less_than *regexp.Regexp = regexp.MustCompile(`<`)
var regexp_greater_than *regexp.Regexp = regexp.MustCompile(`>`)
var regexp_less_or_equal *regexp.Regexp = regexp.MustCompile(`<=`)
var regexp_greater_or_equal *regexp.Regexp = regexp.MustCompile(`>=`)
var regexp_equal *regexp.Regexp = regexp.MustCompile(`=`)
var regexp_if *regexp.Regexp = regexp.MustCompile(`if`)
var regexp_else *regexp.Regexp = regexp.MustCompile(`else`)
var regexp_question *regexp.Regexp = regexp.MustCompile(`\?`)
var regexp_colon *regexp.Regexp = regexp.MustCompile(`:`)
var regexp_do *regexp.Regexp = regexp.MustCompile(`do`)
var regexp_while *regexp.Regexp = regexp.MustCompile(`while`)
var regexp_for *regexp.Regexp = regexp.MustCompile(`for`)
var regexp_break *regexp.Regexp = regexp.MustCompile(`break`)
var regexp_continue *regexp.Regexp = regexp.MustCompile(`continue`)

type TokenEnum int

const (
	NONE_TOKEN TokenEnum = iota
	IDENTIFIER_TOKEN
	INT_CONSTANT_TOKEN
	INT_KEYWORD_TOKEN
	VOID_KEYWORD_TOKEN
	RETURN_KEYWORD_TOKEN
	OPEN_PARENTHESIS_TOKEN
	CLOSE_PARENTHESIS_TOKEN
	OPEN_BRACE_TOKEN
	CLOSE_BRACE_TOKEN
	SEMICOLON_TOKEN
	TILDE_TOKEN
	HYPHEN_TOKEN
	TWO_HYPHENS_TOKEN
	PLUS_TOKEN
	ASTERISK_TOKEN
	FORWARD_SLASH_TOKEN
	PERCENT_TOKEN
	EXCLAMATION_TOKEN
	TWO_AMPERSANDS_TOKEN
	TWO_VERTICAL_BARS_TOKEN
	TWO_EQUAL_SIGNS_TOKEN
	EXCLAMATION_EQUAL_TOKEN
	LESS_THAN_TOKEN
	GREATER_THAN_TOKEN
	LESS_OR_EQUAL_TOKEN
	GREATER_OR_EQUAL_TOKEN
	EQUAL_TOKEN
	IF_KEYWORD_TOKEN
	ELSE_KEYWORD_TOKEN
	QUESTION_TOKEN
	COLON_TOKEN
	DO_KEYWORD_TOKEN
	WHILE_KEYWORD_TOKEN
	FOR_KEYWORD_TOKEN
	BREAK_KEYWORD_TOKEN
	CONTINUE_KEYWORD_TOKEN
)

/////////////////////////////////////////////////////////////////////////////////

var allRegexp = map[TokenEnum]*regexp.Regexp{
	IDENTIFIER_TOKEN:        regexp_identifier,
	INT_CONSTANT_TOKEN:      regexp_int_constant,
	INT_KEYWORD_TOKEN:       regexp_int_keyword,
	VOID_KEYWORD_TOKEN:      regexp_void_keyword,
	RETURN_KEYWORD_TOKEN:    regexp_return_keyword,
	OPEN_PARENTHESIS_TOKEN:  regexp_open_parenthesis,
	CLOSE_PARENTHESIS_TOKEN: regexp_close_parenthesis,
	OPEN_BRACE_TOKEN:        regexp_open_brace,
	CLOSE_BRACE_TOKEN:       regexp_close_brace,
	SEMICOLON_TOKEN:         regexp_semicolon,
	TILDE_TOKEN:             regexp_tilde,
	HYPHEN_TOKEN:            regexp_hyphen,
	TWO_HYPHENS_TOKEN:       regexp_two_hyphens,
	PLUS_TOKEN:              regexp_plus,
	ASTERISK_TOKEN:          regexp_asterisk,
	FORWARD_SLASH_TOKEN:     regexp_forward_slash,
	PERCENT_TOKEN:           regexp_percent,
	EXCLAMATION_TOKEN:       regexp_exclamation,
	TWO_AMPERSANDS_TOKEN:    regexp_two_ampersands,
	TWO_VERTICAL_BARS_TOKEN: regexp_two_vertical_bars,
	TWO_EQUAL_SIGNS_TOKEN:   regexp_two_equal_signs,
	EXCLAMATION_EQUAL_TOKEN: regexp_exclamation_equal,
	LESS_THAN_TOKEN:         regexp_less_than,
	GREATER_THAN_TOKEN:      regexp_greater_than,
	LESS_OR_EQUAL_TOKEN:     regexp_less_or_equal,
	GREATER_OR_EQUAL_TOKEN:  regexp_greater_or_equal,
	EQUAL_TOKEN:             regexp_equal,
	IF_KEYWORD_TOKEN:        regexp_if,
	ELSE_KEYWORD_TOKEN:      regexp_else,
	QUESTION_TOKEN:          regexp_question,
	COLON_TOKEN:             regexp_colon,
	DO_KEYWORD_TOKEN:        regexp_do,
	WHILE_KEYWORD_TOKEN:     regexp_while,
	FOR_KEYWORD_TOKEN:       regexp_for,
	BREAK_KEYWORD_TOKEN:     regexp_break,
	CONTINUE_KEYWORD_TOKEN:  regexp_continue,
}

var allKeywordRegexp = map[TokenEnum]*regexp.Regexp{
	INT_KEYWORD_TOKEN:      regexp_int_keyword,
	VOID_KEYWORD_TOKEN:     regexp_void_keyword,
	RETURN_KEYWORD_TOKEN:   regexp_return_keyword,
	IF_KEYWORD_TOKEN:       regexp_if,
	ELSE_KEYWORD_TOKEN:     regexp_else,
	DO_KEYWORD_TOKEN:       regexp_do,
	WHILE_KEYWORD_TOKEN:    regexp_while,
	FOR_KEYWORD_TOKEN:      regexp_for,
	BREAK_KEYWORD_TOKEN:    regexp_break,
	CONTINUE_KEYWORD_TOKEN: regexp_continue,
}

/////////////////////////////////////////////////////////////////////////////////

type Token struct {
	word      string
	tokenType TokenEnum
}

func doLexer(fileContents string) []Token {
	var allTokens []Token

	fileContents, token := getNextToken(fileContents)
	for token.tokenType != NONE_TOKEN {
		allTokens = append(allTokens, token)
		fileContents, token = getNextToken(fileContents)
	}

	// if there is still data in fileContents then we have data that didn't match a regexp, so generate error
	finalContents := strings.TrimLeft(fileContents, " \n\r\t")
	if len(finalContents) > 0 {
		fmt.Println("some data could not be tokenized:", finalContents)
		os.Exit(1)
	}

	return allTokens
}

/////////////////////////////////////////////////////////////////////////////////

func getNextToken(contents string) (newContents string, token Token) {
	// remove whitespace from beginning
	newContents = strings.TrimLeft(contents, " \n\r\t")

	// use the regexp to find the longest match at beginning
	enum, start, end := longestMatchAtStart(newContents)
	//fmt.Printf("%v %v %v\n", enum, start, end)

	// gather the token info
	token = Token{word: newContents[start:end], tokenType: enum}

	// remove the word from the beginning of the string
	newContents = strings.TrimPrefix(newContents, token.word)

	// the regexp for identifiers might also match keywords, keywords should take priority
	// so switch the enum from identifier to keyword if necessary
	if token.tokenType == IDENTIFIER_TOKEN {
		token.tokenType = keywordOrIdentifier(token.word)
	}

	return newContents, token
}

/////////////////////////////////////////////////////////////////////////////////

func longestMatchAtStart(contents string) (TokenEnum, int, int) {
	start, end := 0, 0
	var detectedTokenType TokenEnum = NONE_TOKEN

	for enum, re := range allRegexp {
		result := re.FindStringIndex(contents)

		// if it found something and if it's at the beginning of the string
		if result != nil && result[0] == 0 {
			// save it if the word is longer than what we've found before
			if (result[1] - result[0]) > (end - start) {
				detectedTokenType = enum
				start = result[0]
				end = result[1]
			}
		}
	}

	return detectedTokenType, start, end
}

/////////////////////////////////////////////////////////////////////////////////

func keywordOrIdentifier(word string) TokenEnum {
	// This function assumes the word has already been labeled an IDENTIFIER.
	// Loop through all keyword regexp's to see if it matches a keyword.
	for enum, re := range allKeywordRegexp {
		result := re.FindStringIndex(word)
		if result != nil {
			return enum
		}
	}
	return IDENTIFIER_TOKEN
}

/////////////////////////////////////////////////////////////////////////////////
