package service

import (
	"context"
	"github.com/sirupsen/logrus"
	"quiz_backend_core/internal/dto"
	"quiz_backend_core/internal/model"
)

type quizzesService struct {
	storage model.QuizzesStorage
	logger  *logrus.Logger
}

func NewQuizzesService(deps Deps) model.Quizzes {
	var svc model.Quizzes = quizzesService{
		storage: deps.Storages.Quizzes,
		logger:  deps.Logger,
	}

	//TODO
	// middleware services
	//svc = middleware.LoggingQuestionsMiddleware(deps.Logger)(svc)
	//svc = middleware.InstrumentingQuestionsMiddleware(deps.RequestCounter, deps.RequestLatencyMeter)(svc)

	return svc
}

func (s quizzesService) GetQuizzes(ctx context.Context, creatorUserID int64) ([]dto.Quiz, error) {
	return s.storage.GetQuizzes(ctx, creatorUserID)
}

func (s quizzesService) GetQuestionsByQuizID(ctx context.Context, quizID int64) ([]dto.Question, error) {
	return s.storage.GetQuestionsByQuizID(ctx, quizID)
}

func (s quizzesService) GetQuizByID(ctx context.Context, quizID int64) (dto.Quiz, error) {
	return s.storage.GetQuizByID(ctx, quizID)
}

func (s quizzesService) AddQuiz(ctx context.Context, quiz dto.InputQuiz) (int64, error) {
	quiz.CreatorUserID = ctx.Value(ContextVariablesUserID).(int64)
	return s.storage.AddQuiz(ctx, quiz)
}

func (s quizzesService) DeleteQuizByID(ctx context.Context, quizID int64) error {
	return s.storage.DeleteQuizByID(ctx, quizID)
}
