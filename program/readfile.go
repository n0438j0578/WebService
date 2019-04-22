package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func main() {

	dat, err := ioutil.ReadFile("program/ignoreWord.txt")
	if err != nil {
		fmt.Print(err)
	}

	str := string(dat)
	fmt.Println(str)
	fmt.Println()


	input := "และ"
	fmt.Print("Result : ")
	if strings.Contains(str, input) {
		fmt.Println("Matched")
	} else {
		fmt.Println("Not match")
	}


	
}
