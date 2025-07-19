package lexer

import (
	"Chapter_1/token"
	// "fmt"
)

type Lexer struct {
	input     string
	curIndex  int
	nextIndex int
	curChar   byte
}

// A function a new lexer
func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.ReadChar()
	return l
}

// A function a new token
func NewToken(tokenType token.TokenType, tokenLiteral byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(tokenLiteral)}
}

// A function to read the current character of a lexer and move on
func (l *Lexer) ReadChar() {
	if l.nextIndex >= len(l.input) {
		l.curChar = 0
	} else {
		l.curChar = l.input[l.nextIndex]
	}
	l.curIndex = l.nextIndex
	l.nextIndex += 1
}

func (l *Lexer) peekChar() byte {
	if l.nextIndex >= len(l.input) {
		return 0
	} else {
		return l.input[l.nextIndex]
	}
}

// A function to
// return the token of the current character of a lexer
// advance to the next character
func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.SkipWhiteSpace()
	switch l.curChar {
	case '+':
		tok = NewToken(token.PLUS, l.curChar)
	case '-':
		tok = NewToken(token.MINUS, l.curChar)
	case '*':
		tok = NewToken(token.MULT, l.curChar)
	case '/':
		tok = NewToken(token.DIV, l.curChar)
	case '=':
		if l.peekChar() == '=' {
			firstEq := l.curChar
			l.ReadChar()
			secondEq := l.curChar
			tok = token.Token{Type: token.EQ, Literal: string(firstEq) + string(secondEq)}
		} else {
			tok = NewToken(token.ASSIGN, l.curChar)
		}
	case '<':
		tok = NewToken(token.LT, l.curChar)
	case '>':
		tok = NewToken(token.GT, l.curChar)
	case '!':
		if l.peekChar() == '=' {
			firstNot := l.curChar
			l.ReadChar()
			secondEq := l.curChar
			tok = token.Token{Type: token.NEQ, Literal: string(firstNot) + string(secondEq)}
		} else {
			tok = NewToken(token.EXCLAMATION, l.curChar)
		}
	case ',':
		tok = NewToken(token.COMMA, l.curChar)
	case ';':
		tok = NewToken(token.SEMICOLON, l.curChar)
	case '[':
		tok = NewToken(token.SLBRACKET, l.curChar)
	case ']':
		tok = NewToken(token.SRBRACKET, l.curChar)
	case '(':
		tok = NewToken(token.RLBRACKET, l.curChar)
	case ')':
		tok = NewToken(token.RRBRACKET, l.curChar)
	case '{':
		tok = NewToken(token.PLBRACKET, l.curChar)
	case '}':
		tok = NewToken(token.PRBRACKET, l.curChar)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if IsLetter(l.curChar) {
			tok.Literal = l.ReadVaraible()
			tok.Type = token.LookUpVariable(tok.Literal)
			return tok
		} else if IsDigit(l.curChar) {
			tok.Literal = l.ReadNumber()
			tok.Type = token.INT
			return tok
		} else {
			tok = NewToken(token.ILLEGAL, l.curChar)
		}
	}
	l.ReadChar()
	return tok
}

func (l *Lexer) SkipWhiteSpace() {
	for l.curChar == ' ' || l.curChar == '\t' || l.curChar == '\n' || l.curChar == '\r' {
		l.ReadChar()
	}
}

func IsLetter(curChar byte) bool {
	return ('a' <= curChar && curChar <= 'z') || ('A' <= curChar && curChar <= 'Z') || curChar == '_'
}

func (l *Lexer) ReadVaraible() string {
	startIndex := l.curIndex
	for IsLetter(l.curChar) {
		l.ReadChar()
	}
	return l.input[startIndex:l.curIndex]
}

func IsDigit(curChar byte) bool {
	return '0' <= curChar && curChar <= '9'
}

func (l *Lexer) ReadNumber() string {
	startIndex := l.curIndex
	for IsDigit(l.curChar) {
		l.ReadChar()
	}
	return l.input[startIndex:l.curIndex]
}
