package test

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"

	gothaiwordcut "github.com/narongdejsrn/go-thaiwordcut"
)

type Id struct {
	id int
}

func TestAll() int {
	db, err := sql.Open("mysql", DATABASE)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	selectMessages, err := db.Query("SELECT id FROM collections ")
	var test []int
	for selectMessages.Next() {
		var tag Id
		err = selectMessages.Scan(&tag.id)
		if err != nil {
			panic(err.Error())
		}
		test = append(test, tag.id)
	}

	featuregreeting := Selectfeature("greeting")
	featureproblem := Selectfeature("problem")
	featureorders := Selectfeature("order")
	featuresearch := Selectfeature("search")

	var wg sync.WaitGroup
	wg.Add(len(test))
	for _, index := range test {
		go TestoneByone(index, &wg, featuregreeting, featureproblem, featureorders, featuresearch)
	}
	wg.Wait()

	if err != nil {
		return 0
	} else {
		return 1
	}

}
func TestoneByone(index int, wg *sync.WaitGroup, featuregreeting []string, featureproblem []string, featureorders []string, featuresearch []string) int {
	db, err := sql.Open("mysql", DATABASE)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	var ctx = context.Background()
	selectMessages, err := db.QueryContext(ctx, "SELECT message FROM collections WHERE id=?", index)
	//fmt.Println(selectMessages)
	rawText := ""

	for selectMessages.Next() {
		var tag Tag
		err = selectMessages.Scan(&tag.Feature)
		if err != nil {
			panic(err.Error())
		}
		rawText += tag.Feature
	}

	segmenter := gothaiwordcut.Wordcut()
	segmenter.LoadDefaultDict()
	res := segmenter.Segment(rawText)

	fmt.Print("Cutdata : ")
	fmt.Println(res)

	result := ""

	for i := 0; i < len(res); i++ {
		result += res[i] + " "
	}

	updateToFeatures, err := db.Prepare("UPDATE collections SET sub_feature=? WHERE id=?")
	if err != nil {
		panic(err.Error())
		return 0
	}
	updateToFeatures.Exec(result, index)

	greeting := 0
	problem := 0
	orders := 0
	search := 0

	for i := 0; i < len(res); i++ {

		if Findfeaturesonebyone(res[i], featuregreeting) == 1 {
			greeting++
		}
		if Findfeaturesonebyone(res[i], featureproblem) == 1 {
			problem++
		}
		if Findfeaturesonebyone(res[i], featureorders) == 1 {
			orders++
		}
		if Findfeaturesonebyone(res[i], featuresearch) == 1 {
			search++
		}
	}

	updateToFeatures, err = db.Prepare("UPDATE collections SET greeting=?,problem=?,orders=?,search=? WHERE id=?")
	if err != nil {
		panic(err.Error())
		return 0
		wg.Done()
	}
	updateToFeatures.Exec(greeting, problem, orders, search, index)
	wg.Done()

	return 1

}
func Findfeaturesonebyone(input string, cut []string) int {
	check := 2
	for i := 0; i < len(cut); i++ {
		check = strings.Compare(input, cut[i])
		if check == 0 {
			return 1
		}
	}

	return 0

}

func Selectfeature(types string) []string {
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

	return cut

}
