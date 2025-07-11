package http

//
//import (
//	"context"
//	httptransport "github.com/go-kit/kit/transport/http"
//	"github.com/gorilla/mux"
//	"net/http"
//	"quiz_backend_core/internal/service"
//	"quiz_backend_core/internal/transport"
//	"strconv"
//)
//
//func makeExamHTTPHandler(s *service.Services, r *mux.Router, options []httptransport.ServerOption) {
//	e := transport.MakeExamEndpoints(s.Quizzes)
//
//	r.Methods("OPTIONS", "POST").Path("/start").Handler(httptransport.NewServer(
//		e.PostExamStartEndpoint,
//		decodePostExamStartRequest,
//		encodeResponse,
//		options...,
//	))
//}
//
//func decodePostExamStartRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
//	//defaults
//	dRequest := transport.GetExamStartRequest{
//		CreatorUserID: -1,
//	}
//
//	query := r.URL.Query()
//
//	// TODO check errors
//	if userIdStr := query.Get("creator_user_id"); userIdStr != "" {
//		if dRequest.CreatorUserID, err = strconv.ParseInt(userIdStr, 10, 64); err != nil {
//			return dRequest, err
//		}
//	}
//
//	return dRequest, nil
//}
