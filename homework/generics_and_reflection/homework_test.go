package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Person struct {
	Name    string `properties:"name"`
	Address string `properties:"address,omitempty"`
	Age     int    `properties:"age"`
	Married bool   `properties:"married"`
}

func Serialize(v any) string {
	val := reflect.ValueOf(v)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	typ := val.Type()
	var lines []string

	for i := 0; i < val.NumField(); i++ {
		fieldVal := val.Field(i)
		fieldType := typ.Field(i)

		if !fieldType.IsExported() {
			continue
		}

		tag := fieldType.Tag.Get("properties")
		if tag == "" || tag == "-" {
			continue
		}

		tagParts := strings.Split(tag, ",")
		key := strings.TrimSpace(tagParts[0])

		omitempty := false
		if len(tagParts) > 1 && strings.TrimSpace(tagParts[1]) == "omitempty" {
			omitempty = true
		}

		if omitempty && fieldVal.IsZero() {
			continue
		}

		var strVal string
		switch fieldVal.Kind() {
		case reflect.String:
			strVal = fieldVal.String()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			strVal = fmt.Sprintf("%d", fieldVal.Int())
		case reflect.Bool:
			strVal = fmt.Sprintf("%t", fieldVal.Bool())
		case reflect.Float32, reflect.Float64:
			strVal = fmt.Sprintf("%g", fieldVal.Float())
		default:
			strVal = fmt.Sprintf("%v", fieldVal.Interface())
		}

		key = escapeKey(key)

		lines = append(lines, fmt.Sprintf("%s=%s", key, strVal))
	}

	return strings.Join(lines, "\n")
}

// escapeKey экранирует служебные символы в ключах
func escapeKey(key string) string {
	r := strings.NewReplacer(
		" ", "\\ ",
		"=", "\\=",
		":", "\\:",
	)
	return r.Replace(key)
}

func TestSerialization(t *testing.T) {
	tests := map[string]struct {
		person Person
		result string
	}{
		"test case with empty fields": {
			result: "name=\nage=0\nmarried=false",
		},
		"test case with fields": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
			},
			result: "name=John Doe\nage=30\nmarried=true",
		},
		"test case with omitempty field": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
				Address: "Paris",
			},
			result: "name=John Doe\naddress=Paris\nage=30\nmarried=true",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := Serialize(test.person)
			assert.Equal(t, test.result, result)
		})
	}
}
