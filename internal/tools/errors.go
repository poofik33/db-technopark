package tools

import "errors"

var (
	ErrDoesntExists = errors.New("Record doesn't exists")
	ErrExistWithEmail = errors.New("User with this email already exists")
)
