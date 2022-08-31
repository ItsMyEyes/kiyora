package validate

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

func ParseError(err error) []string {
	var out []string
	var ve validator.ValidationErrors
	errors.As(err, &ve)
	fmt.Printf("%+v\n", ve)
	if len(ve) > 0 {
		for _, e := range ve {
			out = append(out, ParseFieldError(e))
		}
	} else {
		out = append(out, err.Error())
	}

	return out
}

func ParseFieldError(e validator.FieldError) string {
	// workaround to the fact that the `gt|gtfield=Start` gets passed as an entire tag for some reason
	// https://github.com/go-playground/validator/issues/926
	// fieldPrefix := fmt.Sprintf("The field %s", e.Field())
	fieldPrefix := e.Field()
	tag := strings.Split(e.Tag(), "|")[0]
	fmt.Printf("%+v\n", tag)
	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", fieldPrefix)
	case "min":
		return fmt.Sprintf("%s must be at least %s", fieldPrefix, e.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s", fieldPrefix, e.Param())
	case "email":
		return fmt.Sprintf("%s must be a valid email", fieldPrefix)
	default:
		// if it's a tag for which we don't have a good format string yet we'll try using the default english translator
		english := en.New()
		translator := ut.New(english, english)
		if translatorInstance, found := translator.GetTranslator("en"); found {
			return e.Translate(translatorInstance)
		} else {
			return fmt.Errorf("%v", e).Error()
		}
	}
}
