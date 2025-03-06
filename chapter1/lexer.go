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
}

var allKeywordRegexp = map[TokenEnum]*regexp.Regexp{
	INT_KEYWORD_TOKEN:    regexp_int_keyword,
	VOID_KEYWORD_TOKEN:   regexp_void_keyword,
	RETURN_KEYWORD_TOKEN: regexp_return_keyword,
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
	fmt.Printf("%v %v %v\n", enum, start, end)

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
