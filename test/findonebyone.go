package test

import (
	"WebService/model"
	"bytes"
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/joho/sqltocsv"

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
func TestoneByoneNormal(input string, featuregreeting []string, featureproblem []string, featureorders []string, featuresearch []string) (string, []model.ProductRow) {
	db, err := sql.Open("mysql", DATABASE)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	fmt.Println(input)

	segmenter := gothaiwordcut.Wordcut()
	segmenter.LoadDefaultDict()
	res := segmenter.Segment(input)

	fmt.Print("Cutdata : ")
	fmt.Println(res)

	result := ""

	for i := 0; i < len(res); i++ {
		result += res[i] + " "
	}

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

	boo := (greeting == problem && greeting == orders && greeting == search)
	if boo {
		return "", []model.ProductRow{}
	}

	fmt.Println(greeting, problem, orders, search)

	updateToFeatures, err := db.Prepare("UPDATE collections SET greeting=?,problem=?,orders=?,search=? WHERE  message=?")
	if err != nil {
		panic(err.Error())
		return "", []model.ProductRow{}
	}
	updateToFeatures.Exec(greeting, problem, orders, search, input)

	rows, _ := db.Query("SELECT greeting ,problem ,orders ,search,types  FROM collections")

	err = sqltocsv.WriteFile("test/report.csv", rows)

	//fname := `report.csv`
	f, err := os.OpenFile("test/report.csv", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return "", []model.ProductRow{}
	}
	defer f.Close()
	_, err = popLine(f)
	if err != nil {
		fmt.Println(err)
		return "", []model.ProductRow{}
	}
	//fmt.Print("pop:", string(line))

	if err != nil {
		panic(err)
	} else {
		//fmt.Println("เสร็จแล้ว")
	}

	//เริ่มทำการ knn

	irisMatrix := [][]string{}
	iris, err := os.Open("test/report.csv")
	if err != nil {
		panic(err)
	}
	defer iris.Close()

	reader := csv.NewReader(iris)
	reader.Comma = ','
	reader.LazyQuotes = true
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		irisMatrix = append(irisMatrix, record)
	}

	//split data into explaining and explained variables
	X := [][]float64{}
	Y := []string{}
	for _, data := range irisMatrix {

		//convert str slice data into float slice data
		temp := []float64{}
		for _, i := range data[:4] {
			parsedValue, err := strconv.ParseFloat(i, 64)
			if err != nil {
				panic(err)
			}
			temp = append(temp, parsedValue)
		}
		//fmt.Println(temp)
		//explaining variables
		X = append(X, temp)
		//explained variables
		Y = append(Y, data[4])

	}
	//fmt.Println(X)
	//fmt.Println(Y)

	//split data into training and test
	var (
		trainX [][]float64
		trainY []string
		testX  [][]float64
	)
	for i, _ := range X {
		trainX = append(trainX, X[i])
		trainY = append(trainY, Y[i])

	}

	temp := []float64{float64(greeting), float64(problem), float64(orders), float64(search)}
	//fmt.Println(temp)

	testX = append(testX, temp)
	//training
	knn := KNN{}
	knn.k = 3
	knn.fit(trainX, trainY)

	predicted := knn.predict(testX)

	fmt.Println(predicted[0])

	if strings.Compare(predicted[0], "search") == 0 {
		product := ProductMatching(result)
		return "", product
	} else {

		index := questionMatching(input, predicted[0])
		ans := ""
		var ctx = context.Background()
		selectMessages, err := db.QueryContext(ctx, "SELECT answer FROM collections WHERE id=?", index)
		for selectMessages.Next() {
			var tag Tag
			err = selectMessages.Scan(&tag.Feature)
			if err != nil {
				panic(err.Error())
			}
			ans = ans + tag.Feature
		}
		rawText := ""
		//rawtest:=""

		//fmt.Println(rawText)
		if strings.Compare(ans, "") != 0 {
			cut := strings.Split(ans, ":;")

			fmt.Println(cut, len(cut))
			if len(cut) != 1 {
				rawText = cut[rand.Intn(len(cut)-1)]
				for {
					if strings.Compare(rawText, "") == 0 {
						fmt.Println("เจอด้วยหรอวะ")
						cut := strings.Split(ans, ":;")
						rawText = cut[rand.Intn(len(cut)-1)]
					} else {
						break
					}
				}
			} else {
				cut := strings.Split(ans, ":;")
				rawText = cut[0]
				//ans = rawText
				fmt.Println(rawText)
			}

		}

		insForm, err := db.Prepare("INSERT INTO collections(message,types,answer,sub_feature,count,greeting,problem,orders,search) VALUES (?,?,?,?,?,?,?,?,?)")
		if err != nil {
			panic(err.Error())
		}
		_, err = insForm.Exec(input, predicted[0], ans, result, 0, greeting, problem, orders, search)

		return rawText, []model.ProductRow{}
	}

}

//calculate euclidean distance betwee two slices
func Dist(source, dest []float64) float64 {
	val := 0.0
	for i, _ := range source {
		val += math.Pow(source[i]-dest[i], 2)
	}
	return math.Sqrt(val)
}

//argument sort
type Slice struct {
	sort.Interface
	idx []int
}

func (s Slice) Swap(i, j int) {
	s.Interface.Swap(i, j)
	s.idx[i], s.idx[j] = s.idx[j], s.idx[i]
}

func NewSlice(n sort.Interface) *Slice {
	s := &Slice{Interface: n, idx: make([]int, n.Len())}
	for i := range s.idx {
		s.idx[i] = i
	}
	return s
}

func NewFloat64Slice(n []float64) *Slice { return NewSlice(sort.Float64Slice(n)) }

//map sort
type Entry struct {
	name  string
	value int
}
type List []Entry

func (l List) Len() int {
	return len(l)
}

func (l List) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l List) Less(i, j int) bool {
	if l[i].value == l[j].value {
		return l[i].name < l[j].name
	} else {
		return l[i].value > l[j].value
	}
}

//count item frequence in slice
func Counter(target []string) map[string]int {
	counter := map[string]int{}
	for _, elem := range target {
		counter[elem] += 1
	}
	return counter
}

type KNN struct {
	k      int
	data   [][]float64
	labels []string
}

func (knn *KNN) fit(X [][]float64, Y []string) {
	//read data
	knn.data = X
	knn.labels = Y
}

func (knn *KNN) predict(X [][]float64) []string {

	predictedLabel := []string{}
	for _, source := range X {
		var (
			distList   []float64
			nearLabels []string
		)
		//calculate distance between predict target data and surpervised data
		for _, dest := range knn.data {
			distList = append(distList, Dist(source, dest))
		}
		//take top k nearest item's index
		s := NewFloat64Slice(distList)
		sort.Sort(s)
		targetIndex := s.idx[:knn.k]

		//get the index's label
		for _, ind := range targetIndex {
			nearLabels = append(nearLabels, knn.labels[ind])
		}

		//get label frequency
		labelFreq := Counter(nearLabels)

		//the most frequent label is the predict target label
		a := List{}
		for k, v := range labelFreq {
			e := Entry{k, v}
			a = append(a, e)
		}
		sort.Sort(a)
		predictedLabel = append(predictedLabel, a[0].name)
	}
	return predictedLabel

}

func popLine(f *os.File) ([]byte, error) {
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(make([]byte, 0, fi.Size()))

	_, err = f.Seek(0, os.SEEK_SET)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(buf, f)
	if err != nil {
		return nil, err
	}
	line, err := buf.ReadString('\n')
	if err != nil && err != io.EOF {
		return nil, err
	}

	_, err = f.Seek(0, os.SEEK_SET)
	if err != nil {
		return nil, err
	}
	nw, err := io.Copy(f, buf)
	if err != nil {
		return nil, err
	}
	err = f.Truncate(nw)
	if err != nil {
		return nil, err
	}
	err = f.Sync()
	if err != nil {
		return nil, err
	}

	_, err = f.Seek(0, os.SEEK_SET)
	if err != nil {
		return nil, err
	}
	return []byte(line), nil
}
