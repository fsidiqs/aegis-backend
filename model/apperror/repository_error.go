package apperror

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

// "Set" of gorm error
// How to use: Check whether gorm error contains the following substrings
const (
	RecordExists        string = "SQLSTATE 23505"
	ForeignKeyNotExists string = "SQLSTATE 23503"
)

func NewRepoError(err error) *Error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return NewResourceNotFound()
	}

	if strings.Contains(err.Error(), RecordExists) {
		return NewConflictSimple()
	} else if strings.Contains(err.Error(), ForeignKeyNotExists) {
		return NewResourceNotFound()
	}
	return NewInternal()
}

func NewRepoErrorMsg(err error, msg string) *Error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return NewResourceNotFoundMsg(msg)
	}
	if err != nil {
		if strings.Contains(err.Error(), RecordExists) {
			return NewConflictMsg(msg)
		} else if strings.Contains(err.Error(), ForeignKeyNotExists) {
			return NewResourceNotFoundMsg(msg)
		}
	}
	return NewInternalWrap(msg)
}
