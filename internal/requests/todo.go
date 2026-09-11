package requests

type TodoRequest struct {
	Title       string `json:"title" binding:"required"`
	Assignee    string `json:"assignee" binding:"required"`
	Description string `json:"description" binding:"required"`
	Priority    int    `json:"priority" binding:"required,todo_priority"`
}
type TodoUpdateRequest struct {
	Title       *string `json:"title" binding:"omitempty"`
	Assignee    *string `json:"assignee" binding:"omitempty"`
	Description *string `json:"description" binding:"omitempty"`
	Priority    *int    `json:"priority" binding:"omitempty,todo_priority"`
	Status      *int    `json:"status" binding:"omitempty"`
}
