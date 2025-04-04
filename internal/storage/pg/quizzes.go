package pg

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"quiz_backend_core/internal/dto"
	storage_errors "quiz_backend_core/internal/storage/errors"
	"strings"
)

func NewQuizzesStorage(conn *pgxpool.Pool) *QuizzesStorage {
	return &QuizzesStorage{
		conn: conn,
	}
}

type QuizzesStorage struct {
	conn *pgxpool.Pool
}

func (q QuizzesStorage) GetQuizzes(ctx context.Context, creatorUserID int64) ([]dto.Quiz, error) {
	query := `
        SELECT json_build_object(
			'id', q.id::TEXT,
			'name', q.name,
			'description', q.description,
			'creator', json_build_object(
				'id', cua.id::TEXT,
				'login', cua.login,
				'user_name', cua.user_name
			),
			'created_at', q.created_at,
			'updated_at', q.updated_at,
		    'question_ids', (
				SELECT json_agg(question_id::TEXT) FROM quizzes_questions qq WHERE qq.quiz_id = q.id
		    )
		)
		FROM quiz q
		LEFT JOIN user_account cua on cua.id = q.creator_user_id
        %s
	`

	//conditions
	var conditions []string
	if creatorUserID != -1 {
		conditions = append(conditions, fmt.Sprintf("q.creator_user_id=%d", creatorUserID))
	}

	allConditions := ""
	if len(conditions) != 0 {
		allConditions = fmt.Sprintf("WHERE %s", strings.Join(conditions, " and "))
	}
	query = fmt.Sprintf(query, allConditions)

	//do request
	var quizzes []dto.Quiz
	rows, err := q.conn.Query(ctx, query)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows): //it is not error
			return quizzes, nil
		default:
			return quizzes, &storage_errors.ExecutionPSQLError{Err: err}
		}
	}
	defer rows.Close()

	for rows.Next() {
		var res string
		if err := rows.Scan(&res); err != nil {
			return quizzes, &storage_errors.ScanPSQLResultsError{Err: err}
		}

		var result dto.Quiz
		if err := json.Unmarshal([]byte(res), &result); err != nil {
			return quizzes, &storage_errors.UnmarshalPSQLResultsError{Err: err}
		}
		quizzes = append(quizzes, result)
	}

	return quizzes, nil
}

func (q QuizzesStorage) GetQuestionsByQuizID(ctx context.Context, quizID int64) ([]dto.Question, error) {
	query := `
		SELECT
			json_build_object(
				'id', q.id::TEXT,
				'text', q.text,
				'code', q.code,
				'variants', q.variants,
				'answer', q.answer ,
				'type', json_build_object(
					'id', qt.id::TEXT,
					'name', qt.name
				),
				'status', qs.name,
				'subject_id', q.subject_id::TEXT,
				'subject_name', s.name,
				'creator', json_build_object(
					'id', cua.id::TEXT,
					'login', cua.login,
					'user_name', cua.user_name
				),
				'moderator', json_build_object(
					'id', mua.id::TEXT,
					'login', mua.login,
					'user_name', mua.user_name
				),
				'moderated_at', q.moderated_at,
				'created_at', q.created_at
			)
		FROM quizzes_questions qq
		JOIN question q on qq.question_id = q.id
		LEFT JOIN question_type qt on qt.id = q.type_id
		LEFT JOIN question_status qs on qs.id = q.status_id
		LEFT JOIN subject s on s.id = q.subject_id
		LEFT JOIN user_account cua on cua.id = q.creator_user_id
		LEFT JOIN user_account mua on mua.id = q.moderator_user_id
		WHERE qq.quiz_id = $1
	`
	var questions = []dto.Question{}
	rows, err := q.conn.Query(ctx, query, quizID)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows): //it is not error
			return questions, nil
		default:
			return questions, &storage_errors.ExecutionPSQLError{Err: err}
		}
	}
	defer rows.Close()

	for rows.Next() {
		var res string
		if err := rows.Scan(&res); err != nil {
			return questions, &storage_errors.ScanPSQLResultsError{Err: err}
		}

		var result dto.Question
		if err := json.Unmarshal([]byte(res), &result); err != nil {
			return questions, &storage_errors.UnmarshalPSQLResultsError{Err: err}
		}
		questions = append(questions, result)
	}
	return questions, nil
}

func (q QuizzesStorage) GetQuizByID(ctx context.Context, quizID int64) (dto.Quiz, error) {
	query := `
        SELECT json_build_object(
			'id', q.id::TEXT,
			'name', q.name,
			'description', q.description,
			'creator', json_build_object(
				'id', cua.id::TEXT,
				'login', cua.login,
				'user_name', cua.user_name
			),
			'created_at', q.created_at,
			'updated_at', q.updated_at,
		    'question_ids', (
				SELECT json_agg(question_id) FROM quizzes_questions qq WHERE qq.quiz_id = q.id
		    )
		)
		FROM quiz q
		LEFT JOIN user_account cua on cua.id = q.creator_user_id
	    WHERE q.id = $1
	`

	var quiz dto.Quiz
	err := q.conn.QueryRow(ctx, query, quizID).Scan(&quiz)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows): //it is not error
			return quiz, nil
		default:
			return quiz, &storage_errors.ExecutionPSQLError{Err: err}
		}
	}

	return quiz, nil
}

func (q QuizzesStorage) AddQuiz(ctx context.Context, quiz dto.InputQuiz) (int64, error) {
	addQuizQuery := `
		INSERT INTO
		    quiz (
		    	name,
		     	description,
		     	creator_user_id
		    )
		VALUES (
		    $1, $2, $3
		) 
		RETURNING ID
    `

	//transaction
	tx, err := q.conn.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.RepeatableRead, //TODO check
		AccessMode: pgx.ReadWrite,      //TODO check
	})
	if err != nil {
		return -1, &storage_errors.ExecutionPSQLError{Err: fmt.Errorf("BeginTx failed: %v\n", err)}
	}
	defer tx.Rollback(ctx)

	//first request
	args := []interface{}{
		quiz.Name,
		quiz.Description,
		quiz.CreatorUserID,
	}

	var quizID int64 = -1
	if err = tx.QueryRow(ctx, addQuizQuery, args...).Scan(&quizID); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			return quizID, &storage_errors.AlreadyExistsError{Err: pgErr}
		} else {
			return quizID, &storage_errors.ExecutionPSQLError{Err: err}
		}
	}

	if quizID == -1 {
		return quizID, &storage_errors.ExecutionPSQLError{Err: errors.New("quiz id does not exist")}
	}

	//second request
	rows := make([][]interface{}, len(quiz.QuestionIDs))
	for i, questionID := range quiz.QuestionIDs {
		rows[i] = []interface{}{quizID, questionID}
	}

	copyCount, err := tx.CopyFrom(
		ctx,
		pgx.Identifier{"quizzes_questions"}, // Имя таблицы
		[]string{"quiz_id", "question_id"},  // Имена столбцов
		pgx.CopyFromRows(rows),              // Данные
	)
	if err != nil {
		return -1, &storage_errors.ExecutionPSQLError{Err: fmt.Errorf("CopyFrom failed: %v\n", err)}
	}

	if int(copyCount) != len(rows) {
		return -1, &storage_errors.ExecutionPSQLError{Err: fmt.Errorf("question to add: %v, question added: %v \n", len(rows), copyCount)}
	}

	if err = tx.Commit(ctx); err != nil {
		return -1, &storage_errors.ExecutionPSQLError{Err: err}
	}

	return quizID, nil
}

func (q QuizzesStorage) DeleteQuizByID(ctx context.Context, quizID int64) error {
	removeQuizzesQuestionsQuery := `
		DELETE FROM 
		   quizzes_questions
	    WHERE quiz_id = $1	
	`

	removeQuizQuery := `
		DELETE FROM 
		   quiz
	    WHERE id = $1	
	`

	//transaction
	tx, err := q.conn.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.RepeatableRead, //TODO check
		AccessMode: pgx.ReadWrite,      //TODO check
	})
	if err != nil {
		return &storage_errors.ExecutionPSQLError{Err: fmt.Errorf("BeginTx failed: %v\n", err)}
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, removeQuizzesQuestionsQuery, quizID)
	if err != nil {
		return &storage_errors.ExecutionPSQLError{Err: fmt.Errorf("Exec failed: %v\n", err)}
	}

	_, err = tx.Exec(ctx, removeQuizQuery, quizID)
	if err != nil {
		return &storage_errors.ExecutionPSQLError{Err: fmt.Errorf("Exec failed: %v\n", err)}
	}

	if err = tx.Commit(ctx); err != nil {
		return &storage_errors.ExecutionPSQLError{Err: err}
	}

	return nil
}
