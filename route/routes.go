package route

import (
	"faucet/controller"
	"faucet/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	r.Use(middleware.Cors())
	r.POST("/sign", controller.HandleSign)
	r.GET("/user", middleware.JWTAuthMiddleware(), controller.HandleGetUser)
	r.POST("/faucet", middleware.JWTAuthMiddleware(), controller.HandleFaucet)
}
