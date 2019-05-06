package test

import (
	"context"
	"database/sql"
	"strconv"
	"strings"
)

func ProductCal(msg string, Idcustomer string) string {
	text := strings.Fields(msg)
	var cut [2]string

	if(len(text)==1){
		txt :=strings.Split(msg, ":")
		cut[0]=txt[0]
		cut[1]=txt[1]
	}else{
		cut[0]= text[0]
		cut[1] =text[2]
	}

	db, err := sql.Open("mysql", "root:P@ssword@tcp(35.220.204.174:3306)/N&N_Cafe?charset=utf8")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var ctx = context.Background()
	selectMessages, err := db.QueryContext(ctx, "SELECT name, price FROM menu WHERE amount>0 AND id=? AND amount>=?", cut[0],cut[1])

	for selectMessages.Next() {
		type Product struct {
			Name  string
			Price float64
		}
		var pro Product

		err = selectMessages.Scan(&pro.Name, &pro.Price)
		if err != nil {
			panic(err.Error())
		}
		//แปลงตัวคุณจำนวนสินค้ากับราคา
		count, _ := strconv.ParseFloat(cut[1], 64)
		price := count * pro.Price

		cutprice := strings.Split(strconv.FormatFloat(price, 'f', -1, 64), ".")

		answer := "ชื่อสินค้า : " + pro.Name + "\n" + "เป็นจำนวน : " + cut[1] + "\n" + "ราคารวมทั้งหมด " + cutprice[0]

		selectMessages, err := db.QueryContext(ctx, "SELECT orderold FROM oldmsg WHERE id=?", Idcustomer)
		rawText := ""

		for selectMessages.Next() {
			var tag Tag
			err = selectMessages.Scan(&tag.Feature)
			if err != nil {
				panic(err.Error())
			}
			rawText += tag.Feature
		}

		insForm, _ := db.Prepare("UPDATE oldmsg SET orderold=? WHERE id=? ")
		insForm.Exec(answer, Idcustomer)

		return answer
	}
	return "ไม่มีสินค้าอยู่ในระบบหรือสินค้าหมด"

}

