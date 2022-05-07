package errors

import (
	"github.com/go-playground/validator/v10"
)

type ValidationRule struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value"`
}

var validate = validator.New()

func ValidateStruct(inputStruct interface{}) []*ValidationRule {
	var errors []*ValidationRule
	if err := validate.Struct(inputStruct); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			element := ValidationRule{
				Field: err.StructNamespace(),
				Tag:   err.Tag(),
				Value: err.Param(),
			}
			errors = append(errors, &element)
		}
	}
	return errors
}
