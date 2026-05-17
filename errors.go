package dotenv

import (
	"errors"
	"fmt"
)

var (
	ErrFileName     = errors.New("dotenv: invalid filename — must contain \".env\"")
	ErrFileNotFound = errors.New("dotenv: file not found")
)

type ParseError struct {
	File string
	Line int
	Msg  string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("dotenv: parse error in %s (line %d): %s", e.File, e.Line, e.Msg)
}
