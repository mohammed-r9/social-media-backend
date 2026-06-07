package apperrors

const (
	MissingAuthModeHeader NetworkError = iota
	InvalidAuthModeHeader
	MissingAccessToken
	InvalidAccessToken
	MissingUserClaimsInContext
	MissingAuthorizationHeader
	InvalidRequestBody
)

func (e NetworkError) Error() string {
	switch e {
	case MissingAuthModeHeader:
		return "missing X-Auth-Mode header"
	case InvalidAuthModeHeader:
		return "invalid X-Auth-Mode header"
	case MissingAccessToken:
		return "missing access token"
	case InvalidAccessToken:
		return "invalid access token"
	case MissingUserClaimsInContext:
		return "missing user claims in context"
	case MissingAuthorizationHeader:
		return "missing authorization header"
	case InvalidRequestBody:
		return "invalid request body"
	default:
		return "unhandled network error"
	}
}
