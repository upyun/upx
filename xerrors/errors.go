package xerrors

import "errors"

var (
	ErrInvalidCommand = errors.New("invalid command")
	ErrRequireLogin   = errors.New("log in to UpYun first")
)
