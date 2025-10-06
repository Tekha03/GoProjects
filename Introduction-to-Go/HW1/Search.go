package main

import "fmt"

func main() {
	var lenth int
	fmt.Scan(&lenth)
	arr := make([]int, lenth)
	for i := 0; i < lenth; i++ {
		fmt.Scan(&arr[i])
	}
	var x int
	fmt.Scan(&x)
	fmt.Println(Search(arr, lenth, x))
}

func Search(arr []int, lenth int, x int) int {
	for i := 0; i < lenth; i++ {
		if arr[i] == x {
			return i + 1
		}
	}
	return -1
}
