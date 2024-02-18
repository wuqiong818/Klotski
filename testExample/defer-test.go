package main

import "fmt"

func main() {
	num := 5

	if num == 3 {
		defer func() {

			fmt.Println("defer", 333)
		}()
	} else {
		fmt.Println("未执行")
	}
}
