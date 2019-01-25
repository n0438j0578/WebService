package main

import (
	"fmt"
	"github.com/akreal/knn"
)

func main() {
	knn := knn.NewKNN()

	knn.Train("Hello hello!", "class1")
	knn.Train("Hello, hello.", "class2")
	knn.Train("Hello.", "class3")

	k := 2

	predictedClass := knn.Predict("Say hello!", k)

	fmt.Println(predictedClass)
}