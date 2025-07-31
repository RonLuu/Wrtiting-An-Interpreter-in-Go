package token

import "fmt"

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers
	VARIABLE = "VAR"
	INT      = "INT"

	// Operators
	PLUS   = "+"
	MINUS  = "-"
	MULT   = "*"
	DIV    = "/"
	ASSIGN = "="
	LT     = "<"
	GT     = ">"
	EQ     = "=="
	NEQ    = "!="

	EXCLAMATION = "!"

	// Delimiter
	COMMA     = ","
	SEMICOLON = ";"

	// Bracket
	SLBRACKET = "["
	SRBRACKET = "]"
	RLBRACKET = "("
	RRBRACKET = ")"
	PLBRACKET = "{"
	PRBRACKET = "}"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	IF       = "IF"
	ELSE     = "ELSE"
	IFElSE   = "IFELSE"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	RETURN   = "RETURN"
)

// Alias for string
type TokenType string

type Token struct {
	Type    TokenType // A token contains a TokenType type (what it represents)
	Literal string    // A token contains the string representing it
}

// A map of keywords to its tokentype
var Keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"if":     IF,
	"else":   ELSE,
	"true":   TRUE,
	"false":  FALSE,
	"return": RETURN,
}

// A function to print the token type and its literal
func (token *Token) PrintToken() {
	fmt.Printf("Token:\nType: %s\nLiteral: %s\n", token.Type, token.Literal)
}

// A function to look up all the current keyword and return tokentype
// return token VARIABLE if not a keyword
func LookUpKeyword(word string) TokenType {
	if tok, ok := Keywords[word]; ok {
		return tok
	}
	return VARIABLE
}
