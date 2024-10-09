package rest

import (
	"github.com/gin-gonic/gin"

	"github.com/rubberduckkk/ducker/internal/delivery/rest/aidoc"
	"github.com/rubberduckkk/ducker/internal/delivery/rest/personaltask"
)

func SetupGin(router *gin.Engine) {
	registerValidation()

	api := router.Group("/api")
	v1 := api.Group("/v1")
	{
		const tasks = "tasks"
		v1.GET(tasks, personaltask.Deliver().GetTasks)
		v1.POST(tasks, personaltask.Deliver().AddTask)

		const aidocs = "aidocs"
		v1.POST(aidocs+"/search", aidoc.Deliver().QueryDocument)
		v1.POST(aidocs, aidoc.Deliver().AddDocument)
	}
}
