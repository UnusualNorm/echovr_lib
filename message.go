package echovr

import (
	"fmt"
	"reflect"
)

type Message interface {
	Symbol() uint64
	Stream(s *EasyStream) error
}

type Serializable interface {
	Stream(s *EasyStream) error
}

func Deserialize(obj Serializable, b []byte) error {
	s := NewEasyStream(0, b)
	return obj.Stream(s)
}

func Serialize(obj Serializable) ([]byte, error) {
	s := NewEasyStream(1, []byte{})
	err := obj.Stream(s)
	if err != nil {
		return nil, err
	}
	return s.Bytes(), nil
}

func RunErrorFunctions(funcs []func() error) error {
	for _, f := range funcs {
		if err := f(); err != nil {
			return err
		}
	}
	return nil
}

func StringifyStruct(s interface{}) string {
	valueOf := reflect.ValueOf(s)

	// Check if the input is a struct
	if valueOf.Kind() != reflect.Struct {
		return ""
	}

	// Get the type of the struct
	structType := valueOf.Type()

	// Initialize an empty string to store the result
	result := structType.Name() + "{"

	// Iterate through the fields of the struct
	for i := 0; i < valueOf.NumField(); i++ {
		field := valueOf.Field(i)
		fieldType := structType.Field(i)

		// Add the field name to the result
		result += fieldType.Name + ": "

		// Check if the field is a struct and recursively stringify it
		if field.Kind() == reflect.Struct {
			result += StringifyStruct(field.Interface())
		} else {
			// Convert the field value to a string and add it to the result
			result += fmt.Sprintf("%v", field.Interface())
		}

		// Add a comma and space unless it's the last field
		if i < valueOf.NumField()-1 {
			result += ", "
		}
	}

	// Close the struct representation
	result += "}"

	return result
}
