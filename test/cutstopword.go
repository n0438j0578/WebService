package test

import (
	"encoding/csv"
	"github.com/narongdejsrn/go-thaiwordcut"
	"io"
	"os"
	"strings"
)

func CutStopWord(text string)string   {
	//start := time.Now()
	segmenter := gothaiwordcut.Wordcut()
	segmenter.LoadDefaultDict()
	input := segmenter.Segment(text)
	stopword := []string{}
	file, err := os.Open("test/stopword.txt")
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
	//fmt.Println(stopword)
	result :=[]string{}

	count :=0
	for i:=0;i< len(input);i++  {
		count =0
		for j:=0;j<len(stopword);j++ {
			if!(strings.Compare(input[i],stopword[j])==0){
				count++;
			}
		}
		if(count==len(stopword)){
			result= append(result,input[i])
		}
	}
	//fmt.Println(result)
	real:=""
	for i:=0;i< len(result);i++  {
		real=real +result[i] +" "
	}

	return real
	//elapsed := time.Since(start)
	//fmt.Print("Elapse time :")
	//fmt.Println(elapsed)

}