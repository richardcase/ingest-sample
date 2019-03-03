package phone

import (
	"fmt"
)

// ErrFailedPhoneFormatting is a error that represents an issue formatting a phone number
type ErrFailedPhoneFormatting struct {
	original string
	reason   string
}

// NewErrFailedPhoneFormatting creates a new ErrFailedPhoneFormatting with a given reason
func NewErrFailedPhoneFormatting(original string, reason string) *ErrFailedPhoneFormatting {
	return &ErrFailedPhoneFormatting{
		original: original,
		reason:   reason,
	}
}

// Error returns the error message
func (e *ErrFailedPhoneFormatting) Error() string {
	return fmt.Sprintf("enable to format %s as phone number because %s", e.original, e.reason)
}
