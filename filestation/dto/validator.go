package dto

import (
	"reflect"
	"strings"

	"gopkg.in/go-playground/validator.v8"
)

func UrlValidator(
	v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value,
	field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string,
) bool {
	if urlStr, ok := field.Interface().(string); ok {
		return strings.HasPrefix(urlStr, "http://")
	}

	return false
}
