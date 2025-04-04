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

func makeQuestionsHTTPHandler(s *service.Services, r *mux.Router, options []httptransport.ServerOption) {
	e := transport.MakeQuestionsEndpoints(s.Questions)

	r.Methods("OPTIONS", "GET").Path("/questions").Handler(httptransport.NewServer(
		e.GetQuestionsEndpoint,
		decodeGetQuestionsRequest,
		encodeResponse,
		options...,
	))

	r.Methods("OPTIONS", "GET").Path("/types").Handler(httptransport.NewServer(
		e.GetQuestionTypesEndpoint,
		decodeGetTypesRequest,
		encodeResponse,
		options...,
	))

	r.Methods("OPTIONS", "GET").Path("/statuses").Handler(httptransport.NewServer(
		e.GetQuestionStatusesEndpoint,
		decodeGetStatusesRequest,
		encodeResponse,
		options...,
	))

	r.Methods("OPTIONS", "POST").Path("/question").Handler(httptransport.NewServer(
		e.PostQuestionEndpoint,
		decodePostQuestionRequest,
		encodeResponse,
		options...,
	))

	r.Methods("OPTIONS", "PUT").Path("/question/{id}").Handler(httptransport.NewServer(
		e.PutQuestionEndpoint,
		decodePutQuestionRequest,
		encodeResponse,
		options...,
	))

	r.Methods("OPTIONS", "PUT").Path("/question/{id}/moderate").Handler(httptransport.NewServer( //TODO
		e.PutQuestionModerateEndpoint,
		decodePutQuestionModerateRequest,
		encodeResponse,
		options...,
	))

	r.Methods("OPTIONS", "DELETE").Path("/question/{id}").Handler(httptransport.NewServer(
		e.DeleteQuestionEndpoint,
		decodeDeleteQuestionRequest,
		encodeResponse,
		options...,
	))
}

func decodeGetQuestionsRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	//defaults
	dRequest := transport.GetQuestionsRequest{
		SubjectID:     -1,
		CreatorUserID: -1,
		StatusID:      -1,
	}

	query := r.URL.Query()

	// TODO check errors
	if subjectIdStr := query.Get("subject_id"); subjectIdStr != "" {
		if dRequest.SubjectID, err = strconv.ParseInt(subjectIdStr, 10, 64); err != nil {
			return dRequest, err
		}
	}

	if userIdStr := query.Get("creator_user_id"); userIdStr != "" {
		if dRequest.CreatorUserID, err = strconv.ParseInt(userIdStr, 10, 64); err != nil {
			return dRequest, err
		}
	}

	if statusIdStr := query.Get("status_id"); statusIdStr != "" {
		if dRequest.StatusID, err = strconv.ParseInt(statusIdStr, 10, 64); err != nil {
			return dRequest, err
		}
	}

	return dRequest, nil
}

func decodeGetTypesRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	return transport.GetQuestionTypesResponse{}, nil
}

func decodeGetStatusesRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	return transport.GetQuestionStatusesRequest{}, nil
}

func decodePostQuestionRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var question dto.InputQuestion
	if err := json.NewDecoder(r.Body).Decode(&question); err != nil {
		return nil, err
	}

	return transport.PostQuestionRequest{
		Question: question,
	}, nil
}

func decodePutQuestionRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var question dto.InputQuestion
	if err := json.NewDecoder(r.Body).Decode(&question); err != nil {
		return nil, err
	}

	vars := mux.Vars(r)
	questionIdStr, ok := vars["id"]
	if !ok {
		return nil, dto.ErrBadRouting
	}

	questionId, err := strconv.ParseInt(questionIdStr, 10, 64)
	if err != nil {
		return nil, err //TODO wrap error with dto.ErrBadRouting?
	}

	return transport.PutQuestionRequest{
		ID:       questionId,
		Question: question,
	}, nil
}

func decodePutQuestionModerateRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var moderate = struct {
		Approve bool   `json:"approve"`
		Comment string `json:"comment,omitempty"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&moderate); err != nil {
		return nil, err
	}

	vars := mux.Vars(r)
	questionIdStr, ok := vars["id"]
	if !ok {
		return nil, dto.ErrBadRouting
	}

	questionId, err := strconv.ParseInt(questionIdStr, 10, 64)
	if err != nil {
		return nil, err //TODO wrap error with dto.ErrBadRouting?
	}

	return transport.PutQuestionModerateRequest{
		ID:      questionId,
		Approve: moderate.Approve,
		Comment: moderate.Comment,
	}, nil
}

func decodeDeleteQuestionRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	questionIdStr, ok := vars["id"]
	if !ok {
		return nil, dto.ErrBadRouting
	}

	questionId, err := strconv.ParseInt(questionIdStr, 10, 64)
	if err != nil {
		return nil, err //TODO wrap error with dto.ErrBadRouting?
	}

	return transport.DeleteQuestionRequest{
		ID: questionId,
	}, nil
}
