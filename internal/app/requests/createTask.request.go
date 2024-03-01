package requests

type CreateTaskRequest struct {
	Name   string `json:"name" binding:"required"`
	Detail string `json:"detail"`
}
