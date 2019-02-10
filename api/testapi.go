package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	response,err:=http.Get("https://njmessengerbot.herokuapp.com/test")
	if err != nil {
		fmt.Println(err)
	}else{
		data,_:=ioutil.ReadAll(response.Body)
		fmt.Println(string(data))
	}
	
}
