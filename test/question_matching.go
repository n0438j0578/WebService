package test

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/narongdejsrn/go-thaiwordcut"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type QuestionRow struct {
	/*Msg string `json:"message"`*/
	Answer string `json:"answer"`
	ID int `json:"id"`
	SubFeature string `json:"sub_feature"`
	count int
}

func questionMatching(msg string, strType string) int{

	msgFeatures := subFeature(msg)

	var qsArr []QuestionRow

	db, err := sql.Open("mysql", DATABASE)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	var ctx = context.Background()
	selectMessages, err := db.QueryContext(ctx, "SELECT answer, id, sub_feature FROM collections WHERE types=?",strType)

	qsArr = nil

	for selectMessages.Next() {

		var qsTmp QuestionRow
		err = selectMessages.Scan(&qsTmp.Answer, &qsTmp.ID, &qsTmp.SubFeature)
		if err != nil {
			panic(err.Error())
		}
		qsTmp.count = 0

		qsArr = append(qsArr,qsTmp)
	}

	//fmt.Println(qsArr)

	for i:=0; i< len(msgFeatures); i++  {
		for j:=0; j< len(qsArr); j++  {
			if strings.Contains(qsArr[j].SubFeature,msgFeatures[i]) {
				qsArr[j].count++
				//fmt.Println("qsArr= "+qsArr[j].SubFeature+" ,ID= "+strconv.Itoa(qsArr[j].ID)+" , msgFeature= "+msgFeatures[i])
			}

		}

	}

	max := MinMax(qsArr)
	fmt.Println("Maximum count :"+strconv.Itoa(max))

	var idSet []int

	for i:=0; i< len(qsArr); i++  {
		if qsArr[i].count == max {
			fmt.Println("ID : "+strconv.Itoa(qsArr[i].ID)+", Count :"+strconv.Itoa(qsArr[i].count)+", Ans :"+qsArr[i].Answer)
			idSet = append(idSet, qsArr[i].ID)
		}
	}

	fmt.Print("ID set : ")
	fmt.Println(idSet)

	rand.Seed(time.Now().UnixNano())
	fmt.Print("Random ID :")
	randomId := idSet[rand.Intn(len(idSet))]
	fmt.Println(randomId)


	return randomId
}

func subFeature(str string) []string{

	segmenter := gothaiwordcut.Wordcut()
	segmenter.LoadDefaultDict()
	res := segmenter.Segment(str)

	//make counter and store cut word to map (not store duplicate word)
	keyCount := make(map[string]int)
	for i:=0; i< len(res); i++{
		keyCount[res[i]] = 0
	}

	//Define key
	keys := make([]string, 0)
	for key := range keyCount{
		keys = append(keys, key)
	}


	return keys
}

func MinMax(array []QuestionRow) /*(int,*/ int/*)*/ {
	var max = array[0].count
	//var min int
	for i:=0; i< len(array); i++ {
		if max < array[i].count {
			max = array[i].count
		}
		//if min > array[i].count {
		//	min = array[i].count
		//}
	}
	return /*min,*/ max
}


func main() {
	start := time.Now()
	questionMatching("สวัสดี", "greeting")
	//questionMatching("เจ้าเน็คไม่ได้", "problem")
	//questionMatching("มีไหม", "search")
	//questionMatching("อยากได้", "order")

	fmt.Println()
	elapsed := time.Since(start)
	fmt.Print("Elapse time :")
	fmt.Println(elapsed)
}
