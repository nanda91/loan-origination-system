package validator

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateLoanApplication(err error) ([]string, error) {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		var messages []string
		for _, fe := range ve {
			field := fe.Field()
			switch fe.Tag() {
			case "required":
				messages = append(messages, fmt.Sprintf("%s is a required field", field))
			case "min":
				messages = append(messages, fmt.Sprintf("%s must be at least %s", field, fe.Param()))
			case "max":
				messages = append(messages, fmt.Sprintf("%s must be at most %s", field, fe.Param()))
			case "len":
				messages = append(messages, fmt.Sprintf("%s must be %s characters long", field, fe.Param()))
			default:
				messages = append(messages, fmt.Sprintf("%s is not valid", field))
			}
		}
		return messages, errors.New("Invalid input")
	}

	return []string{"Invalid input"}, nil
}
