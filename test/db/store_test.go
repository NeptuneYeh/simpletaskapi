package db

import (
	"github.com/NeptuneYeh/simpletask/internal/infra/db"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestMemStore_CreateAndGet(t *testing.T) {
	store := db.NewStore()
	task := db.Task{
		ID:        "task1",
		Name:      "Test Task",
		Status:    0,
		Detail:    "Test detail",
		CreatedAt: time.Now(),
	}

	createdTask, err := store.Create(task)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	if !reflect.DeepEqual(createdTask, task) {
		t.Errorf("Expected created task to be %+v, got %+v", task, createdTask)
	}

	retrievedTask, err := store.Get(task.ID)
	if err != nil {
		t.Fatalf("Failed to retrieve task: %v", err)
	}

	if !reflect.DeepEqual(retrievedTask, task) {
		t.Errorf("Expected retrieved task to be %+v, got %+v", task, retrievedTask)
	}
}

func TestMemStore_Update(t *testing.T) {
	store := db.NewStore()
	originalTask := db.Task{
		ID:        "task1",
		Name:      "Original Task",
		Status:    0,
		Detail:    "Original detail",
		CreatedAt: time.Now(),
	}

	_, err := store.Create(originalTask)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	updatedTask := db.Task{
		ID:     "task1",
		Name:   "Updated Task",
		Status: 1,
		Detail: "Updated detail",
	}

	_, err = store.Update(updatedTask.ID, updatedTask)
	if err != nil {
		t.Fatalf("Failed to update task: %v", err)
	}

	taskAfterUpdate, err := store.Get(updatedTask.ID)
	if err != nil {
		t.Fatalf("Failed to get task after update: %v", err)
	}

	// 我不想比較 CreatedAt 的值 因為設計的不會動到
	taskAfterUpdate.CreatedAt = updatedTask.CreatedAt
	if !reflect.DeepEqual(taskAfterUpdate, updatedTask) {
		t.Errorf("Expected task after update to be %+v, got %+v", updatedTask, taskAfterUpdate)
	}
}

func TestMemStore_Delete(t *testing.T) {
	store := db.NewStore()
	task := db.Task{
		ID:        "task1",
		Name:      "Test Task",
		Status:    0,
		Detail:    "Test detail",
		CreatedAt: time.Now(),
	}

	_, err := store.Create(task)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	deletedTask, err := store.Delete(task.ID)
	if err != nil {
		t.Fatalf("Failed to delete task: %v", err)
	}

	if !reflect.DeepEqual(deletedTask, task) {
		t.Errorf("Expected deleted task to be %+v, got %+v", task, deletedTask)
	}

	_, err = store.Get(task.ID)
	if err != db.ErrTaskNotFound {
		t.Fatalf("Expected ErrTaskNotFound after task deletion, got %v", err)
	}
}

func TestMemStore_GetAll(t *testing.T) {
	store := db.NewStore()
	task1 := db.Task{
		ID:        "task1",
		Name:      "Test Task 1",
		Status:    0,
		Detail:    "Test detail 1",
		CreatedAt: time.Now(),
	}
	task2 := db.Task{
		ID:        "task2",
		Name:      "Test Task 2",
		Status:    1,
		Detail:    "Test detail 2",
		CreatedAt: time.Now(),
	}

	_, err := store.Create(task1)
	if err != nil {
		t.Fatalf("Failed to create task1: %v", err)
	}
	_, err = store.Create(task2)
	if err != nil {
		t.Fatalf("Failed to create task2: %v", err)
	}

	tasks := store.GetAll()
	if len(tasks) != 2 {
		t.Fatalf("Expected 2 tasks, got %d", len(tasks))
	}

	expectedTasks := []db.Task{task1, task2}
	if !reflect.DeepEqual(tasks, expectedTasks) {
		t.Errorf("Expected tasks to be %+v, got %+v", expectedTasks, tasks)
	}
}

func TestMemStore_UpdateNonExistent(t *testing.T) {
	store := db.NewStore()
	task := db.Task{
		ID:     "nonexistent",
		Name:   "Nonexistent Task",
		Status: 0,
		Detail: "This task does not exist",
	}

	_, err := store.Update(task.ID, task)
	if err != db.ErrTaskNotFound {
		t.Fatalf("Expected ErrTaskNotFound for non-existent task update, got %v", err)
	}
}

func TestMemStore_DeleteNonExistent(t *testing.T) {
	store := db.NewStore()

	_, err := store.Delete("nonexistent")
	if err != db.ErrTaskNotFound {
		t.Fatalf("Expected ErrTaskNotFound for non-existent task deletion, got %v", err)
	}
}

func TestMemStore_ConcurrentUpdates(t *testing.T) {
	store := db.NewStore()
	task := db.Task{
		ID:        "concurrent",
		Name:      "Concurrent Task",
		Status:    0,
		Detail:    "Test detail",
		CreatedAt: time.Now(),
	}

	_, err := store.Create(task)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			update := db.Task{
				ID:     "concurrent",
				Name:   "Updated Concurrent Task",
				Status: 1,
				Detail: "Updated detail",
			}
			_, err := store.Update(update.ID, update)
			if err != nil {
				t.Errorf("Failed to update task: %v", err)
			}
		}(i)
	}
	wg.Wait()

	updatedTask, err := store.Get("concurrent")
	if err != nil {
		t.Fatalf("Failed to get task after concurrent updates: %v", err)
	}

	if updatedTask.Status != 1 {
		t.Errorf("Expected task status to be 1 after concurrent updates, got %d", updatedTask.Status)
	}
}
