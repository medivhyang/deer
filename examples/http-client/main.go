package main

import (
	"fmt"
	"github.com/medivhyang/deer"
)

func main() {
	text, err := deer.GetText("https://baidu.com")
	if err != nil {
		panic(err)
	}
	fmt.Println(text)
}
