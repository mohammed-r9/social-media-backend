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

func (e UserError) Error() string {
	switch e {
	case UserNotFound:
		return "user not found"
	case InvalidUserID:
		return "invalid user id"
	case InvalidUserEmail:
		return "invalid user email"
	case InvalidUserName:
		return "invalid user name"
	case InvalidUsername:
		return "invalid username"
	case ProfileNotFound:
		return "profile not found"
	case EmailAlreadyTaken:
		return "email already taken"
	default:
		return "unhandled user error"
	}
}

func (e UserError) Code() string {
	switch e {
	case UserNotFound:
		return "user_not_found"
	case InvalidUserID:
		return "invalid_user_id"
	case InvalidUserEmail:
		return "invalid_user_email"
	case InvalidUserName:
		return "invalid_user_name"
	case InvalidUsername:
		return "invalid_username"
	case ProfileNotFound:
		return "profile_not_found"
	case EmailAlreadyTaken:
		return "email_taken"
	default:
		return "unhandled_user_error"
	}
}
func (e UserError) Status() int {
	switch e {
	case UserNotFound:
		return http.StatusNotFound
	case InvalidUserID:
		return http.StatusBadRequest
	case InvalidUserEmail:
		return http.StatusBadRequest
	case InvalidUserName:
		return http.StatusBadRequest
	case InvalidUsername:
		return http.StatusBadRequest
	case ProfileNotFound:
		return http.StatusNotFound
	case EmailAlreadyTaken:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
