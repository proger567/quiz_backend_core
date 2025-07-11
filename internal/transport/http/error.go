package http

import (
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

func errorEncoder(logger *logrus.Logger) func(context.Context, error, http.ResponseWriter) {
	return func(_ context.Context, err error, w http.ResponseWriter) {
		if err == nil {
			panic("encodeError with nil error")
		}
		response := ErrorResponse{
			Code:    codeFrom(err),
			Message: err.Error(),
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(response.Code)

		if err := json.NewEncoder(w).Encode(response); err != nil {
			logger.WithFields(logrus.Fields{
				"error": err,
			}).Info("method", "error encoder")
		}
	}
}

// TODO TODO TODO
func codeFrom(err error) int {
	switch {
	//case errors.Is(err, dto.ErrUserNotFound):
	//	return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

//func codeFrom(err error) int {
//	switch {
//	case errors.Is(err, ErrNotFound):
//		return http.StatusNotFound
//	case errors.Is(err, ErrAlreadyExists), errors.Is(err, ErrInconsistentIDs):
//		return http.StatusBadRequest
//	case errors.Is(err, ErrForbidden):
//		return http.StatusForbidden
//	case errors.Is(err, ErrPreconditionRequired):
//		return http.StatusPreconditionRequired
//	default:
//		return http.StatusInternalServerError
//	}
//}
