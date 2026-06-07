package apperrors

import "net/http"

const (
	UserNotFound UserError = iota
	InvalidUserID
	InvalidUserEmail
	InvalidUserName
	InvalidUsername
	ProfileNotFound
	EmailAlreadyTaken
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

func (e UserError) Error() string {
	return e.info().message
}

func (e UserError) Code() string {
	return e.info().code
}

func (e UserError) Status() int {
	return e.info().status
}

func (e UserError) info() errorInfo {
	if int(e) < 0 || int(e) >= len(userErrorTable) {
		return errorInfo{
			message: "unhandled user error",
			code:    "unhandled_user_error",
			status:  http.StatusInternalServerError,
		}
	}
	return userErrorTable[e]
}
