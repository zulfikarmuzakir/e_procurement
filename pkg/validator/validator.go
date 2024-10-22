package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

// Initialize the validator instance globally when the package is loaded.
func init() {
	validate = validator.New()
}

// ValidateStruct validates the struct based on the tags using the global validator instance.
func ValidateStruct(data interface{}) error {
	err := validate.Struct(data)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return fmt.Errorf("validation error: %v", err)
		}
		return err
	}
	return nil
}

// formatValidationErrors converts validation errors into a readable string format.
func formatValidationErrors(ves validator.ValidationErrors) error {
	errorMessages := ""
	for _, ve := range ves {
		errorMessages += fmt.Sprintf("Field '%s' failed validation for '%s' rule\n", ve.Field(), ve.Tag())
	}
	return fmt.Errorf(errorMessages)
}

func GetValidationErrors(err error) []map[string]string {
	if err == nil {
		return nil
	}
	var validationErrors []map[string]string
	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, e := range ve {
			validationError := map[string]string{
				"field": e.Field(),
				"error": fmt.Sprintf("Field validation for '%s' failed on the '%s' tag", e.Field(), e.Tag()),
			}
			validationErrors = append(validationErrors, validationError)
		}
	}

	return validationErrors
}
