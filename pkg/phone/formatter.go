package phone

import (
	"strings"

	"github.com/ttacon/libphonenumber"
)

// Formatter is an interface that can be used to implement formatters
type Formatter interface {
	Format(phone string) (string, error)
}

// UKToE164Formatter is a formatter that will convert a UK phone number
// into E.164 format
type UKToE164Formatter struct{}

// NewUKToE164Formatter creates a new UKToE164Formatter
func NewUKToE164Formatter() *UKToE164Formatter {
	return &UKToE164Formatter{}
}

// Format will format a given UK phone number into E.164 format
func (u *UKToE164Formatter) Format(phone string) (string, error) {
	if len(phone) < 1 {
		return "", NewErrFailedPhoneFormatting(phone, "phone must be greater than 1")
	}
	// Remove any ( )
	phone = strings.Replace(phone, "(", "", -1)
	phone = strings.Replace(phone, ")", "", -1)

	num, err := libphonenumber.Parse(phone, "GB")
	if err != nil {
		return "", NewErrFailedPhoneFormatting(phone, err.Error())
	}

	formatted := libphonenumber.Format(num, libphonenumber.E164)
	return formatted, nil
}
