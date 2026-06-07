package apperrors

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
