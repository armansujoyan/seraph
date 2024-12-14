package utils

type (
	Stack[T any] struct {
		top    *node[T]
		length int
	}
	node[T any] struct {
		value T
		prev  *node[T]
	}
)

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{nil, 0}
}

func (this *Stack[T]) Len() int {
	return this.length
}

func (this *Stack[T]) Peek() T {
	if this.length == 0 {
    var zero T
		return zero
	}
	return this.top.value
}

func (this *Stack[T]) Pop() T {
	if this.length == 0 {
    var zero T
		return zero
	}

	n := this.top
	this.top = n.prev
	this.length--
	return n.value
}

func (this *Stack[T]) Push(value T) {
	n := &node[T]{value, this.top}
	this.top = n
	this.length++
}
