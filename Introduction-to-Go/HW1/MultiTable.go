package main

import "fmt"

func main() {
	var n int
	fmt.Scan(&n)
	MultiTable(n)
}

func MultiTable(n int) {
	for i := 1; i <= 10; i++ {
		fmt.Println(i, "*", n, "=", i * n)
	}
}
