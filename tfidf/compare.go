package main

import (
	"fmt"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/xrash/smetrics"
)

const (
	text1 = "Lorem ipsum dolor."
	text2 = "Lorem dolor sit amet."
)

func main() {
	dmp := diffmatchpatch.New()

	diffs := dmp.DiffMain(text1, text2, false)

	fmt.Println(dmp.DiffPrettyText(diffs))
	fmt.Println(smetrics.JaroWinkler("AL", "AL", 0.7, 4))

}