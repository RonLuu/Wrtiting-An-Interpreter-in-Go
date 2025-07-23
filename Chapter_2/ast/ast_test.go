package ast

import (
	"Chapter_2/token"
	"testing"
)

func TestSTring(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token:      token.Token{Type: token.LET, Literal: "let"},
				Variable:   &Variable{Token: token.Token{Type: token.VARIABLE, Literal: "myVar"}, Literal: "myVar"},
				Expression: &Variable{Token: token.Token{Type: token.VARIABLE, Literal: "anotherVar"}, Literal: "anotherVar"},
			},
		},
	}

	if program.String() != "let myVar = anotherVar;" {
		t.Errorf("program.String() wrong, got = %q", program.String())
	}
}
