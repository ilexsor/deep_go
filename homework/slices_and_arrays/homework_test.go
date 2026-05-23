package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type signedInt interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type circularQueue[T signedInt] struct {
	values   []T
	head     int
	tail     int
	size     int
	capacity int
}

// NewCircularQueue create new circularQueue.
func NewCircularQueue[T signedInt](size int) *circularQueue[T] {
	return &circularQueue[T]{
		values:   make([]T, size),
		capacity: size,
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
	q.tail = (q.tail + 1) % q.capacity
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
	q.head = (q.head + 1) % q.capacity
	q.size--

	return true
}

// Front returns the value at the front of the queue without removing it.
// If the queue is empty, it returns the -1 value of type T.
func (q *circularQueue[T]) Front() T {
	if q.Empty() {
		var n T = -1

		return n
	}
	return q.values[q.head]
}

// Back returns the value at the end of the queue without removing it.
// If the queue is empty, it returns the -1 value of type T.
func (q *circularQueue[T]) Back() T {
	if q.Empty() {
		var n T = -1

		return n
	}
	lastIdx := (q.tail - 1 + q.capacity) % q.capacity

	return q.values[lastIdx]
}

// Empty reports whether the queue contains no elements.
func (q *circularQueue[T]) Empty() bool {
	return q.size == 0
}

// Full reports whether the queue has reached its maximum buffer capacity.
func (q *circularQueue[T]) Full() bool {
	return q.size == q.capacity
}

func TestCircularQueue(t *testing.T) {
	const queueSize = 3
	queue := NewCircularQueue[int](queueSize)

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())

	assert.Equal(t, -1, queue.Front())
	assert.Equal(t, -1, queue.Back())
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
