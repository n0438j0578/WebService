package controller
//aa
import (
	"WebService/model"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"strings"
	"context"
)
var ctx = context.Background()

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
		test:="http://35.220.204.174/WebProject/img/"
		for i:=1;i< (len(input)-2);  i++{
				test+=input[i]+"%20"
				name +=input[i]+" "

		}
		test += input[len(input)-2]+".jpg"
		name +=  input[len(input)-2]
		fmt.Println(test)

		type Tag struct {
			Amount string `json:"amount"`
		}
		//var tag Tag
		//
		//db, err := sql.Open("mysql", "root:P@ssword@tcp(35.220.204.174:3306)/N&N_Cafe?charset=utf8")
		//if err != nil {
		//	panic(err.Error())
		//}
		//defer db.Close()
		////SELECT * FROM Customers
		////WHERE Country='Mexico';
		//
		//insert, err := db.QueryContext(ctx,"SELECT amount FROM menu WHERE name='?'",name)
		//err = insert.Scan(&tag.Amount)
		//rawText := "ตอนนี้เหลือ "+tag.Amount
		//
		//defer insert.Close()

		ans = model.Answer{
			name,
			test,
			"",
			"",
		}

	} else {
		fmt.Println(input[1])
		test := "http://35.220.204.174/WebProject/img/" + input[1] + ".jpg"

		fmt.Println(test)

		type Tag struct {
			Amount string `json:"amount"`
		}
		//var tag Tag
		//
		//db, err := sql.Open("mysql", "root:P@ssword@tcp(35.220.204.174:3306)/N&N_Cafe?charset=utf8")
		//if err != nil {
		//	panic(err.Error())
		//}
		//defer db.Close()
		////SELECT * FROM Customers
		////WHERE Country='Mexico';
		//
		//insert, err := db.QueryContext(ctx,"SELECT amount FROM menu WHERE name='?'",test)
		//err = insert.Scan(&tag.Amount)
		//rawText := "ตอนนี้เหลือ "+tag.Amount
		//
		//defer insert.Close()

		ans = model.Answer{
			input[1],
			test,
			"",
			"",
		}
	}

	//input := strings.Split(request.Text, " ")

	response.Status = "success"
	response.Answer = ans
	context.JSON(http.StatusOK, response)
}
