package apperrors

import "net/http"

const (
	SessionExpired sessionError = iota
	SessionRevoked
	SessionAlreadyExists
	SessionNotFound

	// new errors should be added above this line
	sessionErrorCount
)

var sessionErrorTable = [sessionErrorCount]errorInfo{
	SessionExpired: {
		message: "session expired",
		code:    "session_expired",
		status:  http.StatusUnauthorized,
	},
	SessionRevoked: {
		message: "session revoked",
		code:    "session_revoked",
		status:  http.StatusUnauthorized,
	},
	SessionAlreadyExists: {
		message: "session already exists",
		code:    "session_exists",
		status:  http.StatusConflict,
	},
	SessionNotFound: {
		message: "session not found",
		code:    "session_not_found",
		status:  http.StatusNotFound,
	},
}

func (e sessionError) Error() string {
	return e.info().message
}

func (e sessionError) Code() string {
	return e.info().code
}

func (e sessionError) Status() int {
	return e.info().status
}

func (e sessionError) info() errorInfo {
	if int(e) < 0 || int(e) >= len(sessionErrorTable) {
		return unhandledSessionError
	}
	return sessionErrorTable[e]
}

var unhandledSessionError = errorInfo{
	message: "unhandled session error",
	code:    "unhandled_session_error",
	status:  http.StatusInternalServerError,
}
