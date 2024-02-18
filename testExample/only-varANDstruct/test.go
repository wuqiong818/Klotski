package main

import (
	"fmt"
)

var Hub = make(map[string]User, 10)

func main() {
	user1 := User{
		UserName: "111",
		Users:    &Hub,
	}
	user2 := User{
		UserName: "222",
		Users:    &Hub,
	}
	user3 := User{
		UserName: "333",
		Users:    &Hub,
	}

	Hub["1"] = user1
	Hub["2"] = user2
	Hub["3"] = user3
	var address = *user1.Users
	fmt.Println(address)
	fmt.Println(address["1"].Users)
}
