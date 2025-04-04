package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"quiz_backend_core/internal/dto"
	storage_errors "quiz_backend_core/internal/storage/errors"
)

func NewSubjectsStorage(conn *pgxpool.Pool) *SubjectsStorage {
	return &SubjectsStorage{
		conn: conn,
	}
}

type SubjectsStorage struct {
	conn *pgxpool.Pool
}

func (s SubjectsStorage) GetSubjects(ctx context.Context) ([]dto.Subject, error) {
	query := `		
		SELECT json_build_object(
			'id', s.id::TEXT,
			'name', s.name,
			'description', s.description,
			'creator_user_id', s.creator_user_id::TEXT,
			'active', s.active,
			'parent_id', s.parent_id::TEXT,
			'created_at', s.created_at,
			'updated_at', s.updated_at,
			'question_count', count(q.id),
			'approved_question_count', count(q.id) FILTER ( WHERE q.status_id = (SELECT id from question_status WHERE name='Одобрен') )
		)
		FROM subject s
		LEFT JOIN public.question q ON s.id = q.subject_id
		GROUP BY s.id, s.name, s.description, s.creator_user_id, s.active, s.parent_id, s.created_at, s.updated_at
		ORDER BY s.path;
	`

	var subjects []dto.Subject
	rows, errRows := s.conn.Query(ctx, query)
	if errRows != nil {
		switch {
		case errors.Is(errRows, pgx.ErrNoRows): //it is not error
			return subjects, nil
		default:
			return subjects, &storage_errors.ExecutionPSQLError{Err: errRows}
		}
	}
	defer rows.Close()

	for rows.Next() {
		//var res string
		//if err := rows.Scan(&res); err != nil {
		//	return subjects, &storage_errors.ScanPSQLResultsError{Err: err}
		//}

		//var result dto.Subject
		//if err := json.Unmarshal([]byte(res), &result); err != nil {
		//	return subjects, &storage_errors.UnmarshalPSQLResultsError{Err: err}
		//}

		var result dto.Subject
		if err := rows.Scan(&result); err != nil {
			return subjects, &storage_errors.ScanPSQLResultsError{Err: err}
		}

		subjects = append(subjects, result)
	}
	return subjects, nil
}

func (s SubjectsStorage) AddSubject(ctx context.Context, subject dto.Subject) (int64, error) {
	var subjectID int64 = -1
	query := `
		INSERT INTO subject (
							 name,
							 description,
							 creator_user_id,
							 active,
							 parent_id
							 )
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;
		`

	//preparedStmt, err := s.conn.Prepare(ctx, "AddSubject", query)
	//if err != nil {
	//	return subjectID, &storage_errors.StatementPSQLError{Err: err}
	//}

	tx, err := s.conn.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.RepeatableRead,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return subjectID, &storage_errors.ExecutionPSQLError{Err: fmt.Errorf("BeginTx failed: %v\n", err)}
	}
	defer tx.Rollback(ctx)

	var ParentId sql.NullInt32
	ParentId.Int32 = int32(subject.ParentId)
	ParentId.Valid = subject.ParentId > 0

	args := []interface{}{
		subject.Name,
		subject.Description,
		subject.CreatorUserId,
		subject.Active,
		ParentId,
	}

	if err = tx.QueryRow(ctx /*preparedStmt.Name*/, query, args...).Scan(&subjectID); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			return subjectID, &storage_errors.AlreadyExistsError{Err: dto.ErrSubjectAlreadyExists}
		} else {
			return subjectID, &storage_errors.ExecutionPSQLError{Err: err}
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return subjectID, &storage_errors.ExecutionPSQLError{Err: err}
	}

	return subjectID, nil
}

// TODO нет проверки на существование userId
//func (s SubjectsStorage) GetStatistic(ctx context.Context, userId int64) (dto.Statistic, error) {
//	getTestCountQuery := "SELECT COUNT(*) FROM quiz;"
//	getSubjectCountQuery := "SELECT COUNT(*) FROM subject;"
//	getQuestionCountQuery := "SELECT COUNT(*) FROM question;"
//	getQuestionToModerateCountQuery := "SELECT COUNT(*) FROM question WHERE status_id=(SELECT id from question_status WHERE name='Cоздан');" //TODO use constants
//
//	getTestCountCreatedByCurrentUserQuery := "SELECT COUNT(*) FROM quiz WHERE  creator_user_id = $1;"
//	getQuestionCountCreatedByCurrentUserQuery := "SELECT COUNT(*) FROM question WHERE creator_user_id = $1;"
//
//	getQuestionCountsBySubjectQuery := `
//	SELECT
//		s.id AS subject_id,
//		s.name AS subject_name,
//		COUNT(q.id) AS total_questions
//	FROM
//		subject s
//    LEFT JOIN
//		question q ON EXISTS (SELECT 1 FROM subject s2 WHERE s2.id = q.subject_id AND s2.path <@ s.path)
//	WHERE parent_id IS NULL
//	GROUP BY
//		s.id, s.name
//	ORDER BY
//		s.path;`
//
//	testCount := 0
//	subjectCount := 0
//	questionCount := 0
//	questionToModerateCount := 0
//
//	testCountCreatedByCurrentUser := 0
//	questionCountCreatedByCurrentUser := 0
//
//	if err := s.conn.QueryRow(ctx, getTestCountQuery).Scan(&testCount); err != nil {
//		return dto.Statistic{}, fmt.Errorf("GetStatistic QueryRow: %v\n", err) //TODO refactor errors
//	}
//
//	if err := s.conn.QueryRow(ctx, getSubjectCountQuery).Scan(&subjectCount); err != nil {
//		return dto.Statistic{}, fmt.Errorf("GetStatistic QueryRow: %v\n", err)
//	}
//
//	if err := s.conn.QueryRow(ctx, getQuestionCountQuery).Scan(&questionCount); err != nil {
//		return dto.Statistic{}, fmt.Errorf("GetStatistic QueryRow: %v\n", err)
//	}
//
//	if err := s.conn.QueryRow(ctx, getQuestionToModerateCountQuery).Scan(&questionToModerateCount); err != nil {
//		return dto.Statistic{}, fmt.Errorf("GetStatistic QueryRow: %v\n", err)
//	}
//
//	if err := s.conn.QueryRow(ctx, getTestCountCreatedByCurrentUserQuery, userId).Scan(&testCountCreatedByCurrentUser); err != nil {
//		return dto.Statistic{}, fmt.Errorf("GetStatistic QueryRow: %v\n", err)
//	}
//
//	if err := s.conn.QueryRow(ctx, getQuestionCountCreatedByCurrentUserQuery, userId).Scan(&questionCountCreatedByCurrentUser); err != nil {
//		return dto.Statistic{}, fmt.Errorf("GetStatistic QueryRow: %v\n", err)
//	}
//
//	rows, err := s.conn.Query(ctx, getQuestionCountsBySubjectQuery)
//	if err != nil {
//		return dto.Statistic{}, fmt.Errorf("GetStatistic QueryRow: %v\n", err)
//	}
//	defer rows.Close()
//
//	//todo json?
//	var questionCountsBySubject []dto.SubjectStatisticItem
//	for rows.Next() {
//		var res dto.SubjectStatisticItem
//		errScan := rows.Scan(&res.SubjectId, &res.SubjectName, &res.QuestionCount)
//		if errScan != nil {
//			erRet := fmt.Errorf("GetStatistic rows.Scan: %v\n", errScan)
//			return dto.Statistic{}, erRet
//		}
//		questionCountsBySubject = append(questionCountsBySubject, res)
//	}
//	return dto.Statistic{
//		TestsCount:                        testCount,
//		SubjectCount:                      subjectCount,
//		QuestionCount:                     questionCount,
//		QuestionToModerateCount:           questionToModerateCount,
//		TestCountCreatedByCurrentUser:     testCountCreatedByCurrentUser,
//		QuestionCountCreatedByCurrentUser: questionCountCreatedByCurrentUser,
//		QuestionCountsBySubject:           questionCountsBySubject,
//	}, nil
//}

// TODO use status as constants
func (s SubjectsStorage) GetStatistic(ctx context.Context, userId int64) (dto.Statistic, error) {
	query := `
		SELECT json_build_object(
			'tests_count', COUNT(DISTINCT q.id),
			'subject_count', COUNT(DISTINCT s.id),
			'question_count', COUNT(DISTINCT qu.id),
			'question_to_moderate_count', COUNT(DISTINCT CASE WHEN qu.status_id = (SELECT id from question_status WHERE name='Создан') THEN qu.id END),  -- Using constant
			'test_count_created_by_current_user', COUNT(DISTINCT CASE WHEN q.creator_user_id = $1 THEN q.id END),
			'question_count_created_by_current_user', COUNT(DISTINCT CASE WHEN qu.creator_user_id = $1 THEN qu.id END),
			'question_counts_by_subject', COALESCE((SELECT json_agg(row_to_json(t)) FROM (
					SELECT
						s2.id::TEXT AS subject_id,
						s2.name AS subject_name,
						COUNT(q2.id) AS question_count
					FROM
						subject s2
					LEFT JOIN
						question q2 ON EXISTS (SELECT 1 FROM subject s3 WHERE s3.id = q2.subject_id AND s3.path <@ s2.path)
					WHERE s2.parent_id IS NULL
					GROUP BY s2.id, s2.name, s2.path
					ORDER BY s2.path
				) AS t), '[]'::json)  -- Add subject data
		) AS counts
		FROM (SELECT 1 as dummy) AS one_row
		LEFT JOIN quiz q ON TRUE
		LEFT JOIN subject s ON TRUE
		LEFT JOIN question qu ON TRUE
	`

	var statistic dto.Statistic
	err := s.conn.QueryRow(ctx, query, userId).Scan(&statistic)
	if err != nil {
		//TODO Not found??
		return dto.Statistic{}, fmt.Errorf("GetStatistic QueryRow: %v\n", err)
	}

	return statistic, nil
}

func (s SubjectsStorage) UpdateSubject(ctx context.Context, subject dto.Subject) error {
	query := `		
		UPDATE 
		    subject
		SET
		    name=$1,
		    description=$2,
		    creator_user_id=$3,
		    active=$4,
		    parent_id=$5
		WHERE
		    id = $6
		`

	//preparedStmt, err := s.conn.Prepare(ctx, "UpdateSubject", query);
	//if err != nil {
	//	return &storage_errors.StatementPSQLError{Err: err}
	//}

	tx, err := s.conn.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.RepeatableRead,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return &storage_errors.ExecutionPSQLError{Err: fmt.Errorf("BeginTx failed: %v\n", err)}
	}
	defer tx.Rollback(ctx)

	var ParentId sql.NullInt32
	ParentId.Int32 = int32(subject.ParentId)
	ParentId.Valid = subject.ParentId > 0

	if _, err = tx.Exec(ctx /*preparedStmt.Name*/, query, subject.Name, subject.Description, subject.CreatorUserId, subject.Active, ParentId, subject.ID); err != nil {
		//if pgErr, ok := err.(*pgconn.PgError); ok && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) { //TODO IntegrityConstraintViolation
		//	return &storage_errors.AlreadyExistsError{Err: dto.ErrSubjectAlreadyExists}
		//} else {
		return &storage_errors.ExecutionPSQLError{Err: err}
		//}
	}

	if err = tx.Commit(ctx); err != nil {
		return &storage_errors.ExecutionPSQLError{Err: err}
	}

	return nil
}

func (s SubjectsStorage) DeleteSubjectByID(ctx context.Context, subjectID int64) error {
	query := `		
		DELETE FROM 
		    subject
		WHERE
		    id = $1
		`

	//preparedStmt, err := s.conn.Prepare(ctx, "DeleteSubjectByID", query);
	//if err != nil {
	//	return &storage_errors.StatementPSQLError{Err: err}
	//}

	tx, err := s.conn.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.RepeatableRead,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return &storage_errors.ExecutionPSQLError{Err: fmt.Errorf("BeginTx failed: %v\n", err)}
	}
	defer tx.Rollback(ctx)

	if _, err = tx.Exec(ctx /*preparedStmt.Name*/, query, subjectID); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgerrcode.IsCaseNotFound(pgErr.Code) {
			return &storage_errors.AlreadyExistsError{Err: dto.ErrSubjectNotFound}
		} else {
			return &storage_errors.ExecutionPSQLError{Err: err}
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return &storage_errors.ExecutionPSQLError{Err: err}
	}

	return nil
}
