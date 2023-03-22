package main

import (
	router "CutVido_api/Router"
	"CutVido_api/common"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	db := common.InitDB()
	fmt.Println("db", db)
	r := gin.Default()
	r = router.GetRouter(r)
	panic(r.Run())

	//result := test.CutVideo("/Users/zyq/Desktop/test.mp4", "/Users/zyq/Desktop", "0:01", 2)
	//fmt.Println(result)

	//router := gin.Default()
	//router.GET("/cut-video", test.HandleCutVideo)
	//router.Run(":8080")
}
