package stack

import "errors"

type FerriteStack struct {
	stack []any
}

func (fs *FerriteStack) Push(a any) {
	fs.stack = append(fs.stack, a)
}

func (fs *FerriteStack) Pop() (any, error) {
	if len(fs.stack) == 0 {
		return nil, errors.New("stack underflow")
	}
	a := fs.stack[len(fs.stack)-1]

	fs.stack = fs.stack[:len(fs.stack)-1]
	return a, nil
}

func (fs *FerriteStack) Bool(b bool) {
	if b {
		fs.Push(float32(-1))
	} else {
		fs.Push(float32(0))
	}
}

func (fs *FerriteStack) Empty() bool {
	return len(fs.stack) == 0
}

func (fs *FerriteStack) Dup() error {
	if a, err := fs.Pop(); err != nil {
		return err
	} else {
		fs.Push(a)
		fs.Push(a)
		return nil
	}
}
