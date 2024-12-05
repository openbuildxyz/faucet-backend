package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化Gin引擎
	r := gin.Default()

	// 创建一个简单的路由
	r.GET("/faucet", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to the faucet backend!",
		})
	})

	// 启动后端服务，监听在 8080 端口
	r.Run(":8080")
}
