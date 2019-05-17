package data

import (
	"WebService/model"
	"WebService/test"
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"regexp"
	"strings"

	"github.com/narongdejsrn/go-thaiwordcut"
)

type Id struct {
	id int
}

type Tag struct {
	Feature string `json:"des"`
	Count   int
}

const DATABASE = "root:n0438@j0578@tcp(35.220.204.174:3306)/N&N_Cafe?charset=utf8"

func WordSet(text string, types string, ans string) int {
	db, err := sql.Open("mysql", DATABASE)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	text = strings.ToLower(text)

	segmenter := gothaiwordcut.Wordcut()
	segmenter.LoadDefaultDict()
	res := segmenter.Segment(text)

	fmt.Print("Cutdata : ")
	fmt.Println(res)

	result := ""

	for i := 0; i < len(res); i++ {
		result += res[i] + " "
	}

	insForm, err := db.Prepare("INSERT INTO collections(message,types,answer,sub_feature,count) VALUES (?,?,?,?,?)")
	if err != nil {
		fmt.Print("Cutdata1 : ")
		fmt.Println(res)
		panic(err.Error())
		return 0
	}
	_, err = insForm.Exec(text, types, ans, result, 0)
	if err != nil {
		return 0
	}

	//num := test.Test(types)
	//
	//if num == 0 {
	//	return 0
	//}
	//num = test.TestAll()
	//if num == 0 {
	//	return 0
	//}
	return 1

}

func WordCome(text string, Idcustomer string) (int, string, []model.ProductRow) {

	text = strings.ToLower(text)
	//ทำการต่อฐานข้อมูล
	db, err := sql.Open("mysql", DATABASE)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	var ctx = context.Background()
	//fmt.Println(text)
	selectMessages, err := db.QueryContext(ctx, "SELECT answer,count FROM collections WHERE message=?", text)
	//fmt.Println(selectMessages)
	rawText := ""
	rawtest := ""
	count := 1

	for selectMessages.Next() {
		var tag Tag
		err = selectMessages.Scan(&tag.Feature, &tag.Count)
		if err != nil {
			panic(err.Error())
		}
		rawtest += tag.Feature
		count = count + tag.Count
	}
	//fmt.Println(rawtest)
	if (strings.Compare(rawtest, "") != 0) {
		cut := strings.Split(rawtest, ":;")
		//ลูบนี้จะเช็คประมานว่าถ้าตัดแล้วเจอข้อความเปล่าๆ ให้มันแรนดอมเอาใหม่อีกครั้ง
		if (len(cut) != 1) {
			rawText = cut[rand.Intn(len(cut)-1)]
			for ; ; {
				if (strings.Compare(rawText, "") == 0) {
					fmt.Println("เจอด้วยหรอวะ")
					cut := strings.Split(rawtest, ":;")
					rawText = cut[rand.Intn(len(cut)-1)]
				} else {
					break
				}
			}
		} else {
			rawText = rawtest

		}

	}
	cutstopword := test.CutStopWord(text)

	//ถ้าไม่เจอ
	if (strings.Compare(rawText, "") == 0) {
		//ทำการหาความถี่แต่ละตัวแล้วเอาไปเข้าฟังชั่น
		featuregreeting := test.Selectfeature("greeting")
		featureproblem := test.Selectfeature("problem")
		featureorders := test.Selectfeature("order")
		featuresearch := test.Selectfeature("search")

		//เซฟว่าผู้ใช้คนนี้เคยส่งอะไรมา
		SaveWord(text, Idcustomer)

		rawText, product := test.TestoneByoneNormal(cutstopword, featuregreeting, featureproblem, featureorders, featuresearch)
		fmt.Println(rawText)

		if (len(product) > 0) {
			//ถ้าเจอของ
			return 3, "", product
		} else if (strings.Compare(rawText, "") != 0 && len(product) == 0) {
			//ไม่เจอของแต่ว่าเป็นข้อความที่สามารถตอบกลับไปได้
			return 2, rawText, []model.ProductRow{}
		} else if (strings.Compare(rawText, "") == 0 || len(product) == 0) {
			//ไม่เจออะไรทั้งนั้น
			return 0, "", []model.ProductRow{}
		} else {
			fmt.Println("Test")
			return 2, rawText, []model.ProductRow{}
		}
	} else {
		//เจอข้อความแล้วให้เซฟแล้วค่อยรีเทินออกไปว่าเจอ
		insForm, _ := db.Prepare("UPDATE collections SET count=? WHERE message=? ")
		insForm.Exec(count, text)
		SaveWord(text, Idcustomer)
		return 1, rawText, []model.ProductRow{}
	}

}

func WordComeCosine(text string, Idcustomer string) (int, string, []model.ProductRow) {
	//ทำการโชว์ข้อความที่เข้ามาพร้อมเลขผู้ใช้
	fmt.Println("ผู้ใช้เลขที่ :",Idcustomer)
	fmt.Println("ข้อความ :",text)
	fmt.Println()

	//ให้มันทำเป็นตัวเล็กให้หมดก่อน
	text = strings.ToLower(text)

	//เซฟข้อความก่อนหน้าแล้วก็เซฟเลขไอดีของผู้ใช้เฟสบุคด้วย
	SaveWord(text, Idcustomer)

	//ทำการเช็ค ตัวแปรว่าเข้ามาโดยเป็นการสั่งซื้อสินค้าหรือไม่ โดยใช้ regular expression

	var validID = regexp.MustCompile(`[0-9]\s*\.*\**\\*\t*\n*\r*:\s*\.*\**\\*\t*\n*\r*[0-9]`)

	//ทำการเช็คว่่าเข้าสู่กรณีออเดอร์ซ้ำเหมือนเดิมไหม
	segmenter := gothaiwordcut.Wordcut()
	segmenter.LoadDefaultDict()

	var checkrepeat = segmenter.Segment(text)

	for i:=0;i< len(checkrepeat);i++  {
		if(test.CheckRepeat(checkrepeat[i])){
			fmt.Println("เป็นเข้าสู่กรณีออเดอร์ซ้ำเหมือนเดิม")
			return 1,test.SendRepeat(Idcustomer),[]model.ProductRow{}
		}
	}


	//ถ้าเข้าเงื่อนไข โปรแกรมจะรับไปคำนวนและคืนผลลัพธ์ที่เป็นข้อความแสดงรายละเอียดจำนวนสินค้าและราคาทั้งหมด
	if test.CheckBank(text){
		fmt.Println("เป็นกรณีธนาคาร")
		return 4,"เป็นกรณีธนาคาร",[]model.ProductRow{}

	}else if validID.MatchString(text){
		//ส่งข้อความกลับไป ถ้ากรณีแรกคือ 1 คือเจอเลย
		answer := test.ProductCal(text,Idcustomer)
		fmt.Println("กรณีเป็นรูปแบบการสั่งค้าทำการส่งข้อความแสดงรายละเอียดจำนวนสินค้าและราคาทั้งหมด")
		return 1, answer, []model.ProductRow{}

	}else {
		//ถ้าไม้เข้าเงื่อนไข

		//ทำการต่อฐานข้อมูล
		db, err := sql.Open("mysql", DATABASE)
		if err != nil {
			panic(err.Error())
		}
		defer db.Close()

		//เหมือนเป็นตัวแปรเอาไว้ใช้ในการใส่ฐานข้อมูลโดยสามารถส่งผ่านตัวแปรได้
		var ctx= context.Background()

		//ทำการลองดูว่าข้อความที่ได้ทำการเข้ามานั้นเคยเข้ามาหรือยัง ถ้าเคยแล้วเราจะสามารถเช็คได้แล้วให้ส่งผลกลับไปเลย
		selectMessages, err := db.QueryContext(ctx, "SELECT answer,count FROM collections WHERE message=?", text)

		//ข้อความเอาไว้ส่ง
		rawText := ""
		//ข้อความที่รับเข้ามาจากฐานข้อมูล
		rawtest := ""
		//เอาไว้บวกกับจำนวนครั้งถ้าเจอเลยแบบข้างล่าง
		count := 1

		for selectMessages.Next() {
			var tag Tag
			err = selectMessages.Scan(&tag.Feature, &tag.Count)
			if err != nil {
				panic(err.Error())
			}
			rawtest += tag.Feature
			count = count + tag.Count
		}

		//ถ้าเจอ สตริงมันจะต้องไม่ว่าง แล้วเราก็ต้องเอาตัดมาเลือกด้วย
		if (strings.Compare(rawtest, "") != 0) {

			cut := strings.Split(rawtest, ":;")

			//ถ้าขนาดไม่เท่ากับหนึ่งแปลว่ามีหลายคำตอบ

			if (len(cut) != 1) {
				rawText = cut[rand.Intn(len(cut)-1)]
				for ; ; {
					//ลูบนี้จะเช็คประมานว่าถ้าตัดแล้วเจอข้อความเปล่าๆ ให้มันแรนดอมเอาใหม่อีกครั้ง
					if strings.Compare(rawText, "") == 0 {
						cut := strings.Split(rawtest, ":;")
						rawText = cut[rand.Intn(len(cut)-1)]
					} else {
						break
					}
				}
			} else {
				rawText = rawtest
			}

			//ทำการเพิ่มจำนวนการเรียกใช้ของคำถามนี้
			insForm, _ := db.Prepare("UPDATE collections SET count=? WHERE message=? ")
			insForm.Exec(count, text)

			//ส่งข้อความกลับไป ถ้ากรณีแรกคือ 1 คือเจอเลย
			fmt.Println("กรณีที่เจอข้อความในฐานข้อมูล")
			return 1, rawText, []model.ProductRow{}
		} else {
			//ในกรณีที่ไม่เจอ
			fmt.Println("ไม่เจอข้อความกำลังเข้าสู่กระบวนการหาคำตอบ...")
			fmt.Println()

			//ทำการตัดคำที่ไม่จำเป็นเช่น ครับ ค่ะ จะเก็บไว้ที่ test/stopword.txt
			text := test.CutStopWord(text)
			fmt.Println("ข้อความหลังการตัด Stop Word :​ ", text)
			fmt.Println()

			rawText, product := test.WordCosine(text)

			if (len(product) > 0) {
				//ถ้าเจอของ
				return 3, "", product
			} else if (strings.Compare(rawText, "") != 0 && len(product) == 0) {
				//ไม่เจอของแต่ว่าเป็นข้อความที่สามารถตอบกลับไปได้
				return 2, rawText, []model.ProductRow{}
			} else if (strings.Compare(rawText, "") == 0 || len(product) == 0) {
				//ไม่เจออะไรทั้งนั้น
				fmt.Println("ชิบหายแล้วววว")
				return 0, "", []model.ProductRow{}
			} else {
				fmt.Println("Test")
				return 2, rawText, []model.ProductRow{}
			}
		}
	}

}

func SaveWord(text string, Idcustomer string) {
	db, err := sql.Open("mysql", DATABASE)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	var ctx = context.Background()
	//fmt.Println(text)
	selectMessages, err := db.QueryContext(ctx, "SELECT message FROM oldmsg WHERE id=?", Idcustomer)
	rawText := ""

	for selectMessages.Next() {
		var tag Tag
		err = selectMessages.Scan(&tag.Feature)
		if err != nil {
			panic(err.Error())
		}
		rawText += tag.Feature
	}

	if (strings.Compare(rawText, "") == 0) {
		insForm, _ := db.Prepare("INSERT INTO oldmsg(id,message) VALUES (?,?)")
		_, err = insForm.Exec(Idcustomer, text)
	} else {
		insForm, _ := db.Prepare("UPDATE oldmsg SET message=? WHERE id=? ")
		insForm.Exec(text, Idcustomer)
	}

}

func AddProduct(name string, des string) int {
	db, err := sql.Open("mysql", DATABASE)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	addname, err := db.Prepare("INSERT INTO collections(message,types,answer,sub_feature,count) VALUES (?,?,?,?,?)")
	_, err = addname.Exec(name, "search","test","",0)

	adddes, err := db.Prepare("INSERT INTO collections(message,types,answer,sub_feature,count) VALUES (?,?,?,?,?)")
	_, err = adddes.Exec(des, "search","test","",0)

	if err != nil {
		return -1
	}else{
		return 1
	}
}
