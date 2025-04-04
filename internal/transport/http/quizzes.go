package http

import (
	"context"
	"encoding/json"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
	"quiz_backend_core/internal/dto"
	"quiz_backend_core/internal/service"
	"quiz_backend_core/internal/transport"
	"strconv"
)

func makeQuizzesHTTPHandler(s *service.Services, r *mux.Router, options []httptransport.ServerOption) {
	e := transport.MakeQuizzesEndpoints(s.Quizzes)

	r.Methods("OPTIONS", "GET").Path("/quizzes").Handler(httptransport.NewServer(
		e.GetQuizzesEndpoint,
		decodeGetQuizzesRequest,
		encodeResponse,
		options...,
	))

	r.Methods("OPTIONS", "GET").Path("/{id}/questions").Handler(httptransport.NewServer(
		e.GetQuestionsByQuizIDEndpoint,
		decodeGetQuestionsByQuizIDRequest,
		encodeResponse,
		options...,
	))

	r.Methods("OPTIONS", "GET").Path("/{id}").Handler(httptransport.NewServer(
		e.GetQuizByIDEndpoint,
		decodeGetQuizByIDRequest,
		encodeResponse,
		options...,
	))

	r.Methods("OPTIONS", "POST").Path("/quiz").Handler(httptransport.NewServer(
		e.PostQuizEndpoint,
		decodePostQuizRequest,
		encodeResponse,
		options...,
	))

	r.Methods("OPTIONS", "DELETE").Path("/{id}").Handler(httptransport.NewServer(
		e.DeleteQuizEndpoint,
		decodeDeleteQuizByIDRequest,
		encodeResponse,
		options...,
	))
}

func decodeGetQuizzesRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	//defaults
	dRequest := transport.GetQuizzesRequest{
		CreatorUserID: -1,
	}

	query := r.URL.Query()

	// TODO check errors
	if userIdStr := query.Get("creator_user_id"); userIdStr != "" {
		if dRequest.CreatorUserID, err = strconv.ParseInt(userIdStr, 10, 64); err != nil {
			return dRequest, err
		}
	}

	return dRequest, nil
}

func decodeGetQuestionsByQuizIDRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	quizIdStr, ok := vars["id"]
	if !ok {
		return nil, dto.ErrBadRouting
	}

	quizId, err := strconv.ParseInt(quizIdStr, 10, 64)
	if err != nil {
		return nil, err //TODO wrap error with dto.ErrBadRouting?
	}

	return transport.GetQuestionsByQuizIDRequest{
		QuizID: quizId,
	}, nil
}

func decodeGetQuizByIDRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	quizIdStr, ok := vars["id"]
	if !ok {
		return nil, dto.ErrBadRouting
	}

	quizId, err := strconv.ParseInt(quizIdStr, 10, 64)
	if err != nil {
		return nil, err //TODO wrap error with dto.ErrBadRouting?
	}

	return transport.GetQuizByIDRequest{
		QuizID: quizId,
	}, nil
}

func decodePostQuizRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var quiz dto.InputQuiz
	if err := json.NewDecoder(r.Body).Decode(&quiz); err != nil {
		return nil, err
	}

	return transport.PostQuizRequest{
		Quiz: quiz,
	}, nil
}

func decodeDeleteQuizByIDRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	quizIdStr, ok := vars["id"]
	if !ok {
		return nil, dto.ErrBadRouting
	}

	quizId, err := strconv.ParseInt(quizIdStr, 10, 64)
	if err != nil {
		return nil, err //TODO wrap error with dto.ErrBadRouting?
	}

	return transport.DeleteQuizByIDRequest{
		QuizID: quizId,
	}, nil
}
