package error

import "net/http"

var ErrInvalidGeneral = APIError{
	StatusCode: http.StatusBadRequest,
	Code:       "G001",
	Message:    "We cannot proceed your request",
}

var ErrInvalidInput = APIError{
	StatusCode: http.StatusBadRequest,
	Code:       "0001",
	Message:    "Please check the input again",
}

var ErrUnauthorized = APIError{
	StatusCode: http.StatusUnauthorized,
	Code:       "0002",
	Message:    "You are not allowed to access this",
}
