package repositories

import "errors"

var (
	SqlStatementError = errors.New("failed to build sql statement")
	SqlInsertError    = errors.New("failed to insert entity")
	SqlReadError      = errors.New("failed to read entity")
	SqlUpdateError    = errors.New("failed to update entity")
	SqlDeleteError    = errors.New("failed to delete entity")
	SqlScanError      = errors.New("failed to scan entities")
)
