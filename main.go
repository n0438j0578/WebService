package main

import (
	"WebService/controller"
	"github.com/gin-gonic/gin"
)


func main() {
	gin:=gin.Default()
	api := gin.Group("/api")
	api.POST("/word", controller.Word)
	api.GET("/findfeature", controller.FindFeature)
	api.POST("/example", controller.Example)
	gin.Run(":20000")

}
