package controller


import (
	"WebService/model"
	"WebService/test"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type input struct {
	Text string
}


func FindProduct(context *gin.Context) {
	var request struct {
		input
	}
	var response struct {
		Status        string `json:",omitempty"` //"success | error | inactive"
		StatusMessage string `json:",omitempty"`
		Result []model.ProductRow
	}

	err := context.BindJSON(&request)

	fmt.Println(request)

	if err != nil {
		response.Status = "error"
		response.StatusMessage = err.Error()
		context.JSON(http.StatusInternalServerError, response)
		return
	}

	response.Result = test.ProductMatching(request.input.Text)


	response.Status = "success"
	response.StatusMessage = "Insert example"
	context.JSON(http.StatusOK, response)


}
