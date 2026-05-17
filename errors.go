package dotenv

import (
	"errors"
	"fmt"
)

var (
	ErrFileName     = errors.New("dotenv: invalid filename — must contain \".env\"")
	ErrFileNotFound = errors.New("dotenv: file not found")
)

// ParseError is returned when a .env file contains a malformed line.
// It carries the filename and line number so callers can report the exact location.
type ParseError struct {
	File string
	Line int
	Msg  string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("dotenv: parse error in %s (line %d): %s", e.File, e.Line, e.Msg)
}
