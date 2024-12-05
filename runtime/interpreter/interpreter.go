package interpreter

import (
	"fmt"
	"strconv"

	"github.com/voidwyrm-2/ferrite/lexer"
	"github.com/voidwyrm-2/ferrite/runtime/dict"
	"github.com/voidwyrm-2/ferrite/runtime/stack"
	"github.com/voidwyrm-2/ferrite/tokens"
)

type Interpreter struct {
	stack stack.FerriteStack
	dict  dict.FerriteDict
}

func New() Interpreter {
	return Interpreter{stack: stack.FerriteStack{}, dict: dict.New(dict.StdDict)}
}

func CustomNew(dict dict.FerriteDict) Interpreter {
	return Interpreter{stack: stack.FerriteStack{}, dict: dict}
}

func (i *Interpreter) InterpretTokens(tokns []tokens.Token) error {
	var register any = nil

	for index := 0; index < len(tokns); index++ {
		tok := tokns[index]
		if tok.Ist(tokens.WORD) {
			if w, err := i.dict.GetWord(tok.Lit()); err != nil {
				if err.Error() == "BYE" {
					return err
				}
				return tok.E(fmt.Sprintf("%s\n>>>%s<<<", err.Error(), tok.Lit()))
			} else if err = w.Run(i.Interpret, &index, &(i.stack), &register, tokns[index:]); err != nil {
				if err.Error() == "BYE" {
					return err
				}
				return tok.E(fmt.Sprintf("%s\n>>>%s<<<", err.Error(), tok.Lit()))
			}
		} else if tok.Ist(tokens.STRING) {
			i.stack.Push(tok.Lit())
		} else if tok.Ist(tokens.CHAR_LITERAL) {
			i.stack.Push(tok.Lit()[0])
		} else if tok.Ist(tokens.NUMBER) {
			f, _ := strconv.ParseFloat(tok.Lit(), 32)
			i.stack.Push(float32(f))
		} else if !tok.Ist(tokens.COMMENT) {
			return tok.E(fmt.Sprintf("unexpected token '%s'", tok.Lit()))
		}
	}

	return nil
}

func (i *Interpreter) Interpret(text string) error {
	l := lexer.New(text)
	if tokens, err := l.Lex(); err != nil {
		return err
	} else {
		return i.InterpretTokens(tokens)
	}
}
