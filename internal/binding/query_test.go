package binding

import (
	"fmt"
	"time"
)

func ExampleBindWithQuery() {
	type User struct {
		ID          string
		Name        string
		Age         int
		Likes       []string
		CreatedTime time.Time `time_format:"unix"`
	}
	var user User
	if err := BindWithQuery(&user, map[string][]string{
		"id":           {"1"},
		"name":         {"Alice"},
		"age":          {"20"},
		"likes":        {"cat", "dog", "deer"},
		"created_time": {fmt.Sprintf("%d", time.Date(2020, 9, 18, 0, 0, 0, 0, time.UTC).Unix())},
	}); err != nil {
		panic(err)
	}
	fmt.Printf("%+v", user)

	// Output:
	// {ID:1 Name:Alice Age:20 Likes:[cat dog deer] CreatedTime:2020-09-18 08:00:00 +0800 CST}
}
