package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type circularQueue[T any] struct {
	values []T
	head   int
	tail   int
	size   int
}

// NewCircularQueue create new circularQueue.
func NewCircularQueue[T any](size int) *circularQueue[T] {
	return &circularQueue[T]{
		values: make([]T, size),
	}
}

// Push adds a new value to the end of the queue.
// It returns false if the queue is full; otherwise, it inserts the value
// and advances the tail pointer.
func (q *circularQueue[T]) Push(value T) bool {
	if q.Full() {
		return false
	}

	q.values[q.tail] = value
	q.tail = (q.tail + 1) % len(q.values)
	q.size++

	return true
}

// Pop removes the oldest value from the front of the queue.
// It returns false if the queue is empty; otherwise, it removes the value,
// advances the head pointer, and returns true.
func (q *circularQueue[T]) Pop() bool {
	if q.Empty() {
		return false
	}

	var zero T
	q.values[q.head] = zero

	q.head = (q.head + 1) % cap(q.values)
	q.size--

	return true
}

// Front returns the value at the front of the queue without removing it.
// If the queue is empty, it returns the -1 value of type T.
func (q *circularQueue[T]) Front() T {
	if q.Empty() {
		var zero T

		return zero
	}
	return q.values[q.head]
}

// Back returns the value at the end of the queue without removing it.
// If the queue is empty, it returns the -1 value of type T.
func (q *circularQueue[T]) Back() T {
	if q.Empty() {
		var zero T

		return zero
	}
	lastIdx := (q.tail - 1 + len(q.values)) % len(q.values)

	return q.values[lastIdx]
}

// Empty reports whether the queue contains no elements.
func (q *circularQueue[T]) Empty() bool {
	return q.size == 0
}

// Full reports whether the queue has reached its maximum buffer capacity.
func (q *circularQueue[T]) Full() bool {
	return q.size == len(q.values)
}

func TestCircularQueueInt(t *testing.T) {
	const queueSize = 3
	queue := NewCircularQueue[int](queueSize)

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())

	assert.Equal(t, 0, queue.Front())
	assert.Equal(t, 0, queue.Back())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Push(1))
	assert.True(t, queue.Push(2))
	assert.True(t, queue.Push(3))
	assert.False(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int{1, 2, 3}, queue.values))

	assert.False(t, queue.Empty())
	assert.True(t, queue.Full())

	assert.Equal(t, 1, queue.Front())
	assert.Equal(t, 3, queue.Back())

	assert.True(t, queue.Pop())
	assert.False(t, queue.Empty())
	assert.False(t, queue.Full())
	assert.True(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int{4, 2, 3}, queue.values))

	assert.Equal(t, 2, queue.Front())
	assert.Equal(t, 4, queue.Back())

	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())
}

func TestCircularQueueString(t *testing.T) {
	const queueSize = 3
	queue := NewCircularQueue[string](queueSize)

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())

	assert.Equal(t, "", queue.Front())
	assert.Equal(t, "", queue.Back())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Push("A"))
	assert.True(t, queue.Push("B"))
	assert.True(t, queue.Push("C"))
	assert.False(t, queue.Push("D"))

	assert.True(t, reflect.DeepEqual([]string{"A", "B", "C"}, queue.values))

	assert.False(t, queue.Empty())
	assert.True(t, queue.Full())

	assert.Equal(t, "A", queue.Front())
	assert.Equal(t, "C", queue.Back())

	assert.True(t, queue.Pop())
	assert.False(t, queue.Empty())
	assert.False(t, queue.Full())
	assert.True(t, queue.Push("D"))

	assert.True(t, reflect.DeepEqual([]string{"D", "B", "C"}, queue.values))

	assert.Equal(t, "B", queue.Front())
	assert.Equal(t, "D", queue.Back())

	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())
}

func TestCircularQueueFloat(t *testing.T) {
	const queueSize = 3
	queue := NewCircularQueue[float64](queueSize)

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())

	assert.Equal(t, 0.0, queue.Front())
	assert.Equal(t, 0.0, queue.Back())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Push(1.1))
	assert.True(t, queue.Push(1.2))
	assert.True(t, queue.Push(1.3))
	assert.False(t, queue.Push(1.4))

	assert.True(t, reflect.DeepEqual([]float64{1.1, 1.2, 1.3}, queue.values))

	assert.False(t, queue.Empty())
	assert.True(t, queue.Full())

	assert.Equal(t, 1.1, queue.Front())
	assert.Equal(t, 1.3, queue.Back())

	assert.True(t, queue.Pop())
	assert.False(t, queue.Empty())
	assert.False(t, queue.Full())
	assert.True(t, queue.Push(1.4))

	assert.True(t, reflect.DeepEqual([]float64{1.4, 1.2, 1.3}, queue.values))

	assert.Equal(t, 1.2, queue.Front())
	assert.Equal(t, 1.4, queue.Back())

	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())
}
