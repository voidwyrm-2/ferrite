package word

import (
	"github.com/voidwyrm-2/ferrite/runtime/stack"
	"github.com/voidwyrm-2/ferrite/tokens"
)

type WordBuiltin func(i *int, stack *stack.FerriteStack, register *any, tokens []tokens.Token) error

type Word struct {
	effect string
	text   *string
	fn     *WordBuiltin
}

func New(text string, effect string) Word {
	return Word{text: &text, fn: nil, effect: effect}
}

func NewBuiltin(fn WordBuiltin, effect string) Word {
	return Word{text: nil, fn: &fn, effect: effect}
}

func (w Word) Run(interpret func(string) error, i *int, stack *stack.FerriteStack, register *any, tokens []tokens.Token) error {
	if w.fn != nil {
		return (*w.fn)(i, stack, register, tokens)
	} else if w.text != nil {
		return interpret(*w.text)
	}
	return nil
}
