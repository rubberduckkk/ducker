package rest

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ReError(c *gin.Context, httpCode, bizCode int, err error) {
	c.AbortWithStatusJSON(httpCode, gin.H{
		"code": bizCode,
		"msg":  err.Error(),
	})
}

func ReData(c *gin.Context, data interface{}) {
	jsonData, _ := json.Marshal(data)
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
		"data": jsonData,
	})
}
