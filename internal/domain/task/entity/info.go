package entity

import (
	"time"

	"github.com/google/uuid"

	"github.com/rubberduckkk/ducker/internal/domain/task/valueobj"
)

type TaskInfo struct {
	ID        string
	Detail    valueobj.TaskDetail
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewTaskInfo(detail valueobj.TaskDetail) *TaskInfo {
	return &TaskInfo{
		ID:        uuid.NewString(),
		Detail:    detail,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
