package dict

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/voidwyrm-2/ferrite/runtime/stack"
	"github.com/voidwyrm-2/ferrite/runtime/word"
	"github.com/voidwyrm-2/ferrite/tokens"
)

type FerriteDict struct {
	dict map[string]word.Word
}

func New(dict map[string]word.Word) FerriteDict {
	return FerriteDict{dict: dict}
}

func (fd *FerriteDict) AddWord(name string, word word.Word) error {
	if _, ok := fd.dict[name]; ok {
		return fmt.Errorf("word '%s' already exists", name)
	}
	fd.dict[name] = word
	return nil
}

func (fd *FerriteDict) GetWord(name string) (word.Word, error) {
	if w, ok := fd.dict[name]; ok {
		return w, nil
	} else {
		return word.Word{}, fmt.Errorf("word '%s' does not exist", name)
	}
}

var StdDict = func() map[string]word.Word {
	d := map[string]word.Word{
		"bye": word.NewBuiltin(func(i *int, stack *stack.FerriteStack, register *any, tokens []tokens.Token) error {
			return errors.New("BYE")
		}, "--"),
		"emit": word.NewBuiltin(func(i *int, stack *stack.FerriteStack, register *any, tokens []tokens.Token) error {
			v, e := stack.Pop()
			if e != nil {
				return e
			}

			if i, ok := v.(float32); ok {
				fmt.Print(string(int32(i)))
			} else {
				return fmt.Errorf("invalid type for emitting '%s'", reflect.TypeOf(v).Name())
			}

			return nil
		}, "i -- "),
		"cr": word.New("10 emit", "--"),
		".": word.NewBuiltin(func(i *int, stack *stack.FerriteStack, register *any, tokens []tokens.Token) error {
			v, e := stack.Pop()
			if e != nil {
				return e
			}

			fmt.Print(v)

			return nil
		}, "v -- "),
	}

	return d
}()
