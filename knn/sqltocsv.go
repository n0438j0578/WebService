package main

import (
	"WebService/controller"
	"bytes"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // or the driver of your choice
	"github.com/joho/sqltocsv"
	"io"
	"os"
)


func popLine(f *os.File) ([]byte, error) {
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(make([]byte, 0, fi.Size()))

	_, err = f.Seek(0, os.SEEK_SET)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(buf, f)
	if err != nil {
		return nil, err
	}
	line, err := buf.ReadString('\n')
	if err != nil && err != io.EOF {
		return nil, err
	}

	_, err = f.Seek(0, os.SEEK_SET)
	if err != nil {
		return nil, err
	}
	nw, err := io.Copy(f, buf)
	if err != nil {
		return nil, err
	}
	err = f.Truncate(nw)
	if err != nil {
		return nil, err
	}
	err = f.Sync()
	if err != nil {
		return nil, err
	}

	_, err = f.Seek(0, os.SEEK_SET)
	if err != nil {
		return nil, err
	}
	return []byte(line), nil
}
func main(){
	db, err := sql.Open("mysql", controller.DATABASE)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	rows, _ := db.Query("SELECT greeting ,problem ,orders ,search,types  FROM collections")

	err = sqltocsv.WriteFile("./test/report.csv", rows)

	//fname := `report.csv`
	f, err := os.OpenFile("./test/report.csv", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	line, err := popLine(f)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("pop:", string(line))

	if err != nil {
		panic(err)
	}else{
		fmt.Println("เสร็จแล้ว")
	}

}
func SqltoCsv(){
	db, err := sql.Open("mysql", controller.DATABASE)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	rows, _ := db.Query("SELECT greeting ,problem ,orders ,search,types  FROM collections")

	err = sqltocsv.WriteFile("./knn/report.csv", rows)

	//fname := `report.csv`
	f, err := os.OpenFile("./knn/report.csv", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	line, err := popLine(f)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("pop:", string(line))

	if err != nil {
		panic(err)
	}else{
		fmt.Println("เสร็จแล้ว")
	}

}