package utils

import (
	"fmt"
	"reflect"
)

// structToMap converts a struct to a map[string]interface{}
// structToNonZeroMap converts a struct to a map[string]interface{}, skipping zero-value fields
func StructToMap(input interface{}) (map[string]interface{}, error) {
	// Create the result map
	result := make(map[string]interface{})

	// Get the reflect.Value of the input
	v := reflect.ValueOf(input)

	// Ensure we have a struct
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("input is not a struct")
	}

	// Loop through struct fields
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		// Skip unexported (private) fields
		// See http://www.golangpatterns.info/object-oriented/unexported-fields
		if !fieldValue.CanInterface() {
			continue
		}

		// Skip zero-value fields (e.g., "", 0, nil)
		// See https://github.com/golang/go/wiki/CodeReviewComments#declaring-empty-slices
		if isZeroValue(fieldValue) {
			continue
		}

		// Use field name as key
		result[field.Name] = fieldValue.Interface()
	}
	return result, nil
}

// isZeroValue checks if a field has a zero value (e.g., "", 0, nil)
func isZeroValue(v reflect.Value) bool {
	zero := reflect.Zero(v.Type()).Interface()
	return reflect.DeepEqual(v.Interface(), zero)
}
