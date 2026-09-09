package handler

import (
	"net/http"

	"github.com/amirhossein-karimi/todo/internal/errors"
	"github.com/amirhossein-karimi/todo/internal/requests"
	"github.com/amirhossein-karimi/todo/internal/response"
	"github.com/amirhossein-karimi/todo/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler interface {
	Create(ctx *gin.Context)
}

type handler struct {
	todoSvc service.TodoService
	logger  *zap.Logger
}

func NewHandler(todoSvc service.TodoService, logger *zap.Logger) Handler {
	return &handler{
		todoSvc: todoSvc,
		logger:  logger,
	}
}

// CreateTodo godoc
//
//	@Summary		Create a new todo
//	@Description	Create a new todo item
//	@Tags			Todos
//	@Accept			json
//	@Produce		json
//	@Param			todo	body	requests.TodoRequest	true	"Todo data"
//	@Success		200	{object}	map[string]string
//	@Failure		400	{object}	map[string]string
//	@Router			/api/v1/todo/create [post]
func (h *handler) Create(ctx *gin.Context) {

	var req requests.TodoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind JSON", zap.Error(err))
		errors.HandleError(ctx, err)
		return
	}

	uuid, err := h.todoSvc.Create(ctx, &req)
	if err != nil {
		h.logger.Error("Failed to create todo", zap.Error(err))
		errors.HandleError(ctx, err)
		return
	}

	res := response.CreateResponse(
		http.StatusCreated,
		"Todo created successfully",
		"success",
		map[string]interface{}{
			"uuid": uuid,
		},
		nil,
	)

	ctx.JSON(http.StatusCreated, res)
}
