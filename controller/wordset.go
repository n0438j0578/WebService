package controller

import (
	"WebService/data"
	"WebService/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func WordSet(context *gin.Context) {
	var request struct {
		*model.WordSet
	}
	var response struct {
		Status        string `json:",omitempty"` //"success | error | inactive"
		StatusMessage string `json:",omitempty"`
		Result string
	}
	err := context.BindJSON(&request)

	if err != nil {
		response.Status = "error"
		response.StatusMessage = err.Error()
		context.JSON(http.StatusInternalServerError, response)
		return
	}



	result := data.WordSet(request.Text,request.Type,request.Answer)


	if result ==1{
		response.Status = "success"
		response.StatusMessage = "Insert example"
		response.Result=request.Text
		context.JSON(http.StatusOK, response)
	}else{
		response.Status = "failed"
		response.StatusMessage = "Insert example"
		response.Result=request.Text
		context.JSON(http.StatusOK, response)
	}
}