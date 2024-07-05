package validators

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	validator *validator.Validate
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

/*
/account
/account/<id>
/account/<id>/funds/all/transactions
/account/<id>/funds/all/transfers
/account/<id>/nominee
/account/<id>/income-target
/banks
/banks/providers
/orgs/<id>
/orgs/<id>/members
*/
func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Customize the error messages
		var errors []ValidationError
		for _, err := range err.(validator.ValidationErrors) {
			var element ValidationError
			element.Field = cv.getJSONField(i, err.Field())
			element.Message = cv.getErrorMessage(err)
			errors = append(errors, element)
		}
		return echo.NewHTTPError(400, errors)
	}
	return nil
}

func (cv *CustomValidator) getJSONField(i interface{}, fieldName string) string {
	stType := reflect.TypeOf(i)
	field, _ := stType.FieldByName(fieldName)
	jsonTag := field.Tag.Get("json")
	jsonField := strings.TrimSuffix(strings.Split(jsonTag, ",")[0], ",omitempty")
	if jsonField == "" {
		return fieldName
	}
	return jsonField
}

func (cv *CustomValidator) getErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return fmt.Sprintf("Should be at least %s characters long", err.Param())
	case "max":
		return fmt.Sprintf("Should be at most %s characters long", err.Param())
	default:
		return fmt.Sprintf("Invalid value: %v", err.Value())
	}
}

func NewCustomValidator() *CustomValidator {
	v := validator.New()
	// Use JSON tag name
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	return &CustomValidator{validator: v}
}
