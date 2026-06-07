package apperrors

type (
	databaseError uint8
	userError     uint8
	authError     uint8
	sessionError  uint8
	networkError  uint8
	envError      uint8
)

type AppError interface {
	error
	Code() string
	Status() int
}

type errorInfo struct {
	message string
	code    string
	status  int
}
