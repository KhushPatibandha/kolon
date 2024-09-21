package lexer

import (
	"fmt"
	"math"
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

			{regexp.MustCompile(`print\s*\((.*?)\)`), printHandler},
			{regexp.MustCompile(`println\s*\((.*?)\)`), printlnHandler},

			{regexp.MustCompile(`[a-zA-Z_][a-zA-Z0-9_]*`), identifierHandler},

			{regexp.MustCompile(`[0-9]+[lL]`), intLongHandler(LONG)},
			{regexp.MustCompile(`[0-9]+\.[0-9]+[dD]`), floatDoubleHandler(DOUBLE)},
			{regexp.MustCompile(`[0-9]+\.[0-9]+[eE][+-]?[0-9]+`), floatDoubleHandler(FLOAT)},
			{regexp.MustCompile(`[0-9]+\.[0-9]+`), floatDoubleHandler(FLOAT)},
			{regexp.MustCompile(`[0-9]+`), intLongHandler(INT)},

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
			{regexp.MustCompile(`\.\.`), defaultHandler(DOT_DOT, "..")},
			{regexp.MustCompile(`\.`), defaultHandler(DOT, ".")},
			{regexp.MustCompile(`,`), defaultHandler(COMMA, ",")},
			{regexp.MustCompile(`\?`), defaultHandler(QUESTION_MARK, "?")},

			// {regexp.MustCompile(`[eE]`), defaultHandler(EXPONENT, "e")},
		},
	}
}

func skipHandler(Lexer *Lexer, regex *regexp.Regexp) {
	Lexer.advanceN(regex.FindStringIndex(Lexer.remainder())[1])
}

func defaultHandler(k TokenKind, v string) regexHandler {
	return func(lex *Lexer, regex *regexp.Regexp) {
		lex.push(GetNewToken(k, v))
		lex.advanceN(len(v))
	}
}

func printHandler(lex *Lexer, regex *regexp.Regexp) {
	match := regex.FindStringSubmatch(lex.remainder())
	lex.push(GetNewToken(PRINT, match[1]))
	lex.advanceN(len(match[0]))
}

func printlnHandler(lex *Lexer, regex *regexp.Regexp) {
	match := regex.FindStringSubmatch(lex.remainder())
	lex.push(GetNewToken(PRINTLN, match[1]))
	lex.advanceN(len(match[0]))
}

func floatDoubleHandler(_ TokenKind) regexHandler {
	return func(lex *Lexer, regex *regexp.Regexp) {
		string := regex.FindString(lex.remainder())
		if string[len(string)-1] == 'd' || string[len(string)-1] == 'D' {
			lex.push(GetNewToken(DOUBLE, string))
		} else {
			value, err := strconv.ParseFloat(string, 64)
			if err != nil {
				panic(fmt.Sprintf("Number Handler Error: %v", err))
			}

			if value <= math.MaxFloat32 && value >= -math.MaxFloat32 {
				lex.push(GetNewToken(FLOAT, string))
			} else {
				lex.push(GetNewToken(DOUBLE, string))
			}
		}
		lex.advanceN(len(string))
	}
}

func intLongHandler(k TokenKind) regexHandler {
	return func(lex *Lexer, regex *regexp.Regexp) {
		string := regex.FindString(lex.remainder())
		if string[len(string)-1] == 'l' || string[len(string)-1] == 'L' {
			lex.push(GetNewToken(LONG, string))
		} else {
			number, err := strconv.ParseInt(string, 10, 64)
			if err != nil {
				panic(fmt.Sprintf("Number Handler Error: %v", err))
			}

			if number <= math.MaxInt32 && number >= math.MinInt32 {
				lex.push(GetNewToken(k, string))
			} else {
				lex.push(GetNewToken(LONG, string))
			}
		}
		lex.advanceN(len(string))
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
	var kind, ok = reservedWords[value]
	if ok {
		lex.push(GetNewToken(kind, value))
	} else {
		lex.push(GetNewToken(IDENTIFIER, value))
	}
	lex.advanceN(len(value))
}

func (Lexer *Lexer) advanceN(n int) {
	Lexer.position += n
}

func (Lexer *Lexer) remainder() string {
	return Lexer.source[Lexer.position:]
}

func (Lexer *Lexer) push(token Token) {
	Lexer.Tokens = append(Lexer.Tokens, token)
}

func (Lexer *Lexer) atEOF() bool {
	return Lexer.position >= len(Lexer.source)
}
