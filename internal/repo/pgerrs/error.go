package pgerrs

import "errors"

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
	ErrForeignKey    = errors.New("incorrect foreign key")
)
