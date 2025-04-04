package middleware

import (
	"context"
	"github.com/sirupsen/logrus"
	"quiz_backend_core/internal/dto"
	"quiz_backend_core/internal/model"
	"time"
)

func LoggingQuestionsMiddleware(logger *logrus.Logger) model.QuestionsMiddleware {
	return func(next model.Questions) model.Questions {
		return &loggingQuestionsMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingQuestionsMiddleware struct {
	next   model.Questions
	logger *logrus.Logger
}

func (mw loggingQuestionsMiddleware) GetQuestions(ctx context.Context, creatorUserID, subjectId, statusID int64) (question []dto.Question, err error) {
	defer func(begin time.Time) {
		mw.logger.WithFields(logrus.Fields{
			"took":  time.Since(begin).Milliseconds(),
			"error": err,
		}).Info("method == GetQuestions")
	}(time.Now())
	return mw.next.GetQuestions(ctx, creatorUserID, subjectId, statusID)
}

func (mw loggingQuestionsMiddleware) GetQuestionTypes(ctx context.Context) (types []dto.QuestionType, err error) {
	defer func(begin time.Time) {
		mw.logger.WithFields(logrus.Fields{
			"took":  time.Since(begin).Milliseconds(),
			"error": err,
		}).Info("method == GetTypes")
	}(time.Now())
	return mw.next.GetQuestionTypes(ctx)
}

func (mw loggingQuestionsMiddleware) GetQuestionStatuses(ctx context.Context) (types []dto.QuestionStatus, err error) {
	defer func(begin time.Time) {
		mw.logger.WithFields(logrus.Fields{
			"took":  time.Since(begin).Milliseconds(),
			"error": err,
		}).Info("method == GetStatuses")
	}(time.Now())
	return mw.next.GetQuestionStatuses(ctx)
}

func (mw loggingQuestionsMiddleware) AddQuestion(ctx context.Context, question dto.InputQuestion) (id int64, err error) {
	defer func(begin time.Time) {
		mw.logger.WithFields(logrus.Fields{
			"took":  time.Since(begin).Milliseconds(),
			"error": err,
			//"userID":   ctx.Value(ContextVariablesUserID),
			//"userRole": ctx.Value(ContextVariablesUserRole),
			"question": question,
		}).Info("method == AddQuestion")
	}(time.Now())
	return mw.next.AddQuestion(ctx, question)
}

func (mw loggingQuestionsMiddleware) UpdateQuestionByID(ctx context.Context, questionID int64, question dto.InputQuestion) (err error) {
	defer func(begin time.Time) {
		mw.logger.WithFields(logrus.Fields{
			"took":     time.Since(begin).Milliseconds(),
			"error":    err,
			"id":       questionID,
			"question": question,
		}).Info("method == UpdateQuestionByID")
	}(time.Now())
	return mw.next.UpdateQuestionByID(ctx, questionID, question)
}

func (mw loggingQuestionsMiddleware) ModerateQuestion(ctx context.Context, questionID int64, approve bool, comment string) (err error) {
	defer func(begin time.Time) {
		mw.logger.WithFields(logrus.Fields{
			"took":  time.Since(begin).Milliseconds(),
			"error": err,
		}).Info("method == ApproveQuestionByID")
	}(time.Now())
	return mw.next.ModerateQuestion(ctx, questionID, approve, comment)
}

func (mw loggingQuestionsMiddleware) DeleteQuestion(ctx context.Context, questionID int64) (err error) {
	defer func(begin time.Time) {
		mw.logger.WithFields(logrus.Fields{
			"took":  time.Since(begin).Milliseconds(),
			"error": err,
		}).Info("method == DeleteQuestion")
	}(time.Now())
	return mw.next.DeleteQuestion(ctx, questionID)
}
