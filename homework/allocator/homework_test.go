package main

import (
	"reflect"
	"sort"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

func Defragment(memory []byte, pointers []unsafe.Pointer) {
	if len(pointers) == 0 || len(memory) == 0 {
		return
	}

	sort.Slice(pointers, func(i, j int) bool {
		return uintptr(pointers[i]) < uintptr(pointers[j])
	})

	writeIndex := 0

	for i := 0; i < len(pointers); i++ {
		currentAddr := uintptr(pointers[i])

		targetAddr := uintptr(unsafe.Pointer(&memory[writeIndex]))

		if currentAddr != targetAddr {
			val := *(*byte)(pointers[i])
			memory[writeIndex] = val
			pointers[i] = unsafe.Pointer(&memory[writeIndex])
		}

		writeIndex++
	}

	for i := writeIndex; i < len(memory); i++ {
		memory[i] = 0x00
	}
}

func TestDefragmentation(t *testing.T) {
	var fragmentedMemory = []byte{
		0xFF, 0x00, 0x00, 0x00,
		0x00, 0xFF, 0x00, 0x00,
		0x00, 0x00, 0xFF, 0x00,
		0x00, 0x00, 0x00, 0xFF,
	}

	var fragmentedPointers = []unsafe.Pointer{
		unsafe.Pointer(&fragmentedMemory[0]),
		unsafe.Pointer(&fragmentedMemory[5]),
		unsafe.Pointer(&fragmentedMemory[10]),
		unsafe.Pointer(&fragmentedMemory[15]),
	}

	var defragmentedPointers = []unsafe.Pointer{
		unsafe.Pointer(&fragmentedMemory[0]),
		unsafe.Pointer(&fragmentedMemory[1]),
		unsafe.Pointer(&fragmentedMemory[2]),
		unsafe.Pointer(&fragmentedMemory[3]),
	}

	var defragmentedMemory = []byte{
		0xFF, 0xFF, 0xFF, 0xFF,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	}

	Defragment(fragmentedMemory, fragmentedPointers)
	assert.True(t, reflect.DeepEqual(defragmentedMemory, fragmentedMemory))
	assert.True(t, reflect.DeepEqual(defragmentedPointers, fragmentedPointers))
}
