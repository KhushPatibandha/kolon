package lexer

import (
	"fmt"
	"regexp"
	"strconv"
)

type regexHandler func(lex *Lexer, regex *regexp.Regexp)

type regexPattern struct {
	regex   *regexp.Regexp
	handler regexHandler
}

type Lexer struct {
	patterns []regexPattern
	Tokens   []Token
	source   string
	position int
}

func Tokenizer(source string) []Token {
	lexer := createLexer(source)
	for !lexer.atEOF() {
		matched := false
		for _, pattern := range lexer.patterns {
			lineOfCode := pattern.regex.FindStringIndex(lexer.remainder())
			if lineOfCode != nil && lineOfCode[0] == 0 {
				pattern.handler(lexer, pattern.regex)
				matched = true
				break
			}
		}
		if !matched {
			panic(fmt.Sprintf("lexer error: unrecognized token '%v' near --> '%v'", lexer.remainder()[:1], lexer.remainder()))
		}
	}
	lexer.push(GetNewToken(EOF, "EOF"))
	return lexer.Tokens
}

func createLexer(source string) *Lexer {
	return &Lexer{
		position: 0,
		source:   source,
		Tokens:   make([]Token, 0),
		patterns: []regexPattern{
			{regexp.MustCompile(`\t+`), skipHandler},
			{regexp.MustCompile(`\s+`), skipHandler},
			{regexp.MustCompile(`\/\/.*`), skipHandler},

			{regexp.MustCompile(`else\s+if`), identifierHandler},

			{regexp.MustCompile(`[a-zA-Z_][a-zA-Z0-9_]*`), identifierHandler},

			{regexp.MustCompile(`\d+\.\d+`), floatHandler(FLOAT)},
			{regexp.MustCompile(`\d+`), intHandler(INT)},

			{regexp.MustCompile(`"[^"]*"`), stringHandler(STRING)},
			{regexp.MustCompile(`'[^']'`), stringHandler(CHAR)},
			{regexp.MustCompile(`true`), defaultHandler(BOOL, "true")},
			{regexp.MustCompile(`false`), defaultHandler(BOOL, "false")},

			{regexp.MustCompile(`\[`), defaultHandler(OPEN_SQUARE_BRACKET, "[")},
			{regexp.MustCompile(`\]`), defaultHandler(CLOSE_SQUARE_BRACKET, "]")},
			{regexp.MustCompile(`\{`), defaultHandler(OPEN_CURLY_BRACKET, "{")},
			{regexp.MustCompile(`\}`), defaultHandler(CLOSE_CURLY_BRACKET, "}")},
			{regexp.MustCompile(`\(`), defaultHandler(OPEN_BRACKET, "(")},
			{regexp.MustCompile(`\)`), defaultHandler(CLOSE_BRACKET, ")")},

			{regexp.MustCompile(`==`), defaultHandler(DOUBLE_EQUAL, "==")},
			{regexp.MustCompile(`!=`), defaultHandler(NOT_EQUAL, "!=")},
			{regexp.MustCompile(`<=`), defaultHandler(LESS_THAN_EQUAL, "<=")},
			{regexp.MustCompile(`>=`), defaultHandler(GREATER_THAN_EQUAL, ">=")},
			{regexp.MustCompile(`<`), defaultHandler(LESS_THAN, "<")},
			{regexp.MustCompile(`>`), defaultHandler(GREATER_THAN, ">")},

			{regexp.MustCompile(`\+\+`), defaultHandler(PLUS_PLUS, "++")},
			{regexp.MustCompile(`\+=`), defaultHandler(PLUS_EQUAL, "+=")},
			{regexp.MustCompile(`\+`), defaultHandler(PLUS, "+")},
			{regexp.MustCompile(`--`), defaultHandler(MINUS_MINUS, "--")},
			{regexp.MustCompile(`-=`), defaultHandler(MINUS_EQUAL, "-=")},
			{regexp.MustCompile(`-`), defaultHandler(DASH, "-")},
			{regexp.MustCompile(`\*=`), defaultHandler(STAR_EQUAL, "*=")},
			{regexp.MustCompile(`\*`), defaultHandler(STAR, "*")},
			{regexp.MustCompile(`/=`), defaultHandler(SLASH_EQUAL, "/=")},
			{regexp.MustCompile(`/`), defaultHandler(SLASH, "/")},
			{regexp.MustCompile(`%=`), defaultHandler(PERCENT_EQUAL, "%=")},
			{regexp.MustCompile(`%`), defaultHandler(PERCENT, "%")},

			{regexp.MustCompile(`&&`), defaultHandler(AND_AND, "&&")},
			{regexp.MustCompile(`\|\|`), defaultHandler(OR_OR, "||")},
			{regexp.MustCompile(`&`), defaultHandler(AND, "&")},
			{regexp.MustCompile(`\|`), defaultHandler(OR, "|")},

			{regexp.MustCompile(`=`), defaultHandler(EQUAL_ASSIGN, "=")},
			{regexp.MustCompile(`!`), defaultHandler(NOT, "!")},
			{regexp.MustCompile(`:`), defaultHandler(COLON, ":")},
			{regexp.MustCompile(`;`), defaultHandler(SEMI_COLON, ";")},
			{regexp.MustCompile(`,`), defaultHandler(COMMA, ",")},
		},
	}
}

func skipHandler(lexer *Lexer, regex *regexp.Regexp) {
	lexer.advanceN(regex.FindStringIndex(lexer.remainder())[1])
}

func defaultHandler(k TokenKind, v string) regexHandler {
	return func(lex *Lexer, regex *regexp.Regexp) {
		lex.push(GetNewToken(k, v))
		lex.advanceN(len(v))
	}
}

func floatHandler(k TokenKind) regexHandler {
	return func(lex *Lexer, regex *regexp.Regexp) {
		matchedString := regex.FindString(lex.remainder())
		_, err := strconv.ParseFloat(matchedString, 64)
		if err != nil {
			panic(fmt.Sprintf("Number Handler Error: %v", err))
		}
		lex.push(GetNewToken(k, matchedString))
		lex.advanceN(len(matchedString))
	}
}

func intHandler(k TokenKind) regexHandler {
	return func(lex *Lexer, regex *regexp.Regexp) {
		matchedString := regex.FindString(lex.remainder())
		if _, err := strconv.ParseInt(matchedString, 10, 64); err == nil {
			lex.push(GetNewToken(k, matchedString))
		} else {
			if _, err := strconv.ParseUint(matchedString, 10, 64); err == nil {
				lex.push(GetNewToken(k, matchedString))
			} else {
				panic(fmt.Sprintf("Number Handler Error: %v", err))
			}
		}
		lex.advanceN(len(matchedString))
	}
}

func stringHandler(k TokenKind) regexHandler {
	return func(lex *Lexer, regex *regexp.Regexp) {
		match := regex.FindString(lex.remainder())
		lex.push(GetNewToken(k, match))
		lex.advanceN(len(match))
	}
}

func identifierHandler(lex *Lexer, regex *regexp.Regexp) {
	value := regex.FindString(lex.remainder())
	kind, ok := reservedWords[value]
	if ok {
		lex.push(GetNewToken(kind, value))
	} else {
		lex.push(GetNewToken(IDENTIFIER, value))
	}
	lex.advanceN(len(value))
}

func (lexer *Lexer) advanceN(n int) {
	lexer.position += n
}

func (lexer *Lexer) remainder() string {
	return lexer.source[lexer.position:]
}

func (lexer *Lexer) push(token Token) {
	lexer.Tokens = append(lexer.Tokens, token)
}

func (lexer *Lexer) atEOF() bool {
	return lexer.position >= len(lexer.source)
}
