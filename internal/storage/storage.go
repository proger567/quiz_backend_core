package storage

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"quiz_backend_core/internal/model"
	"quiz_backend_core/internal/storage/pg"
)

type Storages struct {
	Subjects  model.SubjectsStorage
	Questions model.QuestionsStorage
	Quizzes   model.QuizzesStorage

	pool *pgxpool.Pool
}

func NewStorages(ctx context.Context, databaseDSN string) (*Storages, error) {
	pool, err := pgxpool.New(ctx, databaseDSN)
	if err != nil {
		return nil, err
	}

	// TODO
	//if err = storage.DB.Ping(); err != nil {
	//	log.Fatal(err)
	//}

	return &Storages{
		Subjects:  pg.NewSubjectsStorage(pool),
		Questions: pg.NewQuestionsStorage(pool),
		Quizzes:   pg.NewQuizzesStorage(pool),

		pool: pool,
	}, nil
}

func (s *Storages) Close() {
	s.pool.Close()
}
