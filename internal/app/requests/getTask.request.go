package requests

type GetTaskRequest struct {
	ID string `uri:"id" binding:"required"`
}
