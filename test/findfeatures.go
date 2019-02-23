package test

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/narongdejsrn/go-thaiwordcut"
)

type Tag struct {
	Feature string `json:"des"`
}

const DATABASE  = "root:P@ssword@tcp(35.220.204.174:3306)/N&N_Cafe?charset=utf8"

func Test(strType string) int {
	db, err := sql.Open("mysql", DATABASE)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	var ctx = context.Background()
	selectMessages, err := db.QueryContext(ctx, "SELECT message FROM collections WHERE types=?",strType)
	rawText :=""

	for selectMessages.Next() {
		var tag Tag
		err = selectMessages.Scan(&tag.Feature)
		if err != nil {
			panic(err.Error())
		}
		//fmt.Println(tag.Des)
		rawText+=tag.Feature +" "
	}

	fmt.Print("Rawdata : ")
	fmt.Println(rawText)

	segmenter := gothaiwordcut.Wordcut()
	segmenter.LoadDefaultDict()
	//result := segmenter.Segment("ช่วยแนะนำเร้าเตอร์ที่ส่งสัญญาณ5Ghzได้หน่อยค่ะ")
	res := segmenter.Segment(rawText)

	//fmt.Println(rawText)
	fmt.Print("After cut : ")
	fmt.Println(res)

	featureStr := concatString(res)
	updateToFeatures, err := db.Prepare("UPDATE features SET features=? WHERE types=?")
	if err != nil {
		panic(err.Error())
		return 0
	}
	updateToFeatures.Exec(featureStr, strType)


	//Word ranking

	//make counter
	wordrank := make(map[string]int)
	for i:=0; i< len(res); i++{
		wordrank[res[i]] = 0
	}

	//Define key
	keys := make([]string, 0)
	for key := range wordrank{
		keys = append(keys, key)
	}
	fmt.Print("Key is : ")
	fmt.Println(keys)

	//Count features
	for j:=0; j< len(keys); j++  {
		for k:=0; k< len(res); k++  {
			if res[k]  == keys[j]{
				wordrank[keys[j]]++
			}
		}
	}
	fmt.Println(wordrank)
	fmt.Println()

	//strData := ""
	//for l:=0; l< len(keys); l++{
	//	strData += keys[l]+" "
	//}

	subFeaturesStr := concatString(keys)

	updates, err := db.Prepare("UPDATE features SET sub_features=? WHERE types=?")
	if err != nil {
		panic(err.Error())
		return 0
	}
	updates.Exec(subFeaturesStr, strType)

	return 1

	//Sort
	//mapForSort := make(map[int]string)
	//forSort := make([]int, 0)
	//for key, val := range wordrank{
	//	mapForSort[val] = key
	//	forSort = append(forSort, val)
	//}
	//
	//sort.Ints(forSort)
	//
	//for _, val := range wordrank {
	//	fmt.Println(mapForSort[val], val)
	//}



	//for k:=0; k< len(res); k++{
	//	rawText+=res[k]
	//}
	//
	//
	//tr := textrank.NewTextRank()
	//rule := textrank.NewDefaultRule()
	//language := textrank.NewDefaultLanguage()
	//algorithmDef := textrank.NewDefaultAlgorithm()
	//
	//tr.Populate(rawText, language, rule)
	//tr.Ranking(algorithmDef)
	//
	//text :=textrank.FindSingleWords(tr)
	//fmt.Println(text)
	//
	//result:=""
	//
	//for i:=0;i< len(text);i++  {
	//	//fmt.Println(text[i].Word)
	//	result+=text[i].Word+" "
	//}
	//fmt.Println(result)
	//
	//if err != nil {
	//	panic(err.Error())
	//}
	//
	//var ctx = context.Background()
	//
	//
	//tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	//_, execErr := tx.ExecContext(ctx, "INSERT INTO wordrank (word) VALUES (?)",result)
	//if execErr != nil {
	//	if rollbackErr := tx.Rollback(); rollbackErr != nil {
	//		log.Printf("Could not roll back: %v\n", rollbackErr)
	//	}
	//	log.Fatal(execErr)
	//}
	//if err := tx.Commit(); err != nil {
	//	log.Fatal(err)
	//}
	//
	//defer selectMessages.Close()
}

func concatString(arr []string) string{

	strData := ""
	for l:=0; l< len(arr); l++{
		strData += arr[l]+" "
	}

	return strData
}