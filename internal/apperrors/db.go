package apperrors

import "net/http"

const (
	UniqueViolation DatabaseError = iota
	NoRowsAffected
	CacheMiss
)

func (e DatabaseError) Error() string {
	switch e {
	case UniqueViolation:
		return "unique violation"
	case NoRowsAffected:
		return "no rows affected"
	case CacheMiss:
		return "cache miss"
	default:
		return "unhandled database error"
	}
}

func (e DatabaseError) Code() string {
	switch e {
	case UniqueViolation:
		return "unique_violation"
	case NoRowsAffected:
		return "no_rows_affected"
	case CacheMiss:
		return "cache_miss"
	default:
		return "unhandled_database_error"
	}
}
func (e DatabaseError) Status() int {
	switch e {
	case UniqueViolation:
		return http.StatusConflict
	case NoRowsAffected:
		return http.StatusNotFound
	case CacheMiss:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
