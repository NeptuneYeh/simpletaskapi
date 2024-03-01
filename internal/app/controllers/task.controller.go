package controllers

import (
	"errors"
	"github.com/NeptuneYeh/simpletask/internal/app/requests"
	"github.com/NeptuneYeh/simpletask/internal/infra/db"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type TaskController struct {
	store db.Store
}

func NewTaskController(store db.Store) *TaskController {
	return &TaskController{
		store: store,
	}
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func (c *TaskController) CreateTask(ctx *gin.Context) {
	var req requests.CreateTaskRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.Task{
		ID:        uuid.New().String(),
		Name:      req.Name,
		Detail:    req.Detail,
		Status:    0,
		CreatedAt: time.Now(),
	}

	task, err := c.store.Create(arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, task)
}

func (c *TaskController) GetTask(ctx *gin.Context) {
	var req requests.GetTaskRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	task, err := c.store.Get(req.ID)
	if err != nil {
		if errors.Is(err, db.ErrTaskNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, task)
}

func (c *TaskController) ListTask(ctx *gin.Context) {

	tasks := c.store.GetAll()

	ctx.JSON(http.StatusOK, tasks)
}

func (c *TaskController) UpdateTask(ctx *gin.Context) {
	var req requests.UpdateTaskRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.Task{
		ID:     req.ID,
		Name:   req.Name,
		Detail: req.Detail,
		Status: req.Status,
	}

	task, err := c.store.Update(req.ID, arg)
	if err != nil {
		if errors.Is(err, db.ErrTaskNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, task)
}

func (c *TaskController) DeleteTask(ctx *gin.Context) {
	var req requests.DeleteTaskRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	task, err := c.store.Delete(req.ID)
	if err != nil {
		if errors.Is(err, db.ErrTaskNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, task)
}
