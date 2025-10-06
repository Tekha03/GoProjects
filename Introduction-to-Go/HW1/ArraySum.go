package main

import "fmt"

func main() {
	var n int
	fmt.Scan(&n)
	arr := make([]int, n) 
	for i := 0; i < n; i++ {
		var a int
		fmt.Scan(&a)
		arr[i] = a
	}
	fmt.Println(ArraySum(arr, n))
}

func ArraySum(arr []int, size int) int {
	sum := 0
	for i := 0; i < size; i++ {
		sum += arr[i]
	}
	return sum
}
