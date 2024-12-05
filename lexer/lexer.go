package lexer

import (
	"fmt"
	"strconv"
	"unicode"

	"github.com/voidwyrm-2/ferrite/tokens"
)

type Lexer struct {
	text         string
	idx, col, ln int
	ch           rune
}

func New(text string) Lexer {
	l := Lexer{text: text, idx: -1, col: 0, ln: 1}
	l.advance()
	return l
}

func (l *Lexer) advance() {
	l.idx++
	l.col++
	if l.idx < len(l.text) {
		l.ch = rune(l.text[l.idx])
	} else {
		l.ch = -1
	}

	if l.ch == '\n' {
		l.ln++
		l.col = 1
	}
}

func (l Lexer) e(format string, a ...any) error {
	_a := []any{l.ln, l.col}
	_a = append(_a, a...)
	return fmt.Errorf("error on line %d, col %d: "+format, _a...)
}

func (l *Lexer) collectWord() (tokens.Token, error) {
	start := l.col
	origLn := l.ln
	s := []byte{}

	for l.ch != -1 && !unicode.IsSpace(l.ch) {
		s = append(s, byte(l.ch))
		l.advance()
	}

	if _, err := strconv.ParseFloat(string(s), 32); err == nil {
		return tokens.New(tokens.NUMBER, string(s), start, l.col-1, origLn), nil
	}

	return tokens.New(tokens.WORD, string(s), start, l.col-1, l.ln), nil
}

func (l *Lexer) collectString(stringType int) (tokens.Token, error) {
	start := l.col
	origLn := l.ln
	s := []byte{}
	escaped := false
	term := func(c rune) bool {
		if stringType == 1 {
			return c == '\''
		} else if stringType == 2 {
			return c == '`'
		}
		return c == '"'
	}

	l.advance()

	for l.ch != -1 && !term(l.ch) {
		if escaped && stringType != 2 {
			switch l.ch {
			case '\\', '"', '\'':
				s = append(s, byte(l.ch))
			case 'n':
				s = append(s, '\n')
			case 't':
				s = append(s, '\t')
			case '0':
				s = append(s, 0)
			default:
				return tokens.Token{}, l.e("illegal escaped character '%c'", l.ch)
			}
			escaped = false
		} else if l.ch == '\\' && stringType != 2 {
			escaped = true
		} else {
			s = append(s, byte(l.ch))
		}

		if len(s) > 1 && stringType == 1 {
			return tokens.Token{}, l.e("illegal character literal")
		}

		l.advance()
	}

	if len(s) == 0 {
		return tokens.Token{}, l.e("empty character literal")
	}

	l.advance()

	if stringType == 1 {
		return tokens.New(tokens.CHAR_LITERAL, string(s), start, l.col-1, origLn), nil
	}
	return tokens.New(tokens.STRING, string(s), start, l.col-1, origLn), nil
}

func (l *Lexer) Lex() ([]tokens.Token, error) {
	tokensList := []tokens.Token{}

	for l.ch != -1 {
		if unicode.IsSpace(l.ch) {
			l.advance()
		} else if l.ch == '#' {
			for l.ch != -1 && l.ch != '\n' {
				l.advance()
			}
		} else if l.ch == '(' {
			start := l.col
			startln := l.ln
			s := ""
			for l.ch != -1 && l.ch != ')' {
				s += string(l.ch)
				l.advance()
			}
			if l.ch != ')' {
				l.col = start
				l.ln = startln
				return []tokens.Token{}, l.e("unterminated multi-line comment")
			}
			tokensList = append(tokensList, tokens.New(tokens.COMMENT, s, start, l.col, startln))
			l.advance()
		} else if l.ch == '"' {
			t, err := l.collectString(0)
			if err != nil {
				return []tokens.Token{}, err
			}
			tokensList = append(tokensList, t)
		} else if l.ch == '\'' {
			t, err := l.collectString(1)
			if err != nil {
				return []tokens.Token{}, err
			}
			tokensList = append(tokensList, t)

		} else if l.ch == '`' {
			t, err := l.collectString(2)
			if err != nil {
				return []tokens.Token{}, err
			}
			tokensList = append(tokensList, t)

		} else if !unicode.IsSpace(l.ch) {
			t, err := l.collectWord()
			if err != nil {
				return []tokens.Token{}, err
			}
			tokensList = append(tokensList, t)
		} else {
			return []tokens.Token{}, l.e("illegal character '%c'", l.ch)
		}
	}

	return tokensList, nil
}
