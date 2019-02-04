package data

import (
	"WebService/test"
	"database/sql"
	"fmt"
	"github.com/narongdejsrn/go-thaiwordcut"
)

type Id struct {
	id int
}

type Tag struct {
	Feature string `json:"des"`
}
const DATABASE  = "root:P@ssword@tcp(35.220.204.174:3306)/N&N_Cafe?charset=utf8"


func WordSet(text string,types string,ans string) int {
	db, err := sql.Open("mysql", DATABASE)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	//var ctx = context.Background()
	//selectMessages, err := db.QueryContext(ctx, "INSERT INTO collections(message,types,answer) VALUES (?,?,?)", text,types)



	segmenter := gothaiwordcut.Wordcut()
	segmenter.LoadDefaultDict()
	res := segmenter.Segment(text)

	fmt.Print("Cutdata : ")
	fmt.Println(res)

	result := ""

	for i := 0; i < len(res); i++ {
		result += res[i] + " "
	}

	insForm, err := db.Prepare("INSERT INTO collections(message,types,answer,sub_feature) VALUES (?,?,?,?)")
	if err != nil {
		panic(err.Error())
		return 0
	}
	_,err=insForm.Exec(text, types,ans,result)
	if err != nil {
		return 0
	}

	num := test.Test(types)


	//featuregreeting := test.Selectfeature("greeting")
	//featureproblem := test.Selectfeature("problem")
	//featureorders := test.Selectfeature("orders")
	//featuresearch := test.Selectfeature("search")

	if(num==0){
		return 0
	}
	num = test.TestAll()
	if(num==0){
		return 0
	}
	//
	//greeting := 0
	//problem := 0
	//orders := 0
	//search := 0
	//
	//for i := 0; i < len(res); i++ {
	//
	//	if test.Findfeaturesonebyone(res[i], featuregreeting) == 1 {
	//		greeting++
	//	}
	//	if test.Findfeaturesonebyone(res[i], featureproblem) == 1 {
	//		problem++
	//	}
	//	if test.Findfeaturesonebyone(res[i], featureorders) == 1 {
	//		orders++
	//	}
	//	if test.Findfeaturesonebyone(res[i], featuresearch) == 1 {
	//		search++
	//	}
	//}
	//updateToFeatures, err := db.Prepare("UPDATE collections SET greeting=?,problem=?,orders=?,search=? WHERE message=?")
	//if err != nil {
	//	panic(err.Error())
	//	return 0
	//}
	//updateToFeatures.Exec(greeting, problem, orders, search, text)



	return 1






	//INSERT INTO `collections` (`message`, `greeting`, `problem`, `orders`, `search`, `types`, `answer`, `id`, `sub_feature`) VALUES ('', '0', '0', '0', '0', NULL, '', NULL, NULL)


}