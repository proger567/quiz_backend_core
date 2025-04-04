package error_handler

import (
	"context"
	"github.com/sirupsen/logrus"
)

// ErrorHandlerInterface TODO declare interface where it is used
type ErrorHandlerInterface interface {
	Handle(ctx context.Context, err error)
}

type ErrorHandler struct {
	logger logrus.Logger
}

func NewUnitLogHandler(logger logrus.Logger) *ErrorHandler {
	return &ErrorHandler{
		logger: logger,
	}
}

func (uh *ErrorHandler) Handle(ctx context.Context, err error) {
	uh.logger.Log(logrus.ErrorLevel, err)
}
