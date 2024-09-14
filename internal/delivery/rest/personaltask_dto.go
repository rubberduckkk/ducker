package rest

import (
	"time"

	"github.com/rubberduckkk/ducker/internal/domain/task"
)

type GetTasksRequest struct {
	BatchSize int   `form:"batch_size"`
	Cursor    int64 `form:"cursor"`
	Order     Order `form:"order" binding:"query_order"`
}

type GetTasksResponse struct {
	Cursor int64      `json:"cursor"`
	Items  []TaskItem `json:"items"`
}

type AddTaskRequest struct {
	Content string `json:"content" binding:"required"`
}

type TaskItem struct {
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

func BuildTaskItem(t *task.Task) TaskItem {
	return TaskItem{
		Content:   t.Detail.Content,
		CreatedAt: t.CreatedAt.Format(time.RFC3339),
	}
}
