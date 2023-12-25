package ast

import (
	"monkey-pl/token"
	"testing"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "x"},
					Value: "x",
				},
				Value: &Identifier{
					Token: token.Token{
						Type:    token.IDENT,
						Literal: "y",
					},
					Value: "y",
				},
			},
		},
	}

	if program.String() != "let x = y;" {
		t.Errorf("program.String() is not correct. Received %q", program.String())
	}
}
