package main

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type MultiError struct {
	Errors []error
}

func (e *MultiError) Error() string {
	if e == nil || len(e.Errors) == 0 {
		return ""
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "%d errors occured:\n", len(e.Errors))

	for _, err := range e.Errors {
		fmt.Fprintf(&sb, "\t* %s", err.Error())
	}
	sb.WriteString("\n")

	return sb.String()

}

func Append(err error, errs ...error) *MultiError {
	e := make([]error, 0, len(errs))
	if me, ok := err.(*MultiError); ok && me != nil {
		e = append(e, me.Errors...)
	} else if err != nil {
		e = append(e, err)
	}

	for _, er := range errs {
		if er != nil {
			e = append(e, er)
		}
	}

	if len(e) == 0 {
		return nil
	}

	return &MultiError{e}
}

func (e *MultiError) Is(target error) bool {
	if e == nil {
		return false
	}

	for _, err := range e.Errors {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}

func (e *MultiError) As(target interface{}) bool {
	if e == nil {
		return false
	}

	for _, err := range e.Errors {
		if errors.As(err, target) {
			return true
		}
	}

	return false
}

func (e *MultiError) Unwrap() error {
	if e == nil || len(e.Errors) == 0 {
		return nil
	}

	return e.Errors[0]
}

func (e *MultiError) UnwrapAll() []error {
	if e == nil {
		return nil
	}

	return e.Errors
}

func TestMultiError(t *testing.T) {
	var err error
	err = Append(err, errors.New("error 1"))
	err = Append(err, errors.New("error 2"))

	expectedMessage := "2 errors occured:\n\t* error 1\t* error 2\n"
	assert.EqualError(t, err, expectedMessage)
}

func TestMultiErrorNil(t *testing.T) {
	var err error
	err = Append(err, errors.New("error 1"))
	err = Append(err, nil)
	err = Append(err, errors.New("error 2"))

	expectedMessage := "2 errors occured:\n\t* error 1\t* error 2\n"
	assert.EqualError(t, err, expectedMessage)

	e := Append(nil, nil, nil)
	assert.Nil(t, e)
}

type customErr struct {
	msg string
}

func (e customErr) Error() string {
	return e.msg
}

func TestMultiErrorAdvanced(t *testing.T) {
	t.Run("single error", func(t *testing.T) {
		err := Append(nil, errors.New("single error"))
		assert.Equal(t, "1 errors occured:\n\t* single error\n", err.Error())
	})

	t.Run("append nil", func(t *testing.T) {
		err := Append(nil, nil, nil)
		assert.Nil(t, err)
	})

	t.Run("append to existing MultiError", func(t *testing.T) {
		var err1, err2 error
		err1 = Append(nil, errors.New("error 1"))
		err2 = Append(err1, errors.New("error 2"))

		// Для проверки внутренней структуры нужно привести к *MultiError
		me, ok := err2.(*MultiError)
		assert.True(t, ok)
		assert.Len(t, me.Errors, 2)
	})

	t.Run("errors.Is support", func(t *testing.T) {
		targetErr := errors.New("target error")
		err := Append(nil, errors.New("error 1"), targetErr, errors.New("error 2"))

		assert.True(t, err.Is(targetErr))
	})

	t.Run("errors.As support", func(t *testing.T) {
		targetErr := customErr{msg: "custom"}
		err := Append(nil, errors.New("error 1"), targetErr)

		assert.True(t, err.As(&targetErr))
	})
}
