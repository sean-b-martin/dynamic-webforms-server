package validation

import (
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
)

var validate = setupValidate()

func setupValidate() *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

	return v
}

type ErrorResponse struct {
	Field  string        `json:"field"`
	Value  interface{}   `json:"value"`
	Failed ErrorExpected `json:"failed"`
}

type ErrorExpected struct {
	Constraint    string `json:"constraint"`
	Configuration string `json:"configuration"`
}

func Validate(data interface{}) []ErrorResponse {
	var validationErrors []ErrorResponse

	if errs := validate.Struct(data); errs != nil {
		validationErrors = make([]ErrorResponse, len(errs.(validator.ValidationErrors)))
		for i, err := range errs.(validator.ValidationErrors) {
			validationErrors[i] = ErrorResponse{
				Field: err.Field(),
				Value: err.Value(),
				Failed: ErrorExpected{
					Constraint:    err.Tag(),
					Configuration: err.Param(),
				},
			}
		}
	}

	return validationErrors
}
