package requests

type TodoRequest struct {
	Title       string `json:"title" binding:"required"`
	Assignee    string `json:"assignee" binding:"required"`
	Description string `json:"description" binding:"required"`
	Priority    int    `json:"priority" binding:"required,todo_priority"`
}
