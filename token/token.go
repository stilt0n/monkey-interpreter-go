package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

var keywords = map[string]TokenType{
	"fn":  FUNCTION,
	"let": LET,
}

func LookupIdent(identifierLiteral string) TokenType {
	if tok, ok := keywords[identifierLiteral]; ok {
		return tok
	}
	return IDENT
}

const (
	// Invalid syntax and end of file
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
	// Identifiers and literals
	IDENT = "IDENT"
	INT   = "INT"
	// Operators
	ASSIGN = "="
	PLUS   = "+"
	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
)
