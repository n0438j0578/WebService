package controller

import (
	"WebService/data"
	"WebService/model"
	"WebService/test"
	"github.com/gin-gonic/gin"
	"net/http"
)

func WordCome(context *gin.Context) {
	var request struct {
		*model.Word
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



	result,answer := data.WordCome(request.Text,request.Idcustomer)


	if result ==1{
		response.Status = "success"
		response.StatusMessage = "เจอข้อความ"
		response.Result=answer
		context.JSON(http.StatusOK, response)
	}else if(result==2){
		response.Status = "success"
		response.StatusMessage = "ไม่เจอข้อความแต่ตอบได้"
		response.Result=answer
		context.JSON(http.StatusOK, response)
		test.TestAll()
	}else{
		response.Status = "failed"
		response.StatusMessage = "ไม่เจอข้อความ"
		response.Result=""
		context.JSON(http.StatusOK, response)
	}
}