package main

import (
	"regexp"
	"strings"
)

/////////////////////////////////////////////////////////////////////////////////

// For info on regex syntax:  https://github.com/google/re2/wiki/Syntax
// * is 0 or more, \b is word boundary, use raw strings to make it easier
// + is 1 or more, | is or, ? means it's optional
// \w is the same as [0-9A-Za-z_]
// \W is the same as [^0-9A-Za-z_] where ^ means not
// \b means at ASCII word boundary, it won't capture the character that creates this word boundary (like a space),
// \b means a transition from \w to \W for example
var regexp_identifier *regexp.Regexp = regexp.MustCompile(`[a-zA-Z_][0-9A-Za-z_]*\b`)
var regexp_int_constant *regexp.Regexp = regexp.MustCompile(`([0-9]+)[^\w.]`)
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
var regexp_if_keyword *regexp.Regexp = regexp.MustCompile(`if\b`)
var regexp_else_keyword *regexp.Regexp = regexp.MustCompile(`else\b`)
var regexp_question *regexp.Regexp = regexp.MustCompile(`\?`)
var regexp_colon *regexp.Regexp = regexp.MustCompile(`:`)
var regexp_do_keyword *regexp.Regexp = regexp.MustCompile(`do\b`)
var regexp_while_keyword *regexp.Regexp = regexp.MustCompile(`while\b`)
var regexp_for_keyword *regexp.Regexp = regexp.MustCompile(`for\b`)
var regexp_break_keyword *regexp.Regexp = regexp.MustCompile(`break\b`)
var regexp_continue_keyword *regexp.Regexp = regexp.MustCompile(`continue\b`)
var regexp_comma *regexp.Regexp = regexp.MustCompile(`,`)
var regexp_static_keyword *regexp.Regexp = regexp.MustCompile(`static\b`)
var regexp_extern_keyword *regexp.Regexp = regexp.MustCompile(`extern\b`)
var regexp_long_keyword *regexp.Regexp = regexp.MustCompile(`long\b`)
var regexp_long_constant *regexp.Regexp = regexp.MustCompile(`([0-9]+[lL])[^\w.]`)
var regexp_signed_keyword *regexp.Regexp = regexp.MustCompile(`signed\b`)
var regexp_unsigned_keyword *regexp.Regexp = regexp.MustCompile(`unsigned\b`)
var regexp_unsigned_int_constant *regexp.Regexp = regexp.MustCompile(`([0-9]+[uU])[^\w.]`)
var regexp_unsigned_long_constant *regexp.Regexp = regexp.MustCompile(`([0-9]+([lL][uU]|[uU][lL]))[^\w.]`)
var regexp_double_keyword *regexp.Regexp = regexp.MustCompile(`double\b`)
var regexp_double_constant *regexp.Regexp = regexp.MustCompile(`(([0-9]*\.[0-9]+|[0-9]+\.?)[Ee][+-]?[0-9]+|[0-9]*\.[0-9]+|[0-9]+\.)[^\w.]`)

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
	COMMA_TOKEN
	STATIC_KEYWORD_TOKEN
	EXTERN_KEYWORD_TOKEN
	LONG_KEYWORD_TOKEN
	LONG_CONSTANT_TOKEN
	SIGNED_KEYWORD_TOKEN
	UNSIGNED_KEYWORD_TOKEN
	UNSIGNED_INT_CONSTANT_TOKEN
	UNSIGNED_LONG_CONSTANT_TOKEN
	DOUBLE_KEYWORD_TOKEN
	DOUBLE_CONSTANT_TOKEN
)

/////////////////////////////////////////////////////////////////////////////////

var allRegexp = map[TokenEnum]*regexp.Regexp{
	IDENTIFIER_TOKEN:             regexp_identifier,
	INT_CONSTANT_TOKEN:           regexp_int_constant,
	INT_KEYWORD_TOKEN:            regexp_int_keyword,
	VOID_KEYWORD_TOKEN:           regexp_void_keyword,
	RETURN_KEYWORD_TOKEN:         regexp_return_keyword,
	OPEN_PARENTHESIS_TOKEN:       regexp_open_parenthesis,
	CLOSE_PARENTHESIS_TOKEN:      regexp_close_parenthesis,
	OPEN_BRACE_TOKEN:             regexp_open_brace,
	CLOSE_BRACE_TOKEN:            regexp_close_brace,
	SEMICOLON_TOKEN:              regexp_semicolon,
	TILDE_TOKEN:                  regexp_tilde,
	HYPHEN_TOKEN:                 regexp_hyphen,
	TWO_HYPHENS_TOKEN:            regexp_two_hyphens,
	PLUS_TOKEN:                   regexp_plus,
	ASTERISK_TOKEN:               regexp_asterisk,
	FORWARD_SLASH_TOKEN:          regexp_forward_slash,
	PERCENT_TOKEN:                regexp_percent,
	EXCLAMATION_TOKEN:            regexp_exclamation,
	TWO_AMPERSANDS_TOKEN:         regexp_two_ampersands,
	TWO_VERTICAL_BARS_TOKEN:      regexp_two_vertical_bars,
	TWO_EQUAL_SIGNS_TOKEN:        regexp_two_equal_signs,
	EXCLAMATION_EQUAL_TOKEN:      regexp_exclamation_equal,
	LESS_THAN_TOKEN:              regexp_less_than,
	GREATER_THAN_TOKEN:           regexp_greater_than,
	LESS_OR_EQUAL_TOKEN:          regexp_less_or_equal,
	GREATER_OR_EQUAL_TOKEN:       regexp_greater_or_equal,
	EQUAL_TOKEN:                  regexp_equal,
	IF_KEYWORD_TOKEN:             regexp_if_keyword,
	ELSE_KEYWORD_TOKEN:           regexp_else_keyword,
	QUESTION_TOKEN:               regexp_question,
	COLON_TOKEN:                  regexp_colon,
	DO_KEYWORD_TOKEN:             regexp_do_keyword,
	WHILE_KEYWORD_TOKEN:          regexp_while_keyword,
	FOR_KEYWORD_TOKEN:            regexp_for_keyword,
	BREAK_KEYWORD_TOKEN:          regexp_break_keyword,
	CONTINUE_KEYWORD_TOKEN:       regexp_continue_keyword,
	COMMA_TOKEN:                  regexp_comma,
	STATIC_KEYWORD_TOKEN:         regexp_static_keyword,
	EXTERN_KEYWORD_TOKEN:         regexp_extern_keyword,
	LONG_KEYWORD_TOKEN:           regexp_long_keyword,
	LONG_CONSTANT_TOKEN:          regexp_long_constant,
	SIGNED_KEYWORD_TOKEN:         regexp_signed_keyword,
	UNSIGNED_KEYWORD_TOKEN:       regexp_unsigned_keyword,
	UNSIGNED_INT_CONSTANT_TOKEN:  regexp_unsigned_int_constant,
	UNSIGNED_LONG_CONSTANT_TOKEN: regexp_unsigned_long_constant,
	DOUBLE_KEYWORD_TOKEN:         regexp_double_keyword,
	DOUBLE_CONSTANT_TOKEN:        regexp_double_constant,
}

var allKeywordRegexp = map[TokenEnum]*regexp.Regexp{
	INT_KEYWORD_TOKEN:      regexp_int_keyword,
	VOID_KEYWORD_TOKEN:     regexp_void_keyword,
	RETURN_KEYWORD_TOKEN:   regexp_return_keyword,
	IF_KEYWORD_TOKEN:       regexp_if_keyword,
	ELSE_KEYWORD_TOKEN:     regexp_else_keyword,
	DO_KEYWORD_TOKEN:       regexp_do_keyword,
	WHILE_KEYWORD_TOKEN:    regexp_while_keyword,
	FOR_KEYWORD_TOKEN:      regexp_for_keyword,
	BREAK_KEYWORD_TOKEN:    regexp_break_keyword,
	CONTINUE_KEYWORD_TOKEN: regexp_continue_keyword,
	STATIC_KEYWORD_TOKEN:   regexp_static_keyword,
	EXTERN_KEYWORD_TOKEN:   regexp_extern_keyword,
	LONG_KEYWORD_TOKEN:     regexp_long_keyword,
	SIGNED_KEYWORD_TOKEN:   regexp_signed_keyword,
	UNSIGNED_KEYWORD_TOKEN: regexp_unsigned_keyword,
	DOUBLE_KEYWORD_TOKEN:   regexp_double_keyword,
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
		fail("some data could not be tokenized:", finalContents)
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
	newContents = newContents[end:]

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
	groupStart, groupEnd := 0, 0
	var detectedTokenType TokenEnum = NONE_TOKEN

	for enum, re := range allRegexp {
		result := re.FindStringSubmatchIndex(contents)

		// if it found something and if the first index is at the beginning of the string
		if result != nil && result[0] == 0 {
			// save it if the word is longer than what we've found before
			if (result[1] - result[0]) > (end - start) {
				detectedTokenType = enum
				start = result[0]
				end = result[1]

				// save group 1 if it was found
				if len(result) > 2 {
					groupStart = result[2]
					groupEnd = result[3]
				}
			}
		}
	}

	// send back the group if we found it
	if groupEnd > 0 {
		return detectedTokenType, groupStart, groupEnd
	} else {
		return detectedTokenType, start, end
	}
}

/////////////////////////////////////////////////////////////////////////////////

func keywordOrIdentifier(word string) TokenEnum {
	// This function assumes the word has already been labeled an IDENTIFIER.
	// Loop through all keyword regexp's to see if it matches a keyword.
	// But the regex must match at the beginning.
	for enum, re := range allKeywordRegexp {
		result := re.FindStringIndex(word)
		// if it found a match and if the first index is at the beginning of the string
		if result != nil && result[0] == 0 {
			return enum
		}
	}
	return IDENTIFIER_TOKEN
}

/////////////////////////////////////////////////////////////////////////////////
