package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"WebService/model"
	"fmt"
	"time"
)

func GetExample(context *gin.Context) {
	var request struct {
		Id string
		Name string
	}
	var response struct {
		Status        string `json:",omitempty"` //"success | error | inactive"
		StatusMessage string `json:",omitempty"`
		Example       *model.Example
	}
	path:= context.Param("path")
	err := context.BindQuery(&request)
	if err != nil {
		response.Status = "error"
		response.StatusMessage = err.Error()
		context.JSON(http.StatusInternalServerError, response)
		return
	}
	fmt.Println("Param" ,path)
	fmt.Printf("Query %+v\n" ,request)
	example, err := ds.Mongo.FindExample(request.Id)
	if err != nil {
		cookie1 := &http.Cookie{Name: "sample", Value: "sample",Expires:time.Now(), HttpOnly: false}
		http.SetCookie(context.Writer, cookie1)
		response.Status = "error"
		context.JSON(http.StatusBadRequest, response)
		return
	}

	response.Status = "success"
	response.StatusMessage = ""
	response.Example = example
	context.JSON(http.StatusOK, response)
}

func PostExample(context *gin.Context) {
	var request struct {
		*model.Example
	}
	var response struct {
		Status        string `json:",omitempty"` //"success | error | inactive"
		StatusMessage string `json:",omitempty"`
		Example       *model.Example
	}
	err := context.BindJSON(&request)
	if err != nil {
		response.Status = "error"
		response.StatusMessage = err.Error()
		context.JSON(http.StatusInternalServerError, response)
		return
	}
	example, err := ds.Mongo.InsertExample(request.Example)
	if err != nil {
		response.Status = "error"
		response.StatusMessage = err.Error()
		context.JSON(http.StatusInternalServerError, response)
		return
	}
	response.Status = "success"
	response.StatusMessage = "Insert example"
	response.Example = example
	context.JSON(http.StatusOK, response)
}
