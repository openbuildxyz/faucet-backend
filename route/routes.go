package route

import (
	"faucet/controller"
	"faucet/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	r.Use(middleware.Cors())
	r.POST("/faucet", controller.HandleFaucet)
	r.POST("/sign", controller.HandleSign)
	r.GET("/user", controller.HandleGetUser)
}
