package middleware

import (
	"context"
	"github.com/sirupsen/logrus"
	"quiz_backend_core/internal/dto"
	"quiz_backend_core/internal/model"
	"time"
)

func LoggingSubjectsMiddleware(logger *logrus.Logger) model.SubjectsMiddleware {
	return func(next model.Subjects) model.Subjects {
		return &loggingSubjectsMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingSubjectsMiddleware struct {
	next   model.Subjects
	logger *logrus.Logger
}

func (mw loggingSubjectsMiddleware) GetSubjects(ctx context.Context) (subjects []dto.Subject, err error) {
	defer func(begin time.Time) {
		mw.logger.WithFields(logrus.Fields{
			"took":  time.Since(begin).Milliseconds(),
			"error": err,
		}).Info("method == GetSubjects")
	}(time.Now())
	return mw.next.GetSubjects(ctx)
}

func (mw loggingSubjectsMiddleware) AddSubject(ctx context.Context, subject dto.Subject) (id int64, err error) {
	defer func(begin time.Time) {
		mw.logger.WithFields(logrus.Fields{
			"took":  time.Since(begin).Milliseconds(),
			"error": err,
		}).Info("method == AddSubject")
	}(time.Now())
	return mw.next.AddSubject(ctx, subject)
}

func (mw loggingSubjectsMiddleware) GetStatistic(ctx context.Context, userId int64, userRole dto.Role) (statistic dto.Statistic, err error) {
	defer func(begin time.Time) {
		mw.logger.WithFields(logrus.Fields{
			"took":  time.Since(begin).Milliseconds(),
			"error": err,
		}).Info("method == GetStatistic")
	}(time.Now())
	return mw.next.GetStatistic(ctx, userId, userRole)
}

func (mw loggingSubjectsMiddleware) UpdateSubject(ctx context.Context, subject dto.Subject) (err error) {
	defer func(begin time.Time) {
		mw.logger.WithFields(logrus.Fields{
			"took":  time.Since(begin).Milliseconds(),
			"error": err,
		}).Info("method == UpdateSubject")
	}(time.Now())
	return mw.next.UpdateSubject(ctx, subject)
}

func (mw loggingSubjectsMiddleware) DeleteSubjectByID(ctx context.Context, subjectID int64) (err error) {
	defer func(begin time.Time) {
		mw.logger.WithFields(logrus.Fields{
			"took":  time.Since(begin).Milliseconds(),
			"error": err,
		}).Info("method == DeleteSubjectByID")
	}(time.Now())
	return mw.next.DeleteSubjectByID(ctx, subjectID)
}
