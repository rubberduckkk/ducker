package task

import (
	"time"

	"github.com/rubberduckkk/ducker/internal/domain/task"
	"github.com/rubberduckkk/ducker/internal/domain/task/entity"
	"github.com/rubberduckkk/ducker/internal/domain/task/valueobj"
)

type Task struct {
	ID         string `json:"id" gorm:"primaryKey;not null;type:varchar(32)"`
	CustomerID string `json:"customer_id" gorm:"column:customer_id;type:varchar(32);not null"`
	Detail     string `json:"detail" gorm:"column:detail;type:text"`
	CreatedAt  int64  `json:"created_at" gorm:"column:created_at;type:bigint(20);not null"`
	UpdatedAt  int64  `json:"updated_at" gorm:"column:updated_at;type:bigint(20);not null"`
}

func (Task) TableName() string {
	return "tasks"
}

func FromEntity(task *task.Task) *Task {
	m := new(Task)
	m.ID = task.ID
	m.CustomerID = task.CustomerID
	m.Detail = task.Detail.Marshal()
	m.CreatedAt = task.CreatedAt.Unix()
	m.UpdatedAt = task.UpdatedAt.Unix()
	return m
}

func (t Task) ToEntity() *task.Task {
	var detail valueobj.TaskDetail
	_ = detail.Unmarshal(t.Detail)
	return &task.Task{
		TaskInfo: &entity.TaskInfo{
			ID:        t.ID,
			Detail:    detail,
			CreatedAt: time.Unix(t.CreatedAt, 0),
			UpdatedAt: time.Unix(t.UpdatedAt, 0),
		},
		CustomerID: t.CustomerID,
	}
}
