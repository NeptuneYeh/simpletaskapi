package requests

type UpdateTaskRequest struct {
	ID     string `uri:"id" binding:"required"`
	Name   string `json:"name"`
	Status int    `json:"status"` // 0: incomplete, 1: complete
	Detail string `json:"detail"`
}
