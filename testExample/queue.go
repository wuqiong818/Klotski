package main

import (
	"fmt"
	"strconv"
)

func main() {
	atoi, err := strconv.Atoi("_")
	if err == nil {

		fmt.Println("success", atoi)
	} else {
		fmt.Println("err", err)
		fmt.Println("fail", atoi)
	}
}
