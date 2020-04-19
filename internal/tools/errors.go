package tools

import "errors"

var (
	ErrDoesntExists           = errors.New("Record doesn't exists")
	ErrExistWithEmail         = errors.New("User with this email already exists")
	ErrExistWithSlug          = errors.New("Record with this slug alreday exists")
	ErrUserDoesntExists       = errors.New("User doesn't exists")
	ErrForumDoesntExists      = errors.New("Forum with this slug doesn't exists")
	ErrPostDoesntExists       = errors.New("Post with this id doesn't exists")
	ErrParentPostDoesntExists = errors.New("Post with this parent id doesn't exists")
	ErrThreadDoesntExists     = errors.New("Thread with this id doesn't exists")
	ErrIncorrectSlug          = errors.New("Slug is incorrect")
)
