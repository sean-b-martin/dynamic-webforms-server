package validation

import "github.com/go-playground/validator/v10"

type ErrorResponse struct {
	Field string      `json:"field"`
	Tag   string      `json:"tag"`
	Value interface{} `json:"value"`
	Error string      `json:"error"`
}

var validate = validator.New(validator.WithRequiredStructEnabled())

func Validate(data interface{}) []ErrorResponse {
	var validationErrors []ErrorResponse

	if errs := validate.Struct(data); errs != nil {
		validationErrors = make([]ErrorResponse, len(errs.(validator.ValidationErrors)))
		for i, err := range errs.(validator.ValidationErrors) {
			validationErrors[i] = ErrorResponse{
				Field: err.Field(),
				Tag:   err.Tag(),
				Value: err.Value(),
				Error: err.Error(),
			}
		}
	}

	return validationErrors
}
