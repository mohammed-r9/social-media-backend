package apperrors

import "net/http"

const (
	MissingEnvVar EnvError = iota
)

var envErrorTable = [...]errorInfo{
	MissingEnvVar: {
		message: "missing environment variable",
		code:    "missing_env_var",
		status:  http.StatusInternalServerError,
	},
}

func (e EnvError) Error() string {
	return e.info().message
}

func (e EnvError) Code() string {
	return e.info().code
}

func (e EnvError) Status() int {
	return e.info().status
}

func (e EnvError) info() errorInfo {
	if int(e) < 0 || int(e) >= len(envErrorTable) {
		return errorInfo{
			message: "unhandled env error",
			code:    "unhandled_env_error",
			status:  http.StatusInternalServerError,
		}
	}
	return envErrorTable[e]
}
