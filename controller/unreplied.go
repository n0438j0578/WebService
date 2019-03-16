package controller

import (
	"WebService/model"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)


func Unreplied (con *gin.Context) {

	var response struct {
		Status        string `json:",omitempty"` //"success | error | inactive"
		StatusMessage string
		Result	      []model.Unreplied
	}

	db, err := sql.Open("mysql", "root:P@ssword@tcp(35.220.204.174:3306)/N&N_Cafe?charset=utf8")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	selectMessages, err := db.Query("SELECT id, message FROM oldmsg WHERE status=0")
	var arr model.Unreplied

	for selectMessages.Next() {

		err = selectMessages.Scan(&arr.ID, &arr.Msg)
		if err != nil {
			panic(err.Error())
		}
		//fmt.Println(tag.Des)
		response.Result = append(response.Result, arr)
	}



	response.Status = "Success!"
	con.JSON(http.StatusOK, response)


}
