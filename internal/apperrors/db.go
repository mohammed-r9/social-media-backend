package apperrors

import "net/http"

const (
	UniqueViolation databaseError = iota
	NoRowsAffected
	CacheMiss

	// new errors should be added above this line
	dbErrorCount
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

func (e databaseError) Error() string {
	return e.info().message
}

func (e databaseError) Code() string {
	return e.info().code
}

func (e databaseError) Status() int {
	return e.info().status
}

func (e databaseError) info() errorInfo {
	if int(e) < 0 || int(e) >= len(databaseErrorTable) {
		return unhandledDatabseError
	}

	return databaseErrorTable[e]
}

var unhandledDatabseError = errorInfo{
	message: "unhandled database error",
	code:    "unhandled_database_error",
	status:  http.StatusInternalServerError,
}
