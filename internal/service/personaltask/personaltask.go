package personaltask

import (
	"github.com/rubberduckkk/ducker/internal/domain/customer"
	"github.com/rubberduckkk/ducker/internal/domain/task"
	"github.com/rubberduckkk/ducker/internal/domain/task/valueobj"
)

type Service interface {
	AddTask(customerID string, detail valueobj.TaskDetail) error
	GetTasks(customerID string, cursor int64, batchSize int) (tasks []*task.Task, nextCursor int64, err error)
}

type svcImpl struct {
	customers customer.Repository
	tasks     task.Repository
}

func New(customers customer.Repository, tasks task.Repository) Service {
	return &svcImpl{
		customers: customers,
		tasks:     tasks,
	}
}

func (s *svcImpl) AddTask(cid string, detail valueobj.TaskDetail) error {
	c, err := s.customers.Get(cid)
	if err != nil {
		return err
	}
	t := task.NewTask(c.ID, detail)
	return s.tasks.Create(t)
}

func (s *svcImpl) GetTasks(customerID string, cursor int64, batchSize int) (tasks []*task.Task, nextCursor int64, err error) {
	c, err := s.customers.Get(customerID)
	if err != nil {
		return nil, 0, err
	}
	return s.tasks.List(c.ID, cursor, batchSize, valueobj.OrderDESC)
}
