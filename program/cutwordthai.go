package main

import (
	"fmt"
	"github.com/narongdejsrn/go-thaiwordcut"
	"strings"
)

func main() {
	segmenter := gothaiwordcut.Wordcut()
	segmenter.LoadDefaultDict()
	//result := segmenter.Segment("ช่วยแนะนำเร้าเตอร์ที่ส่งสัญญาณ5Ghzได้หน่อยค่ะ")
	result := segmenter.Segment("แบบเดิม")
	fmt.Println(result)
	fmt.Println(strings.ToLower("Linksys LSS-EA9300-AH Max-Stream AC4000 Tri-Band Wi-Fi Router"))
}


