package test

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"strings"
)

func CheckRepeat(input string) bool {
	dat, err := ioutil.ReadFile("program/ignoreWord.txt")
	if err != nil {
		fmt.Print(err)
	}

	str := string(dat)
	fmt.Println(str)
	fmt.Println()

	fmt.Print("Result : ")
	if strings.Contains(str, input) {
		fmt.Println("Matched")
		return true
	} else {
		fmt.Println("Not match")
		return false
	}



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


