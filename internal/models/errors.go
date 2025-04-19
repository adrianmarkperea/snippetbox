package models

import (
	"errors"
)

// Use a custom error code so our controller doesn't have to deal with DB-specific errors
var ErrNoRecord = errors.New("models: no matching record found")
