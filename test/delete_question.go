package main

import (
	"database/sql"
	//"fmt"
	//"strconv"
	"strings"
	_ "github.com/go-sql-driver/mysql"
)

const DATABASE  = "root:P@ssword@tcp(35.220.204.174:3306)/N&N_Cafe?charset=utf8"

type MessagesRow struct {
	Msg string `json:"message"`
	ID int `json:"id"`
	SubFeature string `json:"sub_feature"`
}

func deleteByWord(word string){

	var qsArr []MessagesRow

	db, err := sql.Open("mysql", DATABASE)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	//var ctx = context.Background()
	selectMessages, err := db.Query("SELECT message, id, sub_feature FROM collections WHERE 1")

	for selectMessages.Next() {

		var qsTmp MessagesRow
		err = selectMessages.Scan(&qsTmp.Msg, &qsTmp.ID, &qsTmp.SubFeature)
		if err != nil {
			panic(err.Error())
		}

		qsArr = append(qsArr,qsTmp)
	}

	//count := 1
	for i:=0; i< len(qsArr); i++  {
		if strings.Contains(qsArr[i].SubFeature,word) {
			//fmt.Print(count)
			//fmt.Println(" message: "+qsArr[i].Msg+" ,subFeature: "+qsArr[i].SubFeature+" ,ID: "+strconv.Itoa(qsArr[i].ID))
			//count++
			deleteRow, err := db.Prepare("DELETE FROM collections WHERE id=?")
			if err != nil {
				panic(err.Error())
			}
			deleteRow.Exec(qsArr[i].ID)
		}

	}

}


func main ()  {
	deleteByWord("สวัสดี")
}
