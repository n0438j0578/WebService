package data

import (
	"WebService/test"
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/narongdejsrn/go-thaiwordcut"
)

type Id struct {
	id int
}

type Tag struct {
	Feature string `json:"des"`
	Count int
}

const DATABASE = "root:P@ssword@tcp(35.220.204.174:3306)/N&N_Cafe?charset=utf8"

func WordSet(text string, types string, ans string) int {
	db, err := sql.Open("mysql", DATABASE)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	segmenter := gothaiwordcut.Wordcut()
	segmenter.LoadDefaultDict()
	res := segmenter.Segment(text)

	fmt.Print("Cutdata : ")
	fmt.Println(res)

	result := ""

	for i := 0; i < len(res); i++ {
		result += res[i] + " "
	}

	insForm, err := db.Prepare("INSERT INTO collections(message,types,answer,sub_feature,count) VALUES (?,?,?,?,?)")
	if err != nil {
		fmt.Print("Cutdata1 : ")
		fmt.Println(res)
		panic(err.Error())
		return 0
	}
	_, err = insForm.Exec(text, types, ans, result, 0)
	if err != nil {
		return 0
	}

	num := test.Test(types)

	if num == 0 {
		return 0
	}
	num = test.TestAll()
	if num == 0 {
		return 0
	}
	return 1

}

func WordCome(text string, Idcustomer string) (int, string) {
	db, err := sql.Open("mysql", DATABASE)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	var ctx = context.Background()
	fmt.Println(text)
	selectMessages, err := db.QueryContext(ctx, "SELECT answer,count FROM collections WHERE message=?", text)
	fmt.Println(selectMessages)
	rawText := ""
	count :=1

	for selectMessages.Next() {
		var tag Tag
		err = selectMessages.Scan(&tag.Feature,&tag.Count)
		if err != nil {
			panic(err.Error())
		}
		rawText += tag.Feature
		count = count +tag.Count
	}

	if (strings.Compare(rawText, "") == 0) {
		SaveWord(text,Idcustomer)
		return 0, ""
	} else {
		insForm, _ := db.Prepare("UPDATE collections SET count=? WHERE message=? ")
		insForm.Exec(count, text)
		SaveWord(text,Idcustomer)
		return 1, rawText
	}

}

func SaveWord(text string, Idcustomer string)  {
	db, err := sql.Open("mysql", DATABASE)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	var ctx = context.Background()
	fmt.Println(text)
	selectMessages, err := db.QueryContext(ctx, "SELECT message FROM oldmsg WHERE id=?", Idcustomer)
	rawText := ""

	for selectMessages.Next() {
		var tag Tag
		err = selectMessages.Scan(&tag.Feature)
		if err != nil {
			panic(err.Error())
		}
		rawText += tag.Feature
	}

	if (strings.Compare(rawText, "") == 0) {
		insForm, _ := db.Prepare("INSERT INTO oldmsg(id,message) VALUES (?,?)")
		_, err = insForm.Exec(Idcustomer, text)
	} else {
		insForm, _ := db.Prepare("UPDATE oldmsg SET message=? WHERE id=? ")
		insForm.Exec(text, Idcustomer)
	}

}
