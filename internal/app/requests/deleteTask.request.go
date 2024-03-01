package requests

type DeleteTaskRequest struct {
	ID string `uri:"id" binding:"required"`
}
