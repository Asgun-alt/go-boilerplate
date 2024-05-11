package vo

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}

func NewCustomValidator() *CustomValidator {
	v := validator.New()
	_ = v.RegisterValidation("msisdn", ValidateMsisdn)
	return &CustomValidator{Validator: v}
}

func ValidateMsisdn(fl validator.FieldLevel) bool {
	reg := `^(62|8|08)`
	rgx := regexp.MustCompile(reg)
	return rgx.MatchString(fl.Field().String())
}
