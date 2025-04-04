package service

import (
	"github.com/go-kit/kit/metrics"
	"github.com/sirupsen/logrus"
	"quiz_backend_core/internal/model"
	"quiz_backend_core/internal/storage"
)

// TODO move to better place?
type ContextVariables string

const (
	ContextVariablesUserID   = "user_id"
	ContextVariablesUserRole = "user_role"
)

type Services struct {
	Subjects  model.Subjects
	Questions model.Questions
	Quizzes   model.Quizzes
}

type Deps struct {
	Storages            *storage.Storages
	Logger              *logrus.Logger
	RequestCounter      metrics.Counter
	RequestLatencyMeter metrics.Histogram
}

func NewServices(deps Deps) *Services {
	subjects := NewSubjectsService(deps)
	questions := NewQuestionsService(deps)
	quizzes := NewQuizzesService(deps)
	return &Services{
		Subjects:  subjects,
		Questions: questions,
		Quizzes:   quizzes,
	}
}
