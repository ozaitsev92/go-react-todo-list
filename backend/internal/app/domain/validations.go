package domain

import (
	"errors"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

func requiredIf(cond bool) validation.RuleFunc {
	return func(value interface{}) error {
		if cond {
			return validation.Validate(value, validation.Required)
		}

		return nil
	}
}

func timeNotZero(value interface{}) error {
	t, _ := value.(time.Time)

	if t.IsZero() {
		return errors.New("time must not be zero")
	}

	return nil
}
