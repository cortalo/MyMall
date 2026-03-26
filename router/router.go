package router

import (
	"MyMall/pkg/logger"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLogger())
	r.Use(logger.GinRecovery(true))

	v1 := r.Group("/api/v1")
	{
		v1.GET("/hello", func(c *gin.Context) {
			c.JSON(200, gin.H{"hello": "world"})
		})
		v1.GET("/hellocd", func(c *gin.Context) {
			c.JSON(200, gin.H{"hello": "cd"})
		})
	}

	return r
}
