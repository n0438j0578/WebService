package controller


import (
	"WebService/data"
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
type product struct {
	Name string
	Des string
}

func AddProduct(context *gin.Context) {
	var request struct {
		product
	}
	var response struct {
		Status        string `json:",omitempty"` //"success | error | inactive"
		StatusMessage string `json:",omitempty"`
		Result []model.ProductRow
	}

	err := context.BindJSON(&request)

	if err != nil {
		response.Status = "error"
		response.StatusMessage = err.Error()
		context.JSON(http.StatusInternalServerError, response)
		return
	}

	result := data.AddProduct(request.product.Name,request.product.Des)

	if(result==1){
		response.Status = "success"
		response.StatusMessage = "Insert product"
		context.JSON(http.StatusOK, response)
	}else{
		response.Status = "failed"
		response.StatusMessage = "Insert product"
		context.JSON(http.StatusOK, response)
	}



}
