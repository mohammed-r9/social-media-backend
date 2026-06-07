package apperrors

import "net/http"

const (
	UniqueViolation DatabaseError = iota
	NoRowsAffected
	CacheMiss
)

var databaseErrorTable = [...]errorInfo{
	UniqueViolation: {
		message: "unique violation",
		code:    "unique_violation",
		status:  http.StatusConflict,
	},
	NoRowsAffected: {
		message: "no rows affected",
		code:    "no_rows_affected",
		status:  http.StatusNotFound,
	},
	CacheMiss: {
		message: "cache miss",
		code:    "cache_miss",
		status:  http.StatusNotFound,
	},
}

func (e DatabaseError) Error() string {
	return e.info().message
}

func (e DatabaseError) Code() string {
	return e.info().code
}

func (e DatabaseError) Status() int {
	return e.info().status
}

func (e DatabaseError) info() errorInfo {
	if int(e) < 0 || int(e) >= len(databaseErrorTable) {
		return errorInfo{
			message: "unhandled database error",
			code:    "unhandled_database_error",
			status:  http.StatusInternalServerError,
		}
	}
	return databaseErrorTable[e]
}
