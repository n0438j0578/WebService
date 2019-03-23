package main

import (
	"WebService/controller"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)


func main() {
	gin:=gin.Default()
	gin.Use(cors.Default())
	api := gin.Group("/api")
	api.POST("/word", controller.Word)
	api.POST("/wordset", controller.WordSet)
	api.GET("/findfeature", controller.FindFeature)
	api.POST("/example", controller.Example)
	api.POST("/test",controller.ExampleFindOneByOne)
	api.POST("/wordcome",controller.WordCome)
	api.POST("/wordcomecosine",controller.WordComeCosine)
	api.POST("product",controller.FindProduct)
	api.GET("/unrepiled", controller.Unreplied)

	//api.POST("knntest",controller.KnnTest)
	gin.Run(":20000")

}
