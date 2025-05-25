package main

import "fmt"

func factorial(n int) int {
	if n == 0 || n == 1 {
		return 1
	}
	return n * factorial(n-1)
}

func main() {
	var result int
	var p int
	for i := 0; i < 1_000_000; i++ {
		result = factorial(20)
		p = i
	}
	fmt.Println(result)
	fmt.Println(p)
}
