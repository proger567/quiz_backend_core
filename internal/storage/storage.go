package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"quiz_backend_core/internal/model"
	"quiz_backend_core/internal/storage/pg"
)

type Storages struct {
	Subjects  model.SubjectsStorage
	Questions model.QuestionsStorage
	Quizzes   model.QuizzesStorage
}

// TODO use pool
func NewStorages(ctx context.Context, databaseDSN string) (*Storages, error) {
	pool, err := pgxpool.New(ctx, databaseDSN)
	if err != nil {
		log.Fatal(fmt.Errorf("Unable to connect to database: %v\n", err))
	}

	// TODO
	//if err = storage.DB.Ping(); err != nil {
	//	log.Fatal(err)
	//}

	return &Storages{
		Subjects:  pg.NewSubjectsStorage(pool),
		Questions: pg.NewQuestionsStorage(pool),
		Quizzes:   pg.NewQuizzesStorage(pool),
	}, nil
}
