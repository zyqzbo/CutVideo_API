package router

import (
	"CutVido_api/controller"
	"github.com/gin-gonic/gin"
)

func GetRouter(r *gin.Engine) *gin.Engine {

	// 接口通过中间件用于识别用户身份
	//r.GET("/cut-video",middleware.AuthMiddleware(), controller.CutVideoController)
	r.GET("/cut-video", controller.CutVideoController)
	return r
}
