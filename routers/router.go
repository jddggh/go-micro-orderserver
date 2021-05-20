package routers

import "github.com/gin-gonic/gin"

func InitRouters() *gin.Engine {
	ginRouter := gin.Default()
	ginRouter.POST("/orders", func(context *gin.Context) {
		context.JSON(200, gin.H{"orders": "orders info"})
	})
	return ginRouter
}