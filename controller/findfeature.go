package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/DavidBelicza/TextRank"
	"fmt"
	"database/sql"
	"log"
	"context"
)
type Tag struct {
	Des string `json:"des"`
}


func FindFeature(con *gin.Context) {

	var response struct {
		Status        string `json:",omitempty"` //"success | error | inactive"
		StatusMessage string `json:",omitempty"`
	}

	db, err := sql.Open("mysql", "root:n0438@j0578@tcp(35.220.204.174:3306)/N&N_Cafe?charset=utf8")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	insert, err := db.Query("SELECT des FROM product")
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
		result+=text[i].Word+" "
	}
	fmt.Println(result)

	if err != nil {
		panic(err.Error())
	}

	var ctx = context.Background()

	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	_, execErr := tx.Exec( "DELETE FROM wordrank")
	if execErr != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Printf("Could not roll back: %v\n", rollbackErr)
		}
		log.Fatal(execErr)
	}
	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}



	//insert
	tx, err = db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	_, execErr = tx.ExecContext(ctx, "INSERT INTO wordrank (word) VALUES (?)",result)
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

	response.Status = "ระบบได้ค้นหาคำแล้ว"
	response.StatusMessage = ""
	con.JSON(http.StatusOK, response)
}