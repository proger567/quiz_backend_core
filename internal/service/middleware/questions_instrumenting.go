package middleware

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/metrics"
	"quiz_backend_core/internal/dto"
	"quiz_backend_core/internal/model"
	"time"
)

func InstrumentingQuestionsMiddleware(requestCount metrics.Counter, requestLatency metrics.Histogram) model.QuestionsMiddleware {
	return func(next model.Questions) model.Questions {
		return instrumentingQuestionsMiddleware{
			requestCount,
			requestLatency,
			next,
		}
	}
}

type instrumentingQuestionsMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	next           model.Questions
}

func (im instrumentingQuestionsMiddleware) GetQuestions(ctx context.Context, creatorUserID, subjectID, statusID int64) (questions []dto.Question, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "getQuestions", "error", fmt.Sprint(err != nil)}
		im.requestCount.With(lvs...).Add(1)
		im.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	questions, err = im.next.GetQuestions(ctx, creatorUserID, subjectID, statusID)
	return
}

func (im instrumentingQuestionsMiddleware) GetQuestionTypes(ctx context.Context) (types []dto.QuestionType, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "getTypes", "error", fmt.Sprint(err != nil)}
		im.requestCount.With(lvs...).Add(1)
		im.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	types, err = im.next.GetQuestionTypes(ctx)
	return
}

func (im instrumentingQuestionsMiddleware) GetQuestionStatuses(ctx context.Context) (types []dto.QuestionStatus, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "getTypes", "error", fmt.Sprint(err != nil)}
		im.requestCount.With(lvs...).Add(1)
		im.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	types, err = im.next.GetQuestionStatuses(ctx)
	return
}

func (im instrumentingQuestionsMiddleware) AddQuestion(ctx context.Context, questionAdd dto.InputQuestion) (id int64, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "addQuestion", "error", fmt.Sprint(err != nil)}
		im.requestCount.With(lvs...).Add(1)
		im.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	id, err = im.next.AddQuestion(ctx, questionAdd)
	return
}

func (im instrumentingQuestionsMiddleware) UpdateQuestionByID(ctx context.Context, questionID int64, question dto.InputQuestion) (err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "updateQuestionByID", "error", fmt.Sprint(err != nil)}
		im.requestCount.With(lvs...).Add(1)
		im.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	err = im.next.UpdateQuestionByID(ctx, questionID, question)
	return
}

func (im instrumentingQuestionsMiddleware) ModerateQuestion(ctx context.Context, questionID int64, approve bool, comment string) (err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "approveQuestion", "error", fmt.Sprint(err != nil)}
		im.requestCount.With(lvs...).Add(1)
		im.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	err = im.next.ModerateQuestion(ctx, questionID, approve, comment)
	return
}

func (im instrumentingQuestionsMiddleware) DeleteQuestion(ctx context.Context, questionID int64) (err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "deleteQuestion", "error", fmt.Sprint(err != nil)}
		im.requestCount.With(lvs...).Add(1)
		im.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	err = im.next.DeleteQuestion(ctx, questionID)
	return
}
