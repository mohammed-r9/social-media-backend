package apperrors

type (
	DatabaseError uint8
	UserError     uint8
	AuthError     uint8
	SessionError  uint8
	NetworkError  uint8
	EnvError      uint8
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
