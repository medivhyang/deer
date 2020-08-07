package main

import (
	"fmt"
	"github.com/medivhyang/deer"
)

func main() {
	demoGet()
}

func demoGet() {
	resp, err := deer.Get("https://example.com").Do()
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.Text())
}

func demoPost() {
	value := map[string]string{
		"foo": "bar",
	}
	resp, err := deer.Post("https://example.com").WithJSONBody(value).Do()
	if err != nil {
		panic(err)
	}
	result := map[string]string{}
	if err = resp.BindWithJSON(&result); err != nil {
		panic(err)
	}
	fmt.Println(result)
}
