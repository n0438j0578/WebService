package controller

import (
	"WebService/model"
	"WebService/test"
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime"
	"sync"
)

func Example(context *gin.Context) {
	var request struct {
		*model.Example
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
	//example, err := ds.Mongo.InsertExample(request.Example)
	//if err != nil {
	//	response.Status = "error"
	//	response.StatusMessage = err.Error()
	//	context.JSON(http.StatusInternalServerError, response)
	//	return
	//}



	result := test.Test("greeting")

	result = test.Test("problem")

	result = test.Test("order")

	result = test.Test("search")

	if result ==1{
		response.Status = "success"
		response.StatusMessage = "Insert example"
		response.Result=request.Value
		context.JSON(http.StatusOK, response)
	}else{
		response.Status = "failed"
		response.StatusMessage = "Insert example"
		response.Result=request.Value
		context.JSON(http.StatusOK, response)
	}
}
func ExampleFindOneByOne(context *gin.Context) {
	var mutex = &sync.Mutex{}
	mutex.Lock()
	runtime.GOMAXPROCS(1)
	var request struct {
		*model.Example
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



	result := test.TestAll()

	//result:= test.TestoneByone(5)

	mutex.Unlock()
	if result ==1{
		response.Status = "success"
		response.StatusMessage = "Insert example"
		response.Result=request.Value
		context.JSON(http.StatusOK, response)
	}else{
		response.Status = "failed"
		response.StatusMessage = "Insert example"
		response.Result=request.Value
		context.JSON(http.StatusOK, response)
	}
}
