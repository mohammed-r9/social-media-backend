package apperrors

const (
	MissingEnvVar EnvError = iota
)

func (e EnvError) Error() string {
	switch e {
	case MissingEnvVar:
		return "missing environment variable"
	default:
		return "unhandled env error"
	}
}

func (e EnvError) Code() string {
	switch e {
	case MissingEnvVar:
		return "missing_env_var"
	default:
		return "unhandled_env_error"
	}
}
