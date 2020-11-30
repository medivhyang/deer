package main

import (
	"fmt"
	"github.com/medivhyang/deer"
)

func main() {
	t := &deer.RequestTemplate{}
	content, err := t.NewBuilder().GetText("https://baidu.com")
	if err != nil {
		panic(err)
	}
	fmt.Println(content)
}
