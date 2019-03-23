package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/narongdejsrn/go-thaiwordcut"
	"github.com/wilcosheh/tfidf"
	"github.com/wilcosheh/tfidf/similarity"
)

type Tag struct {
	Des string `json:"des"`
	Id  string
}

const DATABASE = "root:P@ssword@tcp(35.220.204.174:3306)/N&N_Cafe?charset=utf8"

func main() {

	db, err := sql.Open("mysql", DATABASE)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	insert, err := db.Query("SELECT message,id FROM collections")
	rawText := []string{}

	segmenter := gothaiwordcut.Wordcut()
	segmenter.LoadDefaultDict()
	f := tfidf.New()
	index := []Tag{}

	for insert.Next() {
		var tag Tag
		err = insert.Scan(&tag.Des, &tag.Id)
		if err != nil {
			panic(err.Error())
		}
		result := segmenter.Segment(tag.Des)
		text := ""
		for i := 0; i < len(result); i++ {
			text = text + result[i] + " "
		}
		tag.Des = text
		index = append(index, tag)
		f.AddDocs(text)
		rawText = append(rawText, text)
	}

	t1 := "edimax"
	result := segmenter.Segment(t1)
	text := ""
	for i := 0; i < len(result); i++ {
		text = text + result[i] + " "
	}
	w1 := f.Cal(text)
	fmt.Printf("weight of %s is %+v.\n", t1, w1)

	type Answer struct {
		Des    string `json:"des"`
		Id     string
		Cosine float64
	}
	var answer Answer
	answer.Des = ""
	answer.Id = ""
	answer.Cosine = 0.0

	for i := 0; i < len(index); i++ {
		t2 := index[i].Des
		result = segmenter.Segment(t2)
		w2 := f.Cal(index[i].Des)
		sim := similarity.Cosine(w1, w2)
		if (sim > answer.Cosine) {
			answer.Des = t2
			answer.Id = index[i].Id
			answer.Cosine = sim
		}
	}

	fmt.Println(answer)

	//เหมือนเป็นตัวแปรเอาไว้ใช้ในการใส่ฐานข้อมูลโดยสามารถส่งผ่านตัวแปรได้
	var ctx = context.Background()

	//เอาชนิดของข้อความนั้นออกมาดูว่าเป็นเกี่ยวกับ search หรือเปล่าเพราะต้องเอาออกไปเป็นคาลูเซล
	selectMessages, err := db.QueryContext(ctx, "SELECT types,answer FROM collections WHERE id=?", answer.Id)

	var types string
	var answerindata string
	for selectMessages.Next() {
		err = selectMessages.Scan(&types,&answerindata)
		if err != nil {
			panic(err.Error())
		}
		fmt.Println(types,answerindata)

	}


}
