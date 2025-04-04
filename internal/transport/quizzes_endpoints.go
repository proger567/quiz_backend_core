package transport

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"quiz_backend_core/internal/dto"
	"quiz_backend_core/internal/model"
)

type GetQuizzesRequest struct {
	CreatorUserID int64 `json:"creator_user_id"`
}

type GetQuizzesResponse struct {
	Quizzes []dto.Quiz `json:"quizzes"`
	Err     error      `json:"err,omitempty"`
}

// *********************************************************************************************************************

type GetQuestionsByQuizIDRequest struct {
	QuizID int64 `json:"quiz_id"`
}

type GetQuestionsByQuizIDResponse struct {
	Questions []dto.Question `json:"questions"`
	Err       error          `json:"err,omitempty"`
}

// *********************************************************************************************************************

type GetQuizByIDRequest struct {
	QuizID int64 `json:"quiz_id"`
}

type GetQuizByIDResponse struct {
	Quiz dto.Quiz `json:"quiz"`
	Err  error    `json:"err,omitempty"`
}

// *********************************************************************************************************************

type PostQuizRequest struct {
	Quiz dto.InputQuiz `json:"quiz"`
}

type PostQuizResponse struct {
	QuizID int64 `json:"quiz_id"`
	Err    error `json:"err,omitempty"`
}

// *********************************************************************************************************************

type DeleteQuizByIDRequest struct {
	QuizID int64 `json:"quiz_id"`
}

type DeleteQuizByIDResponse struct {
	Err error `json:"err,omitempty"`
}

// *********************************************************************************************************************

type QuizzesEndpoints struct {
	GetQuizzesEndpoint           endpoint.Endpoint
	GetQuestionsByQuizIDEndpoint endpoint.Endpoint
	GetQuizByIDEndpoint          endpoint.Endpoint
	PostQuizEndpoint             endpoint.Endpoint
	DeleteQuizEndpoint           endpoint.Endpoint
}

func MakeQuizzesEndpoints(s model.Quizzes) QuizzesEndpoints {
	return QuizzesEndpoints{
		GetQuizzesEndpoint:           MakeGetQuizzesEndpoint(s),
		GetQuestionsByQuizIDEndpoint: MakeGetQuestionsByQuizIDEndpoint(s),
		GetQuizByIDEndpoint:          MakeGetQuizByIDEndpoint(s),
		PostQuizEndpoint:             MakePostQuizEndpoint(s),
		DeleteQuizEndpoint:           MakeDeleteQuizEndpoint(s),
	}
}

func MakeGetQuizzesEndpoint(s model.Quizzes) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetQuizzesRequest) //TODO check everywhere, return internal server error
		quizzes, err := s.GetQuizzes(ctx, req.CreatorUserID)
		return GetQuizzesResponse{
			Quizzes: quizzes,
			Err:     err,
		}, err
	}
}

func MakeGetQuestionsByQuizIDEndpoint(s model.Quizzes) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetQuestionsByQuizIDRequest) //TODO check everywhere, return internal server error
		questions, err := s.GetQuestionsByQuizID(ctx, req.QuizID)
		return GetQuestionsByQuizIDResponse{
			Questions: questions,
			Err:       err,
		}, err
	}
}

func MakeGetQuizByIDEndpoint(s model.Quizzes) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetQuizByIDRequest) //TODO check everywhere, return internal server error
		quiz, err := s.GetQuizByID(ctx, req.QuizID)
		return GetQuizByIDResponse{
			Quiz: quiz,
			Err:  err,
		}, err
	}
}

func MakePostQuizEndpoint(s model.Quizzes) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(PostQuizRequest) //TODO check everywhere, return internal server error
		quizID, err := s.AddQuiz(ctx, req.Quiz)
		return PostQuizResponse{
			QuizID: quizID,
			Err:    err,
		}, err
	}
}

func MakeDeleteQuizEndpoint(s model.Quizzes) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteQuizByIDRequest) //TODO check everywhere, return internal server error
		err := s.DeleteQuizByID(ctx, req.QuizID)
		return DeleteQuizByIDResponse{
			Err: err,
		}, err
	}
}
