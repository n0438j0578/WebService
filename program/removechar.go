package main

import (
	"fmt"
	"log"
	"regexp"
)

func main() {

	example := "#GoLangCode!$!"

	// Make a Regex to say we only want letters and numbers
	reg, err := regexp.Compile("[^0-9:]+")
	if err != nil {
		log.Fatal(err)
	}

	processedString := reg.ReplaceAllString(example, "")

	fmt.Printf("A string of %s becomes %s \n", example, processedString)

	//เอา 77:1 ค่ะ
	test := "เอา 77 : 1 ค่ะ"

	processedString = reg.ReplaceAllString(test, "")
	fmt.Println(processedString)



	test = "77:1"

	processedString = reg.ReplaceAllString(test, "")
	fmt.Println(processedString)
}
