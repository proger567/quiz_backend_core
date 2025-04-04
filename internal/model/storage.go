package model

import (
	"context"
	"quiz_backend_core/internal/dto"
)

type Pinger interface {
	PingDB() error
}

type SubjectsStorage interface {
	GetSubjects(ctx context.Context) ([]dto.Subject, error)
	AddSubject(ctx context.Context, subject dto.Subject) (int64, error)
	UpdateSubject(ctx context.Context, subject dto.Subject) error
	DeleteSubjectByID(ctx context.Context, id int64) error
	GetStatistic(ctx context.Context, userId int64) (dto.Statistic, error)
}

type QuestionsStorage interface {
	GetQuestions(ctx context.Context, creatorUserID, subjectID, statusID int64) ([]dto.Question, error)
	GetQuestionByID(ctx context.Context, questionID int64) (dto.Question, error)
	GetQuestionTypes(ctx context.Context) ([]dto.QuestionType, error)
	GetQuestionStatuses(ctx context.Context) ([]dto.QuestionStatus, error)
	AddQuestion(ctx context.Context, question dto.InputQuestion) (int64, error)
	UpdateQuestionByID(ctx context.Context, ID int64, question dto.InputQuestion) error
	UpdateQuestionStatus(ctx context.Context, ID int64, status dto.QuestionStatusName) error
	DeleteQuestion(ctx context.Context, ID int64) error
}

type QuizzesStorage interface {
	GetQuizzes(ctx context.Context, creatorUserID int64) ([]dto.Quiz, error)
	GetQuestionsByQuizID(ctx context.Context, quizID int64) ([]dto.Question, error)
	GetQuizByID(ctx context.Context, quizID int64) (dto.Quiz, error)
	AddQuiz(ctx context.Context, quiz dto.InputQuiz) (int64, error)
	DeleteQuizByID(ctx context.Context, quizID int64) error
}
