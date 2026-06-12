package apperrors

import "net/http"

const (
	MissingAuthModeHeader networkError = iota
	InvalidAuthModeHeader
	MissingUserClaimsInContext
	InvalidRequestBody

	// new errors should be added above this line
	networkErrorCount
)

var networkErrorTable = [networkErrorCount]errorInfo{
	MissingAuthModeHeader: {
		message: "missing X-Auth-Mode header",
		code:    "missing_auth_mode_header",
		status:  http.StatusUnauthorized,
	},
	InvalidAuthModeHeader: {
		message: "invalid X-Auth-Mode header",
		code:    "invalid_auth_mode_header",
		status:  http.StatusUnauthorized,
	},
	MissingUserClaimsInContext: {
		message: "missing user claims in context",
		code:    "missing_user_claims_in_context",
		status:  http.StatusUnauthorized,
	},
	InvalidRequestBody: {
		message: "invalid request body",
		code:    "invalid_request_body",
		status:  http.StatusBadRequest,
	},
}

func (e networkError) Error() string {
	return e.info().message
}

func (e networkError) Code() string {
	return e.info().code
}

func (e networkError) Status() int {
	return e.info().status
}

func (e networkError) info() errorInfo {
	if int(e) < 0 || int(e) >= len(networkErrorTable) {
		return unhandledNetworkError
	}
	return networkErrorTable[e]
}

var unhandledNetworkError = errorInfo{
	message: "unhandled network error",
	code:    "unhandled_network_error",
	status:  http.StatusInternalServerError,
}
