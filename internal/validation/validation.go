package validation

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func TodoPriority(fl validator.FieldLevel) bool {
	priority := fl.Field().Int()

	return priority >= 1 && priority <= 3
}

func Register() error {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		return v.RegisterValidation("todo_priority", TodoPriority)
	}

	return nil
}
