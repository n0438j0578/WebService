package test

import (
	"WebService/model"
	"context"
	"database/sql"
	"fmt"
	"github.com/narongdejsrn/go-thaiwordcut"
	"github.com/wilcosheh/tfidf"
	"github.com/wilcosheh/tfidf/similarity"
	"math/rand"
	"strings"

)

type IndexWord struct {
	//คำถามและเลขไอดีในฐานข้อมูล
	Msg string
	Id  string
}

func WordCosine(input string) (string, []model.ProductRow) {

	//ทำการต่อฐานข้อมูล
	db, err := sql.Open("mysql", DATABASE)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	//นำคำถามและเลขไอดีในฐานข้อมูลออกมา
	message, err := db.Query("SELECT message,id FROM collections")

	//เตรียมเครื่องมือในการตัดคำ
	segmenter := gothaiwordcut.Wordcut()
	segmenter.LoadDefaultDict()

	//เริ่มทำกระบวนการในทำค่า tfidf
	f := tfidf.New()

	//สร้างอินเด็กเอาไว้เก็บข้อมูลข้อความและเลขไอดี
	index := []IndexWord{}

	for message.Next() {
		var tag IndexWord

		err = message.Scan(&tag.Msg, &tag.Id)
		if err != nil {
			panic(err.Error())
		}

		//นำมาตัดแบ่งเพื่อให้เอาไปใส่ความถี่ได้เช่น ก็ และ หรือ มันก็จะเอาไปใส่ในเอกสารใหญ่ว่ามีคำนี้เพิ่มเข้ามาแล้วนะถ้าครั้งแรกผลลัพจะประมานว่า
		// ก็ : 1,และ : 2,หรือ : 3 เลขข้างหลังคือความถี่
		result := segmenter.Segment(tag.Msg)

		text := ""
		for i := 0; i < len(result); i++ {
			text = text + result[i] + " "
		}

		tag.Msg = text

		//ยัดใส้อินเด็ก
		index = append(index, tag)

		//ยัดใส่เอกสาร
		f.AddDocs(text)
	}

	//นำข้อความมาตัดและเอามาคำนวนว่าในแต่ละคำที่ตัดนั้นมีคะแนนเท่าไหร่
	t1 := input
	result := segmenter.Segment(t1)
	text := ""
	for i := 0; i < len(result); i++ {
		text = text + result[i] + " "
	}
	weightinput := f.Cal(text)
	fmt.Printf("ค่าน้ำหนักของ ของ %s เป็น %+v.\n", t1, weightinput)

	//ทำการเตรียมคำตอบให้กับข้อความในนี้โดยในตอนนี้ถ้าค่าของ cosine ออกมาเท่ากับ 0.5 จะแปลว่าหาข้อความไม่เจอไปก่อน
	type Answer struct {
		//คำถามและเลขไอดีในฐานข้อมูลและค่า cosine นั้นๆ
		Msg    string
		Id     string
		Cosine float64
	}
	//ใส่ค่าเปล่าให้มันก่อน และการจะหาค่ามากต้องทำให้ค่าเริ่มต้นมันน้อยๆ ก่อน
	var answer Answer
	answer.Msg = ""
	answer.Id = ""
	answer.Cosine = 0.0

	//วนลูปค่าตัวที่ค่า cosine มากที่สุด โดยจะเทียบกับที่รับเข้ามาเรื่อยๆ ตอนนี้ให้มันหาแค่อันเดียวก่อนจริงๆ ต้องหาอันค่ามากสุดสามอันแรก
	for i := 0; i < len(index); i++ {

		weight := f.Cal(index[i].Msg)
		sim := similarity.Cosine(weightinput, weight)
		if (sim > answer.Cosine) {
			answer.Msg = index[i].Msg
			answer.Id = index[i].Id
			answer.Cosine = sim
		}
	}

	//ถ้ามันไม่เจอเลย
	if answer.Cosine==0.50{
		return  "", []model.ProductRow{}
	}else{

		fmt.Println("คำตอบคือ : ", answer.Msg)
		fmt.Println("เลขไอดี : ", answer.Id)
		fmt.Println("Cosine : ", answer.Cosine)

		//เอาชนิดของข้อความนั้นออกมาดูว่าเป็นเกี่ยวกับ search หรือเปล่าเพราะต้องเอาออกไปเป็นคาลูเซล
		//เหมือนเป็นตัวแปรเอาไว้ใช้ในการใส่ฐานข้อมูลโดยสามารถส่งผ่านตัวแปรได้
		var ctx = context.Background()
		selectMessages, err := db.QueryContext(ctx, "SELECT types,answer FROM collections WHERE id=?", answer.Id)

		var types string
		var answerindata string

		for selectMessages.Next() {
			err = selectMessages.Scan(&types, &answerindata)
			if err != nil {
				panic(err.Error())
			}
			if strings.Compare(types, "search") == 0 {
				fmt.Println("เป็นประเภท : ", types)
			} else {
				fmt.Println("เป็นประเภท : ", types)
				fmt.Println("คำตอบคือ : ", answerindata)
			}

		}

		if strings.Compare(types, "search") == 0 {
			product := ProductMatching(input)
			//ส่งข้อมูลออกไปว่าเจอหรือไม่ ถ้าเจอขนาดของ product จะมากกว่า 1
			return "", product
		} else {

			cut := strings.Split(answerindata, ":;")
			//ตัวแปรเปล่าเอาไว้เก็บข้อความที่ตัด
			var rawText string

			//ถ้าขนาดไม่เท่ากับหนึ่งแปลว่ามีหลายคำตอบ
			if (len(cut) != 1) {
				rawText = cut[rand.Intn(len(cut)-1)]
				for ; ; {
					//ลูบนี้จะเช็คประมานว่าถ้าตัดแล้วเจอข้อความเปล่าๆ ให้มันแรนดอมเอาใหม่อีกครั้ง
					if strings.Compare(rawText, "") == 0 {
						cut := strings.Split(answerindata, ":;")
						rawText = cut[rand.Intn(len(cut)-1)]
					} else {
						break
					}
				}
			} else {
				rawText = answerindata
			}

			insForm, err := db.Prepare("INSERT INTO collections(message,types,answer,sub_feature,count,greeting,problem,orders,search) VALUES (?,?,?,?,?,?,?,?,?)")
			if err != nil {
				panic(err.Error())
			}
			//result คือข้อความที่ตัดแล้วดูที่ บรรทัด 71
			_, err = insForm.Exec(input, types, answerindata, result, 0, 0, 0, 0, 0)

			return rawText, []model.ProductRow{}
		}
	}


}
