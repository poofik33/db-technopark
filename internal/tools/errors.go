package tools

import "errors"

var (
	ErrDoesntExists      = errors.New("Record doesn't exists")
	ErrExistWithEmail    = errors.New("User with this email already exists")
	ErrExistWithSlug     = errors.New("Record with this slug alreday exists")
	ErrUserDoesntExists  = errors.New("User doesn't exists")
	ErrForumDoesntExists = errors.New("Forum with this slug doesn't exists")
	ErrIncorrectSlug     = errors.New("Slug is incorrect")
)
