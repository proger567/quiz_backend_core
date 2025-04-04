package pg

import (
	"context"
	"database/sql"
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

func NewQuestionsStorage(conn *pgxpool.Pool) *QuestionsStorage {
	return &QuestionsStorage{
		conn: conn,
	}
}

type QuestionsStorage struct {
	conn *pgxpool.Pool
}

func (q QuestionsStorage) GetQuestions(ctx context.Context, creatorUserID, subjectID, statusID int64) ([]dto.Question, error) {
	var questions = []dto.Question{}
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
		FROM question q
		LEFT JOIN question_type qt on qt.id = q.type_id
		LEFT JOIN question_status qs on qs.id = q.status_id
		LEFT JOIN subject s on s.id = q.subject_id
		LEFT JOIN user_account cua on cua.id = q.creator_user_id
		LEFT JOIN user_account mua on mua.id = q.moderator_user_id
		%s
	`

	//TODO
	//preparedStmt, err := s.conn.Prepare(ctx, "GetQuestions", query)
	//if err != nil {
	//	return questions, &storage_errors.StatementPSQLError{Err: err}
	//}

	var conditions []string
	if creatorUserID != -1 {
		conditions = append(conditions, fmt.Sprintf("q.creator_user_id=%d", creatorUserID))
	}

	if subjectID != -1 {
		conditions = append(conditions, fmt.Sprintf("q.subject_id=%d", subjectID))
	}

	if statusID != -1 {
		conditions = append(conditions, fmt.Sprintf("q.status_id=%d", statusID))
	}

	allConditions := ""
	if len(conditions) != 0 {
		allConditions = fmt.Sprintf("WHERE %s", strings.Join(conditions, " and "))
	}
	query = fmt.Sprintf(query, allConditions)

	rows, err := q.conn.Query(ctx /*preparedStmt.Name*/, query)
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

func (q QuestionsStorage) GetQuestionByID(ctx context.Context, questionID int64) (dto.Question, error) {
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
		FROM question q
		LEFT JOIN question_type qt on qt.id = q.type_id
		LEFT JOIN question_status qs on qs.id = q.status_id
		LEFT JOIN subject s on s.id = q.subject_id
		LEFT JOIN user_account cua on cua.id = q.creator_user_id
		LEFT JOIN user_account mua on mua.id = q.moderator_user_id
		WHERE qt.id = $1
	`

	var question = dto.Question{}
	if err := q.conn.QueryRow(ctx, query, questionID).Scan(&question); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgerrcode.IsCaseNotFound(pgErr.Code) {
			return question, &storage_errors.NotFoundError{Err: dto.ErrQuestionNotFound}
		} else {
			return question, &storage_errors.ExecutionPSQLError{Err: err}
		}
	}
	return question, nil
}

func (q QuestionsStorage) GetQuestionTypes(ctx context.Context) ([]dto.QuestionType, error) {
	var types = []dto.QuestionType{}
	query := `		
		SELECT
			json_build_object(
				'id', id::TEXT,
				'name', name
			)
		FROM question_type
	`

	//preparedStmt, err := q.conn.Prepare(ctx, "GetQuestionTypes", query)
	//if err != nil {
	//	return types, &storage_errors.StatementPSQLError{Err: err}
	//}

	rows, err := q.conn.Query(ctx /*preparedStmt.Name*/, query)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows): //it is not error
			return types, nil
		default:
			return types, &storage_errors.ExecutionPSQLError{Err: err}
		}
	}
	defer rows.Close()

	for rows.Next() {
		var res string
		if err := rows.Scan(&res); err != nil {
			return types, &storage_errors.ScanPSQLResultsError{Err: err}
		}

		var result dto.QuestionType
		if err := json.Unmarshal([]byte(res), &result); err != nil {
			return types, &storage_errors.UnmarshalPSQLResultsError{Err: err}
		}
		types = append(types, result)
	}
	return types, nil
}

func (q QuestionsStorage) GetQuestionStatuses(ctx context.Context) ([]dto.QuestionStatus, error) {
	var statuses = []dto.QuestionStatus{}
	query := `		
		SELECT
			json_build_object(
				'id', id::TEXT,
				'name', name
			)
		FROM question_status
	`

	rows, err := q.conn.Query(ctx, query)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows): //it is not error
			return statuses, nil
		default:
			return statuses, &storage_errors.ExecutionPSQLError{Err: err}
		}
	}
	defer rows.Close()

	for rows.Next() {
		var res string
		if err := rows.Scan(&res); err != nil {
			return statuses, &storage_errors.ScanPSQLResultsError{Err: err}
		}

		var result dto.QuestionStatus
		if err := json.Unmarshal([]byte(res), &result); err != nil {
			return statuses, &storage_errors.UnmarshalPSQLResultsError{Err: err}
		}
		statuses = append(statuses, result)
	}
	return statuses, nil
}

// TODO return id
func (q QuestionsStorage) AddQuestion(ctx context.Context, question dto.InputQuestion) (int64, error) {
	var questionID int64 = -1
	query := `
		INSERT INTO question (
			text,					-- 1
			code,					-- 2
			variants,				-- 3
			answer,					-- 4
			type_id,				-- 5
			status_id,				-- 6
			subject_id,				-- 7
			creator_user_id,		-- 8
			moderator_user_id,		-- 9
			moderated_at 			-- 10
		)
		values (
		    $1,
			$2,
			$3,
			$4,
			$5,
			(SELECT id FROM question_status WHERE name=$6),
			$7,
			$8,
			$9,
			$10
		)
		RETURNING id;
	`

	//preparedStmt, err := q.conn.Prepare(ctx, "AddQuestion", query)
	//if err != nil {
	//	return questionID, &storage_errors.StatementPSQLError{Err: err}
	//}

	tx, err := q.conn.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.RepeatableRead,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return questionID, &storage_errors.ExecutionPSQLError{Err: fmt.Errorf("BeginTx failed: %v\n", err)}
	}
	defer tx.Rollback(ctx)

	var moderatorUserID sql.NullInt64
	moderatorUserID.Int64 = question.ModeratorUserID
	moderatorUserID.Valid = question.ModeratorUserID > 0

	//TODO time instead of string
	var moderatedAt sql.NullString
	moderatedAt.String = question.ModeratedAt
	moderatedAt.Valid = question.ModeratedAt != ""

	// Порядок параметров должен соответствовать порядку в запросе
	args := []interface{}{
		question.Text,
		question.Code,
		question.Variants,
		question.Answer,
		question.TypeID,
		question.StatusName,
		question.SubjectID,
		question.CreatorUserID,
		moderatorUserID,
		moderatedAt,
	}

	if err := tx.QueryRow(ctx /*preparedStmt.Name*/, query, args...).Scan(&questionID); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			return questionID, &storage_errors.AlreadyExistsError{Err: pgErr}
		} else {
			return questionID, &storage_errors.ExecutionPSQLError{Err: err}
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return questionID, &storage_errors.ExecutionPSQLError{Err: err}
	}

	return questionID, nil
}

func (q QuestionsStorage) UpdateQuestionByID(ctx context.Context, ID int64, question dto.InputQuestion) error {
	query := `		
		UPDATE 
		    question
		SET
			text=$1,
			code=$2,
			variants=$3,
			answer=$4,
			type_id=$5,
			status_id=(SELECT id FROM question_status WHERE name=$6),
			subject_id=$7
		WHERE
		    id = $8
		`

	//preparedStmt, err := q.conn.Prepare(ctx, "UpdateQuestion", query)
	//if err != nil {
	//	return &storage_errors.StatementPSQLError{Err: err}
	//}

	tx, err := q.conn.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.RepeatableRead,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return &storage_errors.ExecutionPSQLError{Err: fmt.Errorf("BeginTx failed: %v\n", err)}
	}
	defer tx.Rollback(ctx)

	// Порядок параметров должен соответствовать порядку в запросе
	args := []interface{}{
		question.Text,
		question.Code,
		question.Variants,
		question.Answer,
		question.TypeID,
		question.StatusName,
		question.SubjectID,

		ID,
	}

	if _, err = tx.Exec(ctx /*preparedStmt.Name*/, query, args...); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgerrcode.IsCaseNotFound(pgErr.Code) {
			return &storage_errors.NotFoundError{Err: dto.ErrSubjectNotFound}
		} else {
			return &storage_errors.ExecutionPSQLError{Err: err}
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return &storage_errors.ExecutionPSQLError{Err: err}
	}

	return nil
}

func (q QuestionsStorage) UpdateQuestionStatus(ctx context.Context, ID int64, status dto.QuestionStatusName) error {
	query := `		
		UPDATE 
		    question
		SET
			status_id=(SELECT id FROM question_status WHERE name=$1)
		WHERE
		    id = $2
		`

	//preparedStmt, err := q.conn.Prepare(ctx, "UpdateQuestionStatus", query)
	//if err != nil {
	//	return &storage_errors.StatementPSQLError{Err: err}
	//}

	tx, err := q.conn.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.RepeatableRead,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return &storage_errors.ExecutionPSQLError{Err: fmt.Errorf("BeginTx failed: %v\n", err)}
	}
	defer tx.Rollback(ctx)

	// Порядок параметров должен соответствовать порядку в запросе
	args := []interface{}{
		status,
		ID,
	}

	if _, err = tx.Exec(ctx /*preparedStmt.Name*/, query, args...); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgerrcode.IsCaseNotFound(pgErr.Code) {
			return &storage_errors.NotFoundError{Err: dto.ErrSubjectNotFound}
		} else {
			return &storage_errors.ExecutionPSQLError{Err: err}
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return &storage_errors.ExecutionPSQLError{Err: err}
	}

	return nil
}

func (q QuestionsStorage) DeleteQuestion(ctx context.Context, ID int64) error {
	query := `		
		DELETE FROM 
		    question
		WHERE
		    id = $1
		`

	//preparedStmt, err := q.conn.Prepare(ctx, "DeleteQuestion", query)
	//if err != nil {
	//	return &storage_errors.StatementPSQLError{Err: err}
	//}

	tx, err := q.conn.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.RepeatableRead,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return &storage_errors.ExecutionPSQLError{Err: fmt.Errorf("BeginTx failed: %v\n", err)}
	}
	defer tx.Rollback(ctx)

	// Порядок параметров должен соответствовать порядку в запросе
	args := []interface{}{
		ID,
	}

	if _, err = tx.Exec(ctx /*preparedStmt.Name*/, query, args...); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgerrcode.IsCaseNotFound(pgErr.Code) {
			return &storage_errors.NotFoundError{Err: dto.ErrSubjectNotFound}
		} else {
			return &storage_errors.ExecutionPSQLError{Err: err}
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return &storage_errors.ExecutionPSQLError{Err: err}
	}

	return nil
}
