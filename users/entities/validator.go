package entities

import (
	"log"
	"reflect"
	"regexp"

	"gopkg.in/go-playground/validator.v8"
)

func EmailValidator(
	v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value,
	field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string,
) bool {
	if email, ok := field.Interface().(string); ok {
		r, err := regexp.Compile(`^[a-zA-Z0-9_-{Han}]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$`)
		if err != nil {
			log.Fatal("正则表达式错误")
		}

		return r.MatchString(email)
	}

	return false
}
