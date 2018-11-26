package main

import (
	"fmt"
	"github.com/narongdejsrn/go-thaiwordcut"
)

func main() {
	segmenter := gothaiwordcut.Wordcut()
	segmenter.LoadDefaultDict()
	result := segmenter.Segment("หน้าหนังหี")
	fmt.Println(result)
}
