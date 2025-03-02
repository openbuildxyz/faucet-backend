package route

import (
	"faucet/controller"
	"faucet/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	r.Use(middleware.Cors())
	r.POST("/api/sign", controller.HandleSign)
	r.GET("/api/user", middleware.JWTAuthMiddleware(), controller.HandleGetUser)
	r.POST("/api/faucet", middleware.JWTAuthMiddleware(), controller.HandleFaucet)
}
