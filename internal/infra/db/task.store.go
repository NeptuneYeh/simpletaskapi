package db

func (s *MemStore) Create(task Task) (Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks[task.ID] = task

	return s.tasks[task.ID], nil
}

func (s *MemStore) Update(id string, task Task) (Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.tasks[id]; !ok {
		return Task{}, ErrTaskNotFound
	}
	// CreateAt 不應該被動到
	task.CreatedAt = s.tasks[id].CreatedAt
	s.tasks[id] = task
	return s.tasks[id], nil
}

func (s *MemStore) Delete(id string) (Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.tasks[id]; !ok {
		return Task{}, ErrTaskNotFound
	}
	deleteItem := s.tasks[id]
	delete(s.tasks, id)
	return deleteItem, nil
}

func (s *MemStore) Get(id string) (Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	task, ok := s.tasks[id]
	if !ok {
		return Task{}, ErrTaskNotFound
	}
	return task, nil
}

func (s *MemStore) GetAll() []Task {
	s.mu.Lock()
	defer s.mu.Unlock()
	tasks := []Task{}
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}
	return tasks
}
