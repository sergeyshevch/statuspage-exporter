package utils

import "errors"

// ErrInvalidURL is returned when URL is invalid.
var ErrInvalidURL = errors.New("invalid URL. It won't be parsed. Check that your url contains scheme")
