package main

import (
	"fmt"
	"github.com/medivhyang/deer"
)

func main() {
	resp, err := deer.Get("https://example.com").Do()
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.Text())
}
