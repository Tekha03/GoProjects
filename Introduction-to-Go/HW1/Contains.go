package main

import "fmt"

func main() {
	var length int
	fmt.Scan(&length)
	arr := make([]float32, length)
	for i := 0; i < length; i++ {
		var x float32
		fmt.Scan(&x)
		arr[i] = x
	}
	var n float32
	fmt.Scan(&n)
	if Contains(arr, length, n) {
		fmt.Println("Contains")
	} else {
		fmt.Println("Doesn't contain")
	}
}

func Contains(arr []float32, length int, n float32) bool {
	for i := 0; i < length; i++ {
		if arr[i] == n {
			return true
		}
	}
	return false
}
