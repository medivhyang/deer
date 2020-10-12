package main

import (
	"github.com/medivhyang/deer"
)

func main() {
	//text, err := deer.GetText("https://baidu.com")
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(text)
	demoGetFile()
}

func demoGetFile() {
	if err := deer.GetFile("https://baidu.com", "baidu.html"); err != nil {
		panic(err)
	}
}
