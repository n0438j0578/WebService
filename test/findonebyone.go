package test

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/narongdejsrn/go-thaiwordcut"
	"strings"
	"sync"
)

type Id struct {
	id int
}


func TestAll() int{
	db, err := sql.Open("mysql", DATABASE)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	//var ctx = context.Background()
	selectMessages, err := db.Query("SELECT id FROM collections ")

	//check :=2
	var test []int
	test =append(test,1)
	for selectMessages.Next() {
		var tag Id
		err = selectMessages.Scan(&tag.id)
		if err != nil {
			panic(err.Error())
		}
		//fmt.Println(tag.Des)
		test =append(test,tag.id)
		//check =TestoneByone(tag.id)
		//fmt.Println(tag.id)
		//
		//if(check!=1){
		//	return 0
		//}
	}
	var wg sync.WaitGroup
	wg.Add(len(test))
	for _, index := range test {
		// เราต้องส่ง reference ของ wg ไปด้วย เพื่อที่จะสั่ง Done
		go TestoneByone(index, &wg)
	}
	wg.Wait()

	if err != nil {
		return 0
	}else{
		return 1
	}

}
func TestoneByone(index int,wg *sync.WaitGroup) int {
	db, err := sql.Open("mysql", DATABASE)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	var ctx = context.Background()
	selectMessages, err := db.QueryContext(ctx, "SELECT message FROM collections WHERE id=?", index)
	rawText := ""

	for selectMessages.Next() {
		var tag Tag
		err = selectMessages.Scan(&tag.Feature)
		if err != nil {
			panic(err.Error())
		}
		//fmt.Println(tag.Des)
		rawText += tag.Feature
	}

	//fmt.Print("Rawdata : ")
	//fmt.Println(rawText)

	segmenter := gothaiwordcut.Wordcut()
	segmenter.LoadDefaultDict()
	res := segmenter.Segment(rawText)



	fmt.Print("Cutdata : ")
	fmt.Println(res)

	result:=""

	for i := 0;i< len(res);i++  {
		result+= res[i]+" "
	}


	updateToFeatures, err := db.Prepare("UPDATE collections SET sub_feature=? WHERE id=?")
	if err != nil {
		panic(err.Error())
		return 0
	}
	updateToFeatures.Exec(result,index)

	greeting := 0
	problem := 0
	orders := 0
	search := 0


	for i := 0; i < len(res); i++ {

		if (findfeaturesonebyone(res[i],"greeting") == 1) {
			greeting++
		}
		if (findfeaturesonebyone(res[i],"problem") == 1) {
			problem++
		}
		if (findfeaturesonebyone(res[i],"orders") == 1) {
			orders++
		}
		if (findfeaturesonebyone(res[i],"search") == 1) {
			search++
		}
	}

	updateToFeatures, err = db.Prepare("UPDATE collections SET greeting=?,problem=?,orders=?,search=? WHERE id=?")
	if err != nil {
		panic(err.Error())
		return 0
		wg.Done()
	}
	updateToFeatures.Exec(greeting, problem, orders, search,index)
	wg.Done()

	return 1

}
func findfeaturesonebyone(input string,types string) int {
	db, err := sql.Open("mysql", DATABASE)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	var ctx = context.Background()

	selectMessages, err := db.QueryContext(ctx, "SELECT sub_features FROM features WHERE types=?", types)
	rawText := ""

	for selectMessages.Next() {
		var tag Tag
		err = selectMessages.Scan(&tag.Feature)
		if err != nil {
			panic(err.Error())
		}
		rawText += tag.Feature
	}


	cut := strings.Split(rawText, " ")


	check:=0

    for i := 0; i< len(cut);i++  {
		check =strings.Compare(input,cut[i])
		if(check==0){
			return 1
		}
	}

	return 0

}
