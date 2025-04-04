package errors

import (
	"fmt"
)

type (
	NotFoundError struct {
		Err error
	}
	AlreadyExistsError struct {
		Err error
	}
	StatementPSQLError struct {
		Err error
	}
	ExecutionPSQLError struct {
		Err error
	}

	ScanPSQLResultsError struct {
		Err error
	}

	UnmarshalPSQLResultsError struct {
		Err error
	}

	ViolationPSQLRequestError struct {
		Err error
	}
)

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s: not found in storage\n", e.Err.Error())
}

func (e *AlreadyExistsError) Error() string {
	return fmt.Sprintf("%s: already exists in storage\n", e.Err.Error())
}

func (e *StatementPSQLError) Error() string {
	return fmt.Sprintf("%s: could not compile statement\n", e.Err.Error())
}

func (e *ExecutionPSQLError) Error() string {
	return fmt.Sprintf("%s: could not query", e.Err.Error())
}

func (e *ScanPSQLResultsError) Error() string {
	return fmt.Sprintf("%s: could not scan results", e.Err.Error())
}

func (e *UnmarshalPSQLResultsError) Error() string {
	return fmt.Sprintf("%s: could not unmarshal results", e.Err.Error())
}

func (e *ViolationPSQLRequestError) Error() string {
	return fmt.Sprintf("%s: more then one record affected, whitch is prohibited", e.Err.Error())
}

func (e *NotFoundError) Unwrap() error {
	return e.Err
}

func (e *AlreadyExistsError) Unwrap() error {
	return e.Err
}

func (e *StatementPSQLError) Unwrap() error {
	return e.Err
}

func (e *ExecutionPSQLError) Unwrap() error {
	return e.Err
}

func (e *ScanPSQLResultsError) Unwrap() error {
	return e.Err
}

func (e *UnmarshalPSQLResultsError) Unwrap() error {
	return e.Err
}

func (e *ViolationPSQLRequestError) Unwrap() error {
	return e.Err
}
