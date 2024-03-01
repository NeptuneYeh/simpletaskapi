package db

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrTaskNotFound = errors.New("task not found")
)

type Task struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Status    int       `json:"status"` // 0: incomplete, 1: complete
	Detail    string    `json:"detail"`
	CreatedAt time.Time `json:"created_at"`
}

type MemStore struct {
	mu    sync.Mutex
	tasks map[string]Task
}

func NewStore() Store {
	return &MemStore{
		tasks: make(map[string]Task),
	}
}
