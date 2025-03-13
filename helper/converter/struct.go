package converter

import (
	"fmt"
	"reflect"
)

// Make sure add tag `json:"key"` for each fields in struct
func StructToMSS(obj any) map[string]string {
	result := make(map[string]string)
	v := reflect.ValueOf(obj)
	t := reflect.TypeOf(obj)

	for i := range v.NumField() {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" {
			jsonTag = field.Name
		}
		result[jsonTag] = fmt.Sprintf("%v", v.Field(i).Interface())
	}

	return result
}

// Make sure add tag `json:"key"` for each fields in struct
func StructToMSX(obj any) map[string]any {
	result := make(map[string]any)
	v := reflect.ValueOf(obj)
	t := reflect.TypeOf(obj)

	for i := range v.NumField() {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" {
			jsonTag = field.Name
		}
		result[jsonTag] = v.Field(i).Interface()
	}

	return result
}
