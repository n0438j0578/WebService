package main

import (
	"fmt"
	"github.com/sjwhitworth/golearn/base"
	"github.com/sjwhitworth/golearn/evaluation"
	"github.com/sjwhitworth/golearn/filters"
	"github.com/sjwhitworth/golearn/trees"
	"math/rand"
)

func main() {

	var tree base.Classifier

	rand.Seed(44111342)

	// Load in the iris dataset
	iris, err := base.ParseCSVToInstances("report.csv", false)
	if err != nil {
		panic(err)
	}

	// Discretise the iris dataset with Chi-Merge
	filt := filters.NewChiMergeFilter(iris, 0.999)
	for _, a := range base.NonClassFloatAttributes(iris) {
		filt.AddAttribute(a)
	}
	filt.Train()
	irisf := base.NewLazilyFilteredInstances(iris, filt)

	// Create a 60-40 training-test split
	_, testData := base.InstancesTrainTestSplit(irisf, 0.99)
	//fmt.Println(testData)
	//// Consider two randomly-chosen attributes
	tree = trees.NewRandomTree(2)
	err = tree.Fit(irisf)
	if err != nil {
		panic(err)
	}
	var testX  [][]float64
	temp := []float64{float64(1),float64(2),float64(3),float64(4)}
	testX = append(testX,temp)
	base.NonClassFloatAttributes(testX)
	testData.AddClassAttribute(testX)
	predictions, err := tree.Predict(ir)
	if err != nil {
		panic(err)
	}

	//fmt.Println("RandomTree Performance")
	//cf, err := evaluation.GetConfusionMatrix(testData, predictions)
	//if err != nil {
	//	panic(fmt.Sprintf("Unable to get confusion matrix: %s", err.Error()))
	//}
	//fmt.Println(evaluation.GetSummary(cf))

	//
	// Finally, Random Forests
	//
}
