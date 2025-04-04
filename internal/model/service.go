package model

import (
	"context"
	"quiz_backend_core/internal/dto"
)

type Subjects interface {
	GetSubjects(ctx context.Context) ([]dto.Subject, error)
	AddSubject(ctx context.Context, subject dto.Subject) (int64, error)
	UpdateSubject(ctx context.Context, subject dto.Subject) error
	DeleteSubjectByID(ctx context.Context, id int64) error

	GetStatistic(ctx context.Context, userId int64) (dto.Statistic, error)
}

type Questions interface {
	GetQuestions(ctx context.Context, creatorUserID, subjectId, statusID int64) ([]dto.Question, error)
	GetQuestionTypes(ctx context.Context) ([]dto.QuestionType, error)
	GetQuestionStatuses(ctx context.Context) ([]dto.QuestionStatus, error)
	AddQuestion(ctx context.Context, question dto.InputQuestion) (int64, error)
	UpdateQuestionByID(ctx context.Context, questionID int64, question dto.InputQuestion) error
	ModerateQuestion(ctx context.Context, ID int64, approve bool, comment string) error
	DeleteQuestion(ctx context.Context, ID int64) error
}

type Quizzes interface {
	GetQuizzes(ctx context.Context, creatorUserID int64) ([]dto.Quiz, error)
	GetQuestionsByQuizID(ctx context.Context, quizID int64) ([]dto.Question, error)
	GetQuizByID(ctx context.Context, quizID int64) (dto.Quiz, error)
	AddQuiz(ctx context.Context, quiz dto.InputQuiz) (int64, error)
	DeleteQuizByID(ctx context.Context, quizID int64) error
}

type SubjectsMiddleware func(Subjects) Subjects

type QuestionsMiddleware func(Questions) Questions
