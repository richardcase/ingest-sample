package person

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/pkg/errors"
	"github.com/ttacon/libphonenumber"

	"github.com/richardcase/ingest-sample/pkg/api"
)

// Validator is an interface for a validation components
type Validator interface {
	Validate(*api.Person) error
}

// PersonValidator will validate a person
//nolint
type PersonValidator struct {
}

// Validate will validate a supplied person
func (p *PersonValidator) Validate(person *api.Person) error {
	return validation.ValidateStruct(person,
		validation.Field(&person.Id, validation.Required, validation.By(checkID)),
		validation.Field(&person.Name, validation.Required),
		validation.Field(&person.Email, validation.Required, is.Email),
		validation.Field(&person.MobileNumber, validation.Required, validation.By(checkPhone)),
	)
}

func checkID(value interface{}) error {
	id, ok := value.(int64)
	if !ok {
		return errors.New("id must be an int64")
	}
	if id < 1 {
		return errors.New("id must be greater than 0")
	}
	return nil
}

func checkPhone(value interface{}) error {
	phone, ok := value.(string)
	if !ok {
		return errors.New("number must be a string")
	}

	_, err := libphonenumber.Parse(phone, "GB")
	if err != nil {
		return errors.Wrapf(err, "validating %s as phone number", phone)
	}
	return nil
}
