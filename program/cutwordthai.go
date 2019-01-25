package main

import (
	"fmt"
	"github.com/narongdejsrn/go-thaiwordcut"
)

func main() {
	segmenter := gothaiwordcut.Wordcut()
	segmenter.LoadDefaultDict()
	//result := segmenter.Segment("ช่วยแนะนำเร้าเตอร์ที่ส่งสัญญาณ5Ghzได้หน่อยค่ะ")
	result := segmenter.Segment("เร้าเตอร์มีไฟขึ้นสีแดง ไม่สามารถเข้าอินเตอร์เน็ตได้")
	fmt.Println(result)
}


