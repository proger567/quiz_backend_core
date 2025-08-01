package http

import (
	"context"
	"encoding/json"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/proger567/quiz_backend_middleware"
	"net/http"
	"quiz_backend_core/internal/dto"
	"quiz_backend_core/internal/service"
	"quiz_backend_core/internal/transport"
	"strconv"
)

func makeSubjectsHTTPHandler(s *service.Services, r *mux.Router, options []httptransport.ServerOption) {
	e := transport.MakeSubjectsEndpoints(s.Subjects)

	r.Methods("OPTIONS", "GET").Path("/subjects").Handler(httptransport.NewServer(
		e.GetSubjectsEndpoint,
		decodeGetSubjectsRequest,
		encodeResponse,
		options...,
	))

	r.Methods("OPTIONS", "POST").Path("/subject").Handler(httptransport.NewServer(
		e.PostSubjectEndpoint,
		decodePostSubjectRequest,
		encodeResponse,
		options...,
	))

	r.Methods("OPTIONS", "PUT").Path("/subject").Handler(httptransport.NewServer(
		e.PutSubjectEndpoint,
		decodePutSubjectRequest,
		encodeResponse,
		options...,
	))

	r.Methods("OPTIONS", "DELETE").Path("/subject/{id}").Handler(httptransport.NewServer(
		e.DeleteSubjectEndpoint,
		decodeDeleteSubjectRequest,
		encodeResponse,
		options...,
	))

	r.Methods("OPTIONS", "GET").Path("/statistic").Handler(httptransport.NewServer(
		e.GetStatisticEndpoint,
		decodeGetStatisticRequest,
		encodeResponse,
		options...,
	))
}

func decodeGetSubjectsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return transport.GetSubjectsRequest{}, nil
}

func decodePostSubjectRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var subject dto.Subject
	if err := json.NewDecoder(r.Body).Decode(&subject); err != nil {
		return nil, err
	}

	subject.CreatorUserId, _ = r.Context().Value(quiz_backend_middleware.ContextVariablesUserID).(int64)

	return transport.PostSubjectRequest{
		Subject: subject,
	}, nil
}

func decodePutSubjectRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var subject dto.Subject
	if err := json.NewDecoder(r.Body).Decode(&subject); err != nil {
		return nil, err
	}

	return transport.PutSubjectRequest{
		Subject: subject,
	}, nil
}

func decodeDeleteSubjectRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	idStr, ok := mux.Vars(r)["id"]
	if !ok {
		return nil, dto.ErrBadRouting
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, err
	}

	return transport.DeleteSubjectRequest{
		ID: id,
	}, nil
}

func decodeGetStatisticRequest(_ context.Context, r *http.Request) (interface{}, error) {
	userID, _ := r.Context().Value(quiz_backend_middleware.ContextVariablesUserID).(int64)
	userRole, _ := r.Context().Value(quiz_backend_middleware.ContextVariablesUserRole).(string)

	return transport.GetStatisticRequest{
		UserId:   userID,
		UserRole: dto.Role(userRole),
	}, nil
}
