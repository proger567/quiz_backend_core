package service

import (
	"github.com/go-kit/kit/metrics"
	"github.com/sirupsen/logrus"
	"quiz_backend_core/internal/model"
	"quiz_backend_core/internal/storage"
)

type Services struct {
	Subjects  model.Subjects
	Questions model.Questions
	Quizzes   model.Quizzes
}

type Deps struct {
	Storages            *storage.Storages
	Logger              *logrus.Logger //TODO interface
	RequestCounter      metrics.Counter
	RequestLatencyMeter metrics.Histogram
	Notifier            model.Notifier //TODO interface
}

func NewServices(deps Deps) *Services {
	subjects := NewSubjectsService(deps) //TODO раскрыть deps для каждого сервиса
	questions := NewQuestionsService(deps)
	quizzes := NewQuizzesService(deps)
	return &Services{
		Subjects:  subjects,
		Questions: questions,
		Quizzes:   quizzes,
	}
}
