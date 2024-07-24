package app_error

type AppError struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Errors  []string `json:"errors,omitempty"`
}

func (appErr *AppError) Error() string {
	return appErr.Message
}

func New(code int, message string, errs ...error) *AppError {
	errMsgs := make([]string, 0)
	for _, err := range errs {
		if err != nil {
			errMsgs = append(errMsgs, err.Error())
		}
	}

	return &AppError{
		Code:    code,
		Message: message,
		Errors:  errMsgs,
	}
}

func IsAppError(err error) bool {
	_, ok := err.(*AppError)
	return ok
}
