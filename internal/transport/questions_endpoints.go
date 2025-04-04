package transport

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"quiz_backend_core/internal/dto"
	"quiz_backend_core/internal/model"
)

type GetQuestionsRequest struct {
	SubjectID     int64 `json:"subject_id"`
	CreatorUserID int64 `json:"creator_user_id"`
	StatusID      int64 `json:"status_id"`
}

type GetQuestionsResponse struct {
	Questions []dto.Question `json:"questions"`
	Err       error          `json:"err,omitempty"`
}

// *********************************************************************************************************************

type GetQuestionTypesRequest struct{}

type GetQuestionTypesResponse struct {
	Types []dto.QuestionType `json:"types,omitempty"`
	Err   error              `json:"err,omitempty"`
}

// *********************************************************************************************************************

type GetQuestionStatusesRequest struct{}

type GetQuestionStatusesResponse struct {
	Statuses []dto.QuestionStatus `json:"statuses,omitempty"`
	Err      error                `json:"err,omitempty"`
}

// *********************************************************************************************************************

type PostQuestionRequest struct {
	Question dto.InputQuestion `json:"question"`
}

type PostQuestionResponse struct {
	ID  int64 `json:"id"`
	Err error `json:"err,omitempty"`
}

// *********************************************************************************************************************

type PutQuestionRequest struct {
	ID       int64
	Question dto.InputQuestion

	UserID       int64    `json:"user_id"`
	UserRoleName dto.Role `json:"user_role_name"`
}

type PutQuestionResponse struct {
	Err error `json:"err,omitempty"`
}

// *********************************************************************************************************************

type PutQuestionModerateRequest struct {
	ID      int64
	Approve bool
	Comment string
}

type PutQuestionModerateResponse struct {
	Err error `json:"err,omitempty"`
}

// *********************************************************************************************************************

type DeleteQuestionRequest struct {
	ID int64
}

type DeleteQuestionResponse struct {
	Err error `json:"err,omitempty"`
}

// *********************************************************************************************************************

type QuestionsEndpoints struct {
	GetQuestionsEndpoint        endpoint.Endpoint
	GetQuestionTypesEndpoint    endpoint.Endpoint
	GetQuestionStatusesEndpoint endpoint.Endpoint
	PostQuestionEndpoint        endpoint.Endpoint
	PutQuestionEndpoint         endpoint.Endpoint
	PutQuestionModerateEndpoint endpoint.Endpoint
	DeleteQuestionEndpoint      endpoint.Endpoint
}

func MakeQuestionsEndpoints(s model.Questions) QuestionsEndpoints {
	return QuestionsEndpoints{
		GetQuestionsEndpoint:        MakeGetQuestionsEndpoint(s),
		GetQuestionTypesEndpoint:    MakeGetQuestionTypesEndpoint(s),
		GetQuestionStatusesEndpoint: MakeGetQuestionStatusesEndpoint(s),
		PostQuestionEndpoint:        MakePostQuestionEndpoint(s),
		PutQuestionEndpoint:         MakePutQuestionEndpoint(s),
		PutQuestionModerateEndpoint: MakePutQuestionModerateEndpoint(s),
		DeleteQuestionEndpoint:      MakeDeleteQuestionEndpoint(s),
	}
}

func MakeGetQuestionsEndpoint(s model.Questions) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetQuestionsRequest) //TODO check everywhere, return internal server error
		t, err := s.GetQuestions(ctx, req.CreatorUserID, req.SubjectID, req.StatusID)
		return GetQuestionsResponse{t, err}, err
	}
}

func MakeGetQuestionTypesEndpoint(s model.Questions) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		t, err := s.GetQuestionTypes(ctx)
		return GetQuestionTypesResponse{
			Types: t,
			Err:   err,
		}, err
	}
}

func MakeGetQuestionStatusesEndpoint(s model.Questions) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		st, err := s.GetQuestionStatuses(ctx)
		return GetQuestionStatusesResponse{
			Statuses: st,
			Err:      err,
		}, err
	}
}

func MakePostQuestionEndpoint(s model.Questions) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(PostQuestionRequest)
		id, err := s.AddQuestion(ctx, req.Question)
		return PostQuestionResponse{
			ID:  id,
			Err: err,
		}, err
	}
}

func MakePutQuestionEndpoint(s model.Questions) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(PutQuestionRequest)
		err := s.UpdateQuestionByID(ctx, req.ID, req.Question)
		return PutQuestionResponse{err}, err
	}
}

func MakePutQuestionModerateEndpoint(s model.Questions) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(PutQuestionModerateRequest)
		err := s.ModerateQuestion(ctx, req.ID, req.Approve, req.Comment)
		return PutQuestionModerateResponse{err}, err
	}
}

func MakeDeleteQuestionEndpoint(s model.Questions) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteQuestionRequest)
		err := s.DeleteQuestion(ctx, req.ID)
		return DeleteQuestionResponse{err}, err
	}
}
