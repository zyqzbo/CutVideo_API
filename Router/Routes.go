package router

import (
	"CutVido_api/controller"
	"github.com/gin-gonic/gin"
)

func GetRouter(r *gin.Engine) *gin.Engine {

	r.GET("/cut-video", controller.CutVideoController)
	return r
}
