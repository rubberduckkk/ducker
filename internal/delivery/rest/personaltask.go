package rest

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/rubberduckkk/ducker/internal/domain/task/valueobj"
	"github.com/rubberduckkk/ducker/internal/infra/repository/sql/customer"
	"github.com/rubberduckkk/ducker/internal/infra/repository/sql/task"
	"github.com/rubberduckkk/ducker/internal/service/personaltask"
	"github.com/rubberduckkk/ducker/pkg/mysql"
)

type PersonalTaskDelivery struct {
	svc personaltask.Service
}

var (
	p     *PersonalTaskDelivery
	pOnce sync.Once
)

func PersonalTask() *PersonalTaskDelivery {
	pOnce.Do(func() {
		customerRepo := customer.NewRepository(customer.WithDB(mysql.Instance()))
		taskRepo := task.NewRepository(task.WithDB(mysql.Instance()))
		svc := personaltask.New(customerRepo, taskRepo)
		p = &PersonalTaskDelivery{svc: svc}
	})
	return p
}

func (p *PersonalTaskDelivery) GetTasks(c *gin.Context) {
	req := new(GetTasksRequest)
	if err := c.ShouldBind(req); err != nil {
		reError(c, http.StatusBadRequest, 0, err)
		return
	}

	result, cursor, err := p.svc.GetTasks("", req.Cursor, req.BatchSize)
	if err != nil {
		reError(c, http.StatusInternalServerError, 0, err)
		return
	}

	items := make([]TaskItem, 0, len(result))
	for _, item := range result {
		items = append(items, BuildTaskItem(item))
	}

	resp := new(GetTasksResponse)
	resp.Cursor = cursor
	resp.Items = items
	reData(c, resp)
}

func (p *PersonalTaskDelivery) AddTask(c *gin.Context) {
	req := new(AddTaskRequest)
	if err := c.ShouldBind(req); err != nil {
		reError(c, http.StatusBadRequest, 0, err)
		return
	}

	if err := p.svc.AddTask("", valueobj.TaskDetail{Content: req.Content}); err != nil {
		logrus.WithError(err).Errorf("add task failed")
		reError(c, http.StatusInternalServerError, 0, err)
		return
	}

	reData(c, "")
}
