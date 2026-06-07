package apperrors

const (
	MissingAuthModeHeader NetworkError = iota
	InvalidAuthModeHeader
	MissingUserClaimsInContext
	InvalidRequestBody
)

func (e NetworkError) Error() string {
	switch e {
	case MissingAuthModeHeader:
		return "missing X-Auth-Mode header"
	case InvalidAuthModeHeader:
		return "invalid X-Auth-Mode header"
	case MissingUserClaimsInContext:
		return "missing user claims in context"
	case InvalidRequestBody:
		return "invalid request body"
	default:
		return "unhandled network error"
	}
}

func (e NetworkError) Code() string {
	switch e {
	case MissingAuthModeHeader:
		return "missing_auth_mode_header"
	case InvalidAuthModeHeader:
		return "invalid_auth_mode_header"
	case MissingUserClaimsInContext:
		return "missing_user_claims_in_context"
	case InvalidRequestBody:
		return "invalid_request_body"
	default:
		return "unhandled_network_error"
	}
}
