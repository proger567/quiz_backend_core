package service

import (
	"context"
	"github.com/sirupsen/logrus"
	"quiz_backend_core/internal/dto"
	"quiz_backend_core/internal/model"
	"quiz_backend_core/internal/service/middleware"
)

type subjectsService struct {
	storage model.SubjectsStorage
	logger  *logrus.Logger
}

func NewSubjectsService(deps Deps) model.Subjects {
	var svc model.Subjects = subjectsService{
		storage: deps.Storages.Subjects,
		logger:  deps.Logger,
	}

	// middleware services
	svc = middleware.LoggingSubjectsMiddleware(deps.Logger)(svc)
	svc = middleware.InstrumentingSubjectsMiddleware(deps.RequestCounter, deps.RequestLatencyMeter)(svc)

	return svc
}

func (s subjectsService) GetSubjects(ctx context.Context) ([]dto.Subject, error) {
	return s.storage.GetSubjects(ctx)
}

func (s subjectsService) AddSubject(ctx context.Context, subject dto.Subject) (int64, error) {
	return s.storage.AddSubject(ctx, subject)
}

func (s subjectsService) GetStatistic(ctx context.Context, userId int64) (dto.Statistic, error) {
	return s.storage.GetStatistic(ctx, userId)
}

func (s subjectsService) UpdateSubject(ctx context.Context, subject dto.Subject) error {
	return s.storage.UpdateSubject(ctx, subject)
}

func (s subjectsService) DeleteSubjectByID(ctx context.Context, subjectID int64) error {
	return s.storage.DeleteSubjectByID(ctx, subjectID)
}
