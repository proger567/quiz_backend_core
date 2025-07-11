package transport

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"quiz_backend_core/internal/dto"
	"quiz_backend_core/internal/model"
)

type GetSubjectsRequest struct {
	User string
	Role string
}

type GetSubjectsResponse struct {
	Subjects []dto.Subject `json:"subjects,omitempty"`
	Err      error         `json:"err,omitempty"`
}

//**********************************************************************************************************************

type PostSubjectRequest struct {
	Subject dto.Subject
}

type PostSubjectResponse struct {
	ID  int64 `json:"id"`
	Err error `json:"err,omitempty"`
}

//**********************************************************************************************************************

type GetStatisticRequest struct {
	UserId   int64
	UserRole dto.Role
}

type GetStatisticResponse struct {
	Statistic dto.Statistic `json:"statistic,omitempty"`
	Err       error         `json:"err,omitempty"`
}

//**********************************************************************************************************************

type PutSubjectRequest struct {
	Subject dto.Subject
}

type PutSubjectResponse struct {
	Err error `json:"err,omitempty"`
}

//**********************************************************************************************************************

type DeleteSubjectRequest struct {
	ID int64
}

type DeleteSubjectResponse struct {
	Err error `json:"err,omitempty"`
}

//**********************************************************************************************************************

type SubjectsEndpoints struct {
	GetSubjectsEndpoint   endpoint.Endpoint
	PostSubjectEndpoint   endpoint.Endpoint
	PutSubjectEndpoint    endpoint.Endpoint
	DeleteSubjectEndpoint endpoint.Endpoint

	GetStatisticEndpoint endpoint.Endpoint
}

func MakeSubjectsEndpoints(s model.Subjects) SubjectsEndpoints {
	return SubjectsEndpoints{
		GetSubjectsEndpoint:   MakeGetSubjectsEndpoint(s),
		PostSubjectEndpoint:   MakePostSubjectEndpoint(s),
		PutSubjectEndpoint:    MakePutSubjectEndpoint(s),
		DeleteSubjectEndpoint: MakeDeleteSubjectEndpoint(s),

		GetStatisticEndpoint: MakeGetStatisticEndpoint(s),
	}
}

func MakeGetSubjectsEndpoint(s model.Subjects) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		t, err := s.GetSubjects(ctx)
		return GetSubjectsResponse{t, err}, err
	}
}

func MakePostSubjectEndpoint(s model.Subjects) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(PostSubjectRequest)
		id, err := s.AddSubject(ctx, req.Subject)
		return PostSubjectResponse{
			ID:  id,
			Err: err,
		}, err
	}
}

func MakePutSubjectEndpoint(s model.Subjects) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(PutSubjectRequest)
		err = s.UpdateSubject(ctx, req.Subject)
		return PutSubjectResponse{err}, err
	}
}

func MakeDeleteSubjectEndpoint(s model.Subjects) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(DeleteSubjectRequest)
		err = s.DeleteSubjectByID(ctx, req.ID)
		return DeleteSubjectResponse{err}, err
	}
}

func MakeGetStatisticEndpoint(s model.Subjects) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetStatisticRequest)
		stat, err := s.GetStatistic(ctx, req.UserId, req.UserRole)
		return GetStatisticResponse{stat, err}, err
	}
}
