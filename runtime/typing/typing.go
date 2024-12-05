package typing

type Number interface {
	int32 | float32
}

type Addable interface {
	int32 | float32 | string
}

type Subtractable interface {
	int32 | float32
}

type Multipliable interface {
	int32 | float32 | string
}

type Divisable interface {
	int32 | float32
}

type DifferenceComparable interface {
	int32 | float32
}
