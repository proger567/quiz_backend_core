package http

import (
	"context"
	"encoding/json"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/proger567/quiz_backend_middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"net/http"
	"quiz_backend_core/internal/service"
	"quiz_backend_core/pkg/error_handler"
)

func MakeHTTPHandler(s *service.Services, logger *logrus.Logger, cors bool) http.Handler {
	r := mux.NewRouter()

	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(error_handler.NewUnitLogHandler(logrus.Logger{})), //TODO from deps (service level)?
		httptransport.ServerErrorEncoder(errorEncoder(logger)),
	}

	//if cors {
	r.Use(quiz_backend_middleware.AccessControlMiddleware)
	//}

	r.Use(quiz_backend_middleware.FillContextMiddleware)

	makeSubjectsHTTPHandler(s, r.PathPrefix("/subjects").Subrouter(), options)
	makeQuestionsHTTPHandler(s, r.PathPrefix("/questions").Subrouter(), options)
	makeQuizzesHTTPHandler(s, r.PathPrefix("/quizzes").Subrouter(), options)
	//makeExamHTTPHandler(s, r.PathPrefix("/examination").Subrouter(), options)

	r.Methods("GET").Path("/metrics").Handler(promhttp.Handler())

	return r
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
