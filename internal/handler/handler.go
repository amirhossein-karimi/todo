package handler

import (
	"net/http"
	"strconv"

	"github.com/amirhossein-karimi/todo/internal/errors"
	"github.com/amirhossein-karimi/todo/internal/requests"
	"github.com/amirhossein-karimi/todo/internal/response"
	"github.com/amirhossein-karimi/todo/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Handler interface {
	Create(ctx *gin.Context)
	List(ctx *gin.Context)
	Delete(ctx *gin.Context)
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

// List godoc
// @Summary      List todos
// @Description  Get a paginated list of todos with optional filters
// @Tags         Todos
// @Produce      json
// @Param        page      query  int     false  "Page number"       default(1) minimum(1)
// @Param        status    query  string  false  "Todo status"       default(0)
// @Param        priority  query  string  false  "Todo priority"     default(1)
// @Param        assignee  query  string  false  "Todo assignee"
// @Success      200       {object} response.response
// @Failure      400       {object} response.response
// @Failure      500       {object} response.response
// @Router       /api/v1/todo/list [get]
func (h *handler) List(ctx *gin.Context) {
	page := ctx.DefaultQuery("page", "1")

	status := ctx.DefaultQuery("status", "0")
	priority := ctx.DefaultQuery("priority", "1")
	assignee := ctx.DefaultQuery("assignee", "")

	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		h.logger.Error("failed to parse page parameter", zap.Error(err))
		errors.HandleError(ctx, err)
		return
	}

	todos, total, err := h.todoSvc.List(ctx, pageInt, status, priority, assignee)
	if err != nil {
		h.logger.Error("Failed to list todos", zap.Error(err))
		errors.HandleError(ctx, err)
		return
	}

	res := response.CreateResponse(
		http.StatusOK,
		"Todos retrieved successfully",
		"success",
		map[string]interface{}{
			"todos": todos,
			"total": total,
		},
		nil,
	)

	ctx.JSON(http.StatusOK, res)
}

// Delete godoc
// @Summary      Delete todo
// @Description  Delete a todo by UUID
// @Tags         Todos
// @Accept       json
// @Produce      json
// @Param        uuid  path      string  true  "Todo UUID"
// @Success      200   {object}  response.response
// @Failure      400   {object}  response.response
// @Failure      404   {object}  response.response
// @Failure      500   {object}  response.response
// @Router       /api/v1/todo/delete/{uuid} [delete]
func (h *handler) Delete(ctx *gin.Context) {
	uuidParam := ctx.Param("uuid")
	if uuidParam == "" {
		h.logger.Error("UUID parameter is missing")
		errors.HandleError(ctx, errors.ErrTodoNotFound)
		return
	}

	uuid, err := uuid.Parse(uuidParam)
	if err != nil {
		h.logger.Error("Invalid UUID format", zap.Error(err))
		errors.HandleError(ctx, err)
		return
	}

	err = h.todoSvc.Delete(ctx, uuid)
	if err != nil {
		h.logger.Error("Failed to delete todo", zap.Error(err))
		errors.HandleError(ctx, err)
		return
	}

	res := response.CreateResponse(
		http.StatusOK,
		"Todo deleted successfully",
		"success",
		nil,
		nil,
	)

	ctx.JSON(http.StatusOK, res)
}
