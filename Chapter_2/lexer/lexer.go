package lexer

// A lexer is used for reading token

import (
	"Chapter_2/token"
	"fmt"
)

// A lexer contains
type Lexer struct {
	input     string // The current string it's reading
	curIndex  int    // The current index of that string
	nextIndex int    // The next index of that string
	curChar   byte   // The current char of that string
}

// A function to create a new lexer
func NewLexer(input string) *Lexer {
	// Set the current input
	l := &Lexer{input: input}
	// Read the current character
	l.readChar()
	return l
}

// A function to read the current character of a lexer and move on
func (l *Lexer) readChar() {
	// The nextIndex is 'out of bound'
	if l.nextIndex >= len(l.input) {
		// Set the current character to EOF
		l.curChar = 0
	} else {
		// Set the current chacter to the current character
		l.curChar = l.input[l.nextIndex]
	}
	// Update curIndex and nextIndex
	l.curIndex = l.nextIndex
	l.nextIndex += 1
}

// Debug function
// A function to print the lexer
func (lexer *Lexer) PrintLexer() {
	fmt.Printf("Lexer:\ninput: \"%s\"\ncurIndex: %d\nnextIndex: %d\ncurChar: %c\n", lexer.input, lexer.curIndex, lexer.nextIndex, lexer.curChar)
}

// A function to
// return the token of the current character of a lexer
// advance to the next character
func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.skipWhiteSpace()
	// Depending on the current character,
	// decide how to read the token
	switch l.curChar {
	// Single-char string case
	case '+':
		tok = NewToken(token.PLUS, l.curChar)
	case '-':
		tok = NewToken(token.MINUS, l.curChar)
	case '*':
		tok = NewToken(token.MULT, l.curChar)
	case '/':
		tok = NewToken(token.DIV, l.curChar)
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
	case '<':
		tok = NewToken(token.LT, l.curChar)
	case '>':
		tok = NewToken(token.GT, l.curChar)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	// Special case "="
	case '=':
		// Check if the next char is "="
		if l.peekChar() == '=' {
			firstEq := l.curChar
			l.readChar()
			secondEq := l.curChar
			// Set the token to be EQ, "=="
			tok = token.Token{Type: token.EQ, Literal: string(firstEq) + string(secondEq)}
		} else {
			// Set the token to be ASSIGN, "="
			tok = NewToken(token.ASSIGN, l.curChar)
		}
	case '!':
		// Check if the next char is "="
		if l.peekChar() == '=' {
			firstNot := l.curChar
			l.readChar()
			secondEq := l.curChar
			// Set the token to be NEQ, "!="
			tok = token.Token{Type: token.NEQ, Literal: string(firstNot) + string(secondEq)}
		} else {
			// Set the token to be EXCLAMATION, "!"
			tok = NewToken(token.EXCLAMATION, l.curChar)
		}

	// Whole string case
	default:
		// If it's a letter
		if isLetter(l.curChar) {
			// Read the whole word
			tok.Literal = l.readWord()
			// Decide if the token is variable or a keyword
			tok.Type = token.LookUpKeyword(tok.Literal)
			return tok
		} else if isDigit(l.curChar) { // If it's a number
			// Read the whole number
			tok.Literal = l.readNumber()
			tok.Type = token.INT

			return tok
		} else { // If's something really weird
			tok = NewToken(token.ILLEGAL, l.curChar)
		}
	}
	// Move on to the next token
	l.readChar()
	return tok
}

func (l *Lexer) skipWhiteSpace() {
	// While the current character is any white space
	for l.curChar == ' ' || l.curChar == '\t' || l.curChar == '\n' || l.curChar == '\r' {
		l.readChar()
	}
}

// A function to create a new token
func NewToken(tokenType token.TokenType, tokenLiteral byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(tokenLiteral)}
}

func (l *Lexer) peekChar() byte {
	if l.nextIndex >= len(l.input) {
		return 0
	} else {
		return l.input[l.nextIndex]
	}
}

func isLetter(curChar byte) bool {
	return ('a' <= curChar && curChar <= 'z') || ('A' <= curChar && curChar <= 'Z') || curChar == '_'
}

func (l *Lexer) readWord() string {
	startIndex := l.curIndex
	for isLetter(l.curChar) {
		l.readChar()
	}
	return l.input[startIndex:l.curIndex]
}

func isDigit(curChar byte) bool {
	return '0' <= curChar && curChar <= '9'
}

func (l *Lexer) readNumber() string {
	startIndex := l.curIndex
	for isDigit(l.curChar) {
		l.readChar()
	}
	return l.input[startIndex:l.curIndex]
}
