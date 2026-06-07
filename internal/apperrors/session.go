package apperrors

import "net/http"

const (
	SessionExpired SessionError = iota
	SessionRevoked
	SessionAlreadyExists
	SessionNotFound
)

var sessionErrorTable = [...]errorInfo{
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

func (e SessionError) Error() string {
	return e.info().message
}

func (e SessionError) Code() string {
	return e.info().code
}

func (e SessionError) Status() int {
	return e.info().status
}

func (e SessionError) info() errorInfo {
	if int(e) < 0 || int(e) >= len(sessionErrorTable) {
		return errorInfo{
			message: "unhandled session error",
			code:    "unhandled_session_error",
			status:  http.StatusInternalServerError,
		}
	}
	return sessionErrorTable[e]
}
