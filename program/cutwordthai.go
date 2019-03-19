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
	result := segmenter.Segment("D-Link DAP-1665 Wireless AC1200 MU-MIMO Dual Band Range Extender Access Point (7in1 Multi-Mode Device)มี D Link ไหมครับ")
	fmt.Println(result)
	fmt.Println(strings.ToLower("Linksys LSS-EA9300-AH Max-Stream AC4000 Tri-Band Wi-Fi Router"))
}


