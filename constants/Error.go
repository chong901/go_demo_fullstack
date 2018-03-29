package constants

import (
	"errors"
	"fmt"
	"strings"
)

const (
	ErrMsgInternalServerError = "tInternalServerError"
)

var (
	ErrParameters = errors.New("Wrong parameters")

	ErrIdDuplicated = errors.New("Data is duplicated.")

	IsDuplicatedErr = func(err error) bool {
		return strings.Index(err.Error(), "Duplicate") != -1
	}

	ErrNoAuth = errors.New("tUnauthorized")

	ErrNoTransaction = errors.New("This Transaction Is Empty.")

	Err404NotFounr = errors.New("tPageNotFound")
)

func ErrDataNotExist(dataName string) error {
	return errors.New(fmt.Sprintf("This %s does not exist.", dataName))
}
