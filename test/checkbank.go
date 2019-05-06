package test

import (
	"encoding/csv"
	"github.com/narongdejsrn/go-thaiwordcut"
	"io"
	"os"
	"strings"
)

func CheckBank(text string)bool   {
	//start := time.Now()
	segmenter := gothaiwordcut.Wordcut()
	segmenter.LoadDefaultDict()
	input := segmenter.Segment(text)
	stopword := []string{}
	file, err := os.Open("test/bank.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','
	reader.LazyQuotes = true
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		stopword = record
	}



	for i:=0;i< len(input);i++  {
		for j:=0;j<len(stopword);j++ {
			if(strings.Compare(input[i],stopword[j])==0){
				return true
			}
		}

	}
	return false

}
