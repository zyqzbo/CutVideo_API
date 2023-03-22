package router

import (
	"CutVido_api/controller"
	"CutVido_api/middleware"
	"github.com/gin-gonic/gin"
)

func GetRouter(r *gin.Engine) *gin.Engine {
	r.Use(middleware.CORSMiddleware(), middleware.RecoveryMiddleware()) // 跨域处理 、 报错捕捉显示
	// 接口通过中间件用于识别用户身份
	//r.POST("/cut-video",middleware.AuthMiddleware(), controller.CutVideoController)
	r.POST("/cutVideo", controller.CutVideoController)
	return r
}
