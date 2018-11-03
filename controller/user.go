package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"WebService/model"
)

func Createuser(context *gin.Context) {
	var request struct {
		Username string
		Password string
	}
	var response struct {
		Status        string `json:",omitempty"` //"success | error | inactive"
		StatusMessage string `json:",omitempty"`
		Username string
	}
	err := context.BindJSON(&request)
	if err != nil {
		response.Status = "error"
		response.StatusMessage = err.Error()
		context.JSON(http.StatusInternalServerError, response)
		return
	}
	if request.Password == ""  {
		response.Status = "error"
		response.StatusMessage = "ลิมกรอก"
		context.JSON(http.StatusInternalServerError, response)
		return
	}
	user := new(model.User)
	user.Username = request.Username
	user.Password = request.Password
	user.Role = "user"
	_, err =  ds.Mongo.CreateUser(user)



	if err != nil {
		response.Status = "error"
		response.StatusMessage = err.Error()
		context.JSON(http.StatusInternalServerError, response)
		return
	}
	response.Status = "success"
	response.StatusMessage = "register success"
	response.Username = request.Username
	context.JSON(http.StatusOK, response)
}
