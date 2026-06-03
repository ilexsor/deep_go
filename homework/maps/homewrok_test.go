package main

import (
	"cmp"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type node[K cmp.Ordered, V any] struct {
	key   K
	value V
	left  *node[K, V]
	right *node[K, V]
}

type OrderedMap[K cmp.Ordered, V any] struct {
	root *node[K, V]
	size int
}

func NewOrderedMap[K cmp.Ordered, V any]() OrderedMap[K, V] {
	return OrderedMap[K, V]{}
}

func (m *OrderedMap[K, V]) Insert(key K, value V) {
	m.root = m.insertNode(m.root, key, value, &m.size)
}

func (m *OrderedMap[K, V]) insertNode(n *node[K, V], key K, value V, size *int) *node[K, V] {
	if n == nil {
		*size++

		return &node[K, V]{key: key, value: value}
	}

	if key < n.key {
		n.left = m.insertNode(n.left, key, value, size)
	} else if key > n.key {
		n.right = m.insertNode(n.right, key, value, size)
	} else {
		n.value = value
	}

	return n
}

func (m *OrderedMap[K, V]) Erase(key K) {
	var deleted bool
	m.root = m.eraseNode(m.root, key, &deleted)
	if deleted {
		m.size--
	}
}

func (m *OrderedMap[K, V]) eraseNode(n *node[K, V], key K, deleted *bool) *node[K, V] {
	if n == nil {
		return nil
	}

	if key < n.key {
		n.left = m.eraseNode(n.left, key, deleted)
	} else if key > n.key {
		n.right = m.eraseNode(n.right, key, deleted)
	} else {
		*deleted = true
		if n.left == nil {
			return n.right
		} else if n.right == nil {
			return n.left
		}

		minNode := m.findMin(n.right)
		n.key = minNode.key
		n.value = minNode.value

		var dummy bool
		n.right = m.eraseNode(n.right, minNode.key, &dummy)
	}

	return n
}

func (m *OrderedMap[K, V]) findMin(n *node[K, V]) *node[K, V] {
	curr := n
	for curr.left != nil {
		curr = curr.left
	}

	return curr
}

func (m *OrderedMap[K, V]) Contains(key K) bool {
	curr := m.root
	for curr != nil {
		if key == curr.key {
			return true
		} else if key < curr.key {
			curr = curr.left
		} else {
			curr = curr.right
		}
	}

	return false
}

func (m *OrderedMap[K, V]) Size() int {
	return m.size
}

func (m *OrderedMap[K, V]) ForEach(action func(K, V)) {
	m.forEachNode(m.root, action)
}

func (m *OrderedMap[K, V]) forEachNode(n *node[K, V], action func(K, V)) {
	if n == nil {
		return
	}

	m.forEachNode(n.left, action)
	action(n.key, n.value)
	m.forEachNode(n.right, action)
}

func TestCircularQueue(t *testing.T) {
	data := NewOrderedMap[int, int32]()
	assert.Zero(t, data.Size())

	data.Insert(10, 10)
	data.Insert(5, 5)
	data.Insert(15, 15)
	data.Insert(2, 2)
	data.Insert(4, 4)
	data.Insert(12, 12)
	data.Insert(14, 14)

	assert.Equal(t, 7, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(3))
	assert.False(t, data.Contains(13))

	var keys []int
	expectedKeys := []int{2, 4, 5, 10, 12, 14, 15}
	data.ForEach(func(key int, _ int32) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	data.Erase(15)
	data.Erase(14)
	data.Erase(2)

	assert.Equal(t, 4, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(2))
	assert.False(t, data.Contains(14))

	keys = nil
	expectedKeys = []int{4, 5, 10, 12}
	data.ForEach(func(key int, _ int32) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))
}
