package apperrors

import "net/http"

const (
	UserNotFound userError = iota
	InvalidUserID
	InvalidUserEmail
	InvalidUserName
	InvalidUsername
	ProfileNotFound
	EmailAlreadyTaken

	// new errors should be added above this line
	userErrorCount
)

var userErrorTable = [...]errorInfo{
	UserNotFound: {
		message: "user not found",
		code:    "user_not_found",
		status:  http.StatusNotFound,
	},
	InvalidUserID: {
		message: "invalid user id",
		code:    "invalid_user_id",
		status:  http.StatusBadRequest,
	},
	InvalidUserEmail: {
		message: "invalid user email",
		code:    "invalid_user_email",
		status:  http.StatusBadRequest,
	},
	InvalidUserName: {
		message: "invalid user name",
		code:    "invalid_user_name",
		status:  http.StatusBadRequest,
	},
	InvalidUsername: {
		message: "invalid username",
		code:    "invalid_username",
		status:  http.StatusBadRequest,
	},
	ProfileNotFound: {
		message: "profile not found",
		code:    "profile_not_found",
		status:  http.StatusNotFound,
	},
	EmailAlreadyTaken: {
		message: "email already taken",
		code:    "email_taken",
		status:  http.StatusConflict,
	},
}

func (e userError) Error() string {
	return e.info().message
}

func (e userError) Code() string {
	return e.info().code
}

func (e userError) Status() int {
	return e.info().status
}

func (e userError) info() errorInfo {
	if int(e) < 0 || int(e) >= len(userErrorTable) {
		return unhandledUserError
	}
	return userErrorTable[e]
}

var unhandledUserError = errorInfo{
	message: "unhandled user error",
	code:    "unhandled_user_error",
	status:  http.StatusInternalServerError,
}
