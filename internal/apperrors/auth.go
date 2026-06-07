package apperrors

const (
	InvalidCredentials AuthError = iota
	UnverifiedUserEmail
	InvalidOldPassword
	TokenMissmatch
	InvalidToken
)

func (e AuthError) Error() string {
	switch e {
	case InvalidCredentials:
		return "invalid credentials"
	case UnverifiedUserEmail:
		return "unverified user email"
	case InvalidOldPassword:
		return "invalid user old password"
	case TokenMissmatch:
		return "token missmatch"
	case InvalidToken:
		return "invalid token"
	default:
		return "unhandled auth error"
	}
}
