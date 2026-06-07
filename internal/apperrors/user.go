package apperrors

const (
	UserNotFound UserError = iota
	InvalidUserID
	InvalidUserEmail
	InvalidUserName
	InvalidUsername
	ProfileNotFound
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
	default:
		return "unhandled user error"
	}
}
