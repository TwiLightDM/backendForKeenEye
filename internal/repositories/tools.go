package repositories

import "database/sql"

func validateString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func validateInt(ni sql.NullInt32) int {
	if ni.Valid {
		return int(ni.Int32)
	}
	return 0
}
