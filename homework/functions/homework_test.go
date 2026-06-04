package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Number interface {
	int | int8 | int16 | int32 | int64 | float32 | float64
}

func Map[T Number](data []T, action func(T) T) []T {
	if data == nil {
		return nil
	}

	result := make([]T, len(data))
	for i, v := range data {
		result[i] = action(v)
	}

	return result
}

func Filter[T Number](data []T, action func(T) bool) []T {
	if data == nil {
		return nil
	}

	result := make([]T, 0, len(data))
	for _, v := range data {
		if action(v) {
			result = append(result, v)
		}
	}

	return result
}

func Reduce[T Number](data []T, initial T, action func(T, T) T) T {
	accumulator := initial
	for _, v := range data {
		accumulator = action(accumulator, v)
	}

	return accumulator
}

func TestMap(t *testing.T) {
	tests := map[string]struct {
		data   []int
		action func(int) int
		result []int
	}{
		"nil numbers": {
			action: func(number int) int {
				return -number
			},
		},
		"empty numbers": {
			data: []int{},
			action: func(number int) int {
				return -number
			},
			result: []int{},
		},
		"inc numbers": {
			data: []int{1, 2, 3, 4, 5},
			action: func(number int) int {
				return number + 1
			},
			result: []int{2, 3, 4, 5, 6},
		},
		"double numbers": {
			data: []int{1, 2, 3, 4, 5},
			action: func(number int) int {
				return number * number
			},
			result: []int{1, 4, 9, 16, 25},
		},
	}

	testsFloat := map[string]struct {
		data   []float32
		action func(float32) float32
		result []float32
	}{
		"nil float": {
			action: func(number float32) float32 {
				return -number
			},
		},
		"empty numbers float": {
			data: []float32{},
			action: func(number float32) float32 {
				return -number
			},
			result: []float32{},
		},
		"inc numbers float": {
			data: []float32{1.1, 2.2, 3.3, 4.4, 5.5},
			action: func(number float32) float32 {
				return number + 1
			},
			result: []float32{2.1, 3.2, 4.3, 5.4, 6.5},
		},
		"double numbers float": {
			data: []float32{1.1, 2.2, 3.3, 4.4, 5.5},
			action: func(number float32) float32 {
				return number * number
			},
			result: []float32{1.21, 4.84, 10.889999, 19.36, 30.25},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := Map(test.data, test.action)
			assert.True(t, reflect.DeepEqual(test.result, result))
		})
	}

	for name, test := range testsFloat {
		t.Run(name, func(t *testing.T) {
			result := Map(test.data, test.action)
			assert.True(t, reflect.DeepEqual(test.result, result))
		})
	}
}

func TestFilter(t *testing.T) {
	tests := map[string]struct {
		data   []int
		action func(int) bool
		result []int
	}{
		"nil numbers": {
			action: func(number int) bool {
				return number == 0
			},
		},
		"empty numbers": {
			data: []int{},
			action: func(number int) bool {
				return number == 1
			},
			result: []int{},
		},
		"even numbers": {
			data: []int{1, 2, 3, 4, 5},
			action: func(number int) bool {
				return number%2 == 0
			},
			result: []int{2, 4},
		},
		"positive numbers": {
			data: []int{-1, -2, 1, 2},
			action: func(number int) bool {
				return number > 0
			},
			result: []int{1, 2},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := Filter(test.data, test.action)
			assert.True(t, reflect.DeepEqual(test.result, result))
		})
	}
}

func TestReduce(t *testing.T) {
	tests := map[string]struct {
		initial int
		data    []int
		action  func(int, int) int
		result  int
	}{
		"nil numbers": {
			action: func(lhs, rhs int) int {
				return 0
			},
		},
		"empty numbers": {
			data: []int{},
			action: func(lhs, rhs int) int {
				return 0
			},
		},
		"sum of numbers": {
			data: []int{1, 2, 3, 4, 5},
			action: func(lhs, rhs int) int {
				return lhs + rhs
			},
			result: 15,
		},
		"sum of numbers with initial value": {
			initial: 10,
			data:    []int{1, 2, 3, 4, 5},
			action: func(lhs, rhs int) int {
				return lhs + rhs
			},
			result: 25,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := Reduce(test.data, test.initial, test.action)
			assert.Equal(t, test.result, result)
		})
	}
}
