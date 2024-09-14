package rest

import (
	"github.com/gin-gonic/gin"
)

func SetupGin(router *gin.Engine) {
	registerValidation()

	api := router.Group("/api")
	v1 := api.Group("/v1")
	{
		const tasks = "tasks"
		v1.GET(tasks, PersonalTask().GetTasks)
		v1.POST(tasks, PersonalTask().AddTask)

		const aidocs = "aidocs"
		v1.GET(aidocs, AIDoc().QueryDocument)
		v1.POST(aidocs, AIDoc().AddDocument)
	}
}
