package controller

import (
	"WebService/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func Word(context *gin.Context) {
	var request struct {
		*model.Word
	}
	var response struct {
		Status        string `json:",omitempty"` //"success | error | inactive"
		StatusMessage string `json:",omitempty"`
		Answer        model.Answer
	}
	err := context.BindJSON(&request)
	if err != nil {
		response.Status = "error"
		response.StatusMessage = err.Error()
		context.JSON(http.StatusInternalServerError, response)
		return
	}

	ans :=model.Answer{}

	input := strings.Split(request.Text, " ")
	if (len(input) > 3) {
		name :=""
		test:="http://35.240.208.104/WebProject/food/"
		for i:=1;i< (len(input)-2);  i++{
				test+=input[i]+"%20"
				name +=input[i]+" "

		}
		test += input[len(input)-2]+".jpg"
		name +=  input[len(input)-2]
		fmt.Println(test)
		ans = model.Answer{
			name,
			test,
			"",
		}

	} else {
		fmt.Println(input[1])
		test := "http://35.240.208.104/WebProject/food/" + input[1] + ".jpg"

		fmt.Println(test)

		ans = model.Answer{
			input[1],
			test,
			"",
		}
	}

	//input := strings.Split(request.Text, " ")

	response.Status = "success"
	response.Answer = ans
	context.JSON(http.StatusOK, response)
}
