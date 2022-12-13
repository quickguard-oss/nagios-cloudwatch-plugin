package errors

import (
	"errors"
	"fmt"
)

type ArgumentError struct {
	err   error
	key   string
	value string
}

type CloudWatchError struct {
	err error
}

func NewArgumentErrorWithError(err error, key string, value string) ArgumentError {
	return ArgumentError{
		err:   err,
		key:   key,
		value: value,
	}
}

func NewArgumentErrorWithMessage(msg string, key string, value string) ArgumentError {
	return NewArgumentErrorWithError(errors.New(msg), key, value)
}

func NewCloudWatchError(err error) CloudWatchError {
	return CloudWatchError{
		err: err,
	}
}

func (e ArgumentError) Error() string {
	return fmt.Sprintf(`invalid argument "%s" for %s: %s`, e.value, e.key, e.err)
}

func (e ArgumentError) Unwrap() error {
	return e.err
}

func (e CloudWatchError) Error() string {
	return fmt.Sprintf(`unable to get metrics: %s`, e.err)
}

func (e CloudWatchError) Unwrap() error {
	return e.err
}
