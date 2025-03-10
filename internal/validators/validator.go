package validators

import (
	"github.com/go-playground/validator/v10"
)

// CustomValidator implementa el validador de Echo con go-playground/validator
type CustomValidator struct {
	Validator *validator.Validate
}

// Validate ejecuta la validación en la estructura recibida
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}

// NewValidator crea una nueva instancia de CustomValidator
func NewValidator() *CustomValidator {
	return &CustomValidator{Validator: validator.New()}
}
