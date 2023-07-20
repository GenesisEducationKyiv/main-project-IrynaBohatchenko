package storageerrors

import "errors"

var ErrEmailExists = errors.New("email already exists")
var ErrInvalidEmail = errors.New("invalid email")
