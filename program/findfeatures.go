package main


import (
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/DavidBelicza/TextRank"
	"context"
	"log"
)

type Tag struct {
	Des string `json:"des"`
}

const DATABASE  = "root:n0438@j0578@tcp(35.220.204.174:3306)/N&N_Cafe?charset=utf8"

func main() {
	db, err := sql.Open("mysql", DATABASE)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	insert, err := db.Query("SELECT des FROM menu")
	rawText :=""

	for insert.Next() {
		var tag Tag
		err = insert.Scan(&tag.Des)
		if err != nil {
			panic(err.Error())
		}
		//fmt.Println(tag.Des)
		rawText+=tag.Des
	}

	tr := textrank.NewTextRank()
	rule := textrank.NewDefaultRule()
	language := textrank.NewDefaultLanguage()
	algorithmDef := textrank.NewDefaultAlgorithm()

	tr.Populate(rawText, language, rule)
	tr.Ranking(algorithmDef)

	text :=textrank.FindSingleWords(tr)


	result:=""

	for i:=0;i<10 ;i++  {
		//fmt.Println(text[i].Word)
		result+=text[i].Word+" "
	}
	fmt.Println(result)

	if err != nil {
		panic(err.Error())
	}

	var ctx = context.Background()


	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	_, execErr := tx.ExecContext(ctx, "INSERT INTO wordrank (word) VALUES (?)",result)
	if execErr != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Printf("Could not roll back: %v\n", rollbackErr)
		}
		log.Fatal(execErr)
	}
	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}

	defer insert.Close()



}
