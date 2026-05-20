package main

import (
	"encoding/binary"
	"math/bits"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

// ToLittleEndianBinary swap bytes using binary library.
func ToLittleEndianBinary(number uint32) uint32 {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, number)

	return binary.LittleEndian.Uint32(buf)
}

// ToLittleEndianManual swap bytes using << >> operators.
func ToLittleEndianManual(number uint32) uint32 {
	return ((number & 0xFF000000) >> 24) | // Берём 1-й байт, двигаем в конец
		((number & 0x00FF0000) >> 8) | // Берём 2-й байт, двигаем на 3-е место
		((number & 0x0000FF00) << 8) | // Берём 3-й байт, двигаем на 2-е место
		((number & 0x000000FF) << 24) // Берём 4-й байт, двигаем в начало
}

type Unsigned interface {
	~uint16 | ~uint32 | ~uint64
}

// ToLittleEndianT swap bytes using generics and bits library.
func ToLittleEndianT[T Unsigned](val T) T {
	switch v := any(val).(type) {
	case uint16:
		return any(bits.ReverseBytes16(v)).(T)
	case uint32:
		return any(bits.ReverseBytes32(v)).(T)
	case uint64:
		return any(bits.ReverseBytes64(v)).(T)
	}

	return val
}

func TestСonversion(t *testing.T) {
	tests := map[string]struct {
		number uint32
		result uint32
	}{
		"test case #1": {
			number: 0x00000000,
			result: 0x00000000,
		},
		"test case #2": {
			number: 0xFFFFFFFF,
			result: 0xFFFFFFFF,
		},
		"test case #3": {
			number: 0x00FF00FF,
			result: 0xFF00FF00,
		},
		"test case #4": {
			number: 0x0000FFFF,
			result: 0xFFFF0000,
		},
		"test case #5": {
			number: 0x01020304,
			result: 0x04030201,
		},
		"only highest byte set": {
			number: 0xFF000000,
			result: 0x000000FF,
		},
		"only lowest byte set": {
			number: 0x000000FF,
			result: 0xFF000000,
		},
		"inner left byte set": {
			number: 0x00FF0000,
			result: 0x0000FF00,
		},
		"inner right byte set": {
			number: 0x0000FF00,
			result: 0x00FF0000,
		},
		"asymmetric hex words (DEAD CODE)": {
			number: 0xDEADC0DE,
			result: 0xDEC0ADDE, // DE | C0 | AD | DE
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Run("Binary", func(t *testing.T) {
				r1 := ToLittleEndianBinary(test.number)
				assert.Equal(t, test.result, r1)
			})

			t.Run("Manual", func(t *testing.T) {
				r2 := ToLittleEndianManual(test.number)
				assert.Equal(t, test.result, r2)
			})

			t.Run("Generic", func(t *testing.T) {
				r3 := ToLittleEndianT(test.number)
				assert.Equal(t, test.result, r3)
			})
		})
	}
}
