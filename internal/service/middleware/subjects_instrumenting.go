package middleware

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/metrics"
	"quiz_backend_core/internal/dto"
	"quiz_backend_core/internal/model"
	"time"
)

func InstrumentingSubjectsMiddleware(requestCount metrics.Counter, requestLatency metrics.Histogram) model.SubjectsMiddleware {
	return func(next model.Subjects) model.Subjects {
		return instrumentingSubjectsMiddleware{
			requestCount,
			requestLatency,
			next,
		}
	}
}

type instrumentingSubjectsMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	next           model.Subjects
}

func (im instrumentingSubjectsMiddleware) GetSubjects(ctx context.Context) (subjects []dto.Subject, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "getSubjects", "error", fmt.Sprint(err != nil)}
		im.requestCount.With(lvs...).Add(1)
		im.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	subjects, err = im.next.GetSubjects(ctx)
	return
}

func (im instrumentingSubjectsMiddleware) AddSubject(ctx context.Context, subject dto.Subject) (id int64, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "AddSubject", "error", fmt.Sprint(err != nil)}
		im.requestCount.With(lvs...).Add(1)
		im.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	id, err = im.next.AddSubject(ctx, subject)
	return
}

func (im instrumentingSubjectsMiddleware) GetStatistic(ctx context.Context, userId int64, userRole dto.Role) (statistic dto.Statistic, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "getStatistic", "error", fmt.Sprint(err != nil)}
		im.requestCount.With(lvs...).Add(1)
		im.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	statistic, err = im.next.GetStatistic(ctx, userId, userRole)
	return
}

func (im instrumentingSubjectsMiddleware) UpdateSubject(ctx context.Context, subject dto.Subject) (err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "UpdateSubject", "error", fmt.Sprint(err != nil)}
		im.requestCount.With(lvs...).Add(1)
		im.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	err = im.next.UpdateSubject(ctx, subject)
	return
}

func (im instrumentingSubjectsMiddleware) DeleteSubjectByID(ctx context.Context, subjectID int64) (err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "DeleteSubject", "error", fmt.Sprint(err != nil)}
		im.requestCount.With(lvs...).Add(1)
		im.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	err = im.next.DeleteSubjectByID(ctx, subjectID)
	return
}
