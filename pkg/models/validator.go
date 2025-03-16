package models

import "github.com/go-playground/validator/v10"

// Global instance of validate
// Accessible only from this package
var validate *validator.Validate

func init() {
	validate = validator.New()
}
