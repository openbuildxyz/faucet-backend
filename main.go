package main

import (
	_ "faucet/config"
	"faucet/logger"
	"faucet/middleware"
	"faucet/route"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	// 初始化日志
	logFile := viper.GetString("log.file")
	logLevel := viper.GetString("log.level")
	logger.Init(logFile, logLevel)

	// 初始化 Gin
	r := gin.Default()
	r.Use(middleware.LeakyBucketRateLimiter(30))
	// 使用自定义的路由配置
	route.SetupRouter(r)

	// 启动服务
	port := viper.GetString("server.port")
	logger.Log.Infof("Starting server on port %s", port)
	if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
		logger.Log.Fatalf("Failed to start server: %v", err)
	}
}
