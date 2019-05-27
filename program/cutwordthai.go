package main

import (
	"fmt"
	"github.com/narongdejsrn/go-thaiwordcut"

)

func main() {
	segmenter := gothaiwordcut.Wordcut()
	segmenter.LoadDefaultDict()
	result := segmenter.Segment("รองเท้ากัดค่ะขอเปลี่ยนได้ไหมคะ")
	fmt.Println(result)
}


