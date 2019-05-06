package test

import (
	"context"
	"database/sql"
	"encoding/csv"
	"github.com/narongdejsrn/go-thaiwordcut"
	"io"
	"os"
	"strings"
)

func CheckRepeat(text string) bool {
	segmenter := gothaiwordcut.Wordcut()
	segmenter.LoadDefaultDict()
	input := segmenter.Segment(text)
	stopword := []string{}
	file, err := os.Open("test/ignoreWord.txt")
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


func SendRepeat(Idcustomer string) string {
	//ทำการต่อฐานข้อมูล
	db, err := sql.Open("mysql", DATABASE)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	//เหมือนเป็นตัวแปรเอาไว้ใช้ในการใส่ฐานข้อมูลโดยสามารถส่งผ่านตัวแปรได้
	var ctx= context.Background()

	//ทำการลองดูว่าข้อความที่ได้ทำการเข้ามานั้นเคยเข้ามาหรือยัง ถ้าเคยแล้วเราจะสามารถเช็คได้แล้วให้ส่งผลกลับไปเลย
	selectMessages, err := db.QueryContext(ctx, "SELECT orderold FROM oldmsg WHERE id=?", Idcustomer)

	//ข้อความที่รับเข้ามาจากฐานข้อมูล
	rawText := ""

	for selectMessages.Next() {
		err = selectMessages.Scan(&rawText)
		if err != nil {
			panic(err.Error())
		}
	}
	if (strings.Compare(rawText, "") != 0) {
		return rawText
	}else{
		return "ลูกค้าไม่เคยทำการสั่งสินค้าจากร้านค้า"
	}

}


