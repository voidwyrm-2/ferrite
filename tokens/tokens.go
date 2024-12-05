package tokens

import "fmt"

type TokenType uint8

const (
	NONE TokenType = iota
	COMMENT
	STRING
	CHAR_LITERAL
	NUMBER
	WORD
)

type Token struct {
	_type          TokenType
	lit            string
	start, end, ln int
}

func New(_type TokenType, lit string, start, end, ln int) Token {
	return Token{_type: _type, lit: lit, start: start, end: end, ln: ln}
}

func Newtl(_type TokenType, lit string) Token {
	return New(_type, lit, -1, -1, -1)
}

func (t Token) Istl(_type TokenType, lit string) bool {
	return t.Ist(_type) && t.Isl(lit)
}

func (t Token) Ist(_type TokenType) bool {
	return t._type == _type
}

func (t Token) Isl(lit string) bool {
	return t.lit == lit
}

func (t Token) Lit() string {
	return t.lit
}

func (t Token) E(msg string) error {
	return fmt.Errorf("error on line %d, col %d: %s", t.ln, t.start, msg)
}

func (t Token) Str() string {
	return fmt.Sprintf("Token{%s, '%s', %d..%d, %d}", []string{"NONE", "COMMENT", "STRING", "CHAR_LITERAL", "NUMBER", "WORD"}[t._type], t.lit, t.start, t.end, t.ln)
}
