package main

import (
	"fmt"
	"github.com/narongdejsrn/go-thaiwordcut"
	"regexp"
	"strings"
)

func main() {
	segmenter := gothaiwordcut.Wordcut()
	segmenter.LoadDefaultDict()
	//result := segmenter.Segment("ช่วยแนะนำเร้าเตอร์ที่ส่งสัญญาณ5Ghzได้หน่อยค่ะ")
	result := segmenter.Segment("แบบเดิม")
	fmt.Println(result)
	fmt.Println(strings.ToLower("Linksys LSS-EA9300-AH Max-Stream AC4000 Tri-Band Wi-Fi Router"))
	var validID = regexp.MustCompile(`[0-9]\s*\.*\**\\*\t*\n*\r*:\s*\.*\**\\*\t*\n*\r*[0-9]`)

	fmt.Println(validID.MatchString("123				:  123123"))

	words := strings.Fields("123:123123")
	for i:=0;i< len(words);i++  {
		fmt.Println(words[i])
	}
	fmt.Println(len(strings.Split("123				:  123123", ":")))
}


