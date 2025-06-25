package repositories

import "errors"

var (
	SqlStatementError = errors.New("failed to build sql statement")
)
