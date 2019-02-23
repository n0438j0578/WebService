package test

import (
	"WebService/model"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

func ProductMatching(msg string) []model.ProductRow{

	msgFeatures := subFeature(msg)

	product :=[]model.ProductRow{}

	db, err := sql.Open("mysql", DATABASE)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()



	//var ctx = context.Background()
	selectMessages, err := db.Query( "SELECT name, des, img,id FROM menu WHERE amount>0")


	for selectMessages.Next() {

	var pro model.ProductRow

	err = selectMessages.Scan(&pro.Name, &pro.Des, &pro.Img,&pro.ID)
	if err != nil {
	panic(err.Error())
	}

	product = append(product,pro)
	}


	//fmt.Println(product)

	for i:=0; i< len(msgFeatures); i++  {
		for j:=0; j< len(product); j++  {
			if strings.Contains(product[j].Name,msgFeatures[i]) {
				product[j].Count++
				//fmt.Println("product= "+product[j].SubFeature+" ,ID= "+strconv.Itoa(product[j].ID)+" , msgFeature= "+msgFeatures[i])
			}

		}

	}

	max := minMax(product)
	fmt.Println("Maximum count :"+strconv.Itoa(max))

	if(max==0){
		return []model.ProductRow{}
	}

	//var idSet []int
	var result []model.ProductRow

	for i:=0; i< len(product); i++  {
		if product[i].Count == max {
			//fmt.Println("ID : "+strconv.Itoa(product[i].ID)+", Count :"+strconv.Itoa(product[i].count)+", Ans :"+ product[i].Answer)
			fmt.Println("ID :"+strconv.Itoa(product[i].ID))
			//idSet = append(idSet, product[i].ID)
			result = append(result,product[i])
		}
	}

	fmt.Println("Result :"+result[0].Name)






	return result
}

func minMax(array []model.ProductRow) /*(int,*/ int/*)*/ {
	var max = array[0].Count
	//var min int
	for i:=0; i< len(array); i++ {
		if max < array[i].Count {
			max = array[i].Count
		}
		//if min > array[i].count {
		//	min = array[i].count
		//}
	}
	return /*min,*/ max
}


