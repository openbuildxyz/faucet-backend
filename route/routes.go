package route

import (
	"faucet/controller"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	r.POST("/faucet", controller.HandleFaucet)
}
