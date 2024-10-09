package ginhelper

import (
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
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
		"data": data,
	})
}
