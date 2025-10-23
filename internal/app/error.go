package app

import (
	"fmt"
	"github.com/pkg/errors"
)

var (
	ErrUserDoesNotExist             = errors.New("user does not exist")
	ErrUserDoesNotExistOrOutOfPage  = errors.New("user does not exist or out of page")
	ErrUserSubscriptionDoesNotExist = errors.New("user subscription does not exist")
	ErrDatabaseProblems             = errors.New("database problems")
)

type errorWrapper struct {
	appErr    error
	originErr error
}

func Wrap(applicationError error, originError error) error {
	if originError == nil {
		return nil
	}

	if applicationError == nil {
		return originError
	}

	return errors.WithStack(&errorWrapper{
		appErr:    applicationError,
		originErr: originError,
	})
}

func (e errorWrapper) Error() string {
	return fmt.Sprintf("%s: %s", e.appErr, e.originErr)
}
