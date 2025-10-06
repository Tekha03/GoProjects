package main

import "fmt"

func main() {
	var n int
	fmt.Scan(&n)
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Scan(&arr[i])
	}
	fmt.Println(MinMax(arr, n))
}

func MinMax(arr []int, n int) (int, int) {
	mn := arr[0]
	mx := arr[0]
	for i := 0; i < n; i++ {
		mn = min(mn, arr[i])
		mx = max(mx, arr[i])
	}
	return mn, mx
}
