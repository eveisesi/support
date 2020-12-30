package internal

import (
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

type Level int

const (
	LevelInternal Level = 500
	LevelBad      Level = 400
)

type InternalError struct {
	Level   Level
	Message string
}

func (i InternalError) Error() string {
	return i.Message
}

func NewInternalError(level Level, message string) InternalError {
	return InternalError{Level: level, Message: message}
}

const duplicateKeyError = 11000

func IsUniqueConstrainViolation(exception error) bool {

	var bwe mongo.BulkWriteException
	if errors.As(exception, &bwe) {
		for _, errs := range bwe.WriteErrors {
			if errs.Code == duplicateKeyError {
				return true
			}
		}
		return false
	}
	var we mongo.WriteException
	if errors.As(exception, &we) {
		for _, errs := range we.WriteErrors {
			if errs.Code == duplicateKeyError {
				return true
			}
		}
		return false
	}

	return false
}
