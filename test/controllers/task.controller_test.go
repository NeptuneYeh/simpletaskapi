package controllers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/NeptuneYeh/simpletask/internal/app/controllers"
	"github.com/NeptuneYeh/simpletask/internal/app/requests"
	"github.com/NeptuneYeh/simpletask/internal/infra/db"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestCreateTaskAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	store := db.NewStore()
	taskController := controllers.NewTaskController(store)
	router.POST("/tasks", taskController.CreateTask)
	router.GET("/tasks/:id", taskController.GetTask)
	router.GET("/tasks", taskController.ListTask)
	router.PUT("/tasks/:id", taskController.UpdateTask)
	router.DELETE("/tasks/:id", taskController.DeleteTask)

	t.Run("Success", func(t *testing.T) {
		task := requests.CreateTaskRequest{Name: "New Task", Detail: "New Detail"}
		body, _ := json.Marshal(task)

		req := httptest.NewRequest("POST", "/tasks", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code 200, got %d", w.Code)
		}
	})

	t.Run("Binding Error", func(t *testing.T) {
		body := []byte(`{"name": "New Task", "detail": 123}`)

		req := httptest.NewRequest("POST", "/tasks", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code 400 for binding error, got %d", w.Code)
		}
	})
}

func TestGetTaskAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	store := db.NewStore()
	taskController := controllers.NewTaskController(store)
	router.POST("/tasks", taskController.CreateTask)
	router.GET("/tasks/:id", taskController.GetTask)
	router.GET("/tasks", taskController.ListTask)
	router.PUT("/tasks/:id", taskController.UpdateTask)
	router.DELETE("/tasks/:id", taskController.DeleteTask)

	// create task
	_, err := store.Create(db.Task{
		ID:        "valid-task-id",
		Name:      "Existing Task",
		Detail:    "balabala",
		Status:    0,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return
	}

	t.Run("Success", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/tasks/valid-task-id", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code 200, got %d", w.Code)
		}

		var task db.Task
		if err := json.Unmarshal(w.Body.Bytes(), &task); err != nil {
			t.Fatalf("Failed to unmarshal response body: %v", err)
		}

		if task.ID != "valid-task-id" || task.Name != "Existing Task" {
			t.Errorf("Expected task with ID 'valid-task-id' and Name 'Existing Task', got %+v", task)
		}
	})

	t.Run("Task Not Found", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/tasks/invalid-task-id", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status code 404 for non-existent task, got %d", w.Code)
		}

		var resp map[string]string
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("Failed to unmarshal response body: %v", err)
		}

		if resp["error"] != db.ErrTaskNotFound.Error() {
			t.Errorf("Expected error message '%v', got '%v'", sql.ErrNoRows.Error(), resp["error"])
		}
	})
}

func TestUpdateTaskAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	store := db.NewStore()
	taskController := controllers.NewTaskController(store)
	router.POST("/tasks", taskController.CreateTask)
	router.GET("/tasks/:id", taskController.GetTask)
	router.GET("/tasks", taskController.ListTask)
	router.PUT("/tasks/:id", taskController.UpdateTask)
	router.DELETE("/tasks/:id", taskController.DeleteTask)

	// create task
	_, err := store.Create(db.Task{
		ID:        "valid-task-id",
		Name:      "Existing Task",
		Detail:    "balabala",
		Status:    0,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return
	}

	t.Run("Success", func(t *testing.T) {
		updatedTask := requests.UpdateTaskRequest{
			ID:     "valid-task-id",
			Name:   "Updated Task",
			Detail: "Updated Detail",
			Status: 1,
		}
		body, _ := json.Marshal(updatedTask)

		req := httptest.NewRequest("PUT", "/tasks/valid-task-id", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code 200, got %d", w.Code)
		}
	})

	t.Run("Task Not Found", func(t *testing.T) {
		updatedTask := requests.UpdateTaskRequest{
			ID:     "invalid-task-id",
			Name:   "Updated Task",
			Detail: "Updated Detail",
			Status: 1,
		}
		body, _ := json.Marshal(updatedTask)
		req := httptest.NewRequest("PUT", "/tasks/invalid-task-id", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status code 404 for non-existent task, got %d", w.Code)
		}
	})

	t.Run("Binding Error", func(t *testing.T) {
		updatedTask := requests.UpdateTaskRequest{
			Name:   "Updated Task",
			Detail: "Updated Detail",
			Status: 1,
		}
		body, _ := json.Marshal(updatedTask)
		req := httptest.NewRequest("PUT", "/tasks/invalid-task-id", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code 400 for binding error, got %d", w.Code)
		}
	})
}

func TestDeleteTaskAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	store := db.NewStore()
	taskController := controllers.NewTaskController(store)
	router.POST("/tasks", taskController.CreateTask)
	router.GET("/tasks/:id", taskController.GetTask)
	router.GET("/tasks", taskController.ListTask)
	router.PUT("/tasks/:id", taskController.UpdateTask)
	router.DELETE("/tasks/:id", taskController.DeleteTask)

	// create task
	_, err := store.Create(db.Task{
		ID:        "valid-task-id",
		Name:      "Existing Task",
		Detail:    "balabala",
		Status:    0,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return
	}

	t.Run("Success", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/tasks/valid-task-id", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code 200, got %d", w.Code)
		}
	})

	t.Run("Task Not Found", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/tasks/valid-task-id", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status code 404 for non-existent task, got %d", w.Code)
		}
	})
}

func TestGetTaskALLAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	store := db.NewStore()
	taskController := controllers.NewTaskController(store)
	router.POST("/tasks", taskController.CreateTask)
	router.GET("/tasks/:id", taskController.GetTask)
	router.GET("/tasks", taskController.ListTask)
	router.PUT("/tasks/:id", taskController.UpdateTask)
	router.DELETE("/tasks/:id", taskController.DeleteTask)

	// create task
	_, err := store.Create(db.Task{
		ID:        "valid-task-id",
		Name:      "Existing Task",
		Detail:    "balabala",
		Status:    0,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return
	}

	t.Run("Success", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/tasks", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code 200, got %d", w.Code)
		}
	})
}
