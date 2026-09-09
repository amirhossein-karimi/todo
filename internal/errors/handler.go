package errors

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/amirhossein-karimi/todo/internal/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func HandleError(ctx *gin.Context, err error) {
	var appErr *AppError

	if errors.As(err, &appErr) {
		ctx.JSON(
			appErr.Code,
			response.CreateResponse(
				appErr.Code,
				appErr.Message,
				"error",
				nil,
				nil,
			),
		)
		return
	}
	var validationErrors validator.ValidationErrors

	if errors.As(err, &validationErrors) {

		fieldErrors := make([]map[string]string, 0)

		for _, fieldError := range validationErrors {

			field := strings.ToLower(fieldError.Field())

			fieldErrors = append(fieldErrors, map[string]string{
				field: validationMessage(fieldError),
			})
		}

		ctx.JSON(
			http.StatusBadRequest,
			response.CreateResponse(
				http.StatusBadRequest,
				"Validation failed",
				"error",
				nil,
				fieldErrors,
			),
		)

		return
	}

	ctx.JSON(
		http.StatusInternalServerError,
		response.CreateResponse(
			http.StatusInternalServerError,
			"Internal server error",
			"error",
			nil,
			nil,
		),
	)
}

func validationMessage(err validator.FieldError) string {
	field := strings.ToLower(err.Field())
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "todo_priority":
		return fmt.Sprintf("%s must be between 1 and 3", field)
	default:
		return "Invalid value"
	}
}
