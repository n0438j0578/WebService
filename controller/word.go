package controller

//aa
import (
	"WebService/model"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"strings"
	"context"
	"fmt"
	"database/sql"
	"strconv"
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

	ans := model.Answer{}

	input := strings.Split(request.Text, " ")

	if (strings.Compare(input[0], "มี") == 0 || strings.Compare(input[len(input)-1], "ไหม") == 0) {
		if (len(input) > 3) {
			name := ""
			test := "http://35.220.204.174/WebProject/img/"
			for i := 1; i < (len(input) - 2); i++ {
				test += input[i] + "%20"
				name += input[i] + " "

			}
			test += input[len(input)-2] + ".jpg"
			name += input[len(input)-2]
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
				"มีของไหม",
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
				"มีของไหม",
				input[1],
				test,
				"",
				"",
			}
		}

		//input := strings.Split(request.Text, " ")
	} else if (strings.Compare(input[len(input)-1], "เหลือไหม") == 0) {

		db, err := sql.Open("mysql", "root:P@ssword@tcp(35.220.204.174:3306)/N&N_Cafe?charset=utf8")
		if err != nil {
			panic(err.Error())
		}
		defer db.Close()

		//SELECT Address,ContactName FROM Customers WHERE City='Berlin' ;

		name := ""
		for i := 0; i < len(input); i++ {
			if (i != len(input)-1) {
				if (i != len(input)-2) {
					name += input[i] + " "
				} else {
					name += input[i]
				}

			}
		}

		fmt.Println(name)


		type Tag struct {
			Amount int `json:"amount"`
			ID   int    `json:"id"`
		}

		var tag Tag

		err = db.QueryRow("SELECT amount,id FROM menu WHERE name=?", name).Scan(&tag.Amount,&tag.ID)

		amount := tag.Amount
		//fmt.Println(tag.Amount)
		//fmt.Println(tag.ID)

		//for insert.Next() {
		//	var tag Tag
		//	err = insert.Scan(&tag.Amount)
		//	if err != nil {
		//		panic(err.Error())
		//	}
		//	//fmt.Println(tag.Des)
		//	amount =tag.Amount
		//}
		text := "ของเหลืออยู่ "+ strconv.Itoa(amount)
		ans = model.Answer{
			"เหลือไหม",
			name,
			"",
			"",
			text,
		}

	}else{
		ans = model.Answer{
			"ไม่เข้าพวก",
			"ไม่ทราบจ้า",
			"",
			"",
			"",
		}

	}

	response.Status = "success"
	response.Answer = ans
	context.JSON(http.StatusOK, response)
}
