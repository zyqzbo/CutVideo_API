package middleware

// 为了拦截重复创建分类是出现的报错信息显示出来 for CategoryController Create
//定义的拦截方法  for CategoryController Create 重复创建出现报错的时候配合panic使用
//gin框架也有在main.go 的 r := gin.Default() 的 Default 点进去的 recover()拦截显示方法

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RecoveryMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				Fail(context, nil, fmt.Sprint(err))
			}
		}()

		context.Next()
	}
}

func Response(ctx *gin.Context, httpStatus int, code int, data gin.H, msg string) { // 封装请求格式
	ctx.JSON(httpStatus, gin.H{
		"code": code,
		"data": data,
		"msg":  msg,
	})
}

func Fail(ctx *gin.Context, data gin.H, msg string) { // 根据上面封装方法定义请求失败时的返回字段格式
	Response(ctx, http.StatusOK, 400, data, msg)
}
