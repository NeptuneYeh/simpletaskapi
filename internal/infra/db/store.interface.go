package db

type Store interface {
	Create(task Task) (Task, error)
	Update(id string, task Task) (Task, error)
	Delete(id string) (Task, error)
	Get(id string) (Task, error)
	GetAll() []Task
}
