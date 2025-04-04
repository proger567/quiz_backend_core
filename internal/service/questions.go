package service

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"quiz_backend_core/internal/dto"
	"quiz_backend_core/internal/model"
	"quiz_backend_core/internal/service/middleware"
	"time"
)

type questionsService struct {
	storage model.QuestionsStorage
	logger  *logrus.Logger
}

func NewQuestionsService(deps Deps) model.Questions {
	var svc model.Questions = questionsService{
		storage: deps.Storages.Questions,
		logger:  deps.Logger,
	}

	// middleware services
	svc = middleware.LoggingQuestionsMiddleware(deps.Logger)(svc)
	svc = middleware.InstrumentingQuestionsMiddleware(deps.RequestCounter, deps.RequestLatencyMeter)(svc)

	return svc
}

func (s questionsService) GetQuestions(ctx context.Context, creatorUserID, subjectId, statusID int64) ([]dto.Question, error) {
	return s.storage.GetQuestions(ctx, creatorUserID, subjectId, statusID)
}

func (s questionsService) GetQuestionTypes(ctx context.Context) ([]dto.QuestionType, error) {
	return s.storage.GetQuestionTypes(ctx)
}

func (s questionsService) GetQuestionStatuses(ctx context.Context) ([]dto.QuestionStatus, error) {
	return s.storage.GetQuestionStatuses(ctx)
}

func (s questionsService) AddQuestion(ctx context.Context, question dto.InputQuestion) (int64, error) {
	userID, ok1 := ctx.Value(ContextVariablesUserID).(int64)
	userRole, ok2 := ctx.Value(ContextVariablesUserRole).(dto.Role) // TODO ALWAYS GOOD
	if !ok1 || !ok2 {
		return -1, errors.New("bad user id or role") //TODO
	}

	if userRole == dto.RoleAdmin || userRole == dto.RoleModerator {
		question.StatusName = dto.QuestionStatusNameApproved
		question.ModeratorUserID = userID
		question.ModeratedAt = time.Now().Format("2006-01-02 15:04")
	} else {
		question.StatusName = dto.QuestionStatusNameCreated
		question.ModeratorUserID = -1
		question.ModeratedAt = "" //TODO
	}

	question.CreatorUserID = userID

	return s.storage.AddQuestion(ctx, question)
}

func (s questionsService) UpdateQuestionByID(ctx context.Context, questionID int64, question dto.InputQuestion) error {
	//TODO check if not changed

	//s.GetQuestions()

	userID := ctx.Value(ContextVariablesUserID).(int64)
	userRole := ctx.Value(ContextVariablesUserRole).(dto.Role)

	if userRole == dto.RoleAdmin || userRole == dto.RoleModerator {
		question.StatusName = dto.QuestionStatusNameApproved
		question.ModeratorUserID = userID
		question.ModeratedAt = time.Now().Format("2006-01-02 15:04")
	} else {
		question.StatusName = dto.QuestionStatusNameCreated
		question.ModeratorUserID = -1
		question.ModeratedAt = "" //TODO
	}

	return s.storage.UpdateQuestionByID(ctx, questionID, question)
}

func (s questionsService) ModerateQuestion(ctx context.Context, ID int64, approve bool, comment string) error {
	//TODO create notification with comment
	if approve {
		return s.storage.UpdateQuestionStatus(ctx, ID, dto.QuestionStatusNameApproved)
	} else {
		return s.storage.UpdateQuestionStatus(ctx, ID, dto.QuestionStatusNameDeclined)
	}
}

func (s questionsService) DeleteQuestion(ctx context.Context, ID int64) error {
	return s.storage.DeleteQuestion(ctx, ID)
}
