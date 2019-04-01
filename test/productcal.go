package test

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

func ProductCal(msg string) string {

	cut :=strings.Split(msg,":")

	//product := []model.ProductRow{}

	db, err := sql.Open("mysql", "root:P@ssword@tcp(35.220.204.174:3306)/N&N_Cafe?charset=utf8")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var ctx = context.Background()
	selectMessages, err := db.QueryContext(ctx, "SELECT name, price FROM menu WHERE amount>0 AND id=?", cut[0])

	for selectMessages.Next() {
		type Product struct {
			Name string
			Price float64
		}
		var pro Product

		err = selectMessages.Scan(&pro.Name, &pro.Price)
		if err != nil {
			panic(err.Error())
		}
		count,_ := strconv.ParseFloat(cut[1], 64)
		price := count*pro.Price
		answer := "ชื่อสินค้า : "+pro.Name+"\n"+"เป็นจำนวน : "+cut[1]+"\n"+"ราคารวมทั้งหมด"+fmt.Sprintf("%f", price)
		return answer
	}
	return ""


}
