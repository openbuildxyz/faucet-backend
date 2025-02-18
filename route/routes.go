package route

import (
	"faucet/controller"
	"faucet/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	r.Use(middleware.Cors())
	r.POST("/api/faucet", controller.HandleFaucet)
	r.POST("/api/sign", controller.HandleSign)
	r.GET("/api/user", controller.HandleGetUser)
}
