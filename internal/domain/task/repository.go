package task

import (
	"github.com/rubberduckkk/ducker/internal/domain/task/valueobj"
)

type Repository interface {
	Create(task *Task) error
	Update(task *Task) error
	Get(id string) (*Task, error)
	Remove(id string) error
	List(customerID string, cursor int64, batchSize int, order valueobj.QueryOrder) (res []*Task, nextCursor int64, err error)
}
