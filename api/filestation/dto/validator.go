package dto

import (
	"log"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
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

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("url_validator", UrlValidator)
		if err != nil {
			log.Println("err:", err)
		}
	}
}
