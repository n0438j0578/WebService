package main

import (
	"fmt"
	"github.com/narongdejsrn/go-thaiwordcut"
)

func main() {
	segmenter := gothaiwordcut.Wordcut()
	segmenter.LoadDefaultDict()
	//result := segmenter.Segment("ช่วยแนะนำเร้าเตอร์ที่ส่งสัญญาณ5Ghzได้หน่อยค่ะ")
	result := segmenter.Segment("Linksys LSS-EA9300-AH Max-Stream AC4000 Tri-Band Wi-Fi Router")
	fmt.Println(result)
}


