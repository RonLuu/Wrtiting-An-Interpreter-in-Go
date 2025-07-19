package token

var Keywords = map[string]TokenType{
	"fn":  FUNCTION,
	"let": LET,
}

func LookUpVariable(variable string) TokenType {
	if tok, ok := Keywords[variable]; ok {
		return tok
	}
	return VARIABLE
}

// Alias for string
type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

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
)
