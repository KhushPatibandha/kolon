package main

import "fmt"

func main() {
	var res int
	for i := 1; i <= 1000000; i++ {
		res += i
	}
	fmt.Println(res)
}
