package controller

import (
	"WebService/data"
	"WebService/model"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func WordCome(context *gin.Context) {
	//รับข้อความเข้ามาแล้วทำการส่งต่อ
	var request struct {
		*model.Word
	}
	err := context.BindJSON(&request)

	if err != nil {
		context.JSON(http.StatusInternalServerError, "")
		return
	}

	var product []model.ProductRow

	result, answer, product := data.WordCome(request.Text, request.Idcustomer)

	fmt.Println(answer)

	if result == 1 {
		var response struct {
			Status        string `json:",omitempty"` //"success | error | inactive"
			StatusMessage string `json:",omitempty"`
			Result        string
		}
		response.Status = "success"
		response.StatusMessage = "เจอข้อความ"
		response.Result = answer
		context.JSON(http.StatusOK, response)
	} else if (result == 2) {
		var response struct {
			Status        string `json:",omitempty"` //"success | error | inactive"
			StatusMessage string `json:",omitempty"`
			Result        string
		}
		response.Status = "success"
		response.StatusMessage = "ไม่เจอข้อความแต่ตอบได้"
		response.Result = answer
		context.JSON(http.StatusOK, response)
	} else if (result == 3) {
		var response struct {
			Status        string `json:",omitempty"` //"success | error | inactive"
			StatusMessage string `json:",omitempty"`
			Product       []model.ProductRow
		}
		response.Status = "success search"
		response.StatusMessage = "search"
		response.Product = product
		context.JSON(http.StatusOK, response)
	} else {
		var response struct {
			Status        string `json:",omitempty"` //"success | error | inactive"
			StatusMessage string `json:",omitempty"`
			Result        string
		}
		db, err := sql.Open("mysql", "root:n0438@j0578@tcp(35.220.204.174:3306)/N&N_Cafe?charset=utf8")
		if err != nil {
			panic(err.Error())
		}
		defer db.Close()
		insForm, _ := db.Prepare("UPDATE oldmsg SET status=? WHERE id=? ")
		insForm.Exec(0, request.Idcustomer	)

		response.Status = "failed"
		response.StatusMessage = "ไม่เจอข้อความ"
		response.Result = ""
		context.JSON(http.StatusOK, response)
	}
}
