package storage

import "errors"

var (
	ErrStudentNotFound = errors.New("student not found")
	ErrParsingQuery    = errors.New("error parsing query")
)