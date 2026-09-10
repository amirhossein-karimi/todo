package errors

type AppError struct {
	Message string
	Code    int
}

func (e *AppError) Error() string {
	return e.Message
}

var (
	ErrTodoWithThisTitleAlreadyExists = &AppError{
		Message: "todo with this title already exists",
		Code:    409,
	}
	ErrTodoNotFound = &AppError{
		Message: "todo not found",
		Code:    404,
	}
)
