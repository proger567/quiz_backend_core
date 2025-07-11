package dto

import "errors"

//TODO make common package with errors? at least with common errors?

// common
var (
	ErrInternalServerError = errors.New("internal server error") //TODO
	ErrBadRouting          = errors.New("bad routing")
)

// subject
var (
	ErrSubjectAlreadyExists = errors.New("subject is already exist")
	ErrSubjectNotFound      = errors.New("subject is not found")
)

// question
var (
	ErrQuestionAlreadyExists = errors.New("question is already exist")
	ErrQuestionNotFound      = errors.New("question is not found")
)

// user //TODO delete?
var (
	ErrUserNotFound = errors.New("user is not found")
)
