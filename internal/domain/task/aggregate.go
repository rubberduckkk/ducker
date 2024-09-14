package task

import (
	"github.com/rubberduckkk/ducker/internal/domain/task/entity"
	"github.com/rubberduckkk/ducker/internal/domain/task/valueobj"
)

type Task struct {
	*entity.TaskInfo
	CustomerID string
}

func NewTask(customerID string, detail valueobj.TaskDetail) *Task {
	return &Task{
		CustomerID: customerID,
		TaskInfo:   entity.NewTaskInfo(detail),
	}
}
